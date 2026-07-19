import { TokenizedFetch } from "./tokenizedFetch";

type Test = {
	uuid: string;
	name: string;
	created_at: Date;
	updated_at: Date;
};

type GetTestResponse = {
	test: Test;
};

export async function GetTest(testUUID: string): Promise<GetTestResponse> {
	const rawRes = await TokenizedFetch(
		"http://" +
			import.meta.env.VITE_DOMAIN +
			"/api/v1/test/" +
			testUUID +
			{
				method: "GET",
				headers: {
					Accept: "application/json",
				},
			},
	);

	return (await rawRes.json()) as GetTestResponse;
}
