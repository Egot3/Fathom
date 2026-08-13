<script lang="ts">
  import { type Answer } from "../lib/contracts/totals";

  const {
    content,
  }: {
    content: Answer;
  } = $props();

  let height = $state(0);

  let page = $state(1);
  let pageSize = $state(5);
</script>

<span>peek!</span>
