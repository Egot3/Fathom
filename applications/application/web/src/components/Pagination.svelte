<script lang="ts">
  import { ArrowLeftIcon, ArrowRightIcon } from "@lucide/svelte";
  import { Pagination } from "@skeletonlabs/skeleton-svelte";

  let {
    pageSize,
    page = $bindable(0),
    total,
  }: {
    pageSize: number;
    page: number;
    total: number;
  } = $props();
</script>

<Pagination
  count={total}
  {pageSize}
  {page}
  onPageChange={(event) => (page = event.page)}
  class="rounded-xl"
>
  <Pagination.PrevTrigger>
    <ArrowLeftIcon class="size-4" />
  </Pagination.PrevTrigger>
  <Pagination.Context>
    {#snippet children(pagination)}
      {#each pagination().pages as page, index (page)}
        {#if page.type === "page"}
          <Pagination.Item {...page}>
            {page.value}
          </Pagination.Item>
        {:else}
          <Pagination.Ellipsis {index}>…</Pagination.Ellipsis>
        {/if}
      {/each}
    {/snippet}
  </Pagination.Context>
  <Pagination.NextTrigger>
    <ArrowRightIcon class="size-4" />
  </Pagination.NextTrigger>
</Pagination>
