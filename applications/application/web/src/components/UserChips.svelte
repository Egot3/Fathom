<script lang="ts">
  import ArrowLeftIcon from "@lucide/svelte/icons/arrow-left";
  import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
  import { Pagination } from "@skeletonlabs/skeleton-svelte";
  import { IsJSONError } from "../lib/statuses/jsonerror";
  import { CheckIcon } from "@lucide/svelte";
  import { SvelteSet } from "svelte/reactivity";
  import { FetchUsers, type User } from "../lib/contracts/user";
  import type { Totals } from "../lib/contracts/totals";

  let {
    chosen = $bindable(new SvelteSet<string>()),
    existing = new SvelteSet<User>(),
    maxSelect,
  }: {
    chosen: SvelteSet<string>;
    existing?: SvelteSet<User>;
    maxSelect?: number;
  } = $props();
  $inspect(chosen);

  let page = $state(1);
  let pageSize = $state(5);

  let loading = $state(true);
  let statusMessage = $state("");

  let paginatedUsers: { users: User[]; total: number } = $state(null as never);

  let time: number;
  $effect(() => {
    const p = page;
    const ps = pageSize;

    clearTimeout(time);

    new Promise((resolve) => {
      time = setTimeout(async () => {
        statusMessage = await FetchUsers(p - 1, ps)
          .andTee((r) => {
            paginatedUsers = r;
            loading = false;
          })
          .match(
            (_) => "",
            (err) => err.error,
          );

        resolve(0);
      }, 500);
    });
  });

  let existingUserUUIDs = $derived(
    new Set(Array.from(existing).map((v) => v.uuid)),
  );
</script>

<div class="grid gap-2 w-full place-items-center h-full overflow-auto">
  {#if loading}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:else}
    {#if statusMessage !== ""}
      <div>{statusMessage}</div>
    {:else}
      {#if paginatedUsers.total !== 0}
        {#each paginatedUsers.users as user (user.uuid)}
          <button
            type="button"
            class={`chip capitalize preset-outlined-surface-400-600 ${chosen.has(user.uuid) ? "preset-tonal-primary" : ""}`}
            onclick={() => {
              if (chosen.has(user.uuid)) {
                chosen.delete(user.uuid);
              } else {
                if (maxSelect && chosen.size >= maxSelect) {
                  return;
                }
                chosen.add(user.uuid);
              }
            }}
            disabled={existingUserUUIDs.has(user.uuid)}
          >
            {#if chosen.has(user.uuid)}<CheckIcon size={14} />{/if}
            <span>{user.nickname}</span>
          </button>
        {/each}

        <div class="flex justify-between items-center gap-4 w-full self-end">
          <Pagination
            count={paginatedUsers?.total ?? 0}
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
</div>
