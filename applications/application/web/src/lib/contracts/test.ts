import type { ETagInfo } from "../bgdata/currentlyrunning.svelte";
import type { JSONError } from "../statuses/jsonerror";
import type { Quiz } from "./quiz";
import { maxAgeRegex, TokenizedFetch } from "./tokenizedFetch";

type Test = {
	uuid: string;
	name: string;
	created_at: Date;
	updated_at: Date;
	quizzes: Quiz[];
};

type GetTestResponse = {
	test: Test;
	deadline: string;
	is_paused: boolean;
};

type GetQuizUUIDsResponse = {
	quiz_uuids: string[];
};

// redo
export async function FetchTest(testUUID: string): Promise<Test> {
	const rawRes = await TokenizedFetch(
		"http://" + import.meta.env.VITE_DOMAIN + "/api/v1/test/" + testUUID,
		{
			method: "GET",
			headers: {
				Accept: "application/json",
			},
		},
	);

	return ((await rawRes.json()) as GetTestResponse).test;
}

export async function FetchCurrentlyRunningTestInfo(): Promise<
	{ test: Test; deadline: Date; isPaused: boolean } | JSONError | null
> {
	try {
		const res = await TokenizedFetch(
			"http://" + import.meta.env.VITE_DOMAIN + "/api/v1/test/running",
		);
		if (!res.ok) {
			if (res.status === 423) {
				return null;
			}
			return (await res.json()) as JSONError;
		}

		const respJSON = (await res.json()) as GetTestResponse;
		return {
			test: respJSON.test,
			deadline: new Date(respJSON.deadline),
			isPaused: respJSON.is_paused,
		};
	} catch (e) {
		console.log("couldn't fetch current test info due to unknown error: ", e);
		return {
			error: "got network error while fetching current test",
		} as JSONError;
	}
}

export async function FetchCurrentlyRunningQuizUUIDs(
	ETag?: string,
): Promise<{ Caching: ETagInfo; UUIDs: string[] } | JSONError | null> {
	if (ETag === undefined) {
		try {
			const res = await TokenizedFetch(
				"http://" +
					import.meta.env.VITE_DOMAIN +
					"/api/v1/test/running/quizzes",
			);
			if (!res.ok) {
				if (res.status === 423) {
					return null;
				}
				return (await res.json()) as JSONError;
			}

			const etag = res.headers.get("ETag");
			if (etag == null) {
				return {
					UUIDs: ((await res.json()) as GetQuizUUIDsResponse).quiz_uuids,
					Caching: { ETag: "", ExpiresAt: new Date() },
				};
			}

			const reg = maxAgeRegex.exec(res.headers.get("Cache-Control") ?? "");
			if (reg == null || reg.length < 2) {
				return {
					UUIDs: ((await res.json()) as GetQuizUUIDsResponse).quiz_uuids,
					Caching: { ETag: "", ExpiresAt: new Date() },
				};
			}
			const maxAge = parseInt(reg[1], 10);
			return {
				UUIDs: ((await res.json()) as GetQuizUUIDsResponse).quiz_uuids,
				Caching: {
					ETag: etag,
					ExpiresAt: new Date(Date.now() + maxAge * 1000),
				},
			};
		} catch (e) {
			console.log("couldn't fetch current test info due to unknown error: ", e);
			return {
				error: "got network error while fetching current test",
			} as JSONError;
		}
	}

	try {
		const res = await TokenizedFetch(
			"http://" + import.meta.env.VITE_DOMAIN + "/api/v1/test/running/quizzes",
			{ headers: {} },
		);
		if (!res.ok) {
			if (res.status === 304) {
				return null;
			}
			return (await res.json()) as JSONError;
		}

		const etag = res.headers.get("ETag");
		if (etag == null) {
			return {
				UUIDs: ((await res.json()) as GetQuizUUIDsResponse).quiz_uuids,
				Caching: { ETag: "", ExpiresAt: new Date() },
			};
		}

		const reg = maxAgeRegex.exec(res.headers.get("Cache-Control") ?? "");
		if (reg == null || reg.length < 2) {
			return {
				UUIDs: ((await res.json()) as GetQuizUUIDsResponse).quiz_uuids,
				Caching: { ETag: "", ExpiresAt: new Date() },
			};
		}
		const maxAge = parseInt(reg[1], 10);
		return {
			UUIDs: ((await res.json()) as GetQuizUUIDsResponse).quiz_uuids,
			Caching: {
				ETag: etag,
				ExpiresAt: new Date(Date.now() + maxAge * 1000),
			},
		};
	} catch (e) {
		console.log("couldn't fetch current test info due to unknown error: ", e);
		return {
			error: "got network error while fetching current test",
		} as JSONError;
	}
}

export type TestsOrError = { tests: Test[]; total: number } | JSONError;

export async function FetchAllTests(
	page: number,
	size: number,
): Promise<TestsOrError> {
	try {
		const rawRes = await TokenizedFetch(
			"http://" +
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
		);

		if (!rawRes.ok) {
			return (await rawRes.json()) as JSONError;
		}

		const totals = (await rawRes.json()) as { tests: Test[]; total: number };
		console.log("totals: ", totals);
		return totals;
	} catch (err) {
		console.log(err);
		return {
			error: "network error",
		} as JSONError;
	}
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
			"http://" + import.meta.env.VITE_DOMAIN + "/api/v1/test/",
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
		console.log("Couldn't fetch login for user: ", err);
		if (err instanceof Error) {
			return { error: "couldn't send login because of in-browser error" };
		}

		return { error: "couldn't send login because of unknown error" };
	}
}
