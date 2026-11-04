<script lang="ts">
  import { FetchAnswers } from "../lib/contracts/totals";
  import { CreatePagedQuery } from "../lib/layoututils/createPagedQuery.svelte";
  import { CreateFitRows } from "../lib/layoututils/fitRows.svelte";
  import QuizView from "./QuizView.svelte";
  import Pagination from "./Pagination.svelte";

  const {
    userUUID,
    groupUUID,
    testUUID,
  }: {
    userUUID: string;
    groupUUID: string;
    testUUID: string;
  } = $props();

  const fit = CreateFitRows({ rowHeight: 27, reserve: 1 });
  const q = CreatePagedQuery({
    pageSize: () => fit.pageSize,
    enabled: () => fit.ready,
    params: () => ({}),
    fetcher: (page, size) =>
      FetchAnswers(groupUUID, userUUID, testUUID, page, size),
  });

  let chosenUUID: string = $state("");

  let currentUUID = $derived(chosenUUID || q.data?.answers?.[0]?.quiz_uuid); // huh?
  let current = $derived(
    q.data?.answers?.find((a) => a.quiz_uuid === currentUUID),
  );
  let chosen = $derived(current?.chosen);
</script>

<div class="grid grid-cols-12 w-full flex-1 overflow-scroll">
  <div class="col-start-1 col-end-9 p-3">
    {#if currentUUID !== undefined}
      <QuizView UUID={currentUUID} {chosen} correct={current?.correct} />
    {:else}
      <div
        class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
      ></div>
    {/if}
  </div>
  <div
    class="flex flex-col col-start-9 col-end-13 bg-surface-600-400 rounded-[0.25rem] p-1 gap-2"
    use:fit.measure
  >
    {#if q.loading}
      <div
        class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl col-start-1 col-end-13"
      ></div>
    {:else if q.error !== ""}
      <div>{q.error}</div>
    {:else if q.data === null}
      <div>no totals to list</div>
    {:else}
      {#each q.data?.answers as answer}
        <button
          class="btn text-surface-50-950 bg-surface-900-100 w-full text-xl h-fit"
          onclick={() => {
            chosenUUID = answer.quiz_uuid;
          }}
        >
          <div>
            {answer.quiz_name.slice(14).slice(0, -3)}
          </div>
          <div>
            {answer.score}/{answer.max_score}
          </div>
        </button>
      {/each}
    {/if}
    <footer class="flex justify-between items-center gap-4 w-full mt-auto">
      <Pagination
        bind:page={q.page}
        pageSize={fit.pageSize}
        total={q.data?.total ?? 0}
      />
    </footer>
  </div>
</div>
