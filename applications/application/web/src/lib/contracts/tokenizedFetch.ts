import { Toaster } from "../apiutils/toaster";
import { GetTokenExpiration, SetTokenExpiration } from "../bgdata/user.svelte";
import { ERR_UNAUTHORIZED } from "../carefulness/unauthorized";

export const maxAgeRegex = /max-age=(\d+)/;

export async function TokenizedFetch(
  url: RequestInfo | URL,
  opts?: RequestInit,
  loggedout?: boolean,
): Promise<Response> {
  if (!loggedout || loggedout === undefined) {
    const exp = GetTokenExpiration();
    console.log("got token exp: ", exp);

    if (exp < new Date()) {
      console.log("expired");
      window.location.reload();
      Toaster.info({
        title: "Session expired!",
        description: "login again",
      });
      throw ERR_UNAUTHORIZED;
    }
  }

  const res = await fetch(url, { ...opts, credentials: "include" });

  switch (res.status) {
    case 401:
      SetTokenExpiration(new Date(Date.now() - 1));
    default:
      const sessionControl = res.headers.get("Session-Control");
      if (sessionControl !== null) {
        const reg = maxAgeRegex.exec(sessionControl);
        if (reg === null || reg.length < 2) {
          return res;
        }
        const maxAge = parseInt(reg[1], 10);

        SetTokenExpiration(new Date(Date.now() + maxAge * 1000));
      }
  }

  return res;
}
