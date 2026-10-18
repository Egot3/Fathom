<script lang="ts">
  import { ResultAsync } from "neverthrow";
  import type { JSONError } from "../lib/statuses/jsonerror";
  import ErrorPopover from "./ErrorPopover.svelte";
  import { FileUpload } from "@skeletonlabs/skeleton-svelte";
  import { CircleArrowDown } from "@lucide/svelte";

  const {
    fetcher,
    importName,
    callback,
    accept,
    hint = "",
  }: {
    fetcher: (arg0: File) => ResultAsync<null, JSONError>;
    importName: string;
    callback?: () => void;
    accept?: Record<string, string[]>;
    hint?: string;
  } = $props();

  let file: File | null = $state(null);
  function switchFile(details: { acceptedFiles: File[] }) {
    file = details.acceptedFiles?.[0] ?? null;
  }

  let importing = $state(false);

  let statusMessage = $state("");
  let errorPopoverOpen = $state(false);
  async function push() {
    importing = true;
    statusMessage = await fetcher(file as File)
      .orTee((_) => (errorPopoverOpen = true))
      .andTee((_) => {
        errorPopoverOpen = false;
        callback?.();
      })
      .match(
        (_) => "",
        (e) => e.error,
      );
    importing = false;
  }

  let rejection = $state("");
  function onReject(details: { files: { file: File; errors: string[] }[] }) {
    rejection = details.files
      .map((f) => `${f.file.name}: ${f.errors.join(", ")}`)
      .join("; ");
  }
</script>

<FileUpload
  {accept}
  maxFiles={1}
  maxFileSize={8 * 1024 * 1024}
  onFileChange={(d) => {
    rejection = "";
    switchFile(d);
  }}
  onFileReject={onReject}
>
  <FileUpload.Label>{importName}</FileUpload.Label>
  <FileUpload.Dropzone>
    <CircleArrowDown class="size-10" />
    <span>Select file or drag here.</span>
    {#if hint}<span class="text-xs opacity-75">{hint}</span>{/if}
    <FileUpload.Trigger>Browse Files</FileUpload.Trigger>
    <FileUpload.HiddenInput />
  </FileUpload.Dropzone>
  <FileUpload.ItemGroup>
    <FileUpload.Context>
      {#snippet children(fileUpload)}
        {#each fileUpload().acceptedFiles as file (file.name)}
          <FileUpload.Item {file}>
            <FileUpload.ItemName>{file.name}</FileUpload.ItemName>
            <FileUpload.ItemSizeText>{file.size} bytes</FileUpload.ItemSizeText>
            <FileUpload.ItemDeleteTrigger />
          </FileUpload.Item>
        {/each}
      {/snippet}
    </FileUpload.Context>
  </FileUpload.ItemGroup>
  <FileUpload.ClearTrigger>Clear Files</FileUpload.ClearTrigger>
</FileUpload>
{#if rejection}<p class="text-xs text-error-500">{rejection}</p>{/if}

<ErrorPopover message={statusMessage} bind:open={errorPopoverOpen}>
  <button
    class="btn preset-filled-primary-500"
    onclick={push}
    disabled={file == null || importing}>Import</button
  >
</ErrorPopover>
