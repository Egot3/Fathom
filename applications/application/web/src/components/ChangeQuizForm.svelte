<script lang="ts">
	import { Dialog } from "@skeletonlabs/skeleton-svelte";
	import {
		FetchQuiz,
		FetchQuizPatch,
		FetchQuizPut,
		Kind,
		type Meta,
		type QuizFile,
	} from "../lib/contracts/quiz";
	import { IsJSONError } from "../lib/statuses/jsonerror";
	import UXInput from "./UXInput.svelte";
	import _ from "lodash";

	const {
		callback,
		UUID,
		name,
	}: {
		callback: () => void;

		name: string;
		UUID: string;
	} = $props();

	const pathRegex = /^([a-zA-Z0-9_\-]+)*\/?[a-zA-Z0-9_\-]+\.md$/;

	let statusMessage = $state("");
	let nameMessage: string = $state("");

	let nameReady: boolean = $derived(pathRegex.test(name));

	let newName = $derived(name);
	let body = $state("");
	let allOrNone: boolean = $state(false);
	let randomized: boolean = $state(true);
	let score: number = $state(1);
	let kind: Kind | null = $state(null);

	let loading: boolean = $state(true);
	let error: string = $state("");

	let response = $state<QuizFile>(null as never);

	$effect(() => {
		loading = true;
		error = "";
		FetchQuiz(UUID)
			.then((r) => {
				console.log("response loaded");
				if (IsJSONError(r)) {
					error = r.error;
					return;
				}
				body = r.body;
				allOrNone = r.meta.all_or_none;
				randomized = r.meta.randomized;
				score = r.meta.score;
				kind = r.meta.kind;
				response = r;
			})
			.finally(() => (loading = false));
	});
</script>

{#if loading}
	<div class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"></div>
{:else}
	{#if error != ""}
		<p>{error}</p>
	{:else}
		<form
			onsubmit={async (e) => {
				e.preventDefault();

				const newMeta: Meta = {
					all_or_none: allOrNone,
					score: score,
					randomized: randomized,
					kind: response.meta.kind, // т.к. раньше надо было думать
				};

				if (newName != name || score != response.meta.score) {
					const patchResponse = await FetchQuizPatch(
						UUID,
						score != response.meta.score ? score : undefined,
						newName != name ? newName : undefined,
					);
					if (patchResponse != null) {
						statusMessage = patchResponse.error;
						return;
					}
				}

				if (body != response.body || !_.isEqual(response.meta, newMeta)) {
					const putResponse = await FetchQuizPut(UUID, newMeta, body);
					if (putResponse != null) {
						statusMessage = putResponse.error;
						return;
					}
				}
				callback();
			}}
			class="w-full flex flex-col h-full space-y-2"
		>
			<UXInput
				label="Name/path"
				placeholder="color_theory/sky.md"
				bind:value={newName}
				bind:ready={nameReady}
				message={nameMessage}
				checker={(v: string) => {
					nameMessage =
						"name must have no special symbols and contain .md extension";
					return pathRegex.test(v);
				}}
			/>

			<label class="label">
				<span class="label-text">Quiz</span>
				<textarea class="textarea" bind:value={body} placeholder="" required
				></textarea>
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

			{#if kind != Kind.Input}
				<label class="flex items-center space-x-2">
					<input class="checkbox" type="checkbox" bind:checked={randomized} />
					<p>Randomized</p>
				</label>
			{/if}

			<label class="flex items-center space-x-2">
				<input class="checkbox" type="checkbox" bind:checked={allOrNone} />
				<p>All or none</p>
			</label>

			<footer class="flex justify-end gap-2 w-full">
				{#if statusMessage != ""}
					<span class="justify-self-start">{statusMessage}</span>
				{/if}
				<Dialog.CloseTrigger class="btn preset-tonal"
					>Cancel</Dialog.CloseTrigger
				>
				<button type="submit" class="btn preset-filled">Change</button>
			</footer>
		</form>
	{/if}
{/if}
