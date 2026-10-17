<script lang="ts">
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import { type JSONError } from "../lib/statuses/jsonerror";
  import {
    FetchAllTests,
    FetchTest,
    FetchTestImport,
    type Test,
  } from "../lib/contracts/test";
  import CreateDialog from "./CreateDialog.svelte";
  import CreateTestForm from "./CreateTestForm.svelte";
  import SmallBookCreateDialog from "./SmallBookCreateDialog.svelte";
  import DeleteDialogSquare from "./DeleteDialogSquare.svelte";
  import DeleteTestForm from "./DeleteTestForm.svelte";
  import ChangeDialogSqare from "./ChangeDialogSqare.svelte";
  import ChangeTestForm from "./ChangeTestForm.svelte";
  import PeekDialogSquare from "./PeekDialogSquare.svelte";
  import TestPeek from "./TestPeek.svelte";
  import { CreateFitRows } from "../lib/layoututils/fitRows.svelte";
  import PeekDialogue from "./PeekDialogue.svelte";
  import Importer from "./Importer.svelte";
  import { CircleArrowDown } from "@lucide/svelte";
  import Pagination from "./Pagination.svelte";
  import { CreatePagedQuery } from "../lib/layoututils/createPagedQuery.svelte";

  const fit = CreateFitRows({ rowHeight: 20, reserve: 1 });

  const q = CreatePagedQuery({
    pageSize: () => fit.pageSize,
    enabled: () => fit.ready,
    params: () => ({}),
    fetcher: (page, size, p) => FetchAllTests(page, size),
  });

  let focused = $state("");
  let clickFocused = $state("");
</script>

<div
  class="grid gap-4 w-full place-items-center h-full overflow-auto"
  use:fit.measure
>
  {#if q.loading}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:else if q.error}
    <div>{q.error}</div>
  {:else if !q.data || q.data.total === 0}
    <CreateDialog name="Create your first test" title="Create your first test"
      ><CreateTestForm callback={q.refresh} /></CreateDialog
    >
  {:else}
    <table class="table table-auto self-start">
      <thead>
        <tr class="text-surface-100-900 flex" style="height:20px">
          <th class="w-1/3">Test</th>
          <th class="w-1/3">Quiz count</th>
          <th class="w-1/3"></th>
        </tr>
      </thead>

      <tbody>
        {#each q.data.tests as test}
          <tr
            onmouseenter={() => (focused = test.uuid)}
            onmouseleave={() => (focused = "")}
            style="height:20px"
            class="bg-surface-700-300 rounded-xl flex hover:motion-safe:hover:brightness-125 dark:hover:motion-safe:hover:brightness-75"
          >
            <td class="w-1/3">{test.name}</td>
            <td class="w-1/3">{test.quizzes?.length ?? 0}</td>
            <td class="w-1/3">
              {#if test.uuid === focused || test.uuid === clickFocused}
                <div class="flex flex-row-reverse">
                  <DeleteDialogSquare
                    title="Quiz deleter"
                    callback={() => (clickFocused = test.uuid)}
                  >
                    <DeleteTestForm
                      name={test.name}
                      UUID={test.uuid}
                      callback={q.refresh}
                    />
                  </DeleteDialogSquare>
                  <ChangeDialogSqare
                    callback={() => (clickFocused = test.uuid)}
                    title="Test changer"
                    contentGetter={() => test}
                  >
                    {#snippet children({ content }: { content: Test })}
                      <ChangeTestForm
                        UUID={test.uuid}
                        name={test.name}
                        callback={q.refresh}
                        response={content}
                      />
                    {/snippet}</ChangeDialogSqare
                  >

                  <PeekDialogSquare
                    callback={() => (clickFocused = test.uuid)}
                    title="Test peeker"
                    contentGetter={() => test}
                  >
                    {#snippet children({ content }: { content: Test })}
                      <TestPeek {content} />
                    {/snippet}
                  </PeekDialogSquare>
                </div>
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>

    <div class="flex justify-between items-center gap-4 w-full self-end">
      <Pagination
        bind:page={q.page}
        pageSize={fit.pageSize}
        total={q.data?.total ?? 0}
      />

      <div class="flex ms-auto gap-2">
        <SmallBookCreateDialog title="Test maker"
          ><CreateTestForm callback={q.refresh} /></SmallBookCreateDialog
        >

        <PeekDialogue title="Test importing">
          {#snippet trigger()}
            <Dialog.Trigger class="btn preset-filled-warning-500">
              <CircleArrowDown></CircleArrowDown>
            </Dialog.Trigger>
          {/snippet}

          <Importer
            callback={q.refresh}
            importName="test import"
            fetcher={(file: File) => FetchTestImport(file)}
          />
        </PeekDialogue>
      </div>
    </div>
  {/if}
</div>
