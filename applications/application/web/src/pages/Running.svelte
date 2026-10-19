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
          ({ error }) => (statusMessage = error),
        );
      loading = false;
    })();
  });

  let answers = $state<Record<string, QuizAnswer>>({});

  let me = $derived(GetUserOrRedirect());
  let groupUUID = ""; //placeboo
</script>

<div>
  {#if loading}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:else if statusMessage}
    <div>{statusMessage}</div>
  {:else}
    <QuizView
      UUID={current}
      {runnerKey}
      bind:chosen={answers[current]}
      onanswer={(a) =>
        FetchAnswerPost(groupUUID, me.UUID, runnerKey, current, a)}
    />
  {/if}
</div>
