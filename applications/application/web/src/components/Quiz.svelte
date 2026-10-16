<script lang="ts">
  import DOMPurify from "isomorphic-dompurify";
  import { marked } from "marked";

  let { content }: { content: string } = $props();

  let htmlContent = $derived(
    DOMPurify.sanitize(marked.parse(content) as string),
  );

  $inspect(htmlContent);
</script>

<div>
  {@html htmlContent}
</div>
