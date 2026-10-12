<script lang="ts">
  import { CheckIcon } from "@lucide/svelte";
  import { SvelteSet } from "svelte/reactivity";
  import { FetchUsers, type User } from "../lib/contracts/user";
  import type { Totals } from "../lib/contracts/totals";
  import Pagination from "./Pagination.svelte";

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
          <Pagination bind:page {pageSize} total={paginatedUsers.total} />
        </div>
      {/if}
    {/if}
  {/if}
</div>
