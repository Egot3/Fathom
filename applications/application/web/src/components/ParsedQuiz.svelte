<script lang="ts">
  import {
    FetchParsedQuiz,
    Kind,
    type QuizAnswer,
    type QuizOptions,
  } from "../lib/contracts/quiz";
  import AnsweredQuiz from "./AnsweredQuiz.svelte";

  const {
    UUID,
    answeredValue,
  }: {
    UUID: string;
    answeredValue: QuizAnswer;
  } = $props();

  let loading: boolean = $state(false);
  let statusMessage: string = $state("");

  let options: QuizOptions | null = $state(null);
  let answers: QuizAnswer | undefined = $state(undefined);
  let kind: Kind | null = $state(null);

  let title: string = $state("");
  let body: string = $state("");

  let time: number;
  $effect(() => {
    console.log("looping");
    loading = true;
    const quizUUID: string = UUID;

    clearTimeout(time)
    time = setTimeout(async () => {
      statusMessage = await FetchParsedQuiz(quizUUID).match(
        (r) => {
          console.log(r);

          console.log(title, body);
          title = r.title;
          body = r.body;
          answers = r.answers;
          options = r.options;
          kind = r.meta.kind;

          console.log(title, body);
          return "";
        },
        (e) => e.error,
      );

      loading = false;
    }, 1000);
  });
</script>

{#if loading }
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
{:else}
{#if statusMessage!==""}
 <span class="justify-self-start">{statusMessage}</span>
{:else}
<div class="space-y-1 mt-4">
  <h2 class="text-5xl font-bold">{title}</h2>

  <p class="text-2xl">{body}</p>

  {#if options !== null && kind !== null}
    <AnsweredQuiz
      disabled={true}
      {kind}
      {options}
      answers={answeredValue}
      correct={answers}
    ></AnsweredQuiz>
  {/if}
</div>

{/if}

{/if}
