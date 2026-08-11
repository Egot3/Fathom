export type JSONError = {
  error: string; // no tags?
};

export function IsJSONError(res: Record<string, any>): res is JSONError {
  return typeof res === "object" && res !== null && "error" in res;
}
