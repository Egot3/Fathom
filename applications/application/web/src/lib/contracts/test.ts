import { errAsync, fromPromise, okAsync, ResultAsync } from "neverthrow";
import type { JSONError } from "../statuses/jsonerror";
import type { ParsedQuiz, ParsedQuizResponse, Quiz } from "./quiz";
import {
  NormalizeBlob,
  NormalizeJSON,
  NormalizeResponse,
  TokenizedFetch,
} from "./tokenizedFetch";
import { AcceptTransform, type Accept } from "../apiutils/acceptHeader";
import type { Result } from "postcss";
import {
  EvictQuiz,
  GetCachedQuiz,
  ReadManifest,
  SetCachedQuiz,
} from "../apiutils/quizManifest";

export type Test = {
  uuid: string;
  name: string;
  created_at: Date;
  updated_at: Date;
  quizzes: Quiz[];
};

// re-re-redone
export function FetchTest(testUUID: string): ResultAsync<Test, JSONError> {
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/${testUUID}`,
    {
      method: "GET",
      headers: {
        Accept: "application/json",
      },
    },
  ).andThen(NormalizeJSON<Test>);
}

export type TestInfo = Test & {
  deadline: Date;
  is_paused: boolean;
  key: string;
};

type CurrentlyRunningResponse = {
  tests: TestInfo[];
};

export function FetchCurrentlyRunningTestInfos(): ResultAsync<
  TestInfo[],
  JSONError
> {
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/running/all`,
  ).andThen((r) => {
    if (r.status === 423) {
      return okAsync([]);
    }
    return NormalizeJSON<CurrentlyRunningResponse>(r).andThen((r) =>
      okAsync(r.tests),
    );
  });
}

export type Tests = { tests: Test[]; total: number };

export function FetchAllTests(
  page: number,
  size: number,
): ResultAsync<Tests, JSONError> {
  return TokenizedFetch(
    "https://" +
      import.meta.env.VITE_DOMAIN +
      "/api/v1/test/" +
      "?page=" +
      page +
      "&size=" +
      size,
    {
      method: "GET",
      headers: {
        Accept: "application/json",
      },
    },
  ).andThen(NormalizeJSON<Tests>);
}

type PostTestRequest = {
  name: string;
  quizzes: string[];
};

