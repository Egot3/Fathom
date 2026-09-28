<script lang="ts" generics="T">
  import { ResultAsync } from "neverthrow";
  import type { JSONError } from "../lib/statuses/jsonerror";
  import { Download } from "../lib/apiutils/download";
  import ErrorPopover from "./ErrorPopover.svelte";
  import { CircleArrowUp, LoaderCircle } from "@lucide/svelte";

  const {
    fetcher,
    disabled,
    downloadName,
  }: {
    fetcher: () => ResultAsync<Blob, JSONError>;
    disabled?: boolean;
    downloadName: string;
  } = $props();

  let downloading = $state(false);

  let statusMessage = $state("");
  let errorPopoverOpen = $state(false);
  async function fetchAndDownload() {
    statusMessage = await fetcher()
      .andThen((r) => Download(r, downloadName))
      .orTee((_) => (errorPopoverOpen = true))
      .andTee((_) => (errorPopoverOpen = false))
      .match(
        (_) => "",
        (e) => e.error,
      );
  }
</script>

<ErrorPopover message={statusMessage} bind:open={errorPopoverOpen}>
  <button
    class="btn preset-tonal-warning"
    {disabled}
    onclick={fetchAndDownload}
  >
    {#if downloading}
      <LoaderCircle class="animate-spin" />
    {:else}
      <CircleArrowUp /> Export
    {/if}
  </button>
</ErrorPopover>
