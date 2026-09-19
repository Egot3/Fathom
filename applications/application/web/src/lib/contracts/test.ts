import { errAsync, fromPromise, okAsync, ResultAsync } from "neverthrow";
import {
  GetCurrentlyRunning,
  GetCurrentlyRunningCaching,
  SetCurrentlyRunningCaching,
  type ETagInfo,
} from "../bgdata/currentlyrunning.svelte";
import type { JSONError } from "../statuses/jsonerror";
import type { Quiz } from "./quiz";
import { maxAgeRegex, TokenizedFetch } from "./tokenizedFetch";

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
export async function FetchTest(testUUID: string): Promise<Test | JSONError> {
  try {
    const rawRes = await TokenizedFetch(
      "https://" + import.meta.env.VITE_DOMAIN + "/api/v1/test/" + testUUID,
      {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
      },
    );

    if (!rawRes.ok) {
      return (await rawRes.json()) as JSONError;
    }

    return ((await rawRes.json()) as GetTestResponse).test;
  } catch (err) {
    console.log("couldn't fetch current test info due to unknown error: ", err);
    return {
      error: "got network error while fetching current test",
    } as JSONError;
  }
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
    if (!r.ok) {
      if (r.status === 423) {
        return okAsync([]);
      }
      return ResultAsync.fromPromise(
        r.json() as Promise<JSONError>,
        (err): JSONError => {
          console.log("couldn't parse error body: ", err);
          return { error: "couldn't parse error body" };
        },
      ).andThen((body) => {
        return errAsync<TestInfo[], JSONError>(body as JSONError);
      });
    }
    return ResultAsync.fromPromise(
      r.json() as Promise<CurrentlyRunningResponse>,
      (err): JSONError => {
        console.log("couldn't parse response body: ", err);
        return { error: "couldn't parse response body" };
      },
    ).andThen((r) => okAsync(r.tests));
  });
}

export function FetchCurrentlyRunningQuizUUIDs(): ResultAsync<
  string[],
  JSONError
