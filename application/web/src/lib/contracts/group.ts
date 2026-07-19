import { TokenizedFetch } from "./tokenizedFetch";

type Group = {
	uuid: string;
	name: string;
};

type GetGroupResponse = {
	group: Group;
};

export async function GetGroup(groupUUID: string): Promise<GetGroupResponse> {
	const rawRes = await TokenizedFetch(
		"http://" +
			import.meta.env.VITE_DOMAIN +
			"/api/v1/group/" +
			groupUUID +
			{
				method: "GET",
				headers: {
					Accept: "application/json",
				},
			},
	);

	return (await rawRes.json()) as GetGroupResponse;
}
