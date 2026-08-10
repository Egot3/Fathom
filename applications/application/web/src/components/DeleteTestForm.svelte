<script lang="ts">
  import { FetchTestDelete } from "../lib/contracts/test";
  import DeleteForm from "./DeleteForm.svelte";

  const {
    name,
    UUID,
    callback = () => {},
  }: {
    name: string;
    UUID: string;
    callback: () => void;
  } = $props();
</script>

<DeleteForm {name} {UUID} {callback} deletor={FetchTestDelete}></DeleteForm>
