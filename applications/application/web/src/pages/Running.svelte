<script lang="ts">
  import QuizView from "../components/QuizView.svelte";
  import { FetchAnswerPost, type QuizAnswer } from "../lib/contracts/quiz";
  import { FetchRunningTestQuizzes } from "../lib/contracts/test";
  import { GetUserOrRedirect } from "../lib/bgdata/user.svelte";

  const { runnerKey }: { runnerKey: string } = $props();

  let quizUUIDs: string[] = $state([]);
  let current: string = $state(null as never);

  let statusMessage: string = $state("");
  let loading: boolean = $state(true);

  $effect(() => {
    (async () => {
      loading = true;
      statusMessage = await FetchRunningTestQuizzes(runnerKey)
        .andTee((r) => {
          quizUUIDs = r;
          if (current === null) current = quizUUIDs[0];
        })
        .match(
          (_) => "",
          ({ error }) => error,
        );
      loading = false;
    })();
  });

  let answers = $state<Record<string, QuizAnswer>>({});

  let me = $derived(GetUserOrRedirect());
</script>

<div>
  {#if loading}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:else if statusMessage}
    <div>{statusMessage}</div>
  {:else}
    <div class="grid grid-cols-12 grid-rows-12 w-screen h-screen p-16">
      <div
        class="row-start-3 row-end-11 col-start-1 col-end-10 p-4 rounded-2xl text-surface-50-950 bg-surface-900-100"
      >
        <QuizView
          UUID={current}
          {runnerKey}
          bind:chosen={answers[current]}
          onanswer={(a) => FetchAnswerPost(runnerKey, current, a)}
        />
      </div>
    </div>
  {/if}
</div>
