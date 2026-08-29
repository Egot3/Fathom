<script lang="ts">
  import { X } from "@lucide/svelte";
  import { Dialog, Portal } from "@skeletonlabs/skeleton-svelte";
  import type { Snippet } from "svelte";

  const {
    title,
    children,
    content,
    trigger,
  }: {
    title: string;
    children: any;
    trigger: any;
    content?: Promise<unknown>;
  } = $props();
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

        {#await content}
          <div>loading</div>
        {:then content}
          {@render children?.({ content })}
        {/await}
      </Dialog.Content>
    </Dialog.Positioner>
  </Portal>
</Dialog>
