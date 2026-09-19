<script lang="ts">
  import { CheckIcon } from "@lucide/svelte";

  let {
    options = [] as string[],
    selected = $bindable(0),
    working = false,
  } = $props();

  function setColor(selectedColorId: number) {
    selected = selectedColorId;
  }

  $inspect(working);
</script>

{#each options as c, i (i)}
  <button
    class={`chip capitalize preset-outlined-surface-400-600 ${selected === i && !working ? "preset-tonal-primary" : ""} ${working ? "preset-tonal-warning animate-pulse" : ""}`}
    onclick={() => setColor(i)}
    disabled={working}
  >
    {#if selected === i}<CheckIcon size={14} />{/if}
    <span>{c}</span>
  </button>
{/each}
