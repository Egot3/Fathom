<script lang="ts">
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import UXInput from "./UXInput.svelte";
  import { CopyToClipboard } from "../lib/apiutils/copy";
  import type { JSONError } from "../lib/statuses/jsonerror";

  const {
    name,
    UUID,

    deletor,
    callback = () => {},
  }: {
    name: string;
    UUID: string;

    deletor: (UUID: string) => Promise<null | JSONError>;
    callback: () => void;
  } = $props();

  let insertedName: string = $state("");
  let statusMessage: string = $state("");
  let nameReady: boolean = $state(false);

  async function DeleteQuiz(e: Event) {
    e.preventDefault();

    if (!nameReady) {
      statusMessage = "Inserted name is not equal to requested";
      return;
    }

    const deleteResponse = await deletor(UUID);
    if (deleteResponse != null) {
      statusMessage = deleteResponse.error;
      return;
    }

    callback();
  }
</script>

<form onsubmit={DeleteQuiz} class=" w-full flex flex-col h-full space-y-2">
  <p class="text-xl">
    Please enter

    <button
      type="button"
      class="italic hover:cursor-pointer"
      onclick={() => {
        CopyToClipboard(name);
      }}>{name}</button
    >

    below to confirm deletion.
  </p>
  <p class="text-xs opacity-75 -mt-2">
    You may click the text in italic to copy it
  </p>

  <UXInput
    placeholder={name}
    bind:value={insertedName}
    bind:ready={nameReady}
    checker={(v: string) => {
      return name === v;
    }}
  />

  <footer class="flex justify-end gap-2 w-full">
    {#if statusMessage != ""}
      <span class="justify-self-start">{statusMessage}</span>
    {/if}
    <Dialog.CloseTrigger class="btn preset-tonal">Cancel</Dialog.CloseTrigger>
    <button type="submit" class="btn preset-filled" disabled={nameReady}
      >Delete</button
    >
  </footer>
</form>
