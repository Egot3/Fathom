<script lang="ts">
	import ArrowLeftIcon from "@lucide/svelte/icons/arrow-left";
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import { Pagination } from "@skeletonlabs/skeleton-svelte";
	import { IsJSONError } from "../lib/statuses/jsonerror";
	import {
		FetchAllQuizzes,
		type Quiz,
		type QuizzesOrError,
	} from "../lib/contracts/quiz";
	import { CheckIcon } from "@lucide/svelte";
	import { SvelteSet } from "svelte/reactivity";

	let {
		chosen = $bindable(new SvelteSet<string>()),
		existing = new SvelteSet<Quiz>(),
	}: {
		chosen: SvelteSet<string>;
		existing?: SvelteSet<Quiz>;
	} = $props();
	$inspect(chosen);

	let height = $state(0);

	let page = $state(1);
	let pageSize = $state(5);

	let time: number;
	const paginatedQuizPromises = $derived.by(() => {
		const p = page;
		const ps = pageSize;

		console.log("detected change");
		clearTimeout(time);

		return new Promise<QuizzesOrError>((resolve) => {
			time = setTimeout(() => {
				resolve(FetchAllQuizzes(p - 1, ps));
			}, 500);
		});
	});

	let existingQuizUUIDs = $derived(
		new Set(Array.from(existing).map((v) => v.uuid)),
	);
</script>

<div
	class="grid gap-2 w-full place-items-center h-full overflow-auto"
	bind:clientHeight={height}
>
	{#await paginatedQuizPromises}
		<div
			class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
		></div>
	{:then paginatedQuizzes}
		{#if IsJSONError(paginatedQuizzes)}
			<div>{paginatedQuizzes.error}</div>
		{:else}
			{#if paginatedQuizzes.total != 0}
				{#each paginatedQuizzes.quizzes as quiz (quiz.uuid)}
					<button
						type="button"
						class={`chip capitalize preset-outlined-surface-400-600 ${chosen.has(quiz.uuid) ? "preset-tonal-primary" : ""}`}
						onclick={() => {
							if (chosen.has(quiz.uuid)) {
								chosen.delete(quiz.uuid);
							} else {
								chosen.add(quiz.uuid);
							}
						}}
						disabled={existingQuizUUIDs.has(quiz.uuid)}
					>
						{#if chosen.has(quiz.uuid)}<CheckIcon size={14} />{/if}
						<span>{quiz.path.slice(14).slice(0, -3)}</span>
					</button>
				{/each}

				<div class="flex justify-between items-center gap-4 w-full self-end">
					<Pagination
						count={paginatedQuizzes?.total ?? 0}
						{pageSize}
						{page}
						onPageChange={(event) => (page = event.page)}
						class="rounded-xl"
					>
						<Pagination.PrevTrigger>
							<ArrowLeftIcon class="size-4" />
						</Pagination.PrevTrigger>
						<Pagination.Context>
							{#snippet children(pagination)}
								{#each pagination().pages as page, index (page)}
									{#if page.type === "page"}
										<Pagination.Item {...page}>
											{page.value}
										</Pagination.Item>
									{:else}
										<Pagination.Ellipsis {index}>…</Pagination.Ellipsis>
									{/if}
								{/each}
							{/snippet}
						</Pagination.Context>
						<Pagination.NextTrigger>
							<ArrowRightIcon class="size-4" />
						</Pagination.NextTrigger>
					</Pagination>
				</div>
			{/if}
		{/if}
	{/await}
</div>
