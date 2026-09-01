<script lang="ts">
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import { IsJSONError, type JSONError } from "../lib/statuses/jsonerror";
  import ChipSelector from "./ChipSelector.svelte";
  import PeekDialogue from "./PeekDialogue.svelte";
  import TestStarter from "./TestStarter.svelte";
  import {
    FetchCurrentlyRunningTestInfo,
    type Test,
    type TestInfo,
  } from "../lib/contracts/test";

  let currentlyRunning: TestInfo | null = $state(null);
  let trig = $state(0);
  let loading = $state(true);
  let statusMessage = $state("");

  $effect(() => {
    trig;

    loading = true;

    (async () => {
      currentlyRunning = await FetchCurrentlyRunningTestInfo()
        .andTee((r) => {
          loading = false;
        })
        .match(
          (r) => r,
          (err) => {
            statusMessage = err.error;
            return null;
          },
        );
    })();
  });
</script>

<article class="flex flex-col h-full space-y-5">
  {#if loading}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:else}
    {#if statusMessage !== ""}
      <div>{statusMessage}</div>
      <button
        class="btn preset-filled-warning-500"
        onclick={() => {
          trig++;
        }}>Reload?</button
      >
    {:else}
      {#if currentlyRunning === null}
        <div>NOTHING</div>
      {:else}
        <p>Test {currentlyRunning.test.name}</p>
        <div class="flex space-x-1">
          Deadline: {currentlyRunning.deadline}
          <button class="chip preset-outlined-primary-500">Extend</button>
        </div>
        <ChipSelector
          options={["running", "paused"]}
          selected={currentlyRunning.isPaused ? 1 : 0}
        ></ChipSelector>

        <div class="flex space-x-1">
          <button class="chip preset-outlined-primary-500">Add quizzes</button>
          <button class="chip preset-outlined-error-500">Prune quizzes</button>
        </div>
      {/if}
    {/if}
    <div
      class="mt-auto flex flex-col lg:flex-row space-x-0 lg:space-x-1 space-y-1 lg:space-y-0"
    >
      <PeekDialogue title="Test starter">
        {#snippet trigger()}
          <Dialog.Trigger
            class="btn preset-filled-primary-500"
            onclick={() => {}}>Start new</Dialog.Trigger
          >
        {/snippet}
        <TestStarter
          callback={() => {
            trig++;
          }}
        />
      </PeekDialogue>

      <button class="btn preset-outlined-error-500" disabled={!loading}
        >End test</button
      >
      <!-- not stop as it could be confused for pause -->
    </div>
  {/if}
</article>
