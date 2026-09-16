<script lang="ts">
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import { IsJSONError, type JSONError } from "../lib/statuses/jsonerror";
  import ChipSelector from "./ChipSelector.svelte";
  import PeekDialogue from "./PeekDialogue.svelte";
  import TestStarter from "./TestStarter.svelte";
  import {
    FetchCurrentlyRunningTestInfos,
    type Test,
    type TestInfo,
  } from "../lib/contracts/test";

  let isCurrentlyRunning: boolean = $state(false);
  let currentlyRunning: TestInfo[] = $state(null as never);

  let trig = $state(0);
  let loading = $state(true);
  let statusMessage = $state("");

  let chosenId = $state(0);
  let chosen = $derived(currentlyRunning[chosenId]);

  $effect(() => {
    trig;

    loading = true;

    (async () => {
      currentlyRunning = (
        await FetchCurrentlyRunningTestInfos().andTee((r) => {
          loading = false;
          isCurrentlyRunning = r.length === 0;
        })
      ).match(
        (r) => r,
        (err) => {
          statusMessage = err.error;
          return [];
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
      {#if !isCurrentlyRunning}
        <div>NOTHING</div>
      {:else}
        <span>
          <ChipSelector
            options={currentlyRunning.map((e) => e.test.name)}
            bind:selected={chosenId}
          />
        </span>
        <p>Test {chosen.test.name}</p>
        <div class="flex space-x-1">
          Deadline: {chosen.deadline}
          <button class="chip preset-outlined-primary-500">Extend</button>
        </div>
        <ChipSelector
          options={["running", "paused"]}
          selected={chosen.isPaused ? 1 : 0}
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
