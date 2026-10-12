<script lang="ts">
  import { CheckIcon } from "@lucide/svelte";
  import { SvelteSet } from "svelte/reactivity";
  import {
    FetchGroups,
    type Group,
    type ListGroupResponse,
  } from "../lib/contracts/group";
  import Pagination from "./Pagination.svelte";

  let {
    chosen = $bindable(new SvelteSet<string>()),
    existing = new SvelteSet<Group>(),
  }: {
    chosen: SvelteSet<string>;
    existing?: SvelteSet<Group>;
  } = $props();
  $inspect(chosen);

  let page = $state(1);
  let pageSize = $state(5);

  let paginatedGroups: ListGroupResponse = $state(null as never);
  let statusMessage = $state("");
  let loading = $state(true);

  let time: number;
  $effect(() => {
    const p = page;
    const ps = pageSize;

    loading = true;
    clearTimeout(time);

    new Promise((resolve) => {
      time = setTimeout(async () => {
        statusMessage = await FetchGroups(p - 1, ps)
          .andTee((r) => {
            paginatedGroups = r;
            loading = false;
          })
          .match(
            (_) => "",
            (e) => e.error,
          );
        resolve(0);
      }, 500);
    });
  });

  let existingGroupUUIDs = $derived(
    new Set(Array.from(existing).map((v) => v.uuid)),
  );
</script>

<div class="w-full h-full overflow-auto gap-2">
  {#if loading}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:else}
    {#if statusMessage}
      <div>{statusMessage}</div>
    {:else}
      {#if paginatedGroups.total !== 0}
        {#each paginatedGroups.groups as group (group.uuid)}
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
          <Pagination {pageSize} bind:page total={paginatedGroups.total}
          ></Pagination>
        </div>
      {/if}
    {/if}
  {/if}
</div>
