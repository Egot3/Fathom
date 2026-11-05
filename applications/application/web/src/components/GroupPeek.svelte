<script lang="ts">
  import { SvelteSet } from "svelte/reactivity";
  import { FetchGroupUsers, type Group } from "../lib/contracts/group";
  import ExistingUserChips from "./ExistingUserChips.svelte";
  import { onMount } from "svelte";
  import { CreateFitRows } from "../lib/layoututils/fitRows.svelte";
  import type { User } from "../lib/contracts/user";
  import { CreatePagedQuery } from "../lib/layoututils/createPagedQuery.svelte";
  import Pagination from "./Pagination.svelte";

  let { content }: { content: Group } = $props();

  const fit = CreateFitRows({
    rowHeight: 24,
    reserve: 1,
  });
  const q = CreatePagedQuery({
    pageSize: () => fit.pageSize,
    enabled: () => fit.ready,
    params: () => content.uuid,
    fetcher: FetchGroupUsers,
  });
</script>

<div>
  <label class="label">
    <span class="label-text">Name</span>
    <input
      class="input"
      type="text"
      value={content.name}
      placeholder="go-generics"
      disabled
    />
  </label>

  <label class="label min-h-19" use:fit.measure>
    <span class="label-text">Pupils</span>
    <div>
      {#if q.data}
        <ExistingUserChips disabled users={q.data.users} />
      {/if}
    </div>
  </label>

  <Pagination
    bind:page={q.page}
    pageSize={fit.pageSize}
    total={q.data?.total ?? 0}
  ></Pagination>
</div>
