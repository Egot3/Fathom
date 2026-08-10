<script lang="ts">
  import { SvelteSet } from "svelte/reactivity";
  import { CheckIcon } from "@lucide/svelte";
  import _ from "lodash";
  import type { Group } from "../lib/contracts/group";

  let {
    groups = $bindable(new SvelteSet<Group>()),
    disabled = false,
  }: {
    groups: SvelteSet<Group>;
    disabled?: boolean;
  } = $props();

  const startinggroups = new Set(groups);
</script>

<div class="grid gap-4 m-2 w-full place-items-center h-full overflow-auto">
  {#each startinggroups as group}
    <button
      {disabled}
      type="button"
      class={`chip capitalize preset-outlined-surface-400-600 ${groups.has(group) ? "preset-tonal-primary" : ""}`}
      onclick={() => {
        if (groups.has(group)) {
          groups.delete(group);
        } else {
          groups.add(group);
        }
      }}
    >
      {#if groups.has(group)}<CheckIcon size={14} />{/if}
      <span>{group.name}</span>
    </button>
  {/each}
</div>
