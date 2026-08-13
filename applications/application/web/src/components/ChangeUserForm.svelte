<script lang="ts">
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import UXInput from "./UXInput.svelte";
  import _ from "lodash";
  import { SvelteSet } from "svelte/reactivity";
  import {
    FetchUserPatch,
    type PatchUserStruct,
    type User,
  } from "../lib/contracts/user";
  import { GetUser } from "../lib/bgdata/user.svelte";
  import { CheckPassword } from "../lib/passwordUtils/checkPassword";

  const {
    callback,
    UUID,
    response,
  }: {
    callback: () => void;

    response: User;
    UUID: string;
  } = $props();

  let statusMessage = $state("");
  let nameMessage: string = $state("");
  let passwordMessage: string = $state("");

  let nameReady: boolean = $state(false);
  let passwordReady: boolean = $state(false);

  let newName = $derived(response.nickname);
  let password = $state("");
  let isTeacher = $derived(response.is_teacher);

  async function ChangeUser(e: Event) {
    e.preventDefault();

    const patched: PatchUserStruct = {};
    const requestor = GetUser();
    if (requestor == null) {
      statusMessage = "You are logged out!";
      return;
    }
    if (requestor.UUID !== UUID && !requestor.isTeacher) {
      statusMessage = "You can't change account, which you don't own";
      return;
    }

    if (isTeacher !== response.is_teacher) {
      if (!requestor.isTeacher) {
        statusMessage = "You are not a teacher to change teachery";
        return;
      }
      if (requestor.UUID === UUID) {
        statusMessage = "You can't change your own teachery!"; // so there won't be a situation when only techer de-teachers themselves
        return;
      }

      patched.IsTeacher = isTeacher;
    }

    if (newName !== response.nickname && newName) {
      patched.Name = newName;
    }

    if (password !== "") {
      patched.Password = password;
    }

    if (!_.isEmpty(patched)) {
      const patchResponse = await FetchUserPatch(UUID, patched);
      if (patchResponse !== null) {
        statusMessage = patchResponse.error;
        return;
      }
    }

    callback();
  }
</script>

<form onsubmit={ChangeUser} class="w-full flex flex-col h-full space-y-2">
  <UXInput
    label="Name"
    placeholder="knownUser"
    bind:value={newName}
    bind:ready={nameReady}
    bind:message={nameMessage}
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

  <UXInput
    label="Password"
    type="password"
    placeholder="unbelievab1y_H4RD_TO_brut3_ForcePassw0rd!"
    bind:value={password}
    bind:ready={passwordReady}
    bind:message={passwordMessage}
    checker={(v: string) => {
      passwordMessage = CheckPassword(v);
      return passwordMessage === "";
    }}
  />

  <label class="flex items-center space-x-2 w-fit">
    <input class="checkbox" type="checkbox" bind:checked={isTeacher} />
    <p>Is teacher</p>
  </label>

  <footer class="flex justify-end gap-2 w-full">
    {#if statusMessage !== ""}
      <span class="justify-self-start">{statusMessage}</span>
    {/if}
    <Dialog.CloseTrigger class="btn preset-tonal">Cancel</Dialog.CloseTrigger>
    <button type="submit" class="btn preset-filled">Change</button>
  </footer>
</form>
