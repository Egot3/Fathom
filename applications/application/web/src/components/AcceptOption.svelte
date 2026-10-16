<script lang="ts">
  import type { SvelteMap } from "svelte/reactivity";

  let {
    tag,
    acceptedMap,
    key,
  }: {
    tag: string;
    key: string;
    acceptedMap: SvelteMap<string, number>;
  } = $props();

  let checked = $state(false);
  let best = $state<number>(0);
</script>

<div class="flex space-x-2 space-y-2">
  <label class="grid items-center grid-cols-2">
    <div class="flex items-center space-x-2">
      <input
        class="checkbox"
        bind:checked
        onclick={() => {
          if (!checked) {
            acceptedMap.delete(key);
          }
          acceptedMap.set(key, best);
        }}
        type="checkbox"
      />
      <p>{tag}</p>
    </div>

    <input
      class="input"
      type="range"
      max="1"
      step="0.1"
      bind:value={best}
      onchange={() => acceptedMap.set(key, best)}
      disabled={!checked}
    />
  </label>
</div>
