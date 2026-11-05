<script lang="ts" generics="T">
  import { Pencil, X } from "@lucide/svelte";
  import { Dialog, Portal } from "@skeletonlabs/skeleton-svelte";
  import type { JSONError } from "../lib/statuses/jsonerror";
  import type { Snippet } from "svelte";
  import type { ResultAsync } from "neverthrow";
  import { IsResultAsync } from "../lib/apiutils/isResultAsync";

  const {
    title,
    callback = () => {},
    children,
    contentGetter,
  }: {
    title: string;
    callback?: () => void;
    children: Snippet<[{ content: T }]>;
    contentGetter: () => ResultAsync<T, JSONError> | T;
  } = $props();

  let open = $state(false);
  let content = $state<T | undefined>(undefined);
  let error = $state("");
  let loading = $state(false);
  let token = 0;

  async function load() {
    const mine = ++token;
    error = "";
    content = undefined;
    loading = true;

    const raw = contentGetter();
    if (IsResultAsync(raw)) {
      await raw.match(
        (v) => {
          if (mine === token) content = v;
        },
        (e) => {
          if (mine === token) error = e.error;
        },
      );
    } else {
      content = raw;
    }

    if (mine === token) loading = false;
  }

  function onOpenChange(details: { open: boolean }) {
    open = details.open;
    if (details.open) {
      load();
    } else {
      token++;
      content = undefined;
      loading = false;
    }
  }
</script>

<Dialog {open} {onOpenChange}>
  <Dialog.Trigger
    onclick={callback}
    class="btn mb-0 mt-0 preset-filled-surface-300-700 hover:preset-filled-warning-300-700 aspect-square h-auto w-auto p-1 m-px"
  >
    <Pencil size="16" />
  </Dialog.Trigger>

  <Portal>
    <Dialog.Backdrop class="fixed inset-0 z-50 bg-surface-50-950/50" />
    <Dialog.Positioner
      class="fixed inset-0 z-50 flex items-center justify-center"
    >
      <Dialog.Content
        class="card bg-surface-100-900 w-md p-4 space-y-4 shadow-xl"
      >
        <header class="flex justify-between items-center">
          <Dialog.Title class="text-2xl font-bold">{title}</Dialog.Title>
          <Dialog.CloseTrigger class="btn-icon hover:preset-tonal"
            ><X /></Dialog.CloseTrigger
          >
        </header>

        {#if open}
          {#if loading}
            <div>loading</div>
          {:else if error}
            {error}
          {:else if content !== undefined}
            {@render children({ content })}
          {/if}
        {/if}
      </Dialog.Content>
    </Dialog.Positioner>
  </Portal>
</Dialog>
