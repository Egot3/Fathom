<script lang="ts">
	import NotFoundText from "../lib/texts/NotFoundText.svelte";
</script>

<div class="absolute grid w-full h-full items-center justify-center">
	<div class="flex flex-col items-center justify-center gap-y-5 z-30">
		<p class="text-7xl">Page not found!</p>
		<a
			type="button"
			class="btn text-2xl preset-filled-error-500 text-surface-50 pb-2"
			href="/home">Take me home!</a
		>
	</div>
</div>

<div class="absolute inset-0 grid place-items-center">
	<div
		class="opacity-10 text-center text-pretty w-2/3 pointer-events-none select-none text-surface-500"
	>
		<NotFoundText />
		<NotFoundText />
		<NotFoundText />
		Yes, this is the same text repetead thrise, I got bored a bit
	</div>
</div>
