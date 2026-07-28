<script lang="ts">
	import ArrowLeftIcon from "@lucide/svelte/icons/arrow-left";
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import { Pagination } from "@skeletonlabs/skeleton-svelte";
	import { IsJSONError } from "../lib/statuses/jsonerror";
	import { FetchAllQuizzes, type QuizzesOrError } from "../lib/contracts/quiz";
	import CreateDialog from "./CreateDialog.svelte";
	import CreateQuizForm from "./CreateQuizForm.svelte";

	let height = $state(0);

	let page = $state(1);
	let pageSize = $derived(Math.trunc((height - 39 - 45 - 40) / 39));

	let trigger = $state(0);
	let time: number;
	const paginatedQuizPromises = $derived.by(() => {
		const p = page;
		const ps = pageSize;
		trigger;

		console.log("detected change");
		clearTimeout(time);

		return new Promise<QuizzesOrError>((resolve) => {
			time = setTimeout(() => {
				resolve(FetchAllQuizzes(p - 1, ps));
			}, 500);
		});
	});
</script>

<div
	class="grid gap-4 w-full place-items-center h-full overflow-auto"
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
			{#if paginatedQuizzes.total == 0}
				<CreateDialog title="Quiz creator" name="Create your first quiz"
					><CreateQuizForm
						callback={() => {
							trigger++;
						}}
					/></CreateDialog
				>
			{:else}
				<table class="table table-auto self-start">
					<thead>
						<tr class="text-surface-100-900">
							<th>Quiz path</th>
							<th>Max score</th>
						</tr>
					</thead>

					<tbody>
						{#each paginatedQuizzes.quizzes as quiz}
							<tr>
								<td>{quiz.path}</td>
								<!-- 10.5 rem = 168px -->
								<td>{quiz.score}</td>
							</tr>
						{/each}
					</tbody>
				</table>

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
										<Pagination.Ellipsis {index}>&#8230;</Pagination.Ellipsis>
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
