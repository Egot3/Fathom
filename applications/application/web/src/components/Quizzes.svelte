<script lang="ts">
  import {
    FetchAllQuizzes,
    type Quiz as QuizType,
    type QuizFile,
    FetchQuizExport,
    FetchQuizImport,
    FetchParsedQuiz,
  } from "../lib/contracts/quiz";
  import CreateDialog from "./CreateDialog.svelte";
  import CreateQuizForm from "./CreateQuizForm.svelte";
  import SmallPaperCreateDialog from "./SmallPaperCreateDialog.svelte";
  import ChangeQuizForm from "./ChangeQuizForm.svelte";
  import ChangeDialogSqare from "./ChangeDialogSqare.svelte";
  import PeekDialogSquare from "./PeekDialogSquare.svelte";
  import DeleteDialogSquare from "./DeleteDialogSquare.svelte";
  import DeleteQuizForm from "./DeleteQuizForm.svelte";
  import { CreateFitRows } from "../lib/layoututils/fitRows.svelte";
  import Quiz from "./Quiz.svelte";
  import Pagination from "./Pagination.svelte";
  import Exporter from "./Exporter.svelte";
  import BestAccept from "./BestAccept.svelte";
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import PeekDialogue from "./PeekDialogue.svelte";
  import { CircleArrowDown, CircleArrowUp } from "@lucide/svelte";
  import Importer from "./Importer.svelte";
  import { CreatePagedQuery } from "../lib/layoututils/createPagedQuery.svelte";
  import { SvelteMap } from "svelte/reactivity";
  import type { Accept } from "../lib/apiutils/acceptHeader";
  import { okAsync } from "neverthrow";

  const fit = CreateFitRows({ rowHeight: 20, reserve: 1 });
  const q = CreatePagedQuery({
    pageSize: () => fit.pageSize,
    enabled: () => fit.ready,
    params: () => ({}),
    fetcher: (page, size, p) => FetchAllQuizzes(page, size),
  });

  let focused = $state("");
  let clickFocused = $state("");
  let selected = $state(new SvelteMap<string, string>());
  let best = $state([] as Accept[]);
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
    <CreateDialog title="Quiz creator" name="Create your first quiz"
      ><CreateQuizForm callback={q.refresh} /></CreateDialog
    >
  {:else}
    <div class="flex flex-col h-full min-h-0 w-full">
      <div class="flex-1 min-h-0 overflow-hidden">
        <table class="table table-auto self-start w-full">
          <thead>
            <tr class="text-surface-100-900 flex" style="height:20px">
              <th class="w-1/3">Quiz path</th>
              <th class="w-1/3">Max score</th>
              <th class="w-1/3"></th>
            </tr>
          </thead>

          <tbody>
            {#each q.data.quizzes as quiz (quiz.uuid)}
              {const name = quiz.path.startsWith("/data/quizzes/")
                ? quiz.path.slice(14)
                : quiz.path}
              <tr
                style="height:20px"
                onmouseenter={() => (focused = quiz.uuid)}
                onmouseleave={() => (focused = "")}
                class="bg-surface-700-300 rounded-xl flex hover:motion-safe:hover:brightness-125 dark:hover:motion-safe:hover:brightness-75"
              >
                <td class="w-1/3">
                  <input
                    class="checkbox"
                    type="checkbox"
                    onclick={(e) => {
                      const { checked } = e.target as HTMLInputElement;
                      if (checked) {
                        selected.set(quiz.uuid, name);
                      } else {
                        selected.delete(quiz.uuid);
                      }
                    }}
                  />
                  {name}
                </td>
                <td class="w-1/3">{quiz.score}</td>
                <td class="w-1/3">
                  {#if quiz.uuid === focused || quiz.uuid === clickFocused}
                    <div class="flex flex-row-reverse">
                      <DeleteDialogSquare
                        title="Quiz deleter"
                        callback={() => (clickFocused = quiz.uuid)}
                      >
                        <DeleteQuizForm
                          {name}
                          UUID={quiz.uuid}
                          callback={q.refresh}
                        />
                      </DeleteDialogSquare>
                      <ChangeDialogSqare
                        callback={() => (clickFocused = quiz.uuid)}
                        title="Quiz changer"
                        contentGetter={() => quiz}
                      >
                        {#snippet children({ content }: { content: QuizFile })}
                          <ChangeQuizForm
                            UUID={quiz.uuid}
                            name={name.endsWith(".md")
                              ? name.slice(0, -3)
                              : name}
                            callback={q.refresh}
                            response={content}
                          />
                        {/snippet}</ChangeDialogSqare
                      >

                      <PeekDialogSquare
                        callback={() => (clickFocused = quiz.uuid)}
                        title="Quiz peeker"
                        contentGetter={() =>
                          FetchParsedQuiz(quiz.uuid).andThen((r) =>
                            okAsync(r.body),
                          )}
                      >
                        {#snippet children({ content }: { content: string })}
                          <Quiz content={content ?? ""} />
                        {/snippet}
                      </PeekDialogSquare>
                    </div>
                  {/if}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
      <footer>
        <div class="flex justify-between items-center gap-4 w-full self-end">
          <Pagination
            bind:page={q.page}
            pageSize={fit.pageSize}
            total={q.data?.total ?? 0}
          />

          <div>
            <PeekDialogue title="Quiz exporter">
              {#snippet trigger()}
                <Dialog.Trigger
                  disabled={selected.size == 0}
                  class="btn preset-filled-warning-500"
                >
                  <CircleArrowUp></CircleArrowUp>
                </Dialog.Trigger>
              {/snippet}
              Chosen: {Array.from(selected.values()).join(", ")}
              <BestAccept bind:best></BestAccept>
              <Exporter
                disabled={selected.size == 0}
                fetcher={() =>
                  FetchQuizExport(Array.from(selected.keys()), best)}
                downloadName="quiz"
              />
            </PeekDialogue>

            <SmallPaperCreateDialog title="Quiz maker">
              <CreateQuizForm callback={q.refresh} />
            </SmallPaperCreateDialog>

            <PeekDialogue title="Quiz importing">
              {#snippet trigger()}
                <Dialog.Trigger class="btn preset-filled-warning-500">
                  <CircleArrowDown></CircleArrowDown>
                </Dialog.Trigger>
              {/snippet}

              <Importer
                callback={q.refresh}
                importName="quiz import"
                fetcher={(file: File) => FetchQuizImport(file)}
              />
            </PeekDialogue>
          </div>
        </div>
      </footer>
    </div>
  {/if}
</div>
