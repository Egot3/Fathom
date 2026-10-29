<script lang="ts">
  import { CheckIcon } from "@lucide/svelte";

  let {
    options = [] as string[],
    selected = $bindable(0),
    working = false,
    onchange = () => {},
  } = $props();
</script>

<span class="flex gap-1">
  {#each options as c, i (i)}
    <button
      class="chip capitalize preset-outlined-surface-400-600 {selected === i &&
      !working
        ? 'preset-tonal-primary'
        : ''} {working ? 'preset-tonal-warning animate-pulse' : ''}"
      onclick={() => {
        selected = i;
        onchange?.();
      }}
      disabled={working}
    >
      {#if selected === i}<CheckIcon size={14} />{/if}
      <span>{c}</span>
    </button>
  {/each}
</span>
