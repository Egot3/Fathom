<script lang="ts">
  import { QrCode } from "@skeletonlabs/skeleton-svelte";

  const { key }: { key: string } = $props();
</script>

<QrCode value={`https://${import.meta.env.PUBLIC_BASE_URL}/${key}`}>
  <QrCode.Frame class="size-full max-size-36">
    <QrCode.Pattern />
  </QrCode.Frame>
</QrCode>
