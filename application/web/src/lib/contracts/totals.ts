import { GetGroup } from "./group";
import { GetTest } from "./test";
import { TokenizedFetch } from "./tokenizedFetch";

type rawTotal = {
	test_uuid: string;
	group_uuid: string;
	user_uuid: string;
};

type testTotal = {
	testName: string;
	groupName: string;
	userName: string;
};

export async function FetchTestsTotalsForUser(
	userUUID: string,
	page: number,
	size: number,
): Promise<{ totals: testTotal[]; total: number }> {
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
			} as testTotal;
		},
	);

	return {
		totals: await Promise.all(testsTotalPromises),
		total: rawTotals.total,
	};
}
