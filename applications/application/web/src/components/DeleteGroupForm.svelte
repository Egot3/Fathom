<script lang="ts">
  import { FetchGroupDelete } from "../lib/contracts/group";
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

<Confirm {name} key={UUID} {callback} actor={FetchGroupDelete} value="Delete"
></Confirm>
