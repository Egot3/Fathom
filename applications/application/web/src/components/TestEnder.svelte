<script lang="ts">
  import { FetchTestStop } from "../lib/contracts/test";
  import DeleteForm from "./DeleteForm.svelte";

  const {
    name,
    key,
    callback = () => {},
  }: {
    name: string;
    key: string;
    callback: () => void;
  } = $props();
</script>

<DeleteForm {name} UUID={key} {callback} deletor={FetchTestStop}></DeleteForm>
