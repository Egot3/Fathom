<script lang="ts">
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import Link from "./Link.svelte";
  import PeekDialogue from "./PeekDialogue.svelte";

  const { key }: { key: string } = $props();
</script>

<PeekDialogue content={(async () => key)()} title="QR sharer">
  {#snippet trigger()}
    <Dialog.Trigger class="btn preset-outlined-warning-50-950"
      >Get connection QR</Dialog.Trigger
    >
  {/snippet}

  <Link {key}></Link>
</PeekDialogue>
