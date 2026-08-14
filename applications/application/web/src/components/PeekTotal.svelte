<script lang="ts">
  import { Kind } from "../lib/contracts/quiz";
  import { FetchAnswers } from "../lib/contracts/totals";
  import { IsJSONError } from "../lib/statuses/jsonerror";
  import AnsweredQuiz from "./AnsweredQuiz.svelte";

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
    FetchAnswers(groupUUID, userUUID, testUUID, page, pageSize),
  );
</script>

{#await response}
  <div>loading...</div>
{:then answers}
  {#if IsJSONError(answers)}
    <div>{answers.error}</div>
  {:else}
    <div class="grid grid-cols-12 h-24/25 w-full">
      <div class="col-start-1 col-end-9 bg-amber-800">
        <AnsweredQuiz kind={Kind.Input} />
      </div>
      <div class="col-start-9 col-end-13 bg-green-950"></div>
    </div>
  {/if}
{/await}
