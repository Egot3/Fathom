<script lang="ts">
  import {
    Kind,
    type QuizAnswer,
    type QuizOptions,
  } from "../lib/contracts/quiz";

  const {
    kind,
    options,
    answer,
  }: {
    kind: Kind;
    options?: QuizOptions;
    answer?: QuizAnswer;
  } = $props();
</script>

{#if kind === Kind.Input}
  <div>input</div>
{:else if kind === Kind.Order}
  <div>order</div>
{:else if kind === Kind.Radio}
  <div>radio</div>
{:else if kind === Kind.Check}
  <div>check</div>
{:else if kind === Kind.Accordance}
  <div>accordance</div>
{:else}
  <div>Unknown quiz kind!</div>
{/if}
