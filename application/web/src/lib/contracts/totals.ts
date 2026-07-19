import type { JSONError } from "../statuses/jsonerror";
import { GetGroup } from "./group";
import { GetTest } from "./test";
import { TokenizedFetch } from "./tokenizedFetch";

type rawTotal = {
	test_uuid: string;
	group_uuid: string;
	user_uuid: string;
	score: number;
};

type testTotal = {
	testName: string;
	groupName: string;
	userName: string;
	score: number;
};

export type TotalsOrError = { totals: testTotal[]; total: number } | JSONError;

export async function GetTotalsForUser(
	userUUID: string,
	page: number,
	size: number,
): Promise<TotalsOrError> {
	const rawRes = await TokenizedFetch(
		"http://" +
			import.meta.env.VITE_DOMAIN +
			"/api/v1/total/all/" +
			userUUID +
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

	const rawTotals = (await rawRes.json()) as {
		totals: rawTotal[];
		total: number;
	};
	const testsTotalPromises: Promise<testTotal>[] = rawTotals.totals.map(
		async (v) => {
			const testResp = await GetTest(v.test_uuid);
			const groupResp = await GetGroup(v.group_uuid);
			return {
				testName: testResp.test.name,
				groupName: groupResp.group.name,
				score: v.score,
			} as testTotal;
		},
	);

	return {
		totals: await Promise.all(testsTotalPromises),
		total: rawTotals.total,
	};
}
