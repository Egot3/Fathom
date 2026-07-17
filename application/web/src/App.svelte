<script lang="ts">
	import { onMount } from "svelte";
	import type { Component } from "svelte";
	import Login from "./pages/Login.svelte";
	import Status404 from "./pages/Status404.svelte";
	import Register from "./pages/Register.svelte";
	import Home from "./pages/Home.svelte";
	import { GetUser } from "./lib/bgdata/user.svelte";

	const regex = /^https?:\/\/[^\/]+\/([^?#]+)/;

	let currentURL = $state("");

	let hier_part = $derived(regex.exec(currentURL)?.[1] ?? "");

	$inspect(hier_part);

	onMount(() => {
		currentURL = window.location.href;
	});
	const pages = new Map<string, Component>([
		["login", Login],
		["register", Register],
		["home", Home],
	]);

	const publicRoutes = new Set(["login", "register"]);

	const sendToLogin = () => {
		window.location.href = "/login";
	};

	console.log(GetUser());
</script>

{#if pages.has(hier_part)}
	{#if GetUser() == null && !publicRoutes.has(hier_part)}
		{sendToLogin()}
	{:else}
		{@const ActivePage = pages.get(hier_part)}
		<ActivePage />
	{/if}
{:else}
	<Status404 />
{/if}
