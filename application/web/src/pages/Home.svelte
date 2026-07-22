<script lang="ts">
	import BentoContainer from "../components/BentoContainer.svelte";
	import CurrentlyRunning from "../components/CurrentlyRunning.svelte";
	import PreviousTests from "../components/PreviousTests.svelte";
	import SelfInfo from "../components/SelfInfo.svelte";
</script>

<BentoContainer>
	<div
		id="current"
		class="row-start-1 row-end-8 col-start-1 col-end-5 bg-surface-700-300 rounded-2xl p-4"
	>
		<CurrentlyRunning />
	</div>
	<div
		id="previous_results"
		class="row-start-8 row-end-13 col-start-1 col-end-13 bg-surface-700-300 rounded-2xl p-8 flex flex-col"
	>
		<h2 class="self-center text-2xl font-bold">Previous results</h2>
		<PreviousTests />
	</div>
	<div
		id="self"
		class="row-start-1 row-end-8 col-start-5 col-end-13 flex flex-col bg-surface-700-300 rounded-2xl p-8"
	>
		<h2 class="text-2xl font-bold">You are</h2>
		<SelfInfo />
	</div>
</BentoContainer>
