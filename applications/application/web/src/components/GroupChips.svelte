<script lang="ts">
  import { CheckIcon } from "@lucide/svelte";
  import { SvelteSet } from "svelte/reactivity";
  import {
    FetchGroups,
    type Group,
    type ListGroupResponse,
  } from "../lib/contracts/group";
  import Pagination from "./Pagination.svelte";
  import { CreateFitRows } from "../lib/layoututils/fitRows.svelte";
  import { CreatePagedQuery } from "../lib/layoututils/createPagedQuery.svelte";

  let {
    chosen = $bindable(new SvelteSet<string>()),
    existing = new SvelteSet<Group>(),
  }: {
    chosen: SvelteSet<string>;
    existing?: SvelteSet<Group>;
  } = $props();
  $inspect(chosen);

  const fit = CreateFitRows({ rowHeight: 20, reserve: 1 });

  const q = CreatePagedQuery({
    pageSize: () => fit.pageSize,
    enabled: () => fit.ready,
    params: () => ({}),
    fetcher: (page, size, p) => FetchGroups(page, size),
  });

  let existingGroupUUIDs = $derived(
    new Set(Array.from(existing).map((v) => v.uuid)),
  );
</script>

<div class="w-full h-full overflow-auto gap-2">
  {#if q.loading}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:else if q.error}
    <div>{q.error}</div>
  {:else if !q.data || q.data.total === 0}
    <div>No groups</div>
  {:else}
    {#each q.data.groups as group (group.uuid)}
      <button
        type="button"
        class={`chip preset-outlined-surface-400-600 ${chosen.has(group.uuid) ? "preset-tonal-primary" : ""}`}
        onclick={() => {
          if (chosen.has(group.uuid)) {
            chosen.delete(group.uuid);
          } else {
            chosen.add(group.uuid);
          }
        }}
        disabled={existingGroupUUIDs.has(group.uuid)}
      >
        {#if chosen.has(group.uuid)}<CheckIcon size={14} />{/if}
        <span>{group.name}</span>
      </button>
    {/each}

    <div class="flex justify-between items-center gap-4 w-full self-end">
      <Pagination
        pageSize={fit.pageSize}
        bind:page={q.page}
        total={q.data?.total ?? 0}
      ></Pagination>
    </div>
  {/if}
</div>
