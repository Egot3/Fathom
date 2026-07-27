<script lang="ts">
	import { RotateCw } from "@lucide/svelte";
	import { GetCurrentlyRunning } from "../lib/bgdata/currentlyrunning.svelte";
	import { IsJSONError, type JSONError } from "../lib/statuses/jsonerror";

	let currentlyRunningPromise = $state(GetCurrentlyRunning());

	function updateCurrentlyRunning() {
		currentlyRunningPromise = GetCurrentlyRunning();
	}

	let reload: HTMLButtonElement | null = $state(null);

	let time: number;
	let accumClicks: number = 0;
	/*
		- What does this do?
		- This Shenanigans(TS for short)
	*/
	let isAnimating = $state(false);
	let duration = $state("0s");
	let iterations = $state("0");

	function handleClick() {
		if (!reload) return;

		accumClicks++;
		clearTimeout(time);

		time = setTimeout(() => {
			const capturedAccum = accumClicks;
			accumClicks = 0;

			duration = `${capturedAccum * 0.6}s`;
			iterations = `${capturedAccum}`;
			isAnimating = true;
			console.log(
				`accumulated ${duration} duration and ${iterations} iterations`,
			);

			updateCurrentlyRunning();

			setTimeout(
				() => {
					isAnimating = false;
				},
				capturedAccum * 600 + 100,
			);
		}, 300);
	}
</script>

<div class="card preset-filled-surface-700-300 max-w-md overflow-hidden">
	<article class="space-y-4 p-4">
		<div class="flex gap-2 items-center">
			<h2 class="self-center text-2xl font-bold">Currently running</h2>
			<button
				bind:this={reload}
				class="ml-auto mt-1"
				class:animate-custom-spin={isAnimating}
				style="--spin-duration: {duration}; --spin-iterations: {iterations};"
				onclick={handleClick}
			>
				<RotateCw />
			</button>
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
					<a class="mb-auto btn preset-filled-primary-500" href="/participate"
						>Join?</a
					>
				{/if}
			{/if}
		{/await}
	</article>
</div>

<style>
	@keyframes custom-spin {
		from {
			transform: rotate(0deg);
		}
		to {
			transform: rotate(360deg);
		}
	}

	.animate-custom-spin {
		animation-name: custom-spin;
		animation-duration: calc(0.8s / var(--spin-iterations));
		animation-iteration-count: var(--spin-iterations, 1);

		animation-timing-function: linear(0, 0.1, 0.2, 0.5, 0.9, 1.03, 0.98, 1);
	}
</style>
