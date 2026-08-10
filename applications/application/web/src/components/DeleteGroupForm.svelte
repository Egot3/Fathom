<script lang="ts">
  import { FetchGroupDelete } from "../lib/contracts/group";
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

<DeleteForm {name} {UUID} {callback} deletor={FetchGroupDelete}></DeleteForm>
