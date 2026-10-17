import { errAsync, fromPromise, okAsync, ResultAsync } from "neverthrow";
import {
  GetCurrentlyRunning,
  GetCurrentlyRunningCaching,
  SetCurrentlyRunningCaching,
  type ETagInfo,
} from "../bgdata/currentlyrunning.svelte";
import type { JSONError } from "../statuses/jsonerror";
import type { Quiz } from "./quiz";
import {
  maxAgeRegex,
  NormalizeBlob,
  NormalizeJSON,
  NormalizeResponse,
  TokenizedFetch,
} from "./tokenizedFetch";
import { AcceptTransform, type Accept } from "../apiutils/acceptHeader";
import type { Result } from "postcss";

export type Test = {
  uuid: string;
  name: string;
  created_at: Date;
  updated_at: Date;
  quizzes: Quiz[];
};

type GetTestResponse = {
  test: Test;
};

type GetQuizUUIDsResponse = {
  quiz_uuids: string[];
};

// redone
export function FetchTest(testUUID: string): ResultAsync<Test, JSONError> {
  return ResultAsync.fromPromise(
    TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/${testUUID}`,
      {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
      },
    ),
    (err) => {
      console.log("Couldn't fetch test: ", err);
      if (err instanceof Error) {
        return {
          error: "couldn't fetch test because of in-browser error",
        };
      }
      return { error: "couldn't fetch test because of unknown error" };
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
  return ResultAsync.fromPromise(
    TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/running/all`,
    ),
    (err): JSONError => {
      console.log(
        "couldn't fetch current test info due to unknown error: ",
        err,
      );
      return {
        error: "got network error while fetching current test",
      } as JSONError;
    },
  ).andThen((r) => {
    if (r.status === 423) {
      return okAsync([]);
    }
    return NormalizeJSON<CurrentlyRunningResponse>(r).andThen((r) =>
      okAsync(r.tests),
    );
  });
}

// export function FetchCurrentlyRunningQuizUUIDs(): ResultAsync<
//   string[],
//   JSONError
// > {
//   const ETagInfo = GetCurrentlyRunningCaching();
//   return ResultAsync.fromPromise(
//     TokenizedFetch(
//       `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/running/quizzes`,
//       {
//         headers: ETagInfo
//           ? ETagInfo.ExpiresAt > new Date()
//             ? { "If-None-Match": `"${ETagInfo.ETag}"` }
//             : {}
//           : {},
//       },
//     ),
//     (err): JSONError => {
//       console.log("Couldn't fetch running test: ", err);
//       if (err instanceof Error) {
//         return {
//           error: "couldn't fetch running test because of in-browser error",
//         };
//       }
//       return { error: "couldn't fetch running test because of unknown error" };
//     },
//   ).andThen((r) => {
//     if (r.status === 304 && ETagInfo) {
//       const cached = GetCurrentlyRunning();
//       if (cached !== null) return okAsync(cached);
//     }
//     if (r.status === 423) {
//       SetCurrentlyRunningCaching({
//         ETag: "",
//         ExpiresAt: new Date(Date.now() - 100), // the Date.now() - x is used because I am pretty sure someone is still living in the 50-s
//       });
//     }
//     if (!r.ok) {
//       return ResultAsync.fromPromise(r.json(), (): JSONError => ({
//         error: "couldn't parse error body",
//       })).andThen((body) => errAsync<string[], JSONError>(body as JSONError));
//     }

//     return ResultAsync.fromPromise(r.json(), (): JSONError => ({
//       error: "couldn't parse response body",
//     })).andThen((r) => {
//       const etag = r.headers.get("ETag")?.replace(/"/g, "") ?? "";
//       const cch = r.headers.get("Cache-Control");
//       const reg = maxAgeRegex.exec(cch);
//       if (reg === null || reg.length < 2) {
//         console.log("couldn't cache the response. Got header: ", cch);
//         return okAsync(r);
//       }
//       const offset = parseInt(reg[1], 10);

//       SetCurrentlyRunningCaching({
//         ETag: etag,
//         ExpiresAt: new Date(Date.now() + offset * 1000),
//       });

//       return okAsync(r);
//     });
//   });
// }

export type Tests = { tests: Test[]; total: number };

