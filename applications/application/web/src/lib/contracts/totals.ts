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

	test_uuid: string;
	group_uuid: string;
	user_uuid: string;
};

export type TotalsOrError = { totals: TestTotal[]; total: number } | JSONError;

export async function GetTotalsForUser(
	userUUID: string,
	page: number,
	size: number,
): Promise<TotalsOrError> {
	try {
		const rawRes = await TokenizedFetch(
			"https://" +
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
