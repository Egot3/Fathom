<script lang="ts">
  import ChangeDialogSqare from "./ChangeDialogSqare.svelte";
  import PeekDialogSquare from "./PeekDialogSquare.svelte";
  import DeleteDialogSquare from "./DeleteDialogSquare.svelte";
  import {
    FetchUserImport,
    FetchUsers,
    type User,
  } from "../lib/contracts/user";
  import DeleteUserForm from "./DeleteUserForm.svelte";
  import ChangeUserForm from "./ChangeUserForm.svelte";
  import UserPeek from "./UserPeek.svelte";
  import { CreateFitRows } from "../lib/layoututils/fitRows.svelte";
  import Pagination from "./Pagination.svelte";
  import PeekDialogue from "./PeekDialogue.svelte";
  import { Dialog } from "@skeletonlabs/skeleton-svelte";
  import { CircleArrowDown } from "@lucide/svelte";
  import { CreatePagedQuery } from "../lib/layoututils/createPagedQuery.svelte";
  import Importer from "./Importer.svelte";

  const fit = CreateFitRows({ rowHeight: 20, reserve: 1 });

  let focused = $state("");
  let clickFocused = $state("");

  const q = CreatePagedQuery({
    pageSize: () => fit.pageSize,
    enabled: () => fit.ready,
    params: () => ({}),
    fetcher: FetchUsers,
  });
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
    No user registered
    <!-- probably useless, 1 user is the one, watching it -->
  {:else}
    <table class="table table-auto self-start">
      <thead>
        <tr style="height:20px" class="text-surface-100-900 flex">
          <th class="w-1/2">Nickname</th>
          <th class="w-1/2"></th>
        </tr>
      </thead>

      <tbody>
        {#each q.data.users as user (user.uuid)}
          <tr
            style="height:20px"
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
                      callback={q.refresh}
                    />
                    <!-- need only cofirmation, as user is regenerated on re-deploy -->
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
                        callback={q.refresh}
                        response={content}
                      />
                    {/snippet}</ChangeDialogSqare
                  >

                  <PeekDialogSquare
                    callback={() => {
                      clickFocused = user.uuid;
                    }}
                    title="User peeker"
                    contentGetter={() => user}
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
        pageSize={fit.pageSize}
        bind:page={q.page}
        total={q.data?.total ?? 0}
      ></Pagination>

      <PeekDialogue title="User importing">
        {#snippet trigger()}
          <Dialog.Trigger class="btn preset-filled-warning-500">
            <CircleArrowDown></CircleArrowDown>
          </Dialog.Trigger>
        {/snippet}

        <Importer
          importName="user import"
          fetcher={FetchUserImport}
          accept={{
            "text/csv": [".csv"],
            "text/tab-separated-values": [".tsv"],
            "application/gzip": [".gz"],
          }}
          hint="CSV/TSV with header: name,password (optionally gzipped)"
          callback={q.refresh}
        />
      </PeekDialogue>
    </div>
  {/if}
</div>
