<script lang="ts">
  import { Pagination } from "@skeletonlabs/skeleton-svelte";
  import { FetchAllTests, type Tests } from "../lib/contracts/test";
  import type { JSONError } from "../lib/statuses/jsonerror";
  import { ArrowLeftIcon, ArrowRightIcon, CheckIcon } from "@lucide/svelte";

  let chosen = $state("");
  let page = $state(1);
  let pageSize = $state(5);

  let statusMessage = $state("");
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
      statusMessage = await FetchAllTests(p - 1, ps)
        .map((r) => {
          listTests = r;
          console.log(listTests);
          loading = false;
          return r;
        })
        .match(
          () => "",
          (err: JSONError) => err.error,
        );
    });
  });
</script>

{#if loading}
  <div class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"></div>
{:else}
  {#if statusMessage}
    <div class="h-full w-full rounded-xl">
      <p>{statusMessage}</p>
    </div>
  {:else}
    {#if listTests.total !== 0}
      {#each listTests.tests as test (test.uuid)}
        <button
          type="button"
          class={`chip capitalize preset-outlined-surface-400-600 ${chosen === test.uuid ? "preset-tonal-primary" : ""}`}
          onclick={() => {
            chosen;
          }}
        >
          {#if chosen === test.uuid}<CheckIcon size={14} />{/if}
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
    {/if}
  {/if}
{/if}
