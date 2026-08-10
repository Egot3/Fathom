<script lang="ts">
  import ArrowLeftIcon from "@lucide/svelte/icons/arrow-left";
  import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
  import { Pagination } from "@skeletonlabs/skeleton-svelte";
  import { IsJSONError } from "../lib/statuses/jsonerror";
  import { CheckIcon } from "@lucide/svelte";
  import { SvelteSet } from "svelte/reactivity";
  import {
    FetchUsers,
    type User,
    type UsersOrJSONError,
  } from "../lib/contracts/user";

  let {
    chosen = $bindable(new SvelteSet<string>()),
    existing = new SvelteSet<User>(),
  }: {
    chosen: SvelteSet<string>;
    existing?: SvelteSet<User>;
  } = $props();
  $inspect(chosen);

  let page = $state(1);
  let pageSize = $state(5);

  let time: number;
  const paginatedUserPromises = $derived.by(() => {
    const p = page;
    const ps = pageSize;

    console.log("detected change");
    clearTimeout(time);

    return new Promise<UsersOrJSONError>((resolve) => {
      time = setTimeout(() => {
        resolve(FetchUsers(p - 1, ps));
      }, 500);
    });
  });

  let existingUserUUIDs = $derived(
    new Set(Array.from(existing).map((v) => v.uuid)),
  );
</script>

<div class="grid gap-2 w-full place-items-center h-full overflow-auto">
  {#await paginatedUserPromises}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:then paginatedUsers}
    {#if IsJSONError(paginatedUsers)}
      <div>{paginatedUsers.error}</div>
    {:else}
      {#if paginatedUsers.total != 0}
        {#each paginatedUsers.users as user (user.uuid)}
          <button
            type="button"
            class={`chip capitalize preset-outlined-surface-400-600 ${chosen.has(user.uuid) ? "preset-tonal-primary" : ""}`}
            onclick={() => {
              if (chosen.has(user.uuid)) {
                chosen.delete(user.uuid);
              } else {
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
  {/await}
</div>
