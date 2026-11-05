export type PasswordPolicy = {
  min_length: number;
  max_length: number;
  min_digits: number;
  min_lowercase: number;
  min_uppercase: number;
  allow_whitespace: boolean;
  max_bytes: number;
};

const count = (s: string, re: RegExp) => (s.match(re) ?? []).length;

export const CheckPassword = (pw: string, p: PasswordPolicy): string[] => {
  const out: string[] = [];
  const len = [...pw].length;

  const blen = new TextEncoder().encode(pw).length;

  if (blen > p.max_bytes) {
    out.push(`mustn't be longer than ${p.max_bytes} bytes`);
  }
  if (len < p.min_length)
    out.push(`must contain at least ${p.min_length} characters`);
  else if (len > p.max_length)
    out.push(`mustn't be longer than ${p.max_length} characters`);
  if (count(pw, /[0-9]/g) < p.min_digits)
    out.push(`must contain at least ${p.min_digits} numbers`);
  if (count(pw, /[a-z]/g) < p.min_lowercase)
    out.push(`must contain at least ${p.min_lowercase} lowercase letters`);
  if (count(pw, /[A-Z]/g) < p.min_uppercase)
    out.push(`must contain at least ${p.min_uppercase} uppercase letters`);
  if (!p.allow_whitespace && /\p{White_Space}/u.test(pw))
    out.push("mustn't contain any type of whitespace");

  return out;
};
