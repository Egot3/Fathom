<script lang="ts">
  import ArrowLeftIcon from "@lucide/svelte/icons/arrow-left";
  import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
  import {
    GetTotalsForUser,
    type TotalsOrError,
  } from "../lib/contracts/totals";
  import { GetUserOrRedirect } from "../lib/bgdata/user.svelte";
  import { IsJSONError } from "../lib/statuses/jsonerror";
  import Pagination from "./Pagination.svelte";

  let height = $state(0);

  let page = $state(1);
  let pageSize = $derived(Math.trunc((height - 39 - 45) / 39));

  let time: number;
  const paginatedTotalPromises = $derived.by(() => {
    const p = page;
    const ps = pageSize;

    console.log("detected change");
    clearTimeout(time);

    return new Promise<TotalsOrError>((resolve) => {
      time = setTimeout(() => {
        resolve(GetTotalsForUser(GetUserOrRedirect().UUID, p - 1, ps));
      }, 500);
    });
  });
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
      <table class="table table-auto self-start">
        <thead>
          <tr class="text-surface-100-900">
            <th>Test</th>
            <th>Group</th>
            <th>Score</th>
          </tr>
        </thead>

        <tbody>
          {#each paginatedTotals.totals as total}
            <tr>
              <td>{total.test_name}</td>
              <!-- 10.5 rem = 168px -->
              <td>{total.group_name}</td>
              <td>{total.score}</td>
            </tr>
          {/each}
        </tbody>
      </table>

      <div class="flex justify-between items-center gap-4 w-full self-end">
        <Pagination bind:page {pageSize} total={paginatedTotals.total} />
      </div>
    {/if}
  {/await}
</div>
