<script lang="ts">
	import { onMount } from "svelte";
	import type { Component } from "svelte";
	import Login from "./pages/Login.svelte";
	import Status404 from "./pages/Status404.svelte";
	import Register from "./pages/Register.svelte";
	import Home from "./pages/Home.svelte";
	import { GetUser } from "./lib/bgdata/user.svelte";
	import Treachery from "./pages/Treachery.svelte";

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
		["teachery", Treachery],
	]);

	const publicRoutes = new Set(["login", "register"]);
	const teacherRoutes = new Set(["teachery"]);

	const sendToLogin = () => {
		window.location.href = "/login";
	};
	const sendToHome = () => {
		window.location.href = "/home";
	};

	console.log(GetUser());
</script>

{#if pages.has(hier_part)}
	{@const requestor = GetUser()}
	{#if requestor == null && !publicRoutes.has(hier_part)}
		{sendToLogin()}
	{:else}
		{#if !requestor?.isTeacher && teacherRoutes.has(hier_part)}
			{sendToHome()}
		{:else}
			{@const ActivePage = pages.get(hier_part)}
			<ActivePage />
		{/if}
	{/if}
{:else}
	<Status404 />
{/if}
