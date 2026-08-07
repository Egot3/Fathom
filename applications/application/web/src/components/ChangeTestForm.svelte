<script lang="ts">
	import { Dialog } from "@skeletonlabs/skeleton-svelte";
	import UXInput from "./UXInput.svelte";
	import _ from "lodash";
	import {
		FetchTestBundle,
		FetchTestPatch,
		FetchTestPrune,
		type Test,
	} from "../lib/contracts/test";
	import QuizChips from "./QuizChips.svelte";
	import { SvelteSet } from "svelte/reactivity";
	import ExistingQuizChips from "./ExistingQuizChips.svelte";

	const {
		callback,
		UUID,
		name,
		response,
	}: {
		callback: () => void;

		response: Test;
		name: string;
		UUID: string;
	} = $props();

	let statusMessage = $state("");
	let nameMessage: string = $state("");

	let nameReady: boolean = $state(false);

	let newName = $derived(name);
	let quizzes = $derived(new SvelteSet(response.quizzes));
	let newQuizzes = $state(new SvelteSet<string>());

	async function ChangeQuiz(e: Event) {
		e.preventDefault();

		if (newName != name && newName) {
			const patchResponse = await FetchTestPatch(UUID, newName);
			if (patchResponse != null) {
				statusMessage = patchResponse.error;
				return;
			}
		}

		const quizzesArr = Array.from(quizzes);
		const quizzesRemoved = _.difference(response.quizzes, quizzesArr);
		console.log("quizzes removed+quizzesArr: ", quizzesRemoved, quizzesArr);
		if (quizzesRemoved.length > 0) {
			const deleteResponse = await FetchTestPrune(
				UUID,
				quizzesRemoved.map((v) => v.uuid),
			);
			if (deleteResponse != null) {
				statusMessage = deleteResponse.error;
				return;
			}
		}

		const newQuizzesArr = Array.from(newQuizzes);
		if (newQuizzesArr.length > 0) {
			const postResponse = await FetchTestBundle(UUID, newQuizzesArr);
			if (postResponse != null) {
				statusMessage = postResponse.error;
				return;
			}
		}

		callback();
	}

	const existing = new SvelteSet(quizzes);
</script>

<form onsubmit={ChangeQuiz} class="w-full flex flex-col h-full space-y-2">
	<UXInput
		label="Name"
		placeholder="go-regenerics"
		bind:value={newName}
		bind:ready={nameReady}
		message={nameMessage}
		checker={(v: string) => {
			if (v.length < 3) {
				nameMessage = "name can't be less than 3 characters in length";
				return false;
			}
			if (v.length > 255) {
				nameMessage = "name can't be more than 255 characters in length";
				return false;
			}

			return true;
		}}
	/>

	<div class="w-full flex">
		<div class="h-full w-1/2">
			<label class="label">
				<span class="label-text">Old quizzes</span>
				<ExistingQuizChips bind:quizzes />
			</label>
		</div>
		<div class="h-full w-1/2">
			<label class="label">
				<span class="label-text">New quizzes</span>
				<QuizChips bind:chosen={newQuizzes} {existing} />
			</label>
		</div>
	</div>

	<footer class="flex justify-end gap-2 w-full">
		{#if statusMessage != ""}
			<span class="justify-self-start">{statusMessage}</span>
		{/if}
		<Dialog.CloseTrigger class="btn preset-tonal">Cancel</Dialog.CloseTrigger>
		<button type="submit" class="btn preset-filled">Change</button>
	</footer>
</form>
