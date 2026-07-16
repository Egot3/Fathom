<script lang="ts">
	import { onMount } from "svelte";
	import type { Component } from "svelte";
	import Login from "./pages/Login.svelte";
	import Status404 from "./pages/Status404.svelte";

	const regex = /^https?:\/\/[^\/]+\/([^?#]+)/;

	let currentURL = $state("");

	let hier_part = $derived(regex.exec(currentURL)?.[1] ?? "");

	$inspect(hier_part);

	onMount(() => {
		currentURL = window.location.href;
	});

	const pages = new Map<string, Component>([["login", Login]]);
</script>

{#if pages.has(hier_part)}
	{@const ActivePage = pages.get(hier_part)}
	<ActivePage />
{:else}
	<Status404 />
{/if}
