<script lang="ts">
  import ArrowLeftIcon from "@lucide/svelte/icons/arrow-left";
  import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
  import { IsJSONError, type JSONError } from "../lib/statuses/jsonerror";
  import {
    FetchAllQuizzes,
    FetchQuiz,
    type Quiz as QuizType,
    type QuizFile,
    FetchQuizExport,
    FetchQuizImport,
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
  import { SvelteMap, SvelteSet } from "svelte/reactivity";
  import Exporter from "./Exporter.svelte";
  import BestAccept from "./BestAccept.svelte";
  import type { Accept } from "../lib/apiutils/acceptHeader";
  import { Dialog, Popover, Portal } from "@skeletonlabs/skeleton-svelte";
  import PeekDialogue from "./PeekDialogue.svelte";
  import { CircleArrowDown, CircleArrowUp } from "@lucide/svelte";
  import Importer from "./Importer.svelte";

  const fit = CreateFitRows({ rowHeight: 20, reserve: 1 });
  let pageSize = $derived(fit.pageSize);
  let page = $state(1);

  let loading = $state(true);
  let statusMessage = $state("");
  let quizzes: QuizType[] = $state(null as never);
  let total: number = $state(0);

  let selected = $state(new SvelteMap<string, string>());
  let best = $state([] as Accept[]);

  let trig = $state(0);
  let time: number;
  $effect(() => {
    const p = page,
      ps = pageSize,
      ready = fit.ready;
    trig;
    if (!ready) {
      console.log("not ready yet...");
      return;
    }

    loading = true;
    clearTimeout(time);

    time = setTimeout(async () => {
      statusMessage = await FetchAllQuizzes(p - 1, ps)
        .andTee((r) => {
          ({ quizzes, total } = r);
        })
        .match(
          () => "",
          (err: JSONError) => err.error,
        );

      loading = false;
    }, 500);
  });

  $inspect(loading);

  let focused = $state("");
  let clickFocused = $state("");
</script>

<div
  class="grid gap-4 w-full place-items-center h-full overflow-auto"
  use:fit.measure
>
  {#if loading}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:else}
    {#if statusMessage !== ""}
      <div>{statusMessage}</div>
    {:else}
      {#if total === 0}
        <CreateDialog title="Quiz creator" name="Create your first quiz"
          ><CreateQuizForm
            callback={() => {
              trig++;
            }}
          /></CreateDialog
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
                {#each quizzes as quiz (quiz.uuid)}
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
                    <!-- 10.5 rem = 168px -->
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
                              callback={() => {
                                trig++;
                              }}
                            />
                          </DeleteDialogSquare>
                          <ChangeDialogSqare
                            callback={() => (clickFocused = quiz.uuid)}
                            title="Quiz changer"
                            contentGetter={async () => {
                              return await FetchQuiz(quiz.uuid);
                            }}
                          >
                            {#snippet children({
                              content,
                            }: {
                              content: QuizFile;
                            })}
                              <ChangeQuizForm
                                UUID={quiz.uuid}
                                name={name.endsWith(".md")
                                  ? name.slice(0, -3)
                                  : name}
                                callback={() => {
                                  trig++;
                                }}
                                response={content}
                              />
                            {/snippet}</ChangeDialogSqare
                          >

                          <PeekDialogSquare
                            callback={() => (clickFocused = quiz.uuid)}
                            title="Quiz peeker"
                            contentGetter={async () => {
                              const response = await FetchQuiz(quiz.uuid);
                              if (IsJSONError(response)) {
                                return response.error;
                              }
                              return response.body;
                            }}
                          >
                            {#snippet children({
                              content,
                            }: {
                              content: string;
                            })}
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
            <div
              class="flex justify-between items-center gap-4 w-full self-end"
            >
              <Pagination bind:page {pageSize} {total} />

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
                  <CreateQuizForm
                    callback={() => {
                      trig++;
                    }}
                  />
                </SmallPaperCreateDialog>

                <PeekDialogue title="Quiz importing">
                  {#snippet trigger()}
                    <Dialog.Trigger class="btn preset-filled-warning-500">
                      <CircleArrowDown></CircleArrowDown>
                    </Dialog.Trigger>
                  {/snippet}

                  <Importer
                    callback={() => {
                      trig++;
                    }}
                    importName="quiz import"
                    fetcher={(file: File) => FetchQuizImport(file)}
                  />
                </PeekDialogue>
              </div>
            </div>
          </footer>
        </div>
      {/if}
    {/if}
  {/if}
</div>
