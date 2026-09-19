<script lang="ts">
  import { Pagination } from "@skeletonlabs/skeleton-svelte";
  import {
    FetchAllTests,
    FetchTestStart,
    type Tests,
  } from "../lib/contracts/test";
  import type { JSONError } from "../lib/statuses/jsonerror";
  import { ArrowLeftIcon, ArrowRightIcon, CheckIcon } from "@lucide/svelte";
  import { SvelteSet } from "svelte/reactivity";
  import GroupChips from "./GroupChips.svelte";

  const { callback }: { callback: () => void } = $props();

  let hours: number = $state(0);
  let minutes: number = $state(0);
  let seconds: number = $state(0);

  let chosenTest = $state("");
  let chosenGroups: SvelteSet<string> = $state(new SvelteSet<string>());
  let duration: string = $derived(`${hours}h${minutes}m${seconds}s`);

  let page = $state(1);
  let pageSize = $state(5);

  let statusMessage = $state("");
  let startStatusMessage = $state("");

  let loading = $state(true);
  let listTests: Tests = $state(null as never);

  let trigger = $state(0);
  let time: number;
  $effect(() => {
    const p = page;
    const ps = pageSize;
    trigger;
    loading = true;

    console.log("detected change");
    clearTimeout(time);

    time = setTimeout(async () => {
      startStatusMessage = await FetchAllTests(p - 1, ps)
        .map((r) => {
          listTests = r;
          loading = false;
          return r;
        })
        .match(
          () => "",
          (err: JSONError) => err.error,
        );
    });
  });

  async function startTest(e: SubmitEvent) {
    e.preventDefault();

    if (chosenTest === "") {
      statusMessage = "can't run test with no test selected";
      return;
    }
    if (chosenGroups.size === 0) {
      statusMessage = "can't run test with no groups selected";
      return;
    }
    if (duration === "0h0m0s") {
      statusMessage = "can't run test with no duration";
      return;
    }

    statusMessage = await FetchTestStart(
      chosenTest,
      duration,
      Array.from(chosenGroups),
    ).match(
      (r) => {
        callback();
        return "";
      },
      (e) => e.error,
    );

    console.log("resp got");

    return;
  }
</script>

{#if loading}
  <div class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"></div>
{:else}
  {#if startStatusMessage}
    <div class="h-full w-full rounded-xl">
      <p>{startStatusMessage}</p>
    </div>
  {:else}
    {#if listTests.total !== 0}
      <form onsubmit={startTest}>
        {#each listTests.tests as test (test.uuid)}
          <button
            type="button"
            class={`chip preset-outlined-surface-400-600 ${chosenTest === test.uuid ? "preset-tonal-primary" : ""}`}
            onclick={() => {
              chosenTest = test.uuid;
            }}
          >
            {#if chosenTest === test.uuid}<CheckIcon size={14} />{/if}
            <span>{test.name}</span>
          </button>
        {/each}

        <div class="flex justify-between items-center gap-4 w-full self-end">
          <Pagination
            count={listTests?.total ?? 0}
            {pageSize}
            {page}
            onPageChange={(event) => (page = event.page)}
            class="rounded-xl"
          >
            <Pagination.PrevTrigger>
              <ArrowLeftIcon class="size-4" />
            </Pagination.PrevTrigger>
            <Pagination.Context>
              {#snippet children(pagination)}
                {#each pagination().pages as page, index (page)}
                  {#if page.type === "page"}
                    <Pagination.Item {...page}>
                      {page.value}
                    </Pagination.Item>
                  {:else}
                    <Pagination.Ellipsis {index}>…</Pagination.Ellipsis>
                  {/if}
                {/each}
              {/snippet}
            </Pagination.Context>
            <Pagination.NextTrigger>
              <ArrowRightIcon class="size-4" />
            </Pagination.NextTrigger>
          </Pagination>
        </div>

        <GroupChips bind:chosen={chosenGroups} />

        <fieldset>
          <legend></legend>
          <label class="label">
            <span class="label-text">Hours</span>
            <input
              type="number"
              class="input border-2"
              min="0"
              bind:value={hours}
            />
          </label>
          <label class="label">
            <span class="label-text">Minutes</span>
            <input
              type="number"
              class="input border-2"
              min="0"
              bind:value={minutes}
            />
          </label>
          <label class="label">
            <span class="label-text">Seconds</span>
            <input
              type="number"
              class="input border-2"
              min="0"
              bind:value={seconds}
            />
          </label>
        </fieldset>

        <div class="flex">
          <button
            disabled={chosenTest === "" ||
              chosenGroups.size === 0 ||
              duration === "0h0m0s"}
            class="btn preset-filled-brand"
            type="submit"
          >
            Run
          </button>
          {#if statusMessage !== ""}
            <div class="ml-auto badge preset-filled-error-100-900 text-xs">
              {statusMessage}
            </div>
          {/if}
        </div>
      </form>
    {/if}
  {/if}
{/if}
