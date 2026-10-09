<script lang="ts" generics="T">
  import { ResultAsync } from "neverthrow";
  import type { JSONError } from "../lib/statuses/jsonerror";
  import ErrorPopover from "./ErrorPopover.svelte";
    import { FileUpload } from "@skeletonlabs/skeleton-svelte";
    import { USER_MIMES } from "../lib/apiutils/acceptHeader";
    import { CircleArrowDown } from "@lucide/svelte";

  const {
    fetcher,
    importName,
  }: {
    fetcher: (arg0: File) => ResultAsync<null, JSONError>;
    importName: string;
  } = $props();

  let file: File | null = $state(null)
  function switchFile(details: { acceptedFiles: File[] }) {
    file = details.acceptedFiles?.[0] ?? null
  }

  let importing = $state(false);

  let statusMessage = $state("");
  let errorPopoverOpen = $state(false);
  async function push() {
    importing = true
    statusMessage = await fetcher(file as File)
      .orTee((_) => (errorPopoverOpen = true))
      .andTee((_) => (errorPopoverOpen = false))
      .match(
        (_) => "",
        (e) => e.error,
      );
    importing = false
  }
</script>

<FileUpload validate={(file)=>USER_MIMES.some((v)=>v.mime==file.type) ? null : []} onFileChange={switchFile} maxFiles={1}>
	<FileUpload.Label>{importName}</FileUpload.Label>
	<FileUpload.Dropzone>
		<CircleArrowDown class="size-10" />
		<span>Select file or drag here.</span>
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

<ErrorPopover message={statusMessage} bind:open={errorPopoverOpen}>
    <button class="btn preset-filled-primary-500" onclick={push} disabled={file == null || importing}>Import</button>
</ErrorPopover>
