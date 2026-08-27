<script lang="ts">
  import XIcon from "@lucide/svelte/icons/x";
  import { Dialog, Portal } from "@skeletonlabs/skeleton-svelte";

  let {
    title,
    children,
    open = $bindable(false),
  }: {
    title: string;
    children: any;
    open: boolean;
  } = $props();
</script>

<Dialog
  {open}
  onOpenChange={(details) => {
    open = details.open;
  }}
>
  <Portal>
    <Dialog.Backdrop class="fixed inset-0 z-50 bg-surface-50-950/50" />
    <Dialog.Positioner
      class="fixed inset-0 z-50 flex justify-center items-center p-4"
    >
      <Dialog.Content
        class="card bg-surface-100-900 w-full max-w-[75dvw] h-full max-h-[75dvh] p-4 space-y-4 shadow-xl"
      >
        <div class="h-full">
          <header class="flex justify-between items-center">
            <Dialog.Title class="text-lg font-bold">{title}</Dialog.Title>
            <Dialog.CloseTrigger class="btn-icon hover:preset-tonal">
              <XIcon class="size-4" />
            </Dialog.CloseTrigger>
          </header>
          {#if open}
            {@render children()}
          {/if}
        </div>
      </Dialog.Content>
    </Dialog.Positioner>
  </Portal>
</Dialog>
