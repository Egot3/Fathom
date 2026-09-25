import { ResultAsync } from "neverthrow";
import type { JSONError } from "../statuses/jsonerror";
import { NormalizeJSON, TokenizedFetch } from "./tokenizedFetch";
import type { QuizAnswer } from "./quiz";

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
  chosen: QuizAnswer;
  correct: QuizAnswer;
  submitted_at: Date;

  score: number;
  max_score: number;

  group_name: string;
  quiz_name: string;
  test_name: string;
};

export type Totals = { totals: TestTotal[]; total: number };
export type TotalsOrError = Totals | JSONError;

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

export function FetchTotals(
  page: number,
  size: number,

  userUUID?: string,
  groupUUID?: string,
  testUUID?: string,
): ResultAsync<Totals, JSONError> {
  return ResultAsync.fromPromise(
    TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/total/${groupUUID ?? "all"}/${userUUID ?? "all"}/${testUUID ?? "all"}?page=${page}&size=${size}`,
      {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
      },
    ),
    (err) => {
      console.log(err);
      return {
        error: "network error",
      } as JSONError;
    },
  ).andThen((r) => {
    return NormalizeJSON<Totals>(r);
  });
}

export type Answers = { answers: Answer[]; total: number };

export function FetchAnswers(
  groupUUID: string,
  userUUID: string,
  testUUID: string,
  page: number,
  size: number,
): ResultAsync<Answers, JSONError> {
  return ResultAsync.fromPromise(
    TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/total/${groupUUID}/${userUUID}/${testUUID}/answers?page=${page}&size=${size}`,
      {
        method: "GET",
        headers: {
          Accept: "application/json",
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
  ).andThen((r) => NormalizeJSON<Answers>(r));
}
