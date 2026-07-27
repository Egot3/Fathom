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
			"http://" +
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
