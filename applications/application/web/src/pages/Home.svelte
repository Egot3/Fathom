<script lang="ts">
  import BentoContainer from "../components/BentoContainer.svelte";
  import CurrentlyRunning from "../components/CurrentlyRunning.svelte";
  import SelfInfo from "../components/SelfInfo.svelte";
  import TotalsTable from "../components/TotalsTable.svelte";
  import { GetUserOrRedirect } from "../lib/bgdata/user.svelte";

  const me = GetUserOrRedirect();
</script>

<BentoContainer>
  <div
    id="current"
    class="row-start-1 row-end-8 col-start-1 col-end-5 bg-surface-700-300 rounded-2xl p-4"
  >
    <h2 class="text-2xl font-bold">Now running</h2>
    <hr class="hr mt-5 mb-5" />
    <CurrentlyRunning />
  </div>
  <div
    id="previous_results"
    class="row-start-8 row-end-13 col-start-1 col-end-13 bg-surface-700-300 rounded-2xl p-8 flex flex-col"
  >
    <h2 class="self-center text-2xl font-bold">Testing results</h2>
    <hr class="hr mt-5 mb-5" />
    <TotalsTable userUUID={me.UUID} filters={["test"]} />
  </div>
  <div
    id="self"
    class="row-start-1 row-end-8 col-start-5 col-end-13 flex flex-col bg-surface-700-300 rounded-2xl p-8"
  >
    <h2 class="text-2xl font-bold">You are</h2>
    <hr class="hr mt-5 mb-5" />
    <SelfInfo />
  </div>
</BentoContainer>
