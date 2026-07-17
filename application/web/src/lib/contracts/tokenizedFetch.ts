import { GetUser, SetTokenExpiration } from "../bgdata/user.svelte";

export async function TokenizedFetch(
	url: RequestInfo | URL,
	opts?: RequestInit,
): Promise<Response> {
	const res = await fetch(url, { ...opts, credentials: "include" });

	const sessionControl = res.headers.get("Session-Control");

	console.log("s control: ", sessionControl);
	if (sessionControl) {
		SetTokenExpiration(new Date(Date.now() + parseInt(sessionControl, 10)));
	}

	return res;
}
