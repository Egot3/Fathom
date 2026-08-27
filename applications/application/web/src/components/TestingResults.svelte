<script lang="ts">
  import ArrowLeftIcon from "@lucide/svelte/icons/arrow-left";
  import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
  import { Pagination } from "@skeletonlabs/skeleton-svelte";
  import { IsJSONError } from "../lib/statuses/jsonerror";
  import { FetchTotals, type TotalsOrError } from "../lib/contracts/totals";
  import PeekDialog from "./PeekDialog.svelte";
  import PeekTotal from "./PeekTotal.svelte";

  let height = $state(0);

  let page = $state(1);
  let pageSize = $derived(Math.trunc((height - 29 - 45 - 40) / 49));

  let trigger = $state(0);
  let time: number;
  const paginatedTotalPromises = $derived.by(() => {
    const p = page;
    const ps = pageSize;
    trigger; // well well well, it updates as an int(increment) and derived.by fires

    clearTimeout(time);

    return new Promise<TotalsOrError>((resolve) => {
      time = setTimeout(() => {
        resolve(FetchTotals(p - 1, ps));
      }, 500);
    });
  });

  let focused = $state("");
</script>

<div
  class="grid gap-4 w-full place-items-center h-full overflow-auto"
  bind:clientHeight={height}
>
  {#await paginatedTotalPromises}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:then paginatedTotals}
    {#if IsJSONError(paginatedTotals)}
      <div>{paginatedTotals.error}</div>
    {:else}
      {#if paginatedTotals.total == 0}
        No total registered
        <!-- probably useless, 1 total is the one, watching it -->
      {:else}
        <table class="table table-auto self-start">
          <thead>
            <tr class="text-surface-100-900 flex">
              <th class="w-1/5">User</th>
              <th class="w-1/5">Group</th>
              <th class="w-1/5">Test</th>
              <th class="w-1/5">Score</th>
              <th class="w-1/5">Totalized at</th>
            </tr>
          </thead>

          <tbody>
            {#each paginatedTotals.totals as total ((total.group_uuid, total.test_uuid, total.user_uuid))}
              <!-- Ultimate IDX -->
              {const ultimateUUID =
                total.group_uuid + total.test_uuid + total.user_uuid}
              <tr
                onclick={() => (focused = ultimateUUID)}
                class="bg-surface-700-300 rounded-xl flex hover:motion-safe:hover:brightness-125 dark:hover:motion-safe:hover:brightness-75"
              >
                <td class="w-1/5">{total.user_name}</td>
                <td class="w-1/5">{total.group_name}</td>
                <td class="w-1/5">{total.test_name}</td>
                <td class="w-1/5">{total.score}</td>
                <td class="w-1/5">{total.finalized_at}</td>
                <td class="max-w-0">
                  <PeekDialog
                    bind:open={
                      () => focused === ultimateUUID,
                      (v) => (focused = v ? ultimateUUID : "")
                    }
                    title="Totals peeker"
                  >
                    <div class="flex">
                      <h2 class="text-2xl">{total.test_name}</h2>
                      <div
                        class="ml-auto text-2xl chip preset-outlined-surface-400-600"
                      >
                        {total.score}/{total.max_score}
                      </div>
                    </div>

                    <PeekTotal
                      userUUID={total.user_uuid}
                      groupUUID={total.group_uuid}
                      testUUID={total.test_uuid}
                    />
                  </PeekDialog>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>

        <div class="flex justify-between items-center gap-4 w-full self-end">
          <Pagination
            count={paginatedTotals?.total ?? 0}
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
  {/await}
</div>
