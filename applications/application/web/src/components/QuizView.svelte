<script module lang="ts">
  import "katex/dist/katex.min.css";
  import { marked } from "marked";
  import markedKatex from "marked-katex-extension";

  marked.use(
    markedKatex({
      throwOnError: false,
    }),
  );
</script>

<script lang="ts">
  import { renderToString } from "katex";
  import {
    FetchParsedQuiz,
    FetchRunningQuiz,
    type ParsedQuiz,
    type QuizAnswer,
  } from "../lib/contracts/quiz";
  import QuizAnswerer from "./QuizAnswerer.svelte";
  import DOMPurify from "isomorphic-dompurify";

  let {
    UUID,
    runnerKey,
    chosen = $bindable(),
    correct,
    onanswer,
  }: {
    UUID: string;
    runnerKey?: string;
    chosen?: QuizAnswer;
    correct?: QuizAnswer;
    onanswer?: (a: QuizAnswer) => void;
  } = $props();

  const taking = $derived(runnerKey !== undefined);

  let quiz: ParsedQuiz | null = $state(null);
  let loading = $state(true);
  let statusMessage = $state("");
  let ticket = 0;

  $effect(() => {
    const id = UUID;
    const key = runnerKey;
    const mine = ++ticket;
    loading = true;

    const t = setTimeout(
      async () => {
        const msg = await (
          key !== undefined ? FetchRunningQuiz(key, id) : FetchParsedQuiz(id)
        ).match(
          (r) => {
            if (mine === ticket) quiz = r;
            return "";
          },
          (e) => e.error,
        );
        if (mine !== ticket) return;
        statusMessage = msg;
        loading = false;
      },
      key !== undefined ? 0 : 300,
    );
    return () => clearTimeout(t);
  });
</script>

{#if loading}
  <div class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"></div>
{:else if statusMessage !== ""}
  <span>{statusMessage}</span>
{:else if quiz}
  <div class="space-y-1 h-full flex flex-col">
    <h2 class="text-5xl font-bold">{quiz.title}</h2>
    <p class="text-2xl">
      {@html DOMPurify.sanitize(marked.parse(quiz.body) as string, {
        USE_PROFILES: { html: true, mathMl: true },
        ADD_TAGS: ["semantics", "annotation"],
      })}
    </p>
    <QuizAnswerer
      kind={quiz.meta.kind}
      options={quiz.options}
      bind:value={chosen}
      correct={taking ? undefined : (correct ?? quiz.answers)}
      readonly={!taking}
      {onanswer}
    />
  </div>
{/if}
