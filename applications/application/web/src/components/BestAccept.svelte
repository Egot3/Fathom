<script lang="ts">
  import { USER_MIMES, type Accept } from "../lib/apiutils/acceptHeader";
  import AcceptHeaders from "./AcceptHeaders.svelte";

  let { best = $bindable() }: { best: Accept[] } = $props();
</script>

<details class="disclosure">
  <summary>Export options</summary>
  <div class="disclosure-content">
    <AcceptHeaders available={USER_MIMES} bind:accepted={best} />
  </div>
</details>
