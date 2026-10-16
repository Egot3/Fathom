<script lang="ts">
  import { SvelteSet } from "svelte/reactivity";
  import type { Group } from "../lib/contracts/group";
  import ExistingUserChips from "./ExistingUserChips.svelte";

  let { content }: { content: Group } = $props();
</script>

<label class="label">
  <span class="label-text">Name</span>
  <input
    class="input"
    type="text"
    value={content.name}
    placeholder="go-generics"
    disabled
  />
</label>

{const pupils = new SvelteSet(content.pupils)}
<label class="label">
  <span class="label-text">Pupils</span>
  <ExistingUserChips disabled {pupils} />
</label>
