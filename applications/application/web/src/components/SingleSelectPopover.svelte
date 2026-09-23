<script lang="ts">
  import { XIcon } from "@lucide/svelte";
  import { Popover, usePopover } from "@skeletonlabs/skeleton-svelte";
  import { FetchUsers } from "../lib/contracts/user";

  let {
    label,
    fetcher,
    itemLabel,
    itemKey,
    selected = $bindable(null),
  } = $props();

  const uid = $props.id();
  let popover = usePopover({ id: uid });
  let page = $state(1);
  let time: number;
</script>

<Popover.Provider value={popover}>
  <Popover.Anchor>
    <Popover.Trigger class="chip preset-outlined-surface-400-600">
      {`Filter by ${label.toLowerCase()}`}
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
      {#each items as item (itemKey(item))}
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
    </Popover.Content>
  </Popover.Positioner>
</Popover.Provider>
