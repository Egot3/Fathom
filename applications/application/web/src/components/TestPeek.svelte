<script lang="ts">
  import { SvelteSet } from "svelte/reactivity";
  import {
    FetchTest,
    FetchTestExport,
    FetchTestQuizzes,
    type Test,
  } from "../lib/contracts/test";
  import ExistingQuizChips from "./ExistingQuizChips.svelte";
  import Exporter from "./Exporter.svelte";
  import AcceptHeaders from "./AcceptHeaders.svelte";
  import { USER_MIMES, type Accept } from "../lib/apiutils/acceptHeader";
  import { CreateFitRows } from "../lib/layoututils/fitRows.svelte";
  import { CreatePagedQuery } from "../lib/layoututils/createPagedQuery.svelte";
  import Pagination from "./Pagination.svelte";

  let { content }: { content: Test } = $props();

  let best: Accept[] = $state([]);
  const q = CreatePagedQuery({
    pageSize: () => 5,
    enabled: () => true,
    params: () => content.uuid,
    fetcher: FetchTestQuizzes,
  });
</script>

<label class="label">
  <span class="label-text">Name</span>
  <input
    class="input"
    type="text"
    value={content.name}
    placeholder="go-generics"
    disabled
  />
</label>
{#if q.loading}
  <div
    class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl col-start-1 col-end-13"
  ></div>
{:else if q.error !== ""}
  <div>{q.error}</div>
{:else if q.data === null}
  <div>no quizzes to list</div>
{:else}
  {const quizzes = new SvelteSet(q.data.quizzes)}
  <label class="label">
    <span class="label-text">Quizzes</span>
    <ExistingQuizChips disabled {quizzes} />
  </label>
{/if}
<footer class="flex justify-between items-center gap-4 w-full mt-auto">
  <Pagination bind:page={q.page} pageSize={5} total={q.data?.total ?? 0} />
</footer>

<div class="flex space-x-2 p-2">
  <details class="disclosure">
    <summary>Export options</summary>
    <div class="disclosure-content">
      <AcceptHeaders available={USER_MIMES} bind:accepted={best} />
    </div>
  </details>

  <div class="mt-auto mb-2">
    <Exporter
      downloadName={content.name}
      fetcher={() => FetchTestExport(content.uuid, best)}
    >
      Export
    </Exporter>
  </div>
</div>
