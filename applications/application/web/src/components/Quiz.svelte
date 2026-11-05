<script module lang="ts">
  import "katex/dist/katex.min.css";
  import { marked } from "marked";
  import markedKatex from "marked-katex-extension";

  marked.use(
    markedKatex({
      throwOnError: false,
    }),
  );
</script>

<script lang="ts">
  import DOMPurify from "isomorphic-dompurify";

  let { content }: { content: string } = $props();

  let htmlContent = $derived(
    DOMPurify.sanitize(marked.parse(content) as string, {
      USE_PROFILES: { html: true, mathMl: true },
      ADD_TAGS: ["semantics", "annotation"],
    }),
  );
</script>

<div>
  {@html htmlContent}
</div>
