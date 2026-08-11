<script lang="ts">
  import { Pencil, X } from "@lucide/svelte";
  import { Dialog, Portal } from "@skeletonlabs/skeleton-svelte";
  import { FetchQuiz } from "../lib/contracts/quiz";
  import { IsJSONError } from "../lib/statuses/jsonerror";

  const {
    title,
    callback = () => {},
    children,
    contentGetter,
  }: {
    title: string;
    callback: () => void;
    children: any;
    contentGetter: () => Record<string, any> | Promise<Record<string, any>>;
  } = $props();

  let content: Record<string, any> | Promise<Record<string, any>> | null =
    $state(null);
</script>

<Dialog>
  <Dialog.Trigger
    onclick={() => {
      content = contentGetter();
      callback();
    }}
    class="btn mb-0 mt-0 preset-filled-surface-300-700 hover:preset-filled-warning-300-700 aspect-square h-auto w-auto p-1 m-px"
    ><Pencil size="16" /></Dialog.Trigger
  >

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

        {#if content != null}
          {#await content}
            <div>loading...</div>
          {:then content}
            {#if IsJSONError(content)}
              <div>{content.error}</div>
            {:else}
              {@render children({ content })}
            {/if}
          {/await}
        {/if}
      </Dialog.Content>
    </Dialog.Positioner>
  </Portal>
</Dialog>
