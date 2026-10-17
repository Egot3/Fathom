<script lang="ts" generics="T">
  import { Eye, EyeClosed } from "@lucide/svelte";
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import PeekDialogue from "./PeekDialogue.svelte";
  import type { ResultAsync } from "neverthrow";
  import type { JSONError } from "../lib/statuses/jsonerror";
  import type { Snippet } from "svelte";

  const {
    title,
    children,
    callback = () => {},
    contentGetter,
  }: {
    title: string;
    children: Snippet<[{ content: T }]>;
    callback?: () => void;
    contentGetter: () => ResultAsync<T, JSONError> | T;
  } = $props();

  let moused = $state(false);
  let ct: ResultAsync<T, JSONError> | T | null = $state(null);
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
