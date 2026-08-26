<script lang="ts">
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import {
    FetchQuizPatch,
    Kind,
    type Meta,
    type QuizFile,
  } from "../lib/contracts/quiz";
  import UXInput from "./UXInput.svelte";
  import _ from "lodash";

  const {
    callback,
    UUID,
    name,
    response,
  }: {
    callback: () => void;

    response: QuizFile;
    name: string;
    UUID: string;
  } = $props();

  const pathRegex = /^([a-zA-Z0-9_\-]+)*\/?[a-zA-Z0-9_\-]+\.md$/;

  let statusMessage = $state("");
  let nameMessage: string = $state("");

  let nameReady: boolean = $derived(pathRegex.test(name));

  let newName = $derived(name);
  let body = $derived(response.body);
  let allOrNone: boolean = $derived(response.meta.all_or_none);
  let randomized: boolean = $derived(response.meta.randomized);
  let score: number = $derived(response.meta.score);
  let kind: Kind = $derived(response.meta.kind);

  async function ChangeQuiz(e: Event) {
    e.preventDefault();

    let newMeta: Meta | undefined = undefined;
    if (
      allOrNone !== response.meta.all_or_none ||
      randomized !== response.meta.randomized ||
      score !== response.meta.score ||
      kind !== response.meta.kind
    ) {
      newMeta = {
        all_or_none: allOrNone,
        score: score,
        randomized: randomized,
        kind: kind,
      };
    }

    let newBody: string | undefined = undefined;
    if (body !== response.body) {
      newBody = body;
    }

    statusMessage = await FetchQuizPatch(UUID, {
      name: newName,
      meta: newMeta,
      body: newBody,
    }).match(
      () => {
        return "";
      },
      (err) => {
        return err.error || "error has occured";
      },
    );

    callback();
  }

  $inspect(statusMessage);
</script>

<form onsubmit={ChangeQuiz} class="w-full flex flex-col h-full space-y-2">
  <UXInput
    label="Name/path"
    placeholder="color_theory/sky.md"
    bind:value={newName}
    bind:ready={nameReady}
    message={nameMessage}
    checker={(v: string) => {
      nameMessage =
        "name must have no special symbols and must contain .md extension";
      return pathRegex.test(v);
    }}
  />

  <label class="label">
    <span class="label-text">Quiz</span>
    <textarea class="textarea h-1/2" bind:value={body} placeholder="" required
    ></textarea>
  </label>

  <label class="label">
    <span class="label-text">Default score</span>
    <input
      class="input"
      type="number"
      bind:value={score}
      placeholder="1"
      required
    />
  </label>

  {#if kind !== Kind.Input}
    <label class="flex items-center space-x-2">
      <input class="checkbox" type="checkbox" bind:checked={randomized} />
      <p>Randomized</p>
    </label>
  {/if}

  <label class="flex items-center space-x-2">
    <input class="checkbox" type="checkbox" bind:checked={allOrNone} />
    <p>All or none</p>
  </label>

  <footer class="flex justify-end gap-2 w-full">
    {#if statusMessage !== ""}
      <span class="justify-self-start">{statusMessage}</span>
    {/if}
    <Dialog.CloseTrigger class="btn preset-tonal">Cancel</Dialog.CloseTrigger>
    <button type="submit" class="btn preset-filled">Change</button>
  </footer>
</form>
