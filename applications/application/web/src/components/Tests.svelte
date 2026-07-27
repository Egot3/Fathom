<script lang="ts">
	import ArrowLeftIcon from "@lucide/svelte/icons/arrow-left";
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import { Pagination } from "@skeletonlabs/skeleton-svelte";
	import { IsJSONError } from "../lib/statuses/jsonerror";
	import { FetchAllTests, type TestsOrError } from "../lib/contracts/test";
	import CreateTest from "./CreateDialog.svelte";
	import CreateDialog from "./CreateDialog.svelte";
	import CreateTestForm from "./CreateTestForm.svelte";

	let height = $state(0);

	let page = $state(1);
	let pageSize = $derived(Math.trunc((height - 39 - 45) / 39));

	let time: number;
	const paginatedTestPromises = $derived.by(() => {
		const p = page;
		const ps = pageSize;

		console.log("detected change");
		clearTimeout(time);

		return new Promise<TestsOrError>((resolve) => {
			time = setTimeout(() => {
				resolve(FetchAllTests(p - 1, ps));
			}, 500);
		});
	});

	let creatingTest = $state(false);
</script>

<div
	class="grid gap-4 w-full place-items-center h-full overflow-auto"
	bind:clientHeight={height}
>
	{#await paginatedTestPromises}
		<div
			class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
		></div>
	{:then paginatedTests}
		{#if IsJSONError(paginatedTests)}
			<div>{paginatedTests.error}</div>
		{:else}
			{#if paginatedTests.total == 0}
				<CreateDialog title="Create your first test"
					><CreateTestForm /></CreateDialog
				>
			{:else}
				<table class="table table-auto self-start">
					<thead>
						<tr class="text-surface-100-900">
							<th>Test</th>
							<th>Quiz count</th>
						</tr>
					</thead>

					<tbody>
						{#each paginatedTests.tests as test}
							<tr>
								<td>{test.name}</td>
								<!-- 10.5 rem = 168px -->
								<td>{test.quizzes?.length ?? 0}</td>
							</tr>
						{/each}
					</tbody>
				</table>

				<div class="flex justify-between items-center gap-4 w-full self-end">
					<Pagination
						count={paginatedTests?.total ?? 0}
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
