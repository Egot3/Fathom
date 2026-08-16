const MANIFEST_KEY = "fathom:cache:manifest";
const ENTRY_PREFIX = "fathom:quiz:";
const SOFT_CAP_BYTES = 3 * 1024 * 1024;
const EVICT_TARGET_RATIO = 0.8;

interface ManifestEntry {
  checksum: string;
  size: number;
  lastAccessed: number;
}

type Manifest = Record<string, ManifestEntry>;

export function ReadManifest(): Manifest {
  try {
    const raw = localStorage.getItem(MANIFEST_KEY);
    return raw ? (JSON.parse(raw) as Manifest) : {};
  } catch {
    return {}; // corrupted means... corupted, duh
  }
}

function writeManifest(manifest: Manifest): void {
  localStorage.setItem(MANIFEST_KEY, JSON.stringify(manifest));
}

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
): T | null {
  const manifest = ReadManifest();
  const entry = manifest[quizUUID];
  if (!entry || entry.checksum !== currentChecksum) return null;

  const raw = localStorage.getItem(ENTRY_PREFIX + quizUUID);
  if (raw === null) {
    delete manifest[quizUUID];
    writeManifest(manifest);
    return null;
  }

  entry.lastAccessed = Date.now();
  writeManifest(manifest);
  return JSON.parse(raw) as T;
}

export function SetCachedQuiz(
  quizUUID: string,
  checksum: string,
  data: unknown,
): void {
  const json = JSON.stringify(data);
  const size = new Blob([json]).size;
  const manifest = ReadManifest();

  const write = () => {
    localStorage.setItem(ENTRY_PREFIX + quizUUID, json);
    manifest[quizUUID] = { checksum, size, lastAccessed: Date.now() };
    writeManifest(manifest);
  };

  try {
    write();
  } catch (err) {
    if (err instanceof DOMException && err.name === "QuotaExceededError") {
      evictUntilUnder(manifest, totalSize(manifest) * 0.5);
      write();
    } else {
      throw err;
    }
  }

  if (totalSize(manifest) > SOFT_CAP_BYTES) {
    evictUntilUnder(manifest, SOFT_CAP_BYTES * EVICT_TARGET_RATIO);
  }
}

export function EvictQuiz(quizUUID: string): void {
  const manifest = ReadManifest();
  if (manifest[quizUUID]) {
    localStorage.removeItem(ENTRY_PREFIX + quizUUID);
    delete manifest[quizUUID];
    writeManifest(manifest);
  }
}
