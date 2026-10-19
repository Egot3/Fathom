import { errAsync, okAsync, ResultAsync } from "neverthrow";
import {
  EvictQuiz,
  GetCachedQuiz,
  SetCachedQuiz,
  ReadManifest,
} from "../apiutils/quizManifest";
import { type JSONError } from "../statuses/jsonerror";
import {
  NormalizeBlob,
  NormalizeJSON,
  NormalizeResponse,
  TokenizedFetch,
} from "./tokenizedFetch";
import { AcceptTransform, type Accept } from "../apiutils/acceptHeader";

export type Quiz = {
  uuid: string;
  path: string;
  score: number;
  correct_answer: string;
};
export type QuizFile = {
  meta: Meta;
  body: string;
};
export type ParsedQuiz = {
  meta: Meta;
  title: string;
  body: string;
  options: QuizOptions;
  answers?: QuizAnswer;
};

type AnswerInput = {
  input: { input: string };
};
type AnswerRadio = {
  radio: { chosen: number };
};
type AnswerCheck = {
  check: { chosen: number[] };
};
type AnswerAccordance = {
  accordance: { accorded: number[] };
};
type AnswerOrder = {
  order: { item_indexes: number[] };
};

export type QuizAnswer =
  AnswerAccordance | AnswerCheck | AnswerInput | AnswerOrder | AnswerRadio;
// warning: may be not safe
export function AnswerInput(qa: QuizAnswer): string {
  return (qa as AnswerInput).input.input;
}
export function AnswerRadio(qa: QuizAnswer): number {
  return (qa as AnswerRadio).radio.chosen;
}
export function AnswerCheck(qa: QuizAnswer): number[] {
  return (qa as AnswerCheck).check.chosen;
}
export function AnswerAccordance(qa: QuizAnswer): number[] {
  return (qa as AnswerAccordance).accordance.accorded;
}
export function AnswerOrder(qa: QuizAnswer): number[] {
  return (qa as AnswerOrder).order.item_indexes;
}

type Choice = {
  id: number;
  label: string;
};

type OptionsRadio = {
  radio: { choices: Choice[] };
};
type OptionsCheck = {
  check: { choices: Choice[] };
};
type OptionsAccordance = {
  accordance: {
    static: string[];
    dynamic: string[];
  };
};
type OptionsOrder = {
  orders: { items: string[] };
};

export type QuizOptions =
  OptionsAccordance | OptionsRadio | OptionsCheck | OptionsOrder;
export function OptionRadio(qo: QuizOptions): Choice[] {
  return (qo as OptionsRadio).radio.choices;
}
export function OptionCheck(qo: QuizOptions): Choice[] {
  return (qo as OptionsCheck).check.choices;
}
export function OptionAccordance(qa: QuizOptions): {
  static: string[];
  dynamic: string[];
} {
  return (qa as OptionsAccordance).accordance;
}
export function OptionOrder(qo: QuizOptions): string[] {
  return (qo as OptionsOrder).orders.items;
}

export type Quizzes = { quizzes: Quiz[]; total: number };

export function FetchAllQuizzes(
  page: number,
  size: number,
): ResultAsync<Quizzes, JSONError> {
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/quiz/?page=${page}&size=${size}`,
    {
      method: "GET",
      headers: {
        Accept: "application/json",
      },
    },
  ).andThen((r) => NormalizeJSON<Quizzes>(r));
}

export enum Kind {
  Input = "INPUT",
  Radio = "RADIO",
  Check = "CHECK",
  Accordance = "ACCORDANCE",
  Order = "ORDER",
}

export type Meta = {
  kind: Kind;
  randomized: boolean;
  score: number;
  all_or_none: boolean;
};

type PostQuizRequest = {
  meta: Meta;
  name: string;
  body: string;
};

export function FetchQuizPost(
  meta: Meta,
  path: string,
  contents: string,
): ResultAsync<null, JSONError> {
  const body: PostQuizRequest = {
    meta: meta,
    body: contents,
    name: path,
  };

  const bodyString = JSON.stringify(body);

  return TokenizedFetch(
    "https://" + import.meta.env.VITE_DOMAIN + "/api/v1/quiz/",
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

type PatchQuizRequest = {
  name: string | undefined;
  meta: Meta | undefined;
  body: string | undefined;
};

export function FetchQuizPatch(
  UUID: string,
  req: PatchQuizRequest,
): ResultAsync<null, JSONError> {
  const bodyString = JSON.stringify(req);
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/quiz/${UUID}`,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/merge-patch+json", // we are being very rfc rn
      },
      body: bodyString,
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

export function FetchQuizDelete(UUID: string): ResultAsync<null, JSONError> {
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/quiz/${UUID}`,
    {
      method: "DELETE",
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

export type ParsedQuizResponse = { quiz: ParsedQuiz };

// may not throw
export function FetchParsedQuiz(
  quizUUID: string,
): ResultAsync<ParsedQuiz, JSONError> {
  const manifest = ReadManifest();
  const known = manifest[quizUUID];
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/quiz/${quizUUID}/parsed`,
    {
      headers: known ? { "If-None-Match": `"${known.checksum}"` } : {},
    },
  ).andThen((r) => {
    if (r.status === 304 && known) {
      const cached = GetCachedQuiz<ParsedQuiz>(quizUUID, known.checksum);
      if (cached !== null) return okAsync(cached);
    }
    if (r.status === 404) {
      EvictQuiz(quizUUID);
      return errAsync<ParsedQuiz, JSONError>({
        error: `Quiz ${quizUUID} not found`,
      });
    }

    return NormalizeJSON<ParsedQuizResponse>(r)
      .andTee((body: ParsedQuizResponse) => {
        const etag = r.headers.get("ETag")?.replace(/"/g, "") ?? "";
        console.log("ETag:", etag);
        SetCachedQuiz(quizUUID, etag, body.quiz);
      })
      .andThen((r) => okAsync(r.quiz));
  });
}

export function FetchQuizExport(
  UUIDs: string[],
  accepts: Accept[],
): ResultAsync<Blob, JSONError> {
  const body = JSON.stringify({ uuids: UUIDs });
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/quiz/export`,
    {
      method: "POST",
      body: body,
      headers: {
        Accept: AcceptTransform(accepts),
        "Content-Type": "application/json",
      },
    },
  ).andThen(NormalizeBlob);
}

export function FetchQuizImport(file: File): ResultAsync<null, JSONError> {
  const formData = new FormData();

  formData.append("imported", file, file.name);

  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/quiz/import`,
    {
      method: "POST",
      body: formData,
    },
  )
    .andThen(NormalizeResponse)
    .andThen((_) => okAsync(null));
}

type RawQuiz = {
  meta: Meta;
  body: string;
};
export function FetchRawQuiz(UUID: string): ResultAsync<RawQuiz, JSONError> {
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/quiz/${UUID}`,
    {
      method: "GET",
    },
  ).andThen(NormalizeJSON<RawQuiz>);
}

export function FetchAnswerPost(
  groupUUID: string,
  userUUID: string,
  runnerKey: string,
  quizUUID: string,
  answer: QuizAnswer,
): ResultAsync<null, JSONError> {
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/total/${groupUUID}/${userUUID}/running/${runnerKey}/${quizUUID}`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ answer }),
    },
  )
    .andThen(NormalizeResponse)
    .map(() => null);
}

export function FetchRunningQuiz(
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
