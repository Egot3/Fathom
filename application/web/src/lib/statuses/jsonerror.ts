export type JSONError = {
	error: string; // no tags?
};

export function IsJSONError(res: Record<string, any>): res is JSONError {
	return "error" in res;
}
