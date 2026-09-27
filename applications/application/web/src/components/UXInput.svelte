<script lang="ts">
  import { ClassForStatus, InputStatus } from "../lib/statuses/input";
  import { onMount } from "svelte";
  import ErrorPopover from "./ErrorPopover.svelte";

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
    status = $bindable(InputStatus.Idle),
  } = $props();

  let popoverOpen = $state(false);

  onMount(() => {
    value = value || initValue;
  });

  $effect(() => {
    ready = checker(value);
  });
</script>

<label class="label">
  <span class="label-text">{label}</span>
  <ErrorPopover {message} open={popoverOpen}>
    <input
      onmouseover={() => {
        if (message !== "") {
          popoverOpen = true;
        }
      }}
      onmouseout={() => {
        if (message !== "") {
          popoverOpen = false;
        }
      }}
      {disabled}
      class={"input border-2 " + ClassForStatus(status)}
      {type}
      onblur={() => {
        if (ready) {
          status = InputStatus.Treat;

          return;
        }
        popoverOpen = true;
        status = InputStatus.Punish;
      }}
      onfocus={() => {
        popoverOpen = false;
        status = InputStatus.Idle;
      }}
      onkeydown={() => {
        if (ready) {
          status = InputStatus.Treat;
        }
      }}
      bind:value
      {placeholder}
    />
  </ErrorPopover>
</label>