export function FetchAllTests(
  page: number,
  size: number,
): ResultAsync<Tests, JSONError> {
  return ResultAsync.fromPromise(
    TokenizedFetch(
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
    ),
    (err): JSONError => {
      console.log("Couldn't list tests: ", err);
      if (err instanceof Error) {
        return {
          error: "couldn't list tests because of in-browser error",
        };
      }

      return { error: "couldn't list tests because of unknown error" };
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

  return ResultAsync.fromPromise(
    TokenizedFetch("https://" + import.meta.env.VITE_DOMAIN + "/api/v1/test/", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: bodyString,
    }),
    (err): JSONError => {
      console.log("Couldn't post tests: ", err);
      if (err instanceof Error) {
        return {
          error: "couldn't post tests because of in-browser error",
        };
      }

      return { error: "couldn't post tests because of unknown error" };
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

export function FetchTestDelete(UUID: string): ResultAsync<null, JSONError> {
  return ResultAsync.fromPromise(
    TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/${UUID}`,
      {
        method: "DELETE",
      },
    ),
    (err) => {
      console.log("Couldn't fetch test delete: ", err);
      if (err instanceof Error) {
        return {
          error: "couldn't send test delete because of in-browser error",
        };
      }

      return { error: "couldn't send test delete because of unknown error" };
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

  return ResultAsync.fromPromise(
    TokenizedFetch(
      "https://" + import.meta.env.VITE_DOMAIN + "/api/v1/test/" + UUID,
      {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: bodyString,
      },
    ),
    (err) => {
      console.log("Couldn't fetch test delete: ", err);
      if (err instanceof Error) {
        return {
          error: "couldn't send test patch because of in-browser error",
        };
      }

      return { error: "couldn't send test patch because of unknown error" };
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

  return ResultAsync.fromPromise(
    TokenizedFetch(
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
    ),
    (err) => {
      console.log("Couldn't fetch test bundle: ", err);
      if (err instanceof Error) {
        return {
          error: "couldn't send test bundle because of in-browser error",
        };
      }

      return { error: "couldn't send test bundle because of unknown error" };
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

  return ResultAsync.fromPromise(
    TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/${testUUID}/quizzes`,
      {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: bodyString,
      },
    ),
    (err) => {
      console.log("Couldn't prune test: ", err);
      if (err instanceof Error) {
        return {
          error: "couldn't fetch prune test because of in-browser error",
        };
      }
      return { error: "couldn't fetch prune test because of unknown error" };
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
  return ResultAsync.fromPromise(
    TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/running/start`,
      {
        method: "POST",
        body: body,
        headers: {
          "Content-Type": "application/json",
          Accept: "*/*;q=0", //those who nose
        },
      },
    ),
    (err): JSONError => {
      console.log("Couldn't fetch test start: ", err);
      if (err instanceof Error) {
        return {
          error: "couldn't send test start because of in-browser error",
        };
      }

      return { error: "couldn't send test start because of unknown error" };
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

export function FetchTestPause(key: string): ResultAsync<null, JSONError> {
  return ResultAsync.fromPromise(
    TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/running/${key}/pause`,
      {
        method: "POST",
        headers: {
          Accept: "*/*;q=0", //those who nose
        },
      },
    ),
    (err): JSONError => {
      console.log("Couldn't fetch test pause: ", err);
      if (err instanceof Error) {
        return {
          error: "couldn't send test pause because of in-browser error",
        };
      }

      return { error: "couldn't send test pause because of unknown error" };
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

export function FetchTestResume(key: string): ResultAsync<null, JSONError> {
  return ResultAsync.fromPromise(
    TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/running/${key}/resume`,
      {
        method: "POST",
        headers: {
          Accept: "*/*;q=0", //those who nose
        },
      },
    ),
    (err): JSONError => {
      console.log("Couldn't fetch test resume for quiz: ", err);
      if (err instanceof Error) {
        return {
          error: "couldn't send test resume because of in-browser error",
        };
      }

      return { error: "couldn't send test resume because of unknown error" };
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

export function FetchTestStop(key: string): ResultAsync<null, JSONError> {
  return ResultAsync.fromPromise(
    TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/running/${key}/stop`,
      {
        method: "POST",
        headers: {
          Accept: "*/*;q=0",
        },
      },
    ),
    (err): JSONError => {
      console.log("Couldn't fetch test stop for quiz: ", err);
      if (err instanceof Error) {
        return {
          error: "couldn't send test stop because of in-browser error",
        };
      }

      return { error: "couldn't send test stop because of unknown error" };
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
  return ResultAsync.fromPromise(
    TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/running/${key}/extend`,
      {
        method: "POST",
        body: bodyString,
        headers: {
          Accept: "*/*;q=0",
        },
      },
    ),
    (err): JSONError => {
      console.log("Couldn't fetch test extend for quiz: ", err);
      if (err instanceof Error) {
        return {
          error: "couldn't send test extend because of in-browser error",
        };
      }

      return { error: "couldn't send test extend because of unknown error" };
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

export function FetchTestExport(
  testUUID: string,
  accepts: Accept[],
): ResultAsync<Blob, JSONError> {
  return ResultAsync.fromPromise(
    TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/${testUUID}/export`,
      {
        method: "GET",
        headers: {
          Accept: AcceptTransform(accepts),
          "Content-Type": "application/x-www-form-urlencoded",
        },
      },
    ),
    (err) => {
      console.log("Couldn't list answers: ", err);
      if (err instanceof Error) {
        return {
          error: "couldn't fetch list answers because of in-browser error",
        };
      }
      return { error: "couldn't fetch list answers because of unknown error" };
    },
  ).andThen(NormalizeBlob);
}

export function FetchTestImport(file: File): ResultAsync<null, JSONError> {
  const formData = new FormData();

  formData.append("imported", file, file.name);

  return ResultAsync.fromPromise(
    TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/import`,
      {
        method: "POST",
        body: formData,
      },
    ),
    (err) => {
      console.log("Couldn't import: ", err);
      if (err instanceof Error) {
        return {
          error: "couldn't fetch import test because of in-browser error",
        };
      }
      return { error: "couldn't fetch test import because of unknown error" };
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}
