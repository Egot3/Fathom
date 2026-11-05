<script lang="ts">
  import { QrCode } from "@skeletonlabs/skeleton-svelte";
    import { BaseURL } from "../lib/config";

  const { key }: { key: string } = $props();
</script>

<QrCode value={`${BaseURL}/${key}`}>
  <QrCode.Frame class="size-full max-size-36">
    <QrCode.Pattern />
  </QrCode.Frame>
</QrCode>
