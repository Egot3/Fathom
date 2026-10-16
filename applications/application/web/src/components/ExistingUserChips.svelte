<script lang="ts">
  import { SvelteSet } from "svelte/reactivity";
  import { CheckIcon } from "@lucide/svelte";
  import _ from "lodash";
  import type { User } from "../lib/contracts/user";

  let {
    pupils = $bindable(new SvelteSet<User>()),
    disabled = false,
  }: {
    pupils: SvelteSet<User>;
    disabled?: boolean;
  } = $props();

  const startinggroups = new Set(pupils);
</script>

<div class="grid gap-4 m-2 w-full place-items-center h-full overflow-auto">
  {#each startinggroups as pupil}
    <button
      {disabled}
      type="button"
      class={`chip capitalize preset-outlined-surface-400-600 ${pupils.has(pupil) ? "preset-tonal-primary" : ""}`}
      onclick={() => {
        if (pupils.has(pupil)) {
          pupils.delete(pupil);
        } else {
          pupils.add(pupil);
        }
      }}
    >
      {#if pupils.has(pupil)}<CheckIcon size={14} />{/if}
      <span>{pupil.nickname}</span>
    </button>
  {/each}
</div>
