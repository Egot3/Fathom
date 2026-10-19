<script lang="ts">
  import ParsedQuiz from "../components/ParsedQuiz.svelte";
  import type { QuizAnswer } from "../lib/contracts/quiz";
  import {
    FetchCurrentlyRunningTestQuiz,
    FetchCurrentyRunningTestQuizzes,
  } from "../lib/contracts/test";

  const { runnerKey }: { runnerKey: string } = $props();

  let quizUUIDs: string[] = $state([]);
  let chosenUUID: string = $state(null as never);

  let statusMessage: string = $state("");
  let loading: boolean = $state(true);

  $effect(() => {
    (async () => {
      loading = true;
      statusMessage = await FetchCurrentyRunningTestQuizzes(runnerKey)
        .andTee((r) => {
          quizUUIDs = r;
          if (chosenUUID === null) chosenUUID = quizUUIDs[0];
        })
        .match(
          (_) => "",
          ({ error }) => (statusMessage = error),
        );
      loading = false;
    })();
  });
</script>

<div>
  {#if loading}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:else if statusMessage}
    <div>{statusMessage}</div>
  {:else}
    <ParsedQuiz
      UUID={chosenUUID}
      quizGetter={(UUID) => FetchCurrentlyRunningTestQuiz(runnerKey, UUID)}
      answeredValue={{
        input: { input: "1" },
      } as QuizAnswer}
    />
  {/if}
</div>
