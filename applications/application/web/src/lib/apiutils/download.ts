import { Result } from "neverthrow";
import type { JSONError } from "../statuses/jsonerror";

export function Download(blob: Blob, filename = (Date.now()).toString()): Result<null, JSONError> {
  return Result.fromThrowable(() => {
    const downloadUrl: string = window.URL.createObjectURL(blob);

    const link: HTMLAnchorElement = document.createElement('a');
    link.href = downloadUrl;
    link.download = filename;

    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);

    window.URL.revokeObjectURL(downloadUrl);
    return null
  }, (err): JSONError => {
    console.log("couldn't install ", filename, err)
    return { error: "Unable to install file" }
  })()
}
