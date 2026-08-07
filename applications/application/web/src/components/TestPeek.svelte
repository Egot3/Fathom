<script lang="ts">
	import { SvelteSet } from "svelte/reactivity";
	import type { Test } from "../lib/contracts/test";
	import ExistingQuizChips from "./ExistingQuizChips.svelte";
	import { IsJSONError, type JSONError } from "../lib/statuses/jsonerror";

	let { content }: { content: Test | JSONError } = $props();
</script>

{#if IsJSONError(content)}
	<div>{content.error}</div>
{:else}
	<label class="label">
		<span class="label-text">Name</span>
		<input
			class="input"
			type="text"
			value={content.name}
			placeholder="go-generics"
			disabled
		/>
	</label>

	{@const quizzes = new SvelteSet(content.quizzes)}
	<label class="label">
		<span class="label-text">Quizzes</span>
		<ExistingQuizChips disabled {quizzes} />
	</label>
{/if}
