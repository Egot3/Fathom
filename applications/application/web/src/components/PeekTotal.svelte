<script lang="ts">
    import type { QuizAnswer } from "../lib/contracts/quiz";
  import { FetchAnswers, type Answer, type Answers } from "../lib/contracts/totals";
  import ParsedQuiz from "./ParsedQuiz.svelte";

  const {
    userUUID,
    groupUUID,
    testUUID,
  }: {
    userUUID: string;
    groupUUID: string;
    testUUID: string;
  } = $props();

  let height = $state(0);

  let page = $state(1);
  let pageSize = $state(5);

  let answers: Answers = $state(null as never);
  let loading = $state(true);
  let statusMessage = $state("");
  $effect(() => {
    loading = true;
    (async () => {
      statusMessage = (
        await FetchAnswers(groupUUID, userUUID, testUUID, page - 1, pageSize)
      )
        .andTee((r) => {
          loading = false;
          answers = r;
        })
        .match(
          (_) => "",
          (err) => err.error,
        );
    })();
  });

  let chosenUUID: string = $state("");
  let answerValue: QuizAnswer = $state(null as never);
  $effect(()=>{
    answerValue = answers?.answers?.[0]?.chosen ?? null as never
  })
</script>

<div class="flex flex-col h-23/25">
  {#if loading}
    <div>loading...</div>
  {:else}
    {#if statusMessage !== ""}
      <div>{statusMessage}</div>
    {:else}
      <div class="grid grid-cols-12 w-full h-full mt-2">
        <div class="col-start-1 col-end-9 p-3">
          <ParsedQuiz
              answeredValue={answerValue}
            UUID={chosenUUID || answers.answers[0].quiz_uuid}
          />
        </div>
        <div
          class="col-start-9 col-end-13 bg-surface-600-400 rounded-[0.25rem] p-1"
        >
          {#each answers.answers as answer}
            <button
              class="btn preset-filled-error-50-950 w-full text-xl h-fit"
              onclick={() => {
                answerValue = answer.chosen;
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
        </div>
      </div>
    {/if}
  {/if}
</div>
