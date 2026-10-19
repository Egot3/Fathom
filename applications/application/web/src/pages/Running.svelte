<script lang="ts">
  const { runnerKey }: { runnerKey: string } = $props();

  const quizzes = $derived(1);
</script>

<div>{runnerKey}</div>
