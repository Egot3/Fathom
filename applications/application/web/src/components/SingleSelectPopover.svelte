<script lang="ts" generics="T, R extends { total: number }">
  import { ArrowLeftIcon, ArrowRightIcon, XIcon } from "@lucide/svelte";
  import { Popover, usePopover } from "@skeletonlabs/skeleton-svelte";
  import type { ResultAsync } from "neverthrow";
  import type { JSONError } from "../lib/statuses/jsonerror";
  import Pagination from "./Pagination.svelte";

  let {
    label,
    fetcher,
    itemLabel,
    itemKey,
    selected = $bindable(null),
    items,
  }: {
    label: string;
    fetcher: (page: number, size: number) => ResultAsync<R, JSONError>;
    items: (response: R) => T[];
    itemLabel: (arg: T) => string;
    itemKey: (arg: T) => string;
    selected: string | null;
  } = $props();

  const uid = $props.id();
  let popover = usePopover({ id: uid });

  let loading = $state(true);
  let statusMessage = $state("");
  let paginated: T[] = $state(null as never);
  let total = $state(0);

  let selectedLabel: string | null = $state(null);

  let page = $state(1);
  let time: number;
  $effect(() => {
    const p = page;

    loading = true;
    clearTimeout(time);

    new Promise((resolve) => {
      time = setTimeout(async () => {
        statusMessage = await fetcher(p - 1, 5)
          .andTee((r) => {
            total = r.total;
            paginated = items(r);
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
</script>

<Popover.Provider value={popover}>
  <Popover.Anchor>
    <Popover.Trigger class="chip preset-outlined-surface-400-600">
      {selectedLabel ?? `Filter by ${label.toLowerCase()}`}
      {#if selected}
        <button
          onclick={(e) => {
            e.stopPropagation();
            selected = null;
            selectedLabel = null;
          }}
        >
          <XIcon size={7} />
        </button>
      {/if}
    </Popover.Trigger>
  </Popover.Anchor>
  <Popover.Positioner>
    <Popover.Content
      class="card bg-surface-100-900 p-2 w-64 max-h-72 z-40 overflow-auto space-y-1"
    >
      {#if loading}
        <div
          class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
        ></div>
      {:else}
        {#if statusMessage}
          <div class="h-full w-full rounded-xl">
            <p>{statusMessage}</p>
          </div>
        {:else}
          {#each paginated as item (itemKey(item))}
            <button
              class={`chip w-full justify-start text-primary-950-50 ${selected === itemKey(item) ? "preset-tonal-primary" : "preset-outlined-surface-400-600"}`}
              onclick={() => {
                selected = itemKey(item);
                selectedLabel = itemLabel(item);
                popover().setOpen(false);
              }}
            >
              {itemLabel(item)}
            </button>
          {/each}
        {/if}
      {/if}
      <Pagination bind:page pageSize={5} {total} />
    </Popover.Content>
  </Popover.Positioner>
</Popover.Provider>