> {
  const ETagInfo = GetCurrentlyRunningCaching();
  return ResultAsync.fromPromise(
    TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/running/quizzes`,
      {
        headers: ETagInfo
          ? ETagInfo.ExpiresAt > new Date()
            ? { "If-None-Match": `"${ETagInfo.ETag}"` }
            : {}
          : {},
      },
    ),
    (err): JSONError => {
      console.log("Couldn't fetch running test: ", err);
      if (err instanceof Error) {
        return {
          error: "couldn't fetch running test because of in-browser error",
        };
      }
      return { error: "couldn't fetch running test because of unknown error" };
    },
  ).andThen((r) => {
    if (r.status === 304 && ETagInfo) {
      const cached = GetCurrentlyRunning();
      if (cached !== null) return okAsync(cached);
    }
    if (r.status === 423) {
      SetCurrentlyRunningCaching({
        ETag: "",
        ExpiresAt: new Date(Date.now() - 100), // the Date.now() - x is used because I am pretty sure someone is still living in the 50-s
      });
    }
    if (!r.ok) {
      return ResultAsync.fromPromise(r.json(), (): JSONError => ({
        error: "couldn't parse error body",
      })).andThen((body) => errAsync<string[], JSONError>(body as JSONError));
    }

    return ResultAsync.fromPromise(r.json(), (): JSONError => ({
      error: "couldn't parse response body",
    })).andThen((r) => {
      const etag = r.headers.get("ETag")?.replace(/"/g, "") ?? "";
      const cch = r.headers.get("Cache-Control");
      const reg = maxAgeRegex.exec(cch);
      if (reg === null || reg.length < 2) {
        console.log("couldn't cache the response. Got header: ", cch);
        return okAsync(r);
      }
      const offset = parseInt(reg[1], 10);

      SetCurrentlyRunningCaching({
        ETag: etag,
        ExpiresAt: new Date(Date.now() + offset * 1000),
      });

      return okAsync(r);
    });
  });
}

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
  ).andThen((r) => {
    if (!r.ok) {
      return ResultAsync.fromPromise(
        r.json() as Promise<JSONError>,
        (): JSONError => ({
          error: "couldn't parse error body",
        }),
      ).andThen((body) => errAsync(body as JSONError));
    }

    return ResultAsync.fromPromise(r.json(), (): JSONError => ({
      error: "couldn't parse response body",
    })).andThen((body: Tests) => okAsync(body));
  });
}

type PostTestRequest = {
  name: string;
  quizzes: string[];
};

export async function FetchTestPost(
  name: string,
  quizzes: string[],
): Promise<JSONError | null> {
  const body: PostTestRequest = {
    name: name,
    quizzes: quizzes,
  };
  const bodyString = JSON.stringify(body);
  try {
    const response = await TokenizedFetch(
      "https://" + import.meta.env.VITE_DOMAIN + "/api/v1/test/",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: bodyString,
      },
    );

    if (!response.ok) {
      return (await response.json()) as JSONError;
    }

    return null;
  } catch (err) {
    console.log("Couldn't fetch test post: ", err);
    if (err instanceof Error) {
      return { error: "couldn't send test post because of in-browser error" };
    }

    return { error: "couldn't send test post because of unknown error" };
  }
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
  ).andThen((r) => {
    if (!r.ok) {
      console.log("response is not ok!");
      return ResultAsync.fromPromise(
        r.json() as Promise<JSONError>,
        (err): JSONError => {
          console.log("couldn't parse error's body: ", err);
          return { error: "couldn't parse error's body" };
        },
      ).andThen((e: JSONError) => errAsync(e));
    }
    return okAsync(null);
  });
}

type PatchTestRequest = {
  name: string | undefined;
};

export async function FetchTestPatch(
  UUID: string,
  name: string | undefined,
): Promise<JSONError | null> {
  const body: PatchTestRequest = {
    name: name,
  };

  const bodyString = JSON.stringify(body);

  try {
    const response = await TokenizedFetch(
      "https://" + import.meta.env.VITE_DOMAIN + "/api/v1/test/" + UUID,
      {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: bodyString,
      },
    );

    if (!response.ok) {
      console.log("response is not ok!");
      return (await response.json()) as JSONError;
    }

    return null;
  } catch (err) {
    console.log("Couldn't fetch test patch: ", err);
    if (err instanceof Error) {
      return { error: "couldn't send test patch because of in-browser error" };
    }

    return { error: "couldn't send test patch because of unknown error" };
  }
}

type TestBundleRequest = {
  quiz_uuids: string[];
};

export async function FetchTestBundle(
  testUUID: string,
  quizUUIDs: string[],
): Promise<null | JSONError> {
  const body: TestBundleRequest = {
    quiz_uuids: quizUUIDs,
  };

  const bodyString = JSON.stringify(body);

  try {
    const response = await TokenizedFetch(
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
    );

    if (!response.ok) {
      console.log("response is not ok!");
      return (await response.json()) as JSONError;
    }

    return null;
  } catch (err) {
    console.log("Couldn't fetch test bundle: ", err);
    if (err instanceof Error) {
      return {
        error: "couldn't send test bundle patch because of in-browser error",
      };
    }

    return { error: "couldn't send test bundle because of unknown error" };
  }
}

export async function FetchTestPrune(
  testUUID: string,
  quizUUIDs: string[],
): Promise<null | JSONError> {
  const body: TestBundleRequest = {
    quiz_uuids: quizUUIDs,
  };

  const bodyString = JSON.stringify(body);

  try {
    const response = await TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/${testUUID}/quizzes`,
      {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: bodyString,
      },
    );

    if (!response.ok) {
      console.log("response is not ok!");
      return (await response.json()) as JSONError;
    }

    return null;
  } catch (err) {
    console.log("Couldn't fetch test prune for test: ", err);
    if (err instanceof Error) {
      return { error: "couldn't send test prune because of in-browser error" };
    }

    return { error: "couldn't send test prune because of unknown error" };
  }
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
  ).andThen((r) => {
    if (!r.ok) {
      return ResultAsync.fromPromise(r.json(), (err): JSONError => {
        console.log("couldn't parse error's body: ", err);
        return { error: "couldn't parse error's body" };
      }).andThen((e: JSONError) => errAsync(e));
    }

    return okAsync(null);
  });
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
  ).andThen((r) => {
    if (!r.ok) {
      return ResultAsync.fromPromise(r.json(), (err): JSONError => {
        console.log("couldn't parse error's body: ", err);
        return { error: "couldn't parse error's body" };
      }).andThen((e: JSONError) => errAsync(e));
    }

    return okAsync(null);
  });
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
  ).andThen((r) => {
    if (!r.ok) {
      return ResultAsync.fromPromise(r.json(), (err): JSONError => {
        console.log("couldn't parse error's body: ", err);
        return { error: "couldn't parse error's body" };
      }).andThen((e: JSONError) => errAsync(e));
    }

    return okAsync(null);
  });
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
  ).andThen((r) => {
    if (!r.ok) {
      return ResultAsync.fromPromise(r.json(), (err): JSONError => {
        console.log("couldn't parse error's body: ", err);
        return { error: "couldn't parse error's body" };
      }).andThen((e: JSONError) => errAsync(e));
    }

    return okAsync(null);
  });
}
