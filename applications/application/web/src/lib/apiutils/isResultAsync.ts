import { ResultAsync } from "neverthrow";
import type { JSONError } from "../statuses/jsonerror";

export function IsResultAsync(
  value: unknown,
): value is ResultAsync<unknown, JSONError> {
  return (
    value instanceof ResultAsync ||
    (typeof value === "object" && value !== null && "asyncAndThen" in value)
  );
}
