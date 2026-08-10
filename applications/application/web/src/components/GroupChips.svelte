<script lang="ts">
  import ArrowLeftIcon from "@lucide/svelte/icons/arrow-left";
  import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
  import { Pagination } from "@skeletonlabs/skeleton-svelte";
  import { IsJSONError } from "../lib/statuses/jsonerror";
  import { CheckIcon } from "@lucide/svelte";
  import { SvelteSet } from "svelte/reactivity";
  import {
    FetchGroups,
    type Group,
    type GroupsOrJSONError,
  } from "../lib/contracts/group";

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

  let time: number;
  const paginatedGroupPromises = $derived.by(() => {
    const p = page;
    const ps = pageSize;

    console.log("detected change");
    clearTimeout(time);

    return new Promise<GroupsOrJSONError>((resolve) => {
      time = setTimeout(() => {
        resolve(FetchGroups(p - 1, ps));
      }, 500);
    });
  });

  let existingGroupUUIDs = $derived(
    new Set(Array.from(existing).map((v) => v.uuid)),
  );
</script>

<div class="grid gap-2 w-full place-items-center h-full overflow-auto">
  {#await paginatedGroupPromises}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:then paginatedGroups}
    {#if IsJSONError(paginatedGroups)}
      <div>{paginatedGroups.error}</div>
    {:else}
      {#if paginatedGroups.total != 0}
        {#each paginatedGroups.groups as group (group.uuid)}
          <button
            type="button"
            class={`chip capitalize preset-outlined-surface-400-600 ${chosen.has(group.uuid) ? "preset-tonal-primary" : ""}`}
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
            count={paginatedGroups?.total ?? 0}
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
