<script lang="ts">
  import { FetchTotals } from "../lib/contracts/totals";
  import { CreateFitRows } from "../lib/layoututils/fitRows.svelte";
  import { CreatePagedQuery } from "../lib/layoututils/createPagedQuery.svelte";
  import { FetchUsers, type User } from "../lib/contracts/user";
  import { FetchGroups, type Group } from "../lib/contracts/group";
  import { FetchAllTests, type Test } from "../lib/contracts/test";
  import SingleSelectPopover from "./SingleSelectPopover.svelte";
  import PeekDialogBig from "./PeekDialogBig.svelte";
  import PeekTotal from "./PeekTotal.svelte";
  import Pagination from "./Pagination.svelte";

  let {
    userUUID,
    filters = ["group", "test"],
  }: {
    userUUID?: string;
    filters?: readonly ("group" | "test")[];
  } = $props();

  const showUser = $derived(!userUUID);
  const colClass = $derived(showUser ? "w-1/5" : "w-1/4");

  const fit = CreateFitRows({ rowHeight: 20, reserve: 1 });

  let selected = $state({
    user: null as string | null,
    group: null as string | null,
    test: null as string | null,
  });

  const q = CreatePagedQuery({
    pageSize: () => fit.pageSize,
    enabled: () => fit.ready,
    params: () => ({
      user: userUUID ?? selected.user ?? undefined,
      group: selected.group ?? undefined,
      test: selected.test ?? undefined,
    }),
    fetcher: (page, size, p) =>
      FetchTotals(page, size, p.user, p.group, p.test),
    debounce: 500,
  });

  let focused = $state("");
</script>

<div class="flex flex-col h-full w-full gap-2">
  <div class="flex gap-2">
    {#if showUser}
      <SingleSelectPopover
        label="User"
        fetcher={FetchUsers}
        itemLabel={(u: User) => u.nickname}
        itemKey={(u: User) => u.uuid}
        items={(r) => r.users}
        bind:selected={selected.user}
      />
    {/if}
    {#if filters.includes("group")}
      <SingleSelectPopover
        label="Group"
        fetcher={FetchGroups}
        itemLabel={(g: Group) => g.name}
        itemKey={(g: Group) => g.uuid}
        items={(r) => r.groups}
        bind:selected={selected.group}
      />
    {/if}
    {#if filters.includes("test")}
      <SingleSelectPopover
        label="Test"
        fetcher={FetchAllTests}
        itemLabel={(t: Test) => t.name}
        itemKey={(t: Test) => t.uuid}
        items={(r) => r.tests}
        bind:selected={selected.test}
      />
    {/if}
  </div>

  <div class="flex-1 min-h-0 overflow-auto" use:fit.measure>
    {#if q.loading}
      <div
        class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
      ></div>
    {:else if q.error}
      <div>{q.error}</div>
    {:else if !q.data || q.data.total === 0}
      <div>No results</div>
    {:else}
      <table class="table table-auto w-full">
        <thead>
          <tr style="height:20px" class="text-surface-100-900 flex">
            {#if showUser}<th class={colClass}>User</th>{/if}
            <th class={colClass}>Group</th>
            <th class={colClass}>Test</th>
            <th class={colClass}>Score</th>
            <th class={colClass}>Totalized at</th>
          </tr>
        </thead>
        <tbody>
          {#each q.data.totals as total (total.group_uuid + total.test_uuid + total.user_uuid)}
            {const id = total.group_uuid + total.test_uuid + total.user_uuid}
            <tr
              style="height:20px"
              onclick={() => (focused = id)}
              class="bg-surface-700-300 rounded-xl flex hover:motion-safe:hover:brightness-125 dark:hover:motion-safe:hover:brightness-75"
            >
              {#if showUser}<td class={colClass}>{total.user_name}</td>{/if}
              <td class={colClass}>{total.group_name}</td>
              <td class={colClass}>{total.test_name}</td>
              <td class={colClass}>{total.score}/{total.max_score}</td>
              <td class={colClass}>{total.finalized_at}</td>
              <td class="max-w-0">
                <PeekDialogBig
                  bind:open={
                    () => focused === id, (v) => (focused = v ? id : "")
                  }
                  title="Totals peeker"
                >
                  <div class="flex flex-col h-full w-full gap-2">
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
                  </div>
                </PeekDialogBig>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    {/if}
  </div>

  <footer class="flex justify-between items-center gap-4 w-full">
    <Pagination
      bind:page={q.page}
      pageSize={fit.pageSize}
      total={q.data?.total ?? 0}
    />
  </footer>
</div>
