<script lang="ts">
  import { QrCode } from "@skeletonlabs/skeleton-svelte";

  const { key }: { key: string } = $props();

  const domain = import.meta.env.PUBLIC_BASE_URL;

  // ah, yes a message in qr-code
  const url = $derived(
    domain
      ? `https://${domain}/${key}`
      : "set PUBLIC_BASE_URL in env to use this feature",
  );
</script>

<QrCode value={url}>
  <QrCode.Frame class="size-full max-size-36">
    <QrCode.Pattern />
  </QrCode.Frame>
</QrCode>
