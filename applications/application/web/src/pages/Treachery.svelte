<script lang="ts">
	import BentoContainer from "../components/BentoContainer.svelte";
	import Quizzes from "../components/Quizzes.svelte";
	import TeacherCurrentlyRunning from "../components/TeacherCurrentlyRunning.svelte";
	import Tests from "../components/Tests.svelte";
</script>

<BentoContainer>
	<div
		id="current"
		class="row-start-1 row-end-8 col-start-1 col-end-5 flex flex-col bg-surface-700-300 rounded-2xl p-8"
	>
		<h2 class="text-2xl font-bold">Now running</h2>
		<hr class="hr mt-5 mb-5" />
		<TeacherCurrentlyRunning />
	</div>
	<div
		id="testing_results"
		class="row-start-8 row-end-13 col-start-1 col-end-9 bg-surface-700-300 rounded-2xl p-8 flex flex-col"
	>
		<h2 class="self-center text-2xl font-bold">Testing results</h2>
	</div>
	<div
		id="pupils"
		class="row-start-1 row-end-8 col-start-5 col-end-9 flex flex-col bg-surface-700-300 rounded-2xl p-8"
	>
		<h2 class="text-2xl font-bold">Registered pupils</h2>
	</div>
	<div
		id="tests"
		class="row-start-1 row-end-7 col-start-9 col-end-13 flex flex-col bg-surface-700-300 rounded-2xl p-8"
	>
		<h2 class="text-2xl font-bold">Tests</h2>
		<hr class="hr mt-5 mb-5" />
		<Tests />
	</div>
	<div
		id="quizzes"
		class="row-start-7 row-end-13 col-start-9 col-end-13 flex flex-col bg-surface-700-300 rounded-2xl p-8"
	>
		<h2 class="text-2xl font-bold">Quizzes</h2>
		<hr class="hr mt-5 mb-5" />
		<Quizzes />
	</div>
</BentoContainer>
