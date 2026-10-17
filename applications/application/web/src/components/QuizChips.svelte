<script lang="ts">
  import { FetchAllQuizzes, type Quiz } from "../lib/contracts/quiz";
  import { CheckIcon } from "@lucide/svelte";
  import { SvelteSet } from "svelte/reactivity";
  import Pagination from "./Pagination.svelte";
  import { CreateFitRows } from "../lib/layoututils/fitRows.svelte";
  import { CreatePagedQuery } from "../lib/layoututils/createPagedQuery.svelte";

  let {
    chosen = $bindable(new SvelteSet<string>()),
    existing = new SvelteSet<Quiz>(),
  }: {
    chosen: SvelteSet<string>;
    existing?: SvelteSet<Quiz>;
  } = $props();

  const fit = CreateFitRows({ rowHeight: 20, reserve: 1 });

  const q = CreatePagedQuery({
    pageSize: () => fit.pageSize,
    enabled: () => true,
    params: () => ({}),
    fetcher: FetchAllQuizzes,
  });

  let existingQuizUUIDs = $derived(
    new Set(Array.from(existing).map((v) => v.uuid)),
  );
</script>

<div
  class="grid gap-4 m-2 w-full place-items-center h-full overflow-auto"
  use:fit.measure
>
  {#if q.loading}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:else if q.error}
    <div>{q.error}</div>
  {:else if !q.data || q.data.total === 0}
    <div>No quizzes</div>
  {:else}
    {#each q.data.quizzes as quiz (quiz.uuid)}
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
      <Pagination bind:page={q.page} pageSize={5} total={q.data?.total ?? 0} />
    </div>
  {/if}
</div>
