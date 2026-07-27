<script lang="ts">
	import { CheckIcon } from "@lucide/svelte";

	const { options = [] as string[], selected = 0 } = $props();

	let color = $derived(options[selected]);

	function setColor(selectedColor: string) {
		color = selectedColor;
	}
</script>

{#each options as c (c)}
	<button
		class={`chip capitalize preset-outlined-surface-400-600 ${color === c ? "preset-tonal-primary" : ""}`}
		onclick={() => setColor(c)}
	>
		{#if color === c}<CheckIcon size={14} />{/if}
		<span>{c}</span>
	</button>
{/each}
