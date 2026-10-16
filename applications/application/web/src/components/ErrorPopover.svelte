<script lang="ts">
  import { Popover, usePopover } from "@skeletonlabs/skeleton-svelte";

  const {
    message,
    children,
    open = $bindable(false),
  }: {
    message: string;
    children: any;
    open: boolean;
  } = $props();

  const uid = $props.id();
  const popover = usePopover({ id: uid });

  $effect(() => {
    popover().setOpen(open);
  });
</script>

<Popover.Provider value={popover}>
  <Popover.Anchor>
    {@render children()}
  </Popover.Anchor>

  <Popover.Positioner>
    <Popover.Content
      class="bg-error-50-950 p-2 rounded-[4px] text-surface-950-50"
    >
      <Popover.Title tabindex={-1}>{message}</Popover.Title>
    </Popover.Content>
  </Popover.Positioner>
</Popover.Provider>
