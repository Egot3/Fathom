import { err, ok, Result } from "neverthrow";
import type { JSONError } from "../statuses/jsonerror";

const MANIFEST_KEY = "fathom:cache:manifest";
const ENTRY_PREFIX = "fathom:quiz:";
const SOFT_CAP_BYTES = 3 * 1024 * 1024;
const EVICT_TARGET_RATIO = 0.8;

interface ManifestEntry {
  checksum: string;
  size: number;
  lastAccessed: number;
}

export type Manifest = Record<string, ManifestEntry>;

function safeParse<T = unknown>(raw: string): Result<T, JSONError> {
  return Result.fromThrowable(
    (r): T => JSON.parse(r),
    (): JSONError => ({
      error: "Quiz manifest parse error",
    }),
  )(raw);
}

export function ReadManifest(): Result<Manifest, JSONError> {
  const raw = localStorage.getItem(MANIFEST_KEY);
  if (raw === null) {
    return err({ error: "unable to retrive manifest" } as JSONError);
  }
  return safeParse<Manifest>(raw);
}

type cacheError = { kind: "quota" } | { kind: "storage"; cause: unknown };

const classify = (e: unknown): cacheError =>
  e instanceof DOMException && e.name === "QuotaExceededError"
    ? { kind: "quota" }
    : { kind: "storage", cause: e };

const safeSet = Result.fromThrowable(
  (key: string, value: string) => localStorage.setItem(key, value),
  classify,
);

const writeManifest = (manifest: Manifest): Result<void, cacheError> =>
  safeSet(MANIFEST_KEY, JSON.stringify(manifest));

function totalSize(manifest: Manifest): number {
  return Object.values(manifest).reduce((sum, e) => sum + e.size, 0);
}

function evictUntilUnder(manifest: Manifest, targetBytes: number): void {
  const oldestFirst = Object.entries(manifest).sort(
    (a, b) => a[1].lastAccessed - b[1].lastAccessed,
  );

  let size = totalSize(manifest);
  for (const [quizUUID, entry] of oldestFirst) {
    if (size <= targetBytes) break;
    localStorage.removeItem(ENTRY_PREFIX + quizUUID);
    delete manifest[quizUUID];
    size -= entry.size;
  }
}

export function GetCachedQuiz<T = unknown>(
  quizUUID: string,
  currentChecksum: string,
): Result<T, JSONError> {
  return ReadManifest()
    .andThrough((manifest) => {
      const entry = manifest[quizUUID];
      if (!entry || entry.checksum !== currentChecksum)
        return err({ error: "Needed quiz wasn't found" } as JSONError);
      return ok(undefined);
    })
    .andThen((manifest) => {
      const raw = localStorage.getItem(ENTRY_PREFIX + quizUUID);
      if (raw === null) {
        delete manifest[quizUUID];

        // which is strange, async writes don't happen, and written is smaller, yet golang mask is calling for me. Errors should be handled explicitly
        return writeManifest(manifest)
          .mapErr((e) => {
            console.error(
              "Unable to write quiz because of local storage error",
              e,
            );
            return {
              error: "unable to write quiz because of local storage error",
            } as JSONError;
          })
          .andThen((_) =>
            err({ error: "Needed quiz couldn't be retrieved" } as JSONError),
          );
      }

      manifest[quizUUID].lastAccessed = Date.now();
      return writeManifest(manifest)
        .mapErr((e) => {
          console.error(
            "Unable to write quiz because of local storage error",
            e,
          );
          return {
            error: "unable to write quiz because of local storage error",
          } as JSONError;
        })
        .andThen(() => ok(raw));
    })
    .andThen(safeParse<T>);
}

export function SetCachedQuiz(
  quizUUID: string,
  checksum: string,
  data: unknown,
): Result<void, JSONError> {
  return Result.fromThrowable(() => JSON.stringify(data), classify)()
    .andThen((json) => {
      const size = new Blob([json]).size;

      return ReadManifest()
        .orElse(() => ok({} as Manifest))
        .andThen((manifest) => {
          const attempt = () =>
            safeSet(ENTRY_PREFIX + quizUUID, json)
              .andThen(() => {
                manifest[quizUUID] = {
                  checksum,
                  size,
                  lastAccessed: Date.now(),
                };
                return writeManifest(manifest);
              })
              .map(() => manifest);

          return attempt()
            .orElse((e) => {
              if (e.kind !== "quota") return err(e);
              evictUntilUnder(manifest, totalSize(manifest) / 2);
              return attempt();
            })
            .andThen((m) => {
              if (totalSize(m) <= SOFT_CAP_BYTES) return ok(undefined);
              evictUntilUnder(manifest, SOFT_CAP_BYTES * EVICT_TARGET_RATIO);
              return writeManifest(m);
            });
        });
    })
    .orTee((e) => console.error("Couldn't write quiz to cache:", e))
    .mapErr((): JSONError => ({ error: "Couldn't write quiz to cache" }));
}

export function EvictQuiz(quizUUID: string): Result<void, JSONError> {
  return ReadManifest().andThen((manifest) => {
    if (manifest[quizUUID]) {
      localStorage.removeItem(ENTRY_PREFIX + quizUUID);
      delete manifest[quizUUID];
      return writeManifest(manifest).mapErr((e) => {
        console.error("Unable to write quiz because of local storage error", e);
        return {
          error: "unable to write quiz because of local storage error",
        } as JSONError;
      });
    }
    return ok(undefined);
  });
}
