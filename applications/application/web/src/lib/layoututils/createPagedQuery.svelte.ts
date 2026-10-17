import type { ResultAsync } from "neverthrow";
import type { JSONError } from "../statuses/jsonerror";

export function CreatePagedQuery<R, P>(opts: {
  fetcher: (page: number, size: number, params: P) => ResultAsync<R, JSONError>;
  params: () => P;
  pageSize: () => number;
  enabled?: () => boolean;
  debounce?: number;
}) {
  let data = $state<R | null>(null);
  let error = $state("");
  let loading = $state(true);
  let page = $state(1);
  let trigger = $state(0);

  let ticket = 0;
  let lastKey: string | undefined;

  $effect(() => {
    const params = opts.params();
    const size = opts.pageSize();
    const p = page;
    trigger;

    if (opts.enabled && !opts.enabled()) return;

    const key = JSON.stringify(params);
    if (lastKey !== undefined && key !== lastKey && p !== 1) {
      lastKey = key;
      page = 1;
      return;
    }
    lastKey = key;

    loading = true;
    const mine = ++ticket;

    const t = setTimeout(async () => {
      const res = await opts.fetcher(p - 1, size, params);
      if (mine !== ticket) return;
      console.log("got res", res);
      res.match(
        (r) => {
          data = r;
          error = "";
        },
        (e) => {
          error = e.error;
        },
      );
      loading = false;
    }, opts.debounce ?? 300);

    return () => clearTimeout(t);
  });

  return {
    get data() {
      return data;
    },
    get error() {
      return error;
    },
    get loading() {
      return loading;
    },
    get page() {
      return page;
    },
    set page(v: number) {
      page = v;
    },
    refresh() {
      trigger++;
    },
  };
}
