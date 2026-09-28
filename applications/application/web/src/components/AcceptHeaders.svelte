<script lang="ts">
    import type { SvelteSet } from "svelte/reactivity";
    import type { Accept } from "../lib/apiutils/acceptHeader";

    let {
      available,
      accepted = $bindable()
    }: {
      accepted: SvelteSet<Accept>,
      available: string[]
    } = $props();


</script>

<fieldset class="fieldset space-y-2">
    {#each available as tag }
        <div class="flex space-x-2">
            <label class="flex items-center space-x-2">
                {let checked = false}
                <input class="checkbox" bind:checked={checked} onclick={()=>{
                  if(checked) {
                    accepted.delete(tag)
                  }
                }} type="checkbox" />
                <p>Option 1</p>
                <input class="input" type="range" value="0" max="1" />
            </label>
        </div>
    {/each}
</fieldset>
