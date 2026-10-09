export type Accept = { accept: string; quality?: number };
export function AcceptTransform(accepts: Accept[]): string {
  return accepts
    .map((v) => {
      return v.accept + (v.quality ? ";q=" + v.quality : "");
    })
    .join(", ");
}

const mime = {
  Tar: { name: "Tar", mime: "application/tar" },
  Targz: { name: "Tar+gzip", mime: "application/gzip" },
  Zip: { name: "Zip", mime: "application/zip" },
} as const;

export type USER_MIME = (typeof mime)[keyof typeof mime];
export const USER_MIMES = Object.values(mime) as readonly USER_MIME[];
