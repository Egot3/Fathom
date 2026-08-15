import { GetUser, SetTokenExpiration } from "../bgdata/user.svelte";

export const maxAgeRegex = /max-age=(\d+)/;

export async function TokenizedFetch(
  url: RequestInfo | URL,
  opts?: RequestInit,
): Promise<Response> {
  const res = await fetch(url, { ...opts, credentials: "include" });

  switch (res.status) {
    case 401:
      SetTokenExpiration(new Date(Date.now() - 1));
    default:
      const sessionControl = res.headers.get("Session-Control");
      console.log("session control: ", sessionControl);
      if (sessionControl != null) {
        const reg = maxAgeRegex.exec(sessionControl);
        if (reg == null || reg.length < 2) {
          return res;
        }
        const maxAge = parseInt(reg[1], 10);

        console.log("max age", maxAge);
        SetTokenExpiration(new Date(Date.now() + maxAge * 1000));
      }
  }

  return res;
}
