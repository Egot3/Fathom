
export type Accept = { accept: string, quality?: number }
export function AcceptTransform(accepts: Accept[]): string {
  return accepts.map((v) => {
    return v.accept + (v.quality ? "q=" + v.quality + ";" : "")
  }).join(", ")
}
