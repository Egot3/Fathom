<script lang="ts">
  import { CheckIcon } from "@lucide/svelte";

  const {
    options = [] as string[],
    selected = $bindable(0),
    working = false,
  } = $props();

  let color = $derived(options[selected]);

  function setColor(selectedColor: string) {
    color = selectedColor;
  }

  $inspect(working);
</script>

{#each options as c (c)}
  <button
    class={`chip capitalize preset-outlined-surface-400-600 ${color === c && !working ? "preset-tonal-primary" : ""} ${working ? "preset-tonal-warning animate-pulse" : ""}`}
    onclick={() => setColor(c)}
    disabled={working}
  >
    {#if color === c}<CheckIcon size={14} />{/if}
    <span>{c}</span>
  </button>
{/each}
