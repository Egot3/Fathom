<script lang="ts">
	import { X } from "@lucide/svelte";
	import { Dialog, Portal } from "@skeletonlabs/skeleton-svelte";

	let name: string = $state("");
	let description: string = $state("");
	let quizzess: string[] = $state([]);
</script>

<Dialog>
	<Dialog.Trigger class="btn preset-filled-primary-500 whitespace-normal"
		>Create your first test</Dialog.Trigger
	>

	<Portal>
		<Dialog.Backdrop class="fixed inset-0 z-50 bg-surface-50-950/50" />

		<Dialog.Positioner
			class="fixed inset-0 z-50 flex items-center justify-center"
		>
			<Dialog.Content
				class="card bg-surface-100-900 w-md p-4 space-y-4 shadow-xl"
			>
				<header class="flex justify-between items-center">
					<Dialog.Title class="text-2xl font-bold">Test creation</Dialog.Title>
					<Dialog.CloseTrigger class="btn-icon hover:preset-tonal"
						><X /></Dialog.CloseTrigger
					>
				</header>

				<input class="input border-2" />

				<footer class="flex justify-end gap-2">
					<Dialog.CloseTrigger class="btn preset-tonal"
						>Cancel</Dialog.CloseTrigger
					>
					<button type="button" class="btn preset-filled">Agree</button>
				</footer>
			</Dialog.Content>
		</Dialog.Positioner>
	</Portal>
</Dialog>
