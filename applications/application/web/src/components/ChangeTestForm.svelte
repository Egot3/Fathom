<script lang="ts">
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import UXInput from "./UXInput.svelte";
  import _ from "lodash";
  import {
    FetchTestBundle,
    FetchTestPatch,
    FetchTestPrune,
    type Test,
  } from "../lib/contracts/test";
  import QuizChips from "./QuizChips.svelte";
  import { SvelteSet } from "svelte/reactivity";
  import ExistingQuizChips from "./ExistingQuizChips.svelte";
  import { ok, safeTry } from "neverthrow";

  const {
    callback,
    UUID,
    name,
    response,
  }: {
    callback: () => void;

    response: Test;
    name: string;
    UUID: string;
  } = $props();

  let statusMessage = $state("");
  let nameMessage: string = $state("");

  let nameReady: boolean = $state(false);

  let newName = $derived(name);
  let quizzes = $derived(new SvelteSet(response.quizzes));
  let newQuizzes = $state(new SvelteSet<string>());

  async function ChangeTest(e: Event) {
    e.preventDefault();

    const result = await safeTry(async function* () {
      if (newName && newName !== name) {
        yield* FetchTestPatch(UUID, newName);
      }

      const quizzesRemoved = _.difference(
        response.quizzes,
        Array.from(quizzes),
      );
      if (quizzesRemoved.length > 0) {
        yield* FetchTestPrune(
          UUID,
          quizzesRemoved.map((v) => v.uuid),
        );
      }

      const newQuizzesArr = Array.from(newQuizzes);
      if (newQuizzesArr.length > 0) {
        yield* FetchTestBundle(UUID, newQuizzesArr);
      }

      return ok(undefined); // ok
    });

    result.match(
      () => {
        statusMessage = "";
        callback();
      },
      ({ error }) => {
        statusMessage = error;
      },
    );
  }

  const existing = new SvelteSet((() => quizzes)());
</script>

<form onsubmit={ChangeTest} class="w-full flex flex-col h-full space-y-2">
  <UXInput
    label="Name"
    placeholder="go-regenerics"
    bind:value={newName}
    bind:ready={nameReady}
    message={nameMessage}
    checker={(v: string) => {
      if (v.length < 3) {
        nameMessage = "name can't be less than 3 characters in length";
        return false;
      }
      if (v.length > 255) {
        nameMessage = "name can't be more than 255 characters in length";
        return false;
      }

      return true;
    }}
  />

  <div class="w-full flex">
    <div class="h-full w-1/2">
      <label class="label">
        <span class="label-text">Old quizzes</span>
        <ExistingQuizChips bind:quizzes />
      </label>
    </div>
    <div class="h-full w-1/2">
      <label class="label h-full">
        <span class="label-text">New quizzes</span>
        <div class="h-full">
          <QuizChips bind:chosen={newQuizzes} {existing} />
        </div>
      </label>
    </div>
  </div>

  <footer class="flex justify-end gap-2 w-full">
    {#if statusMessage !== ""}
      <span class="justify-self-start">{statusMessage}</span>
    {/if}
    <Dialog.CloseTrigger class="btn preset-tonal">Cancel</Dialog.CloseTrigger>
    <button type="submit" class="btn preset-filled">Change</button>
  </footer>
</form>
