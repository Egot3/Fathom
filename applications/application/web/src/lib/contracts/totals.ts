import type { JSONError } from "../statuses/jsonerror";
import { TokenizedFetch } from "./tokenizedFetch";

type rawTotal = {
  test_uuid: string;
  group_uuid: string;
  user_uuid: string;
  score: number;
};

export type TestTotal = {
  test_name: string;
  group_name: string;
  user_name: string;
  score: number;
  finalized_at: Date;
  max_score: number;

  test_uuid: string;
  group_uuid: string;
  user_uuid: string;
};

export type Answer = {
  group_uuid: string;
  test_uuid: string;
  user_uuid: string;
  quiz_uuid: string;
  chosen: string;
  correct: string;
  submitted_at: Date;

  score: number;
  max_score: number;

  group_name: string;
  quiz_name: string;
  test_name: string;
};

export type TotalsOrError = { totals: TestTotal[]; total: number } | JSONError;

export async function GetTotalsForUser(
  userUUID: string,
  page: number,
  size: number,
): Promise<TotalsOrError> {
  try {
    const rawRes = await TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/total/all/${userUUID}?page=${page}&size=${size}`,
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

    const totals = (await rawRes.json()) as {
      totals: TestTotal[];
      total: number;
    };
    return totals;
  } catch (err) {
    console.log(err);
    return {
      error: "network error",
    } as JSONError;
  }
}

export async function FetchTotals(
  page: number,
  size: number,
): Promise<TotalsOrError> {
  try {
    const rawRes = await TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/total?page=${page}&size=${size}`,
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

    const totals = (await rawRes.json()) as {
      totals: TestTotal[];
      total: number;
    };
    return totals;
  } catch (err) {
    console.log(err);
    return {
      error: "network error",
    } as JSONError;
  }
}

type AnswersOrJSONError = { answers: Answer[]; total: number } | JSONError;

export async function FetchAnswers(
  groupUUID: string,
  userUUID: string,
  testUUID: string,
  page: number,
  size: number,
): Promise<AnswersOrJSONError> {
  try {
    const rawRes = await TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/total/${groupUUID}/${userUUID}/${testUUID}/answers?page=${page}&size=${size}`,
      {
        method: "GET",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/x-www-form-urlencoded",
        },
      },
    );

    if (!rawRes.ok) {
      return (await rawRes.json()) as JSONError;
    }

    const answers = (await rawRes.json()) as {
      answers: Answer[];
      total: number;
    };
    return answers;
  } catch (err) {
    console.log("Couldn't list answers: ", err);
    if (err instanceof Error) {
      return {
        error: "couldn't fetch list answers because of in-browser error",
      };
    }

    return { error: "couldn't fetch list answers because of unknown error" };
  }
}
