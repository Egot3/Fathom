<script lang="ts">
  import { FetchTestStop } from "../lib/contracts/test";
  import Confirm from "./Confirm.svelte";

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

<Confirm {name} {key} {callback} value={"Stop"} actor={FetchTestStop}></Confirm>
