<script lang="ts">
  import ArrowLeftIcon from "@lucide/svelte/icons/arrow-left";
  import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
  import { Pagination } from "@skeletonlabs/skeleton-svelte";
  import { IsJSONError, type JSONError } from "../lib/statuses/jsonerror";

  import CreateDialog from "./CreateDialog.svelte";
  import SmallBookCreateDialog from "./SmallBookCreateDialog.svelte";
  import DeleteDialogSquare from "./DeleteDialogSquare.svelte";
  import PeekDialogSquare from "./PeekDialogSquare.svelte";
  import {
    FetchGroups,
    type Group,
    type GroupsOrJSONError,
  } from "../lib/contracts/group";
  import CreateGroupForm from "./CreateGroupForm.svelte";
  import DeleteGroupForm from "./DeleteGroupForm.svelte";
  import ChangeDialogSqare from "./ChangeDialogSqare.svelte";
  import ChangeGroupForm from "./ChangeGroupForm.svelte";
  import GroupPeek from "./GroupPeek.svelte";

  let height = $state(0);

  let page = $state(1);
  let pageSize = $derived(Math.trunc((height - 39 - 45) / 39));

  let trigger = $state(0);
  let time: number;
  const paginatedGroupPromises = $derived.by(() => {
    const p = page;
    const ps = pageSize;
    trigger;

    console.log("detected change");
    clearTimeout(time);

    return new Promise<GroupsOrJSONError>((resolve) => {
      time = setTimeout(() => {
        resolve(FetchGroups(p - 1, ps));
      }, 500);
    });
  });

  let focused = $state("");
  let clickFocused = $state("");
</script>

<div
  class="grid gap-4 w-full place-items-center h-full overflow-auto"
  bind:clientHeight={height}
>
  {#await paginatedGroupPromises}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:then paginatedGroups}
    {#if IsJSONError(paginatedGroups)}
      <div>{paginatedGroups.error}</div>
    {:else}
      {#if paginatedGroups.total == 0}
        <CreateDialog
          name="Create your first group"
          title="Create your first group"
          ><CreateGroupForm
            callback={() => {
              trigger++;
            }}
          /></CreateDialog
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
            {#each paginatedGroups.groups as group}
              <tr
                onmouseenter={() => (focused = group.uuid)}
                onmouseleave={() => (focused = "")}
                class="bg-surface-700-300 rounded-xl flex hover:motion-safe:hover:brightness-125 dark:hover:motion-safe:hover:brightness-75"
              >
                <td class="w-1/3">{group.name}</td>
                <!-- 10.5 rem = 168px -->
                <td class="w-1/3">{group.pupils?.length ?? 0}</td>
                <td class="w-1/3">
                  {#if group.uuid === focused || group.uuid === clickFocused}
                    <div class="flex flex-row-reverse">
                      <DeleteDialogSquare
                        title="Quiz deleter"
                        callback={() => (clickFocused = group.uuid)}
                      >
                        <DeleteGroupForm
                          name={group.name}
                          UUID={group.uuid}
                          callback={() => {
                            trigger++;
                          }}
                        />
                      </DeleteDialogSquare>
                      <ChangeDialogSqare
                        callback={() => (clickFocused = group.uuid)}
                        title="Quiz changer"
                        contentGetter={async () => {
                          return group;
                        }}
                      >
                        {#snippet children({ content }: { content: Group })}
                          <ChangeGroupForm
                            UUID={group.uuid}
                            name={group.name}
                            callback={() => {
                              trigger++;
                            }}
                            response={content}
                          />
                        {/snippet}</ChangeDialogSqare
                      >

                      <PeekDialogSquare
                        callback={() => (clickFocused = group.uuid)}
                        title="Test peeker"
                        contentGetter={async () => {
                          return group;
                        }}
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
            count={paginatedGroups?.total ?? 0}
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

          <SmallBookCreateDialog title="Test maker"
            ><CreateGroupForm
              callback={() => {
                trigger++;
              }}
            /></SmallBookCreateDialog
          >
        </div>
      {/if}
    {/if}
  {/await}
</div>
