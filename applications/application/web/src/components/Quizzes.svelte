<script lang="ts">
	import ArrowLeftIcon from "@lucide/svelte/icons/arrow-left";
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import { Pagination } from "@skeletonlabs/skeleton-svelte";
	import { IsJSONError } from "../lib/statuses/jsonerror";
	import { FetchAllQuizzes, type QuizzesOrError } from "../lib/contracts/quiz";
	import CreateDialog from "./CreateDialog.svelte";
	import CreateQuizForm from "./CreateQuizForm.svelte";
	import SmallPaperCreateDialog from "./SmallPaperCreateDialog.svelte";
	import { Eye, EyeClosed, Pencil, Trash, Trash2 } from "@lucide/svelte";

	let height = $state(0);

	let page = $state(1);
	let pageSize = $derived(Math.trunc((height - 29 - 35 - 20) / 39));

	let focused = $state("")
	let moused = $state(false)

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
						<tr class="text-surface-100-900 flex">
							<th class="w-1/3">Quiz path</th>
							<th class="w-1/3">Max score</th>
							<th class="w-1/3"></th>
						</tr>
					</thead>

					<tbody>
						{#each paginatedQuizzes.quizzes as quiz (quiz.uuid)}
							
							<tr
								onmouseenter={()=>focused=quiz.uuid}
								onmouseleave={()=>focused=""}
							 	class="bg-surface-700-300 rounded-xl flex hover:motion-safe:hover:brightness-125 dark:hover:motion-safe:hover:brightness-75">
								<td class="w-1/3">{quiz.path.startsWith("/data/quizzes/") ? quiz.path.slice(14) : quiz.path}</td>
								<!-- 10.5 rem = 168px -->
								<td class="w-1/3">{quiz.score}</td>
								<td class="w-1/3">
									{#if quiz.uuid===focused}
										<div class="flex flex-row-reverse">
											<button class="btn mb-0 mt-0 preset-filled-surface-300-700 hover:preset-filled-error-300-700 aspect-square h-auto w-auto p-1 m-px">
												<Trash2 size=16/>
											</button>
											<button class="btn mb-0 mt-0 preset-filled-surface-300-700 hover:preset-filled-warning-300-700 aspect-square h-auto w-auto p-1 m-px">
												<Pencil size=16/>
											</button>
											<button onmouseenter={()=>moused=true} onmouseleave={()=>moused=false} class="btn mb-0 mt-0 preset-filled-surface-300-700 hover:preset-filled-primary-300-700 aspect-square h-auto w-auto p-1 m-px">
												{#if moused}
													<Eye size=16/>
												{:else}
												
													<EyeClosed size=16/>
												{/if}
												
											</button>
										</div>
									{/if}
								</td>
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

					<SmallPaperCreateDialog title="Quiz maker">
						<CreateQuizForm 
							callback={()=>{
								trigger++
							}}
						/>
					</SmallPaperCreateDialog>
				</div>
			{/if}
		{/if}
	{/await}
</div>
