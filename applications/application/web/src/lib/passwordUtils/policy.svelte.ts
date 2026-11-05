import { NormalizeJSON, TokenizedFetch } from "../contracts/tokenizedFetch";
import type { PasswordPolicy } from "./checkPassword";

let policy = $state<PasswordPolicy | null>(null);

export const GetPasswordPolicy = () => policy;

export async function LoadPasswordPolicy() {
  policy = await TokenizedFetch("/api/password-policy")
    .andThen(NormalizeJSON<PasswordPolicy>)
    .match(
      (r) => r,
      (e) => {
        console.log("couldn't get password policy");
        return null;
      },
    );
}
