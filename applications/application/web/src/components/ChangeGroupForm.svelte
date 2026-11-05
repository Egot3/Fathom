<script lang="ts">
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import UXInput from "./UXInput.svelte";
  import _ from "lodash";
  import { SvelteSet } from "svelte/reactivity";
  import {
    FetchGroupAppend,
    FetchGroupPatch,
    FetchGroupPrune,
    FetchGroupUsers,
    type Group,
  } from "../lib/contracts/group";
  import ExistingUserChips from "./ExistingUserChips.svelte";
  import UserChips from "./UserChips.svelte";
  import { ok, safeTry } from "neverthrow";
  import { CreateFitRows } from "../lib/layoututils/fitRows.svelte";
  import { CreatePagedQuery } from "../lib/layoututils/createPagedQuery.svelte";
  import Pagination from "./Pagination.svelte";

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

  const q = CreatePagedQuery({
    pageSize: () => 5,
    enabled: () => true,
    params: () => UUID,
    fetcher: FetchGroupUsers,
  });

  let newName = $derived(name);
  let removed = $state(new SvelteSet<string>());
  let newPupils = $state(new SvelteSet<string>());

  async function ChangeGroup(e: Event) {
    e.preventDefault();

    const result = await safeTry(async function* () {
      if (newName !== name && newName) {
        yield* FetchGroupPatch(UUID, newName);
      }

      if (removed.size > 0) {
        yield* FetchGroupPrune(UUID, Array.from(removed));
      }

      if (newPupils.size > 0) {
        yield* FetchGroupAppend(UUID, Array.from(newPupils));
      }

      return ok(undefined);
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

  const existing = $derived(new SvelteSet(q.data?.users ?? []));
</script>

<form onsubmit={ChangeGroup} class="w-full flex flex-col h-full space-y-2">
  <UXInput
    label="Name"
    placeholder="ETS' Evil Corp fan club"
    bind:value={newName}
    bind:ready={nameReady}
    message={nameMessage}
    checker={(v: string) => {
      if (v.length === 0 || v === name) {
        nameMessage = "";
        return false;
      }
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

  <div class="w-full flex min-h-19">
    <div class="h-full w-1/2 min-h-19">
      <label class="label min-h-19">
        <span class="label-text">Existing pupils</span>
        <ExistingUserChips users={q.data?.users ?? []} bind:removed />
      </label>
      <Pagination bind:page={q.page} pageSize={5} total={q.data?.total ?? 0} />
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
