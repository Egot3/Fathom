<script lang="ts">
  import { FetchAnswers } from "../lib/contracts/totals";
  import { IsJSONError } from "../lib/statuses/jsonerror";
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

  const response = $derived(
    FetchAnswers(groupUUID, userUUID, testUUID, page - 1, pageSize),
  );
  let chosenUUID: string = $state("");
  let answerValue: string = $state(null as never);
</script>

{#await response}
  <div>loading...</div>
{:then answers}
  {#if IsJSONError(answers)}
    <div>{answers.error}</div>
  {:else}
    <div class="grid grid-cols-12 h-24/25 w-full">
      <div class="col-start-1 col-end-9">
        <ParsedQuiz
          UUID={chosenUUID || answers.answers[0].quiz_uuid}
          answerValue={answerValue || answers.answers[0].chosen}
        />
      </div>
      <div
        class="col-start-9 col-end-13 bg-surface-600-400 rounded-[0.25rem] p-1"
      >
        {#each answers.answers as answer}
          <button
            class="btn preset-filled-error-50-950 w-full text-xl"
            onclick={() => {
              answerValue = answer.chosen;
              chosenUUID = answer.quiz_uuid;
            }}>{answer.quiz_name.slice(14).slice(0, -3)}</button
          >
        {/each}
      </div>
    </div>
  {/if}
{/await}
