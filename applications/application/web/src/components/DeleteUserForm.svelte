<script lang="ts">
  import { FetchUserDelete } from "../lib/contracts/user";
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

<DeleteForm {name} {UUID} {callback} deletor={FetchUserDelete}></DeleteForm>
