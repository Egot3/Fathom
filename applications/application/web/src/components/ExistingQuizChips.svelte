<script lang="ts">
	import { SvelteSet } from "svelte/reactivity";
	import { CheckIcon } from "@lucide/svelte";
	import type { Quiz } from "../lib/contracts/quiz";
	import _ from "lodash";

	let {
		quizzes = $bindable(new SvelteSet<Quiz>()),
		disabled = false,
	}: {
		quizzes: SvelteSet<Quiz>;
		disabled?: boolean;
	} = $props();

	const startingQuizzes = new Set(quizzes);
</script>

<div class="grid gap-4 m-2 w-full place-items-center h-full overflow-auto">
	{#each startingQuizzes as quiz}
		<button
			{disabled}
			type="button"
			class={`chip capitalize preset-outlined-surface-400-600 ${quizzes.has(quiz) ? "preset-tonal-primary" : ""}`}
			onclick={() => {
				if (quizzes.has(quiz)) {
					quizzes.delete(quiz);
				} else {
					quizzes.add(quiz);
				}
			}}
		>
			{#if quizzes.has(quiz)}<CheckIcon size={14} />{/if}
			<span>{quiz.path.slice(14).slice(0, -3)}</span>
		</button>
	{/each}
</div>
