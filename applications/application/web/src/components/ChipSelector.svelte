<script lang="ts">
  import { CheckIcon } from "@lucide/svelte";

  let {
    options = [] as string[],
    selected = $bindable(0),
    working = false,
    onchange = () => {},
  } = $props();

  function setColor(selectedColorId: number) {
    selected = selectedColorId;
  }
</script>

{#each options as c, i (i)}
  <button
    class={`chip capitalize preset-outlined-surface-400-600 ${selected === i && !working ? "preset-tonal-primary" : ""} ${working ? "preset-tonal-warning animate-pulse" : ""}`}
    onclick={() => {
      setColor(i);
      onchange?.();
    }}
    disabled={working}
  >
    {#if selected === i}<CheckIcon size={14} />{/if}
    <span>{c}</span>
  </button>
{/each}
