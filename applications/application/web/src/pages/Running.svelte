<script lang="ts">
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import Confirm from "../components/Confirm.svelte";
  import PeekDialogBig from "../components/PeekDialogBig.svelte";
  import PeekDialogue from "../components/PeekDialogue.svelte";
  import QuizView from "../components/QuizView.svelte";
  import { Toaster } from "../lib/apiutils/toaster";
  import {
    FetchAnswerPost,
    FetchMyAnswers,
    type Quiz,
    type QuizAnswer,
  } from "../lib/contracts/quiz";
  import {
    FetchRunningTestInfo,
    FetchRunningTestQuizzes,
    FetchTestTotalize,
    type TestInfo,
  } from "../lib/contracts/test";

  const { runnerKey }: { runnerKey: string } = $props();

  let quizzes: Quiz[] = $state([]);
  let currentIdx: number = $state(0);
  let current: Quiz | null = $derived(quizzes[currentIdx] ?? null);

  let statusMessage: string = $state("");
  let loadingQuizzes: boolean = $state(true);
  let loadingTestInfo: boolean = $state(true);

  let answers = $state<Record<string, QuizAnswer>>({});

  $effect(() => {
    (async () => {
      loadingQuizzes = true;
      statusMessage = await FetchRunningTestQuizzes(runnerKey)
        .andTee((r) => {
          quizzes = r;
        })
        .match(
          (_) => "",
          ({ error }) => error,
        );

      await FetchMyAnswers(runnerKey).match(
        (r) => {
          answers = { ...r, ...answers };
        },
        (e) => {
          statusMessage = e.error;
        },
      );
      loadingQuizzes = false;
    })();
  });

  let info: TestInfo | null = $state(null);
  $effect(() => {
    currentIdx;
    (async () => {
      loadingTestInfo = true;
      statusMessage = await FetchRunningTestInfo(runnerKey)
        .andTee((r) => {
          info = r;
        })
        .match(
          (_) => "",
          ({ error }) => error,
        );
      loadingTestInfo = false;
    })();
  });
</script>

<div>
  {#if loadingQuizzes || loadingTestInfo}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:else if statusMessage}
    <div>{statusMessage}</div>
  {:else if info?.is_paused}
    <div class="grid grid-cols-12 grid-rows-12 w-screen h-screen p-16 gap-4">
      <h1
        class="row-start-3 row-end-11 col-start-1 col-end-13 p-8 rounded-2xl text-surface-50-950 bg-surface-900-100 font-bold flex items-center justify-center"
      >
        test is currently paused
      </h1>
    </div>
  {:else}
    <div class="grid grid-cols-12 grid-rows-12 w-screen h-screen p-16 gap-4">
      <div
        class="row-start-3 row-end-11 col-start-1 col-end-10 p-8 rounded-2xl text-surface-50-950 bg-surface-900-100"
      >
        {#if current !== null}
          <QuizView
            UUID={current.uuid}
            {runnerKey}
            bind:chosen={answers[current.uuid]}
            onanswer={(a) => {
              current !== null
                ? FetchAnswerPost(runnerKey, current.uuid, a)
                    .andTee((_) => currentIdx++)
                    .orTee((e) =>
                      Toaster.error({
                        title: "Couldn't send answer!",
                        description: e.error,
                      }),
                    )
                : {};
            }}
          />
        {/if}
      </div>
      <div
        class="row-start-3 row-end-11 col-start-10 col-end-13 p-4 rounded-2xl flex flex-col bg-surface-200-800"
      >
        <div class="h-11/12 flex flex-col gap-1">
          {#each quizzes as quiz (quiz.uuid)}
            {const quizName = quiz.path.slice(14).slice(0, -3)}
            <button
              class="items-center justify-center w-full flex bg-surface-300-700"
              onclick={() => {
                current = quiz;
              }}
            >
              {quizName}
            </button>
          {/each}
        </div>

        <PeekDialogue title="Test finisher">
          <Confirm
            value="Finish"
            name="Totalize"
            key={runnerKey}
            actor={FetchTestTotalize}
            callback={() => {
              Toaster.info({
                title: "Finished",
                description: "Your answers were send successfuly",
              });
              window.location.href = "/home";
            }}
          />
          {#snippet trigger()}
            <Dialog.Trigger class="h-1/12 btn preset-filled-brand"
              >Totalize</Dialog.Trigger
            >
          {/snippet}
        </PeekDialogue>
      </div>
    </div>
  {/if}
</div>
