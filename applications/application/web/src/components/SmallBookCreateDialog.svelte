<script lang="ts">
	import { Book, Pencil, X } from "@lucide/svelte";
	import { Dialog, Portal } from "@skeletonlabs/skeleton-svelte";

	const { title, children } = $props();
</script>

<Dialog>
	<Dialog.Trigger class="btn preset-filled-primary-500 whitespace-normal"
		><Book>
			<Pencil size="14" x="5" y="2"></Pencil>
		</Book></Dialog.Trigger
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
					<Dialog.Title class="text-2xl font-bold">{title}</Dialog.Title>
					<Dialog.CloseTrigger class="btn-icon hover:preset-tonal"
						><X /></Dialog.CloseTrigger
					>
				</header>

				{@render children()}
			</Dialog.Content>
		</Dialog.Positioner>
	</Portal>
</Dialog>
