<script lang="ts">
  import { SvelteSet } from "svelte/reactivity";
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import { FetchGroupPost } from "../lib/contracts/group";
  import UserChips from "./UserChips.svelte";

  const { callback }: { callback: () => void } = $props();

  let name: string = $state("");
  let users: SvelteSet<string> = $state(new SvelteSet<string>());
  let statusMessage: string = $state("");

  async function PostGroup(e: Event) {
    e.preventDefault();
    const postResponse = await FetchGroupPost(name, Array.from(users));
    if (postResponse != null) {
      statusMessage = postResponse.error;
      return;
    }

    callback();
  }
</script>

<form onsubmit={PostGroup} class="w-full flex flex-col h-full">
  <label class="label">
    <span class="label-text">Name</span>
    <input
      class="input"
      type="text"
      bind:value={name}
      placeholder="k*3"
      required
    />
  </label>
  <label class="label">
    <span class="label-text">Users</span>
    <UserChips bind:chosen={users} />
  </label>

  <footer class="flex justify-end gap-2">
    {#if statusMessage != ""}
      <span class="justify-self-start">{statusMessage}</span>
    {/if}
    <Dialog.CloseTrigger class="btn preset-tonal">Cancel</Dialog.CloseTrigger>
    <button type="submit" class="btn preset-filled">Create</button>
  </footer>
</form>
