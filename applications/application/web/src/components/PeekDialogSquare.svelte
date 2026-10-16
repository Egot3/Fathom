<script lang="ts">
  import { Eye, EyeClosed, X } from "@lucide/svelte";
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import PeekDialogue from "./PeekDialogue.svelte";

  const {
    title,
    children,
    callback = () => {},
    contentGetter,
  }: {
    title: string;
    children: any;
    callback?: () => void;
    contentGetter: () => Promise<unknown>;
  } = $props();

  let moused = $state(false);
  let ct = $state<Promise<unknown>>(new Promise(() => {}));
</script>

<PeekDialogue content={ct} {title} {children}>
  {#snippet trigger()}
    <Dialog.Trigger
      onmouseenter={() => (moused = true)}
      onmouseleave={() => (moused = false)}
      onclick={() => {
        callback();
        ct = contentGetter();
      }}
      class="btn mb-0 mt-0 preset-filled-surface-300-700 hover:preset-filled-primary-300-700 aspect-square h-auto w-auto p-1 m-px"
    >
      {#if moused}
        <Eye size="16" />
      {:else}
        <EyeClosed size="16" />
      {/if}
    </Dialog.Trigger>
  {/snippet}
</PeekDialogue>
