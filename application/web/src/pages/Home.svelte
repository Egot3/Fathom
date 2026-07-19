<script lang="ts">
	import PreviousTests from "../components/PreviousTests.svelte";
</script>

<div
	class="flex w-full h-full items-center justify-center text-surface-100-900"
>
	<div
		class="h-5/6 w-9/10 grid grid-cols-12 grid-rows-12 bg-surface-300-700 rounded-4xl gap-4 p-4"
	>
		<div
			id="current"
			class="row-start-1 row-end-8 col-start-1 col-end-5 bg-surface-700-300 rounded-2xl grid justify-center pb-4 pt-2"
		>
			<PreviousTests />
		</div>
		<div
			id="previous_results"
			class="row-start-8 row-end-13 col-start-1 col-end-13 bg-surface-700-300 rounded-2xl"
		></div>
		<div
			id="self"
			class="row-start-1 row-end-8 col-start-5 col-end-13 bg-surface-700-300 rounded-2xl"
		></div>
	</div>
</div>
