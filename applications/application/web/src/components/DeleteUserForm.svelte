<script lang="ts">
  import { FetchUserDelete } from "../lib/contracts/user";
  import Confirm from "./Confirm.svelte";

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

<Confirm {name} key={UUID} {callback} actor={FetchUserDelete} value="Delete"
></Confirm>
