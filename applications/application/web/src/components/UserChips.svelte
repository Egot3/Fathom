<script lang="ts">
  import { CheckIcon } from "@lucide/svelte";
  import { SvelteSet } from "svelte/reactivity";
  import { FetchUsers, type User } from "../lib/contracts/user";
  import Pagination from "./Pagination.svelte";
  import { CreatePagedQuery } from "../lib/layoututils/createPagedQuery.svelte";
  import { CreateFitRows } from "../lib/layoututils/fitRows.svelte";

  let {
    chosen = $bindable(new SvelteSet<string>()),
    existing = new SvelteSet<User>(),
    maxSelect,
  }: {
    chosen: SvelteSet<string>;
    existing?: SvelteSet<User>;
    maxSelect?: number;
  } = $props();

  const fit = CreateFitRows({ rowHeight: 20, reserve: 1 });

  const q = CreatePagedQuery({
    pageSize: () => fit.pageSize,
    enabled: () => fit.ready,
    params: () => ({}),
    fetcher: (page, size, p) => FetchUsers(page, size),
  });

  let existingUserUUIDs = $derived(
    new Set(Array.from(existing).map((v) => v.uuid)),
  );
</script>

<div class="grid gap-2 w-full place-items-center h-full overflow-auto">
  {#if q.loading}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:else if q.error}
    <div>{q.error}</div>
  {:else if !q.data || q.data.total === 0}
    <div>No users</div>
  {:else}
    {#each q.data.users as user (user.uuid)}
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
        bind:page={q.page}
        pageSize={fit.pageSize}
        total={q.data?.total ?? 0}
      />
    </div>
  {/if}
</div>
