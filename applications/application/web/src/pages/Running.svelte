<script lang="ts">
  import QuizView from "../components/QuizView.svelte";
  import {
    FetchAnswerPost,
    type Quiz,
    type QuizAnswer,
  } from "../lib/contracts/quiz";
  import {
    FetchRunningTestInfo,
    FetchRunningTestQuizzes,
    type TestInfo,
  } from "../lib/contracts/test";

  const { runnerKey }: { runnerKey: string } = $props();

  let quizzes: Quiz[] = $state([]);
  let current: Quiz | null = $state(null);

  let statusMessage: string = $state("");
  let loadingQuizzes: boolean = $state(true);
  let loadingTestInfo: boolean = $state(true);

  $effect(() => {
    (async () => {
      loadingQuizzes = true;
      statusMessage = await FetchRunningTestQuizzes(runnerKey)
        .andTee((r) => {
          quizzes = r;
          if (current === null) current = quizzes[0];
        })
        .match(
          (_) => "",
          ({ error }) => error,
        );
      loadingQuizzes = false;
    })();
  });

  let info: TestInfo | null = $state(null);
  $effect(() => {
    current;
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

  let answers = $state<Record<string, QuizAnswer>>({});
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
            onanswer={(a) =>
              current !== null
                ? FetchAnswerPost(runnerKey, current!.uuid, a)
                : {}}
          />
        {/if}
      </div>
      <div
        class="row-start-3 row-end-11 col-start-10 col-end-13 p-4 rounded-2xl flex flex-col gap-1 bg-surface-200-800"
      >
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
    </div>
  {/if}
</div>
