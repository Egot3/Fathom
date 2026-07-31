import { type JSONError } from "../statuses/jsonerror";
import { TokenizedFetch } from "./tokenizedFetch";

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
	lable: string;
};

type OptionsRadio = {
	radio: { chosen: Choice[] };
};
type OptionsCheck = {
	check: { chosen: Choice[] };
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
	return (qo as OptionsRadio).radio.chosen;
}
export function OptionCheck(qo: QuizOptions): Choice[] {
	return (qo as OptionsCheck).check.chosen;
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

export type QuizzesOrError = { quizzes: Quiz[]; total: number } | JSONError;

export async function FetchAllQuizzes(
	page: number,
	size: number,
): Promise<QuizzesOrError> {
	try {
		const rawRes = await TokenizedFetch(
			"https://" +
				import.meta.env.VITE_DOMAIN +
				"/api/v1/quiz/" +
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
		);

		if (!rawRes.ok) {
			return (await rawRes.json()) as JSONError;
		}

		const quizzes = (await rawRes.json()) as { quizzes: Quiz[]; total: number };
		console.log("quizzes: ", quizzes);
		return quizzes;
	} catch (err) {
		console.log(err);
		return {
			error: "network error",
		} as JSONError;
	}
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

export async function FetchQuizPost(
	meta: Meta,
	path: string,
	contents: string,
): Promise<null | JSONError> {
	const body: PostQuizRequest = {
		meta: meta,
		body: contents,
		name: path,
	};

	const bodyString = JSON.stringify(body);

	try {
		const response = await TokenizedFetch(
			"https://" + import.meta.env.VITE_DOMAIN + "/api/v1/quiz/",
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
		console.log("Couldn't fetch quiz post: ", err);
		if (err instanceof Error) {
			return { error: "couldn't send quiz post because of in-browser error" };
		}

		return { error: "couldn't send quiz post because of unknown error" };
	}
}

type PutQuizRequest = {
	meta: Meta;
	body: string;
};

export async function FetchQuizPut(
	UUID: string,
	meta: Meta,
	contents: string,
): Promise<JSONError | null> {
	const body: PutQuizRequest = {
		meta: meta,
		body: contents,
	};

	const bodyString = JSON.stringify(body);

	try {
		const response = await TokenizedFetch(
			"https://" + import.meta.env.VITE_DOMAIN + "/api/v1/quiz/" + UUID,
			{
				method: "PUT",
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
		console.log("Couldn't fetch quiz put for quiz: ", err);
		if (err instanceof Error) {
			return { error: "couldn't send quiz put because of in-browser error" };
		}

		return { error: "couldn't send quiz put because of unknown error" };
	}
}

type PatchQuizRequest = {
	name: string | undefined;
	score: number | undefined;
};

export async function FetchQuizPatch(
	UUID: string,
	score: number | undefined,
	name: string | undefined,
): Promise<JSONError | null> {
	const body: PatchQuizRequest = {
		name: name,
		score: score,
	};

	const bodyString = JSON.stringify(body);

	try {
		const response = await TokenizedFetch(
			"https://" + import.meta.env.VITE_DOMAIN + "/api/v1/quiz/" + UUID,
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
		console.log("Couldn't fetch quiz patch for quiz: ", err);
		if (err instanceof Error) {
			return { error: "couldn't send quiz patch because of in-browser error" };
		}

		return { error: "couldn't send quiz patch because of unknown error" };
	}
}

// helper
export async function FetchChangeQuiz(
	UUID: string,
	meta: Meta,
	contents: string,
	name: string | undefined,
): Promise<JSONError | null> {
	const putResp = await FetchQuizPut(UUID, meta, contents);
	if (putResp != null) {
		return putResp;
	}
	return FetchQuizPatch(UUID, meta.score, name);
}

type QuizFileOrError = QuizFile | JSONError;
export async function FetchQuiz(UUID: string): Promise<QuizFileOrError> {
	try {
		const response = await TokenizedFetch(
			"https://" + import.meta.env.VITE_DOMAIN + "/api/v1/quiz/" + UUID,
			{
				method: "GET",
				headers: {
					Accept: "application/json",
				},
			},
		);

		if (!response.ok) {
			console.log("response is not ok!");
			return (await response.json()) as JSONError;
		}

		return await response.json();
	} catch (err) {
		console.log("Couldn't fetch quiz patch for quiz: ", err);
		if (err instanceof Error) {
			return { error: "couldn't send quiz patch because of in-browser error" };
		}

		return { error: "couldn't send quiz patch because of unknown error" };
	}
}