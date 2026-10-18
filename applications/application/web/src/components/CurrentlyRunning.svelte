<script lang="ts">
  import { onMount } from "svelte";
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import { GetUser } from "../lib/bgdata/user.svelte";
  import { CreateRunningTests } from "../lib/bgdata/runningTests.svelte";
  import ChipSelector from "./ChipSelector.svelte";
  import PeekDialogue from "./PeekDialogue.svelte";
  import PeekLink from "./PeekLink.svelte";
  import ReloadButton from "./ReloadButton.svelte";
  import RunningControls from "./RunningControls.svelte";
  import TestEnder from "./TestEnder.svelte";
  import TestExtender from "./TestExtender.svelte";
  import TestStarter from "./TestStarter.svelte";

  const { inTeachery }: { inTeachery?: boolean } = $props();

  const isTeacher = $derived((GetUser()?.isTeacher ?? false) && inTeachery);
  const tests = CreateRunningTests();

  let selectedKey: string | null = $state(null);
  const chosen = $derived(
    selectedKey === null
      ? (tests.data.find((t) => t.key === selectedKey) ?? tests.data[0])
      : tests.data[0],
  ); // something malicious is brewing
  const chosenIdx = $derived(Math.max(0, tests.data.indexOf(chosen)));
  const deadline = $derived(
    chosen?.deadline ? new Date(chosen.deadline).toLocaleString() : "",
  );

  onMount(() => tests.refresh());
</script>

<article class="flex flex-col h-full space-y-5">
  <div class="flex justify-end">
    <ReloadButton onreload={() => tests.refresh()} />
  </div>

  {#if tests.loading}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:else if tests.error}
    <div>{tests.error}</div>
    <button
      class="btn preset-filled-warning-500"
      onclick={() => tests.refresh()}>Reload?</button
    >
  {:else if !chosen}
    <div>NOTHING</div>
  {:else}
    {#if tests.data.length > 1}
      <span>
        <ChipSelector
          options={tests.data.map((t) => t.name)}
          bind:selected={
            () => chosenIdx, (i) => (selectedKey = tests.data[i]?.key ?? null)
          }
        />
      </span>
    {/if}

    <p>Test {chosen.name}</p>

    <div class="flex space-x-1">
      Deadline: {deadline}
      {#if isTeacher}
        <PeekDialogue title="Test extender">
          {#snippet trigger()}
            <Dialog.Trigger class="chip preset-outlined-primary-500 ml-auto"
              >Extend</Dialog.Trigger
            >
          {/snippet}
          <TestExtender callback={() => tests.refresh(true)} key={chosen.key} />
        </PeekDialogue>
      {/if}
    </div>

    {#if isTeacher}
      <RunningControls test={chosen} onchanged={() => tests.refresh(true)} />
      <PeekLink key={chosen.key} />
    {:else if chosen.is_paused}
      <div class="badge preset-tonal-warning w-fit">Paused</div>
    {:else}
      <a class="mb-auto btn preset-filled-primary-500" href={`/${chosen.key}`}
        >Join?</a
      >
    {/if}
  {/if}

  {#if isTeacher}
    <div
      class="mt-auto flex flex-col lg:flex-row space-x-0 lg:space-x-1 space-y-1 lg:space-y-0"
    >
      <PeekDialogue title="Test starter">
        {#snippet trigger()}
          <Dialog.Trigger class="btn preset-filled-primary-500"
            >Start new</Dialog.Trigger
          >
        {/snippet}
        <TestStarter callback={() => tests.refresh(true)} />
      </PeekDialogue>

      <PeekDialogue title="Test ender">
        {#snippet trigger()}
          <Dialog.Trigger
            class="btn preset-outlined-error-500"
            disabled={tests.loading || !chosen}
          >
            End now
          </Dialog.Trigger>
        {/snippet}
        <TestEnder
          key={chosen?.key}
          name={chosen?.name}
          callback={() => tests.refresh(true)}
        />
      </PeekDialogue>
    </div>
  {/if}
</article>
