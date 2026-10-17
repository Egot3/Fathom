<script lang="ts">
  import CreateDialog from "./CreateDialog.svelte";
  import SmallBookCreateDialog from "./SmallBookCreateDialog.svelte";
  import DeleteDialogSquare from "./DeleteDialogSquare.svelte";
  import PeekDialogSquare from "./PeekDialogSquare.svelte";
  import { FetchGroups, type Group } from "../lib/contracts/group";
  import CreateGroupForm from "./CreateGroupForm.svelte";
  import DeleteGroupForm from "./DeleteGroupForm.svelte";
  import ChangeDialogSqare from "./ChangeDialogSqare.svelte";
  import ChangeGroupForm from "./ChangeGroupForm.svelte";
  import GroupPeek from "./GroupPeek.svelte";
  import { CreateFitRows } from "../lib/layoututils/fitRows.svelte";
  import Pagination from "./Pagination.svelte";
  import { CreatePagedQuery } from "../lib/layoututils/createPagedQuery.svelte";

  const fit = CreateFitRows({ rowHeight: 20, reserve: 1 });

  const q = CreatePagedQuery({
    pageSize: () => fit.pageSize,
    enabled: () => fit.ready,
    params: () => ({}),
    fetcher: FetchGroups,
  });

  let focused = $state("");
  let clickFocused = $state("");
</script>

<div
  class="grid gap-4 w-full place-items-center h-full overflow-auto"
  use:fit.measure
>
  {#if q.loading}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:else if q.error}
    <div>{q.error}</div>
  {:else if !q.data || q.data.total === 0}
    <CreateDialog name="Create your first group" title="Create your first group"
      ><CreateGroupForm callback={q.refresh} /></CreateDialog
    >
  {:else}
    <table class="table table-auto self-start">
      <thead>
        <tr class="text-surface-100-900 flex">
          <th class="w-1/3">Group</th>
          <th class="w-1/3">User count</th>
          <th class="w-1/3"></th>
        </tr>
      </thead>

      <tbody>
        {#each q.data.groups as group}
          <tr
            onmouseenter={() => (focused = group.uuid)}
            onmouseleave={() => (focused = "")}
            class="bg-surface-700-300 rounded-xl flex hover:motion-safe:hover:brightness-125 dark:hover:motion-safe:hover:brightness-75"
          >
            <td class="w-1/3">{group.name}</td>
            <td class="w-1/3">{group.pupils?.length ?? 0}</td>
            <td class="w-1/3">
              {#if group.uuid === focused || group.uuid === clickFocused}
                <div class="flex flex-row-reverse">
                  <DeleteDialogSquare
                    title="Group deleter"
                    callback={() => (clickFocused = group.uuid)}
                  >
                    <DeleteGroupForm
                      name={group.name}
                      UUID={group.uuid}
                      callback={q.refresh}
                    />
                  </DeleteDialogSquare>
                  <ChangeDialogSqare
                    callback={() => (clickFocused = group.uuid)}
                    title="Group changer"
                    contentGetter={() => group}
                  >
                    {#snippet children({ content }: { content: Group })}
                      <ChangeGroupForm
                        UUID={group.uuid}
                        name={group.name}
                        callback={q.refresh}
                        response={content}
                      />
                    {/snippet}</ChangeDialogSqare
                  >

                  <PeekDialogSquare
                    callback={() => (clickFocused = group.uuid)}
                    title="Group peeker"
                    contentGetter={() => group}
                  >
                    {#snippet children({ content }: { content: Group })}
                      <GroupPeek {content} />
                    {/snippet}
                  </PeekDialogSquare>
                </div>
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>

    <div class="flex justify-between items-center gap-4 w-full self-end">
      <Pagination
        bind:page={q.page}
        pageSize={fit.pageSize}
        total={q.data?.total ?? 0}
      />

      <SmallBookCreateDialog title="Group maker"
        ><CreateGroupForm callback={q.refresh} /></SmallBookCreateDialog
      >
    </div>
  {/if}
</div>
