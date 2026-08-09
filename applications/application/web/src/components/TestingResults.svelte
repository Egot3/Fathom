<script lang="ts">
  import ArrowLeftIcon from "@lucide/svelte/icons/arrow-left";
  import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
  import { Pagination } from "@skeletonlabs/skeleton-svelte";
  import { IsJSONError } from "../lib/statuses/jsonerror";
  import CreateDialog from "./CreateDialog.svelte";
  import SmallPaperCreateDialog from "./SmallPaperCreateDialog.svelte";
  import ChangeQuizForm from "./ChangeQuizForm.svelte";
  import ChangeDialogSqare from "./ChangeDialogSqare.svelte";
  import PeekDialogSquare from "./PeekDialogSquare.svelte";
  import DeleteDialogSquare from "./DeleteDialogSquare.svelte";
  import {
    FetchUsers,
    type User,
    type UsersOrJSONError,
  } from "../lib/contracts/user";
  import DeleteUserForm from "./DeleteUserForm.svelte";
  import ChangeUserForm from "./ChangeUserForm.svelte";
  import UserPeek from "./UserPeek.svelte";

  let height = $state(0);

  let page = $state(1);
  let pageSize = $derived(Math.trunc((height - 29 - 45 - 40) / 49));

  let focused = $state("");
  let clickFocused = $state("");

  let trigger = $state(0);
  let time: number;
  const paginatedQuizPromises = $derived.by(() => {
    const p = page;
    const ps = pageSize;
    trigger;

    clearTimeout(time);

    return new Promise<UsersOrJSONError>((resolve) => {
      time = setTimeout(() => {
        resolve(FetchUsers(p - 1, ps));
      }, 500);
    });
  });
</script>

<div
  class="grid gap-4 w-full place-items-center h-full overflow-auto"
  bind:clientHeight={height}
>
  {#await paginatedQuizPromises}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:then paginatedUsers}
    {#if IsJSONError(paginatedUsers)}
      <div>{paginatedUsers.error}</div>
    {:else}
      {#if paginatedUsers.total == 0}
        No user registered
        <!-- probably useless, 1 user is the one, watching it -->
      {:else}
        <table class="table table-auto self-start">
          <thead>
            <tr class="text-surface-100-900 flex">
              <th class="w-1/2">Nickname</th>
              <th class="w-1/2"></th>
            </tr>
          </thead>

          <tbody>
            {#each paginatedUsers.users as user (user.uuid)}
              <tr
                onmouseenter={() => (focused = user.uuid)}
                onmouseleave={() => (focused = "")}
                class="bg-surface-700-300 rounded-xl flex hover:motion-safe:hover:brightness-125 dark:hover:motion-safe:hover:brightness-75"
              >
                <td class="w-1/2">{user.nickname}</td>
                <td class="w-1/2">
                  {#if user.uuid === focused || user.uuid === clickFocused}
                    <div class="flex flex-row-reverse">
                      <DeleteDialogSquare
                        title="User deleter"
                        callback={() => (clickFocused = user.uuid)}
                      >
                        <DeleteUserForm
                          name={user.nickname}
                          UUID={user.uuid}
                          callback={() => {
                            trigger++;
                          }}
                        />
                      </DeleteDialogSquare>
                      <ChangeDialogSqare
                        callback={() => (clickFocused = user.uuid)}
                        title="User changer"
                        contentGetter={() => {
                          return user;
                        }}
                      >
                        {#snippet children({ content }: { content: User })}
                          <ChangeUserForm
                            UUID={user.uuid}
                            callback={() => {
                              trigger++;
                            }}
                            response={content}
                          />
                        {/snippet}</ChangeDialogSqare
                      >

                      <PeekDialogSquare
                        callback={() => (clickFocused = user.uuid)}
                        title="User peeker"
                        contentGetter={async () =>
                          new Promise((res) => res(user))}
                      >
                        {#snippet children({ content }: { content: User })}
                          <UserPeek {content} />
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
            count={paginatedUsers?.total ?? 0}
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
        </div>
      {/if}
    {/if}
  {/await}
</div>
