<script lang="ts">
  import { FetchTotals, type Totals } from "../lib/contracts/totals";
  import PeekDialogBig from "./PeekDialogBig.svelte";
  import PeekTotal from "./PeekTotal.svelte";
  import { FetchUsers, type User } from "../lib/contracts/user";
  import { FetchGroups, type Group } from "../lib/contracts/group";
  import { FetchAllTests, type Test } from "../lib/contracts/test";
  import SingleSelectPopover from "./SingleSelectPopover.svelte";
  import { CreateFitRows } from "../lib/layoututils/fitRows.svelte";
  import Pagination from "./Pagination.svelte";

  const fit = CreateFitRows({ rowHeight: 20, reserve: 1 });
  let pageSize = $derived(fit.pageSize);
  let page = $state(1);

  let loading = $state(true);
  let statusMessage = $state("");

  let selectedUser = $state("");
  let selectedGroup = $state("");
  let selectedTest = $state("");

  let trigger = $state(0);
  let time: number;
  let paginatedTotals: Totals = $state(null as never);
  $effect(() => {
    const p = page;
    const ps = pageSize;
    const user = selectedUser;
    const group = selectedGroup;
    const test = selectedTest;
    const ready = fit.ready;
    trigger;

    if (!ready) return;

    clearTimeout(time);
    loading = true;

    new Promise((resolve) => {
      time = setTimeout(async () => {
        statusMessage = await FetchTotals(
          p - 1,
          ps,
          user || "all",
          group || "all",
          test || "all",
        )
          .andTee((r) => {
            paginatedTotals = r;
          })
          .match(
            (_) => "",
            (err) => err.error,
          );

        loading = false;

        resolve(0);
      }, 500);
    });
  });

  let focused = $state("");
</script>

<div
  class="grid gap-4 w-full place-items-center h-full overflow-auto"
  use:fit.measure
>
  {#if loading}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:else}
    {#if statusMessage !== ""}
      <div>{statusMessage}</div>
    {:else}
      {#if paginatedTotals.total === 0}
        No total registered
        <!-- probably useless, 1 total is the one, watching it -->
      {:else}
        <table class="table table-auto self-start">
          <thead>
            <tr style="height:20px" class="text-surface-100-900 flex">
              <th class="w-1/5 flex gap-2 items-center">
                User
                <SingleSelectPopover
                  label="User"
                  fetcher={FetchUsers}
                  itemLabel={(u: User) => u.nickname}
                  itemKey={(u: User) => u.uuid}
                  items={(t) => t.users}
                  bind:selected={selectedUser}
                />
              </th>
              <th class="w-1/5 flex gap-2 items-center">
                Group
                <SingleSelectPopover
                  label="Group"
                  fetcher={FetchGroups}
                  itemLabel={(g: Group) => g.name}
                  itemKey={(g: Group) => g.uuid}
                  items={(g) => g.groups}
                  bind:selected={selectedGroup}
                />
              </th>
              <th class="w-1/5 flex gap-2 items-center">
                Test
                <SingleSelectPopover
                  label="Test"
                  fetcher={FetchAllTests}
                  itemLabel={(t: Test) => t.name}
                  itemKey={(t: Test) => t.uuid}
                  items={(t) => t.tests}
                  bind:selected={selectedTest}
                /></th
              >
              <th class="w-1/5">Score</th>
              <th class="w-1/5">Totalized at</th>
            </tr>
          </thead>

          <tbody>
            {#each paginatedTotals.totals as total (total.group_uuid + total.test_uuid + total.user_uuid)}
              <!-- Ultimate IDX -->
              {const ultimateUUID =
                total.group_uuid + total.test_uuid + total.user_uuid}
              <tr
                style="height:20px"
                onclick={() => (focused = ultimateUUID)}
                class="bg-surface-700-300 rounded-xl flex hover:motion-safe:hover:brightness-125 dark:hover:motion-safe:hover:brightness-75"
              >
                <td class="w-1/5">
                  {total.user_name}
                </td>
                <td class="w-1/5">{total.group_name}</td>
                <td class="w-1/5">{total.test_name}</td>
                <td class="w-1/5">{total.score}</td>
                <td class="w-1/5">{total.finalized_at}</td>
                <td class="max-w-0">
                  <PeekDialogBig
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
                  </PeekDialogBig>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>

        <div class="flex justify-between items-center gap-4 w-full self-end">
          <Pagination bind:page {pageSize} total={paginatedTotals.total} />
        </div>
      {/if}
    {/if}
  {/if}
</div>
