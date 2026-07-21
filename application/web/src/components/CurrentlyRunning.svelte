<script lang="ts">
	import { GetCurrentlyRunning } from "../lib/bgdata/currentlyrunning.svelte";
	import { IsJSONError, type JSONError } from "../lib/statuses/jsonerror";

	let currentlyRunningPromise = $state(GetCurrentlyRunning());

	function updateCurrentlyRunning() {
		currentlyRunningPromise = GetCurrentlyRunning();
	}
</script>

<div class="card preset-filled-surface-700-300 max-w-md overflow-hidden">
	<article class="space-y-4 p-4">
		<div>
			<h2 class="self-center text-2xl">Currently running</h2>
		</div>
		{#await currentlyRunningPromise}
			<div
				class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
			></div>
		{:then currentlyRunning}
			{#if currentlyRunning === null}
				<div>NOTHING</div>
			{:else}
				{#if IsJSONError(currentlyRunning)}
					<div>{currentlyRunning.error}</div>
				{:else}
					<p>Test {currentlyRunning.Name}</p>
				{/if}
			{/if}
		{/await}
	</article>
</div>
