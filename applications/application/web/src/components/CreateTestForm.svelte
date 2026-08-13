<script lang="ts">
  import { SvelteSet } from "svelte/reactivity";
  import QuizChips from "./QuizChips.svelte";
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import { FetchTestPost } from "../lib/contracts/test";

  const { callback }: { callback: () => void } = $props();

  let name: string = $state("");
  let quizzes: SvelteSet<string> = $state(new SvelteSet<string>());
  let statusMessage: string = $state("");

  async function PostTest(e: Event) {
    e.preventDefault();
    const postResponse = await FetchTestPost(name, Array.from(quizzes));
    if (postResponse !== null) {
      statusMessage = postResponse.error;
      return;
    }

    callback();
  }

  $inspect(name);
</script>

<form onsubmit={PostTest} class="w-full flex flex-col h-full">
  <label class="label">
    <span class="label-text">Name</span>
    <input
      class="input"
      type="text"
      bind:value={name}
      placeholder="go-generics"
      required
    />
  </label>
  <label class="label">
    <span class="label-text">Quizzes</span>
    <QuizChips bind:chosen={quizzes} />
  </label>

  <footer class="flex justify-end gap-2">
    {#if statusMessage !== ""}
      <span class="justify-self-start">{statusMessage}</span>
    {/if}
    <Dialog.CloseTrigger class="btn preset-tonal">Cancel</Dialog.CloseTrigger>
    <button type="submit" class="btn preset-filled">Create</button>
  </footer>
</form>
