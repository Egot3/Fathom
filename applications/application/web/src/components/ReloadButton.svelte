<script lang="ts">
  import { RotateCw } from "@lucide/svelte";

  let { onreload }: { onreload: () => void } = $props();
  let time: number;
  let accumClicks: number = 0;
  /*
		- What does this do?
		- Tailwind Shenanigans(TS for short)
	*/
  let isAnimating = $state(false);
  let duration = $state("0s");
  let iterations = $state("0");

  function handleClick() {
    accumClicks++;
    clearTimeout(time);

    time = setTimeout(() => {
      const capturedAccum = accumClicks;
      accumClicks = 0;

      duration = `${capturedAccum * 0.6}s`;
      iterations = `${capturedAccum}`;
      isAnimating = true;

      onreload();

      setTimeout(
        () => {
          isAnimating = false;
        },
        capturedAccum * 600 + 100,
      );
    }, 300);
  }
</script>

<button
  class="ml-auto mt-1"
  class:animate-custom-spin={isAnimating}
  style="--spin-duration: {duration}; --spin-iterations: {iterations};"
  onclick={handleClick}
>
  <RotateCw />
</button>

<style>
  @keyframes custom-spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }

  .animate-custom-spin {
    animation-name: custom-spin;
    animation-duration: calc(0.8s / var(--spin-iterations));
    animation-iteration-count: var(--spin-iterations, 1);

    animation-timing-function: linear(0, 0.1, 0.2, 0.5, 0.9, 1.03, 0.98, 1);
  }
</style>
