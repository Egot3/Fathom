<script lang="ts">
	import { Eye, EyeClosed, X } from "@lucide/svelte";
	import { Dialog, Portal } from "@skeletonlabs/skeleton-svelte";

	const {
		title,
		children,
		callback = () => {},
		contentGetter = () => {},
	}: {
		title: string;
		children: any;
		callback: () => void;
		contentGetter: () => unknown;
	} = $props();

	let moused = $state(false);
	let clicked = $state(false);
	let content: unknown = $state();
</script>

<Dialog>
	<!-- why not just make bind:open gng -->
	<Dialog.Trigger
		onmouseenter={() => (moused = true)}
		onmouseleave={() => (moused = false)}
		onclick={() => {
			callback();
			content = contentGetter();
		}}
		class="btn mb-0 mt-0 preset-filled-surface-300-700 hover:preset-filled-primary-300-700 aspect-square h-auto w-auto p-1 m-px"
	>
		{#if moused}
			<Eye size="16" />
		{:else}
			<EyeClosed size="16" />
		{/if}
	</Dialog.Trigger>

	<Portal>
		<Dialog.Backdrop class="fixed inset-0 z-50 bg-surface-50-950/50" />

		<Dialog.Positioner
			class="fixed inset-0 z-50 flex items-center justify-center"
		>
			<Dialog.Content
				class="card bg-surface-100-900 w-md p-4 space-y-4 shadow-xl"
			>
				<header class="flex justify-between items-center">
					<Dialog.Title class="text-2xl font-bold">{title}</Dialog.Title>
					<Dialog.CloseTrigger class="btn-icon hover:preset-tonal"
						><X /></Dialog.CloseTrigger
					>
				</header>

				{#await content}
					<div>loading</div>
				{:then content}
					{@render children({ content })}
				{/await}
			</Dialog.Content>
		</Dialog.Positioner>
	</Portal>
</Dialog>
