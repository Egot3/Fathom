<script lang="ts">
  import {
    AnswerInput,
    Kind,
    OptionAccordance,
    OptionCheck,
    OptionOrder,
    OptionRadio,
    type QuizAnswer,
    type QuizOptions,
  } from "../lib/contracts/quiz";

  const {
    kind,
    options,
    answers,
    disabled = false,
  }: {
    kind: Kind;
    options: QuizOptions;
    answers?: QuizAnswer;
    disabled: boolean;
  } = $props();

  $inspect(options);
</script>

{#if kind === Kind.Input}
  <input {disabled} value={answers ?? ""} />
{:else if kind === Kind.Order}
  <div class="space-y-2">
    {#each OptionOrder(options) as option}
      <label class="flex items-center space-x-2">
        <input name="radio" class="radio" type="radio" {disabled} />
        <p>{option}</p>
        <!-- for now -->
      </label>
    {/each}
  </div>
{:else if kind === Kind.Radio}
  <div class="space-y-2">
    {#each OptionRadio(options) as option}
      <label class="flex items-center space-x-2">
        <input name="radio" class="radio" type="radio" {disabled} />
        <p>{option.label}</p>
      </label>
    {/each}
  </div>
{:else if kind === Kind.Check}
  <div class="space-y-2">
    {#each OptionCheck(options) as option}
      <label class="flex items-center space-x-2">
        <input class="checkbox" type="checkbox" {disabled} />
        <p>{option.label}</p>
      </label>
    {/each}
  </div>
{:else if kind === Kind.Accordance}
  <div class="space-y-2 flex">
    {let opt = OptionAccordance(options)}
    <div>
      {#each opt.static as st}
        <label class="flex items-center space-x-2">
          <input class="checkbox" type="checkbox" {disabled} />
          <p>{st}</p>
        </label>
      {/each}
    </div>
    <div>
      {#each opt.dynamic as dy}
        <label class="flex items-center space-x-2">
          <input class="checkbox" type="checkbox" {disabled} />
          <p>{dy}</p>
        </label>
      {/each}
    </div>
  </div>
{:else}
  <div>Unknown quiz kind!</div>
{/if}
