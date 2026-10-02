<script lang="ts">
  import { SvelteMap } from "svelte/reactivity";
  import type { Accept, USER_MIME } from "../lib/apiutils/acceptHeader";
  import AcceptOption from "./AcceptOption.svelte";

  let {
    available,
    accepted = $bindable(),
  }: {
    accepted: Accept[];
    available: readonly USER_MIME[];
  } = $props();

  const acceptedMap: SvelteMap<string, number> = $state(
    new SvelteMap<string, number>(),
  );

  $effect(() => {
    acceptedMap.values(); // trigger

    accepted = [...acceptedMap.entries()].map((v): Accept => {
      console.log("map entry:", v);
      return {
        accept: v[0],
        quality: v[1],
      };
    });
  });
</script>

<fieldset class="fieldset space-y-2">
  <legend>Compression type & Want level</legend>
  {#each available as tag}
    <AcceptOption {acceptedMap} key={tag.mime} tag={tag.name} />
  {/each}
</fieldset>
