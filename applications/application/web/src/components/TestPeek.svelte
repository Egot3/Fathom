<script lang="ts">
  import { SvelteSet } from "svelte/reactivity";
  import { FetchTestExport, type Test } from "../lib/contracts/test";
  import ExistingQuizChips from "./ExistingQuizChips.svelte";
  import { IsJSONError, type JSONError } from "../lib/statuses/jsonerror";
    import Exporter from "./Exporter.svelte";

  let { content }: { content: Test | JSONError } = $props();
</script>

{#if IsJSONError(content)}
  <div>{content.error}</div>
{:else}
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

  <div class="flex space-x-2">
      <details class="disclosure">
	<summary>Export options</summary>
	<div class="disclosure-content">
		<p>
			Standard orders ship within 1-2 business days and arrive in 3-5 business days. Expedited shipping is available at checkout for
			next-day delivery in most regions.
		</p>
	</div>
      </details>

      <Exporter downloadName={content.name} fetcher={FetchTestExport.bind(content.uuid), }/>
  </div>

{/if}
