<script lang="ts">
  import { GetCurrentlyRunning } from "../lib/bgdata/currentlyrunning.svelte";
  import { IsJSONError, type JSONError } from "../lib/statuses/jsonerror";
  import ChipSelector from "./ChipSelector.svelte";
  import TestStarter from "./TestStarter.svelte";

  let currentlyRunningPromise = $state(GetCurrentlyRunning());

  let loaded = $state(false);

  let choosing = $state(false);
</script>

<article class="flex flex-col h-full space-y-5">
  {#await currentlyRunningPromise}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:then currentlyRunning}
    {#if currentlyRunning === null}
      <div>NOTHING</div>
    {:else}
      {#if IsJSONError(currentlyRunning)}
        <div>{currentlyRunning.error}</div>
        <button
          class="btn preset-filled-warning-500"
          onclick={() => {
            currentlyRunningPromise = GetCurrentlyRunning();
          }}>Reload?</button
        >
      {:else}
        {(loaded = true)}
        <p>Test {currentlyRunning.Name}</p>
        <div class="flex space-x-1">
          Deadline: {currentlyRunning.Deadline}
          <button class="chip preset-outlined-primary-500">Extend</button>
        </div>
        <ChipSelector
          options={["running", "paused"]}
          selected={currentlyRunning.IsPaused ? 1 : 0}
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
      <button
        class="btn preset-filled-primary-500"
        onclick={() => (choosing = true)}>Start new</button
      >
      <!-- <PeekDialog title="Test chooser" open={choosing}>
        <TestStarter></TestStarter>
      </PeekDialog> -->

      <button class="btn preset-outlined-error-500" disabled={!loaded}
        >End test</button
      >
      <!-- not stop as it could be confused for pause -->
    </div>
  {/await}
</article>
