<script lang="ts">
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import UXInput from "./UXInput.svelte";
  import _ from "lodash";
  import { SvelteSet } from "svelte/reactivity";
  import {
    FetchGroupAppend,
    FetchGroupPatch,
    FetchGroupPrune,
    type Group,
  } from "../lib/contracts/group";
  import ExistingGroupChips from "./ExistingGroupChips.svelte";
  import ExistingUserChips from "./ExistingUserChips.svelte";
  import UserChips from "./UserChips.svelte";

  const {
    callback,
    UUID,
    name,
    response,
  }: {
    callback: () => void;

    response: Group;
    name: string;
    UUID: string;
  } = $props();

  let statusMessage = $state("");
  let nameMessage: string = $state("");

  let nameReady: boolean = $state(false);

  let newName = $derived(name);
  let pupils = $derived(new SvelteSet(response.pupils));
  let newPupils = $state(new SvelteSet<string>());

  async function ChangeGroup(e: Event) {
    e.preventDefault();

    if (newName !== name && newName) {
      const patchResponse = await FetchGroupPatch(UUID, newName);
      if (patchResponse !== null) {
        statusMessage = patchResponse.error;
        return;
      }
    }

    const pupilsArr = Array.from(pupils);
    const pupilsRemoved = _.difference(response.pupils, pupilsArr);
    if (pupilsRemoved.length > 0) {
      const deleteResponse = await FetchGroupPrune(
        UUID,
        pupilsRemoved.map((v) => v.uuid),
      );
      if (deleteResponse !== null) {
        statusMessage = deleteResponse.error;
        return;
      }
    }

    const newPupilsArr = Array.from(newPupils);
    if (newPupilsArr.length > 0) {
      const postResponse = await FetchGroupAppend(UUID, newPupilsArr);
      if (postResponse !== null) {
        statusMessage = postResponse.error;
        return;
      }
    }

    callback();
  }

  const existing = new SvelteSet((() => pupils)());
</script>

<form onsubmit={ChangeGroup} class="w-full flex flex-col h-full space-y-2">
  <UXInput
    label="Name"
    placeholder="ETS' Evil Corp fan club"
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
        <span class="label-text">Existing pupils</span>
        <ExistingUserChips bind:pupils />
      </label>
    </div>
    <div class="h-full w-1/2">
      <label class="label">
        <span class="label-text">New pupils</span>
        <UserChips bind:chosen={newPupils} {existing} />
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
