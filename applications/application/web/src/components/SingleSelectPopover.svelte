<script lang="ts" generics="T, R extends { total: number }">
  import { XIcon } from "@lucide/svelte";
  import { Popover, usePopover } from "@skeletonlabs/skeleton-svelte";
  import type { ResultAsync } from "neverthrow";
  import type { JSONError } from "../lib/statuses/jsonerror";

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
</script>

<Popover.Provider value={popover}>
  <Popover.Anchor>
    <Popover.Trigger class="chip preset-outlined-surface-400-600">
      {selected ?? `Filter by ${label.toLowerCase()}`}
      {#if selected}
        <button
          onclick={(e) => {
            e.stopPropagation();
            selected = null;
          }}
        >
          <XIcon size={12} />
        </button>
      {/if}
    </Popover.Trigger>
  </Popover.Anchor>
  <Popover.Positioner>
    <Popover.Content
      class="card bg-surface-100-900 p-2 w-64 max-h-72 overflow-auto space-y-1"
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
              class={`chip w-full justify-start ${selected === itemKey(item) ? "preset-tonal-primary" : "preset-outlined-surface-400-600"}`}
              onclick={() => {
                selected = itemKey(item);
                popover().setOpen(false);
              }}
            >
              {itemLabel(item)}
            </button>
          {/each}
        {/if}
      {/if}
    </Popover.Content>
  </Popover.Positioner>
</Popover.Provider>
