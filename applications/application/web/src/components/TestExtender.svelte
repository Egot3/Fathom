<script lang="ts">
  import { FetchTestExtend } from "../lib/contracts/test";

  const { callback, key }: { callback: () => void; key: string } = $props();

  let hours: number = $state(0);
  let minutes: number = $state(0);
  let seconds: number = $state(0);

  let duration: string = $derived(`${hours}h${minutes}m${seconds}s`);
  let statusMessage = $state("");

  async function extendTest(e: SubmitEvent) {
    e.preventDefault();

    if (
      duration === "0h0m0s" ||
      typeof hours !== "number" ||
      typeof minutes !== "number" ||
      typeof seconds !== "number"
    ) {
      statusMessage = "can't run test with no duration";
      return;
    }

    statusMessage = await FetchTestExtend(key, duration).match(
      (r) => {
        callback();
        return "";
      },
      (e) => e.error,
    );

    return;
  }
</script>

<form onsubmit={extendTest} class="space-y-4">
  <fieldset>
    <legend></legend>
    <label class="label">
      <span class="label-text">Hours</span>
      <input type="number" class="input border-2" min="0" bind:value={hours} />
    </label>
    <label class="label">
      <span class="label-text">Minutes</span>
      <input
        type="number"
        class="input border-2"
        min="0"
        bind:value={minutes}
      />
    </label>
    <label class="label">
      <span class="label-text">Seconds</span>
      <input
        type="number"
        class="input border-2"
        min="0"
        bind:value={seconds}
      />
    </label>
  </fieldset>
  <div class="flex">
    <button type="submit" class="btn preset-filled-brand"> Extend </button>
    {#if statusMessage !== ""}
      <div class="ml-auto badge preset-filled-error-100-900 text-xs">
        {statusMessage}
      </div>
    {/if}
  </div>
</form>
