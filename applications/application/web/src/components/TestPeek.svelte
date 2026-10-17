<script lang="ts">
  import { SvelteSet } from "svelte/reactivity";
  import { FetchTestExport, type Test } from "../lib/contracts/test";
  import ExistingQuizChips from "./ExistingQuizChips.svelte";
  import Exporter from "./Exporter.svelte";
  import AcceptHeaders from "./AcceptHeaders.svelte";
  import { USER_MIMES, type Accept } from "../lib/apiutils/acceptHeader";

  let { content }: { content: Test } = $props();

  let best: Accept[] = $state([]);
</script>

<label class="label">
  <span class="label-text">Name</span>
  <input
    class="input"
    type="text"
    value={content.name}
    placeholder="go-generics"
    disabled
  />
</label>

{const quizzes = new SvelteSet(content.quizzes)}
<label class="label">
  <span class="label-text">Quizzes</span>
  <ExistingQuizChips disabled {quizzes} />
</label>

<div class="flex space-x-2 p-2">
  <details class="disclosure">
    <summary>Export options</summary>
    <div class="disclosure-content">
      <AcceptHeaders available={USER_MIMES} bind:accepted={best} />
    </div>
  </details>

  <div class="mt-auto mb-2">
    <Exporter
      downloadName={content.name}
      fetcher={() => FetchTestExport(content.uuid, best)}
    >
      Export
    </Exporter>
  </div>
</div>
