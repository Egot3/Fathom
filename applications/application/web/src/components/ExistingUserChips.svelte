<script lang="ts">
  import { SvelteSet } from "svelte/reactivity";
  import { CheckIcon } from "@lucide/svelte";
  import type { User } from "../lib/contracts/user";

  let {
    users,
    removed = $bindable(new SvelteSet<string>()),
    disabled = false,
  }: {
    users: User[];
    removed?: SvelteSet<string>;
    disabled?: boolean;
  } = $props();
</script>

<div class="grid gap-4 m-2 w-full place-items-center h-full overflow-auto">
  {#each users as user (user.uuid)}
    {const kept = !removed.has(user.uuid)}
    <button
      {disabled}
      type="button"
      class={`chip capitalize preset-outlined-surface-400-600 ${kept ? "preset-tonal-primary" : "line-through opacity-60"}`}
      onclick={() => {
        if (kept) {
          removed.add(user.uuid);
        } else {
          removed.delete(user.uuid);
        }
      }}
    >
      {#if kept}<CheckIcon size={14} />{/if}
      <span>{user.nickname}</span>
    </button>
  {/each}
</div>
