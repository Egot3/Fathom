<script lang="ts" generics="T">
  import { X } from "@lucide/svelte";
  import { Dialog, Portal } from "@skeletonlabs/skeleton-svelte";
  import { ResultAsync } from "neverthrow";
  import type { JSONError } from "../lib/statuses/jsonerror";
  import type { Snippet } from "svelte";

  const {
    title,
    children,
    content,
    trigger,
  }: {
    title: string;
    children: Snippet<[{ content: T }]>;
    trigger: Snippet;
    content?: ResultAsync<T, JSONError> | T | null;
  } = $props();

  function isResultAsync(
    value: unknown,
  ): value is ResultAsync<unknown, JSONError> {
    return (
      value instanceof ResultAsync ||
      (typeof value === "object" && value !== null && "asyncAndThen" in value)
    );
  }

  let error = $state("");
  async function getContent(): Promise<T | null | undefined> {
    if (!isResultAsync(content)) {
      return content;
    }

    return await content
      .orTee((r) => (error = r.error))
      .match(
        (r) => r,
        (_) => undefined,
      );
  }
</script>

<Dialog>
  {@render trigger()}

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

        {#await getContent()}
          <div>loading</div>
        {:then content}
          {#if error}
            {error}
          {:else if content}
            {@render children?.({ content })}
          {/if}
        {/await}
      </Dialog.Content>
    </Dialog.Positioner>
  </Portal>
</Dialog>
