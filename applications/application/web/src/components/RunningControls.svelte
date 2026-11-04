<script lang="ts">
  import ChipSelector from "./ChipSelector.svelte";
  import ErrorPopover from "./ErrorPopover.svelte";
  import {
    FetchTestPause,
    FetchTestResume,
    type TestInfo,
  } from "../lib/contracts/test";

  const {
    test,
    onchanged,
  }: {
    test: TestInfo;
    onchanged: () => void | Promise<void>;
  } = $props();

  let busy = $state(false);
  let selected = $state(0);
  let error: string = $state("");

  $effect(() => {
    if (!busy) selected = test.is_paused ? 1 : 0;
  });

  async function toggle(next: number) {
    if (busy) return;
    busy = true;
    try {
      error = await (
        next === 1 ? FetchTestPause(test.key) : FetchTestResume(test.key)
      ).match(
        () => "",
        ({ error }) => error,
      );
      await onchanged();
    } finally {
      busy = false;
    }
  }
</script>

<ErrorPopover message={error} open={error !== ""}>
  <ChipSelector
    options={["running", "paused"]}
    working={busy}
    bind:selected
    onchange={() => toggle(selected)}
  />
</ErrorPopover>
