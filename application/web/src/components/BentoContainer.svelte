<script lang="ts">
	const { children } = $props();
</script>

<div
	class="flex w-full h-full items-center justify-center text-surface-100-900"
>
	<div
		class="h-5/6 w-9/10 grid grid-cols-12 grid-rows-12 bg-surface-300-700 rounded-4xl gap-4 p-4"
	>
		{@render children()}
	</div>
</div>