export function FetchTestPost(
  name: string,
  quizzes: string[],
): ResultAsync<null, JSONError> {
  const body: PostTestRequest = {
    name: name,
    quizzes: quizzes,
  };
  const bodyString = JSON.stringify(body);

  return TokenizedFetch(
    "https://" + import.meta.env.VITE_DOMAIN + "/api/v1/test/",
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: bodyString,
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

export function FetchTestDelete(UUID: string): ResultAsync<null, JSONError> {
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/${UUID}`,
    {
      method: "DELETE",
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

type PatchTestRequest = {
  name: string | undefined;
};

export function FetchTestPatch(
  UUID: string,
  name: string | undefined,
): ResultAsync<null, JSONError> {
  const body: PatchTestRequest = {
    name: name,
  };

  const bodyString = JSON.stringify(body);

  return TokenizedFetch(
    "https://" + import.meta.env.VITE_DOMAIN + "/api/v1/test/" + UUID,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: bodyString,
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

type TestBundleRequest = {
  quiz_uuids: string[];
};

export function FetchTestBundle(
  testUUID: string,
  quizUUIDs: string[],
): ResultAsync<null, JSONError> {
  const body: TestBundleRequest = {
    quiz_uuids: quizUUIDs,
  };

  const bodyString = JSON.stringify(body);

  return TokenizedFetch(
    "https://" +
      import.meta.env.VITE_DOMAIN +
      "/api/v1/test/" +
      testUUID +
      "/quizzes",
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: bodyString,
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

export function FetchTestPrune(
  testUUID: string,
  quizUUIDs: string[],
): ResultAsync<null, JSONError> {
  const body: TestBundleRequest = {
    quiz_uuids: quizUUIDs,
  };

  const bodyString = JSON.stringify(body);

  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/${testUUID}/quizzes`,
    {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: bodyString,
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

type TestStartRequest = {
  duration: string;
  test_uuid: string;
  group_uuids: string[];
};

/** @param duration - time in go's duration string. ex: 1h30m6s. Ex on what it can't understand: 2:05, PT1H30M6S, green
 */
export function FetchTestStart(
  UUID: string,
  duration: string,
  groupUUIDs: string[],
): ResultAsync<null, JSONError> {
  const body: string = JSON.stringify({
    test_uuid: UUID,
    duration: duration,
    group_uuids: groupUUIDs,
  } as TestStartRequest);
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/running/start`,
    {
      method: "POST",
      body: body,
      headers: {
        "Content-Type": "application/json",
        Accept: "*/*;q=0", //those who nose
      },
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

export function FetchTestPause(key: string): ResultAsync<null, JSONError> {
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/running/${key}/pause`,
    {
      method: "POST",
      headers: {
        Accept: "*/*;q=0", //those who nose
      },
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

export function FetchTestResume(key: string): ResultAsync<null, JSONError> {
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/running/${key}/resume`,
    {
      method: "POST",
      headers: {
        Accept: "*/*;q=0", //those who nose
      },
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

export function FetchTestStop(key: string): ResultAsync<null, JSONError> {
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/running/${key}/stop`,
    {
      method: "POST",
      headers: {
        Accept: "*/*;q=0",
      },
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

type TestExtendRequest = {
  extend_by: string;
};

export function FetchTestExtend(
  key: string,
  duration: string,
): ResultAsync<null, JSONError> {
  const body: TestExtendRequest = {
    extend_by: duration,
  };

  const bodyString = JSON.stringify(body);
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/running/${key}/extend`,
    {
      method: "POST",
      body: bodyString,
      headers: {
        Accept: "*/*;q=0",
      },
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

export function FetchTestExport(
  testUUID: string,
  accepts: Accept[],
): ResultAsync<Blob, JSONError> {
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/${testUUID}/export`,
    {
      method: "GET",
      headers: {
        Accept: AcceptTransform(accepts),
        "Content-Type": "application/x-www-form-urlencoded",
      },
    },
  ).andThen(NormalizeBlob);
}

export function FetchTestImport(file: File): ResultAsync<null, JSONError> {
  const formData = new FormData();

  formData.append("imported", file, file.name);

  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/import`,
    {
      method: "POST",
      body: formData,
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

type GetQuizzesUUIDsResponse = {
  quiz_uuids: string[];
};

export function FetchCurrentyRunningTestQuizzes(
  key: string,
): ResultAsync<string[], JSONError> {
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/running/${key}/quizzes`,
  )
    .andThen(NormalizeJSON<GetQuizzesUUIDsResponse>)
    .andThen(({ quiz_uuids }) => okAsync(quiz_uuids));
}

export function FetchCurrentlyRunningTestQuiz(
  key: string,
  UUID: string,
): ResultAsync<ParsedQuiz, JSONError> {
  const manifest = ReadManifest();
  const known = manifest[UUID];
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/running/${key}/${UUID}`,
    {
      headers: known ? { "If-None-Match": `"${known.checksum}"` } : {},
    },
  ).andThen((r) => {
    if (r.status === 304 && known) {
      const cached = GetCachedQuiz<ParsedQuiz>(UUID, known.checksum);
      if (cached !== null) return okAsync(cached);
    }
    if (r.status === 404) {
      EvictQuiz(UUID);
      return errAsync<ParsedQuiz, JSONError>({
        error: `Quiz ${UUID} not found`,
      });
    }

    return NormalizeJSON<ParsedQuizResponse>(r)
      .andTee((body: ParsedQuizResponse) => {
        const etag = r.headers.get("ETag")?.replace(/"/g, "") ?? "";
        console.log("ETag:", etag);
        SetCachedQuiz(UUID, etag, body.quiz);
      })
      .andThen(({ quiz }) => okAsync(quiz));
  });
}
