<script lang="ts">
	import { Dialog } from "@skeletonlabs/skeleton-svelte";
	import UXInput from "./UXInput.svelte";
	import { CopyToClipboard } from "../lib/apiutils/copy";
	import { FetchTestDelete } from "../lib/contracts/test";

	const {
		name,
		UUID,
		callback = () => {},
	}: {
		name: string;
		UUID: string;
		callback: () => void;
	} = $props();

	let insertedName: string = $state("");
	let statusMessage: string = $state("");
	let nameReady: boolean = $state(false);

	async function DeleteTest(e: Event) {
		e.preventDefault();

		if (!nameReady) {
			statusMessage = "Inserted name is not equal to requested";
			return;
		}

		const deleteResponse = await FetchTestDelete(UUID);
		if (deleteResponse != null) {
			statusMessage = deleteResponse.error;
			return;
		}

		callback();
	}
</script>

<form onsubmit={DeleteTest} class=" w-full flex flex-col h-full space-y-2">
	<p class="text-xl">
		Please enter

		<button
			type="button"
			class="italic hover:cursor-pointer"
			onclick={() => {
				CopyToClipboard(name);
			}}>{name}</button
		>

		below to confirm this action.
	</p>
	<p class="text-xs opacity-75 -mt-2">
		You may click the text in italic to copy it
	</p>

	<UXInput
		placeholder={name}
		bind:value={insertedName}
		bind:ready={nameReady}
		checker={(v: string) => {
			return name === v;
		}}
	/>

	<footer class="flex justify-end gap-2 w-full">
		{#if statusMessage != ""}
			<span class="justify-self-start">{statusMessage}</span>
		{/if}
		<Dialog.CloseTrigger class="btn preset-tonal">Cancel</Dialog.CloseTrigger>
		<button type="submit" class="btn preset-filled" disabled={nameReady}
			>Delete</button
		>
	</footer>
</form>
