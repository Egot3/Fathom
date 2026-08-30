<script lang="ts">
  import { Popover, usePopover } from "@skeletonlabs/skeleton-svelte";
  import { ClassForStatus, InputStatus } from "../lib/statuses/input";
  import { onMount } from "svelte";

  let {
    value = $bindable(""),
    label = "",
    type = "text",
    ready = $bindable(true),
    checker = (v: string) => true,
    placeholder = "",
    message = $bindable(""),
    initValue = "",
    disabled = false,
    state = $bindable(InputStatus.Idle),
  } = $props();
  const uid = $props.id();

  const popover = $derived(usePopover({ id: uid }));

  onMount(() => {
    value = value || initValue;
  });

  $effect(() => {
    ready = checker(value);
  });
</script>

<label class="label">
  <span class="label-text">{label}</span>
  <Popover.Provider value={popover}>
    <Popover.Anchor>
      <input
        onmouseover={() => {
          if (message !== "") {
            popover().setOpen(true);
          }
        }}
        onmouseout={() => {
          if (message !== "") {
            popover().setOpen(false);
          }
        }}
        {disabled}
        class={"input border-2 " + ClassForStatus(state)}
        {type}
        onblur={() => {
          if (ready) {
            state = InputStatus.Treat;

            return;
          }
          popover().setOpen(true);
          state = InputStatus.Punish;
        }}
        onfocus={() => {
          popover().setOpen(false);
          state = InputStatus.Idle;
        }}
        onkeydown={() => {
          if (ready) {
            state = InputStatus.Treat;
          }
        }}
        bind:value
        {placeholder}
      />
    </Popover.Anchor>

    <Popover.Positioner>
      <Popover.Content
        class="bg-error-50-950 p-2 rounded-[4px] text-surface-950-50"
      >
        <Popover.Title tabindex={-1}>{message}</Popover.Title>
      </Popover.Content>
    </Popover.Positioner>
  </Popover.Provider>
</label>
