<script lang="ts">
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import { IsJSONError, type JSONError } from "../lib/statuses/jsonerror";
  import ChipSelector from "./ChipSelector.svelte";
  import PeekDialogue from "./PeekDialogue.svelte";
  import TestStarter from "./TestStarter.svelte";
  import {
    FetchCurrentlyRunningTestInfos,
    FetchTestPause,
    FetchTestResume,
    type TestInfo,
  } from "../lib/contracts/test";
  import TestEnder from "./TestEnder.svelte";
  import TestExtender from "./TestExtender.svelte";
  import { onMount } from "svelte";

  let isCurrentlyRunning: boolean = $state(true);
  let currentlyRunning: TestInfo[] = $state([]);

  let trig = $state(0);
  let loading = $state(true);
  let statusMessage = $state("");

  let chosenTestId = $state(0);
  let chosenTest = $derived(currentlyRunning?.[chosenTestId]);
  let deadline = $derived(
    chosenTest?.deadline ? new Date(chosenTest.deadline).toLocaleString() : "",
  );

  let selectedId = $state(0);
  let pauseresuming = $state(false);
  let pauseresumeMessage: null | JSONError = $state(null);

  async function refresh() {
    loading = true;
    (
      await FetchCurrentlyRunningTestInfos().andTee((r) => {
        isCurrentlyRunning = r.length !== 0;
      })
    ).match(
      (r) => (currentlyRunning = r),
      (err) => {
        statusMessage = err.error;
        currentlyRunning = [];
      },
    );
    loading = false;
  }

  $effect(() => {
    if (!pauseresuming && chosenTest !== undefined) {
      selectedId = chosenTest.is_paused ? 1 : 0;
    }
  });

  async function pause(key: string) {
    pauseresumeMessage = await FetchTestPause(key).match(
      (r) => r,
      (e) => e,
    );
  }
  async function resume(key: string) {
    pauseresumeMessage = await FetchTestResume(key).match(
      (r) => r,
      (e) => e,
    );
  }

  async function togglePauseResume(nextSelected: number) {
    if (pauseresuming) return;
    pauseresuming = true;
    try {
      if (nextSelected === 1) {
        await pause(chosenTest.key);
      } else {
        await resume(chosenTest.key);
      }
      await refresh();
    } finally {
      pauseresuming = false;
    }
  }

  onMount(() => {
    refresh();
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
      <button class="btn preset-filled-warning-500" onclick={refresh}
        >Reload?</button
      >
    {:else}
      {#if !isCurrentlyRunning || currentlyRunning.length === 0}
        <div>NOTHING</div>
      {:else}
        <span>
          <ChipSelector
            options={currentlyRunning.map((e) => e.name)}
            bind:selected={chosenTestId}
          />
        </span>
        <p>Test {chosenTest.name}</p>
        <div class="flex space-x-1">
          Deadline: {deadline}

          <PeekDialogue title="Test extender">
            {#snippet trigger()}
              <Dialog.Trigger class="chip preset-outlined-primary-500 ml-auto"
                >Extend</Dialog.Trigger
              >
            {/snippet}
            <TestExtender callback={refresh} key={chosenTest.key} />
          </PeekDialogue>
        </div>
        <ChipSelector
          options={["running", "paused"]}
          working={pauseresuming}
          bind:selected={selectedId}
          onchange={() => togglePauseResume(selectedId)}
        />
      {/if}
    {/if}
    <div
      class="mt-auto flex flex-col lg:flex-row space-x-0 lg:space-x-1 space-y-1 lg:space-y-0"
    >
      <PeekDialogue title="Test starter">
        {#snippet trigger()}
          <Dialog.Trigger class="btn preset-filled-primary-500"
            >Start new</Dialog.Trigger
          >
        {/snippet}
        <TestStarter callback={refresh} />
      </PeekDialogue>

      <PeekDialogue title="Test ender">
        {#snippet trigger()}
          <Dialog.Trigger
            class="btn preset-outlined-error-500"
            disabled={loading || !chosenTest}>End now</Dialog.Trigger
          >
        {/snippet}
        <TestEnder
          key={chosenTest?.key}
          name={chosenTest?.name}
          callback={refresh}
        />
      </PeekDialogue>
      <!-- not stop as it could be confused for pause -->
    </div>
  {/if}
</article>
