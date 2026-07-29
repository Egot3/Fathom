<script lang="ts">
	import DOMPurify from "isomorphic-dompurify";
	import { marked } from "marked";

	const { content = "" } = $props();
	let md = $derived(marked.parse(content));
</script>

<div class="prose">
	{#await md}
		<div
			class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
		></div>
	{:then mark}
		{@html DOMPurify.sanitize(mark)}
	{/await}
</div>
