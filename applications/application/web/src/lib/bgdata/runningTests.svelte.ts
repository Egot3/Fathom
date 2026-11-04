import {
  FetchCurrentlyRunningTestInfos,
  type TestInfo,
} from "../contracts/test";

export function CreateRunningTests() {
  let data = $state<TestInfo[]>([]);
  let error = $state("");
  let loading = $state(true);
  let ticket = 0;

  async function refresh(silent = false) {
    const mine = ++ticket;
    if (!silent) loading = true;

    await FetchCurrentlyRunningTestInfos().match(
      (r) => {
        if (mine !== ticket) return;
        data = r;
        error = "";
      },
      (e) => {
        if (mine !== ticket) return;
        data = [];
        error = e.error;
      },
    );

    if (mine === ticket) loading = false;
  }

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
    refresh,
  };
}
