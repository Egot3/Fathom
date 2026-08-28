<script lang="ts">
  import ArrowLeftIcon from "@lucide/svelte/icons/arrow-left";
  import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
  import { Pagination } from "@skeletonlabs/skeleton-svelte";
  import { IsJSONError, type JSONError } from "../lib/statuses/jsonerror";
  import {
    FetchAllTests,
    FetchTest,
    type Test,
    type Tests,
  } from "../lib/contracts/test";
  import CreateDialog from "./CreateDialog.svelte";
  import CreateTestForm from "./CreateTestForm.svelte";
  import SmallBookCreateDialog from "./SmallBookCreateDialog.svelte";
  import DeleteDialogSquare from "./DeleteDialogSquare.svelte";
  import DeleteTestForm from "./DeleteTestForm.svelte";
  import ChangeDialogSqare from "./ChangeDialogSqare.svelte";
  import ChangeTestForm from "./ChangeTestForm.svelte";
  import PeekDialogSquare from "./PeekDialogSquare.svelte";
  import TestPeek from "./TestPeek.svelte";

  let height = $state(0);

  let page = $state(1);
  let pageSize = $derived(Math.trunc((height - 39 - 45) / 39));

  let statusMessage = $state("");
  let loading = $state(true);
  let listTests: Tests = $state(null as never);

  let trigger = $state(0);
  let time: number;
  $effect(() => {
    const p = page;
    const ps = pageSize;
    trigger;
    loading = true;

    console.log("detected change");
    clearTimeout(time);

    time = setTimeout(async () => {
      statusMessage = await FetchAllTests(p - 1, ps)
        .map((r) => {
          listTests = r;
          loading = false;
          return r;
        })
        .match(
          () => "",
          (err: JSONError) => err.error,
        );
    });
  });

  let focused = $state("");
  let clickFocused = $state("");
</script>

<div
  class="grid gap-4 w-full place-items-center h-full overflow-auto"
  bind:clientHeight={height}
>
  {#if loading}
    <div
      class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl"
    ></div>
  {:else}
    {#if statusMessage !== ""}
      <div>{statusMessage}</div>
    {:else}
      {#if listTests?.total === 0}
        <CreateDialog
          name="Create your first test"
          title="Create your first test"
          ><CreateTestForm
            callback={() => {
              trigger++;
            }}
          /></CreateDialog
        >
      {:else}
        <table class="table table-auto self-start">
          <thead>
            <tr class="text-surface-100-900 flex">
              <th class="w-1/3">Test</th>
              <th class="w-1/3">Quiz count</th>
              <th class="w-1/3"></th>
            </tr>
          </thead>

          <tbody>
            {#each listTests.tests as test}
              <tr
                onmouseenter={() => (focused = test.uuid)}
                onmouseleave={() => (focused = "")}
                class="bg-surface-700-300 rounded-xl flex hover:motion-safe:hover:brightness-125 dark:hover:motion-safe:hover:brightness-75"
              >
                <td class="w-1/3">{test.name}</td>
                <!-- 10.5 rem = 168px -->
                <td class="w-1/3">{test.quizzes?.length ?? 0}</td>
                <td class="w-1/3">
                  {#if test.uuid === focused || test.uuid === clickFocused}
                    <div class="flex flex-row-reverse">
                      <DeleteDialogSquare
                        title="Quiz deleter"
                        callback={() => (clickFocused = test.uuid)}
                      >
                        <DeleteTestForm
                          name={test.name}
                          UUID={test.uuid}
                          callback={() => {
                            trigger++;
                          }}
                        />
                      </DeleteDialogSquare>
                      <ChangeDialogSqare
                        callback={() => (clickFocused = test.uuid)}
                        title="Quiz changer"
                        contentGetter={async () => {
                          return await FetchTest(test.uuid);
                        }}
                      >
                        {#snippet children({ content }: { content: Test })}
                          <ChangeTestForm
                            UUID={test.uuid}
                            name={test.name}
                            callback={() => {
                              trigger++;
                            }}
                            response={content}
                          />
                        {/snippet}</ChangeDialogSqare
                      >

                      <PeekDialogSquare
                        callback={() => (clickFocused = test.uuid)}
                        title="Test peeker"
                        contentGetter={async () => {
                          return FetchTest(test.uuid);
                        }}
                      >
                        {#snippet children({
                          content,
                        }: {
                          content: Test | JSONError;
                        })}
                          <TestPeek {content} />
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
            count={listTests?.total ?? 0}
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
            ><CreateTestForm
              callback={() => {
                trigger++;
              }}
            /></SmallBookCreateDialog
          >
        </div>
      {/if}
    {/if}
  {/if}
</div>
