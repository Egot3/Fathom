<script lang="ts">
  import { FetchTestDelete } from "../lib/contracts/test";
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

<Confirm {name} key={UUID} {callback} actor={FetchTestDelete} value="Delete"
></Confirm>
