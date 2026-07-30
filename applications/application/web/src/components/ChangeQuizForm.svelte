<script lang="ts">
	import { Dialog, Popover, usePopover } from "@skeletonlabs/skeleton-svelte";
	import {
		FetchQuizPost,
		FetchQuizPut,
		Kind,
		type Meta,
	} from "../lib/contracts/quiz";
	import { ClassForStatus, InputStatus } from "../lib/statuses/input";
	import UXInput from "./UXInput.svelte";
	import Quiz from "./Quiz.svelte";
	import type { JSONError } from "../lib/statuses/jsonerror";

	const {
		callback,
		name,
		title,
		question,
		answer,
		meta,

		UUID,
	}: {
		callback: () => void;
		name: string;
		title: string;
		question: string;
		answer: string;
		meta: Meta;

		UUID: string;
	} = $props();

	const pathRegex = /^([a-zA-Z0-9_\-]+)*\/?[a-zA-Z0-9_\-]+$/;

	let newName: string = $derived(name);
	let newTitle: string = $derived(title);
	let newQuestion: string = $derived(question);
	let newAnswer: string = $derived(answer);

	let type: Kind = $derived(meta.kind);
	let allOrNone: boolean = $derived(meta.all_or_none);
	let randomized: boolean = $derived(meta.randomized);
	let score: number = $derived(meta.score);

	let statusMessage: string = $state("");
	let nameMessage: string = $state("");

	let nameReady: boolean = $state(false);
	let content = $derived(`# ${title}
		
${question}
		
${answer}`);

	async function PostQuiz(e: Event) {
		e.preventDefault();

		if (newTitle === "") {
			statusMessage = "Title is required";
			return;
		}

		if (newAnswer === "") {
			statusMessage = "Answer is required";
		}

		let response: null | JSONError = null;
		if (score === meta.score) {
			response = await FetchQuizPut(
				UUID,
				{
					kind: type,
					score: score,
					randomized: randomized,
					all_or_none: allOrNone,
				} as Meta,
				content,
			);
		} else {
			if (
				newTitle === title &&
				newAnswer === answer &&
				newQuestion === question
			) {
			} // compile answers
		}
		const postResponse = await FetchQuizPost(
			{
				all_or_none: allOrNone,
				randomized: randomized,
				score: score,
				kind: type,
			} as Meta,
			name,
			content,
		);
		if (postResponse != null) {
			statusMessage = postResponse.error;
			return;
		}

		callback();
	}
</script>

<form onsubmit={PostQuiz} class="w-full flex flex-col h-full space-y-2">
	<div class="flex">
		<div class="flex flex-col w-1/2">
			<UXInput
				label="Name/path"
				placeholder="color_theory/sky.md"
				bind:value={newName}
				bind:ready={nameReady}
				message={nameMessage}
				checker={(v: string) => {
					nameMessage =
						"name must have no special symbols and doesn't contain any extensions";
					return pathRegex.test(v);
				}}
			/>

			<label class="label">
				<span class="label-text">Title</span>
				<input
					class="input"
					type="text"
					bind:value={newTitle}
					placeholder="Sky color"
					required
				/>
			</label>

			<label class="label">
				<span class="label-text">Question</span>
				<textarea
					class="textarea"
					bind:value={newQuestion}
					placeholder="What colour is the sky?"
					required></textarea>
			</label>

			<label class="label">
				<span class="label-text">Kind of quiz</span>
				<select class="select" bind:value={type}>
					<option selected value={Kind.Input}>Input</option>
					<option value={Kind.Radio}>Radio</option>
					<option value={Kind.Check}>Check</option>
					<option value={Kind.Accordance}>Accordance</option>
					<option value={Kind.Order}>Order</option>
				</select>
			</label>

			<label class="label">
				<span class="label-text">Option</span>
				<textarea
					class="textarea"
					bind:value={newAnswer}
					placeholder="[green]"
					required></textarea>
			</label>

			<label class="label">
				<span class="label-text">Default score</span>
				<input
					class="input"
					type="number"
					bind:value={score}
					placeholder="1"
					required
				/>
			</label>

			{#if type != Kind.Input}
				<label class="flex items-center space-x-2">
					<input class="checkbox" type="checkbox" bind:checked={randomized} />
					<p>Randomized</p>
				</label>
			{/if}

			<label class="flex items-center space-x-2">
				<input class="checkbox" type="checkbox" bind:checked={allOrNone} />
				<p>All or none</p>
			</label>
		</div>
		<div class="w-1/2">
			<Quiz {content} />
		</div>
	</div>

	<footer class="flex justify-end gap-2 w-full">
		{#if statusMessage != ""}
			<span class="justify-self-start">{statusMessage}</span>
		{/if}
		<Dialog.CloseTrigger class="btn preset-tonal">Cancel</Dialog.CloseTrigger>
		<button type="submit" class="btn preset-filled">Create</button>
	</footer>
</form>
