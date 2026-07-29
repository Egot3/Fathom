import type { JSONError } from "../statuses/jsonerror";
import { TokenizedFetch } from "./tokenizedFetch";

export type Quiz = {
	uuid: string;
	path: string;
	score: number;
	correct_answer: string;
};

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
		console.log("Couldn't fetch quiz post for user: ", err);
		if (err instanceof Error) {
			return { error: "couldn't send quiz post because of in-browser error" };
		}

		return { error: "couldn't send quiz post because of unknown error" };
	}
}
