import type { ETagInfo } from "../bgdata/currentlyrunning.svelte";
import type { JSONError } from "../statuses/jsonerror";
import { maxAgeRegex, TokenizedFetch } from "./tokenizedFetch";

type Test = {
	uuid: string;
	name: string;
	created_at: Date;
	updated_at: Date;
};

type GetTestResponse = {
	test: Test;
};

type GetQuizUUIDsResponse = {
	quiz_uuids: string[];
};

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
	Test | JSONError | null
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

		return ((await res.json()) as GetTestResponse).test;
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
