<script lang="ts">
  import { FetchAnswers } from "../lib/contracts/totals";

  const {
    userUUID,
    groupUUID,
    testUUID,
  }: {
    userUUID: string;
    groupUUID: string;
    testUUID: string;
  } = $props();

  let height = $state(0);

  let page = $state(1);
  let pageSize = $state(5);

  const response = $derived(
    FetchAnswers(groupUUID, userUUID, testUUID, page, pageSize),
  );
</script>

{#await response}
  <div>loading...</div>
{:then answers}
  <div class="grid grid-cols-12">
    <div class="col-start-1 col-end-8 bg-amber-800"></div>
    <div class="col-start-8 col-end-13 border-green-950"></div>
  </div>
{/await}
