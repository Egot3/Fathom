<script lang="ts">
  import { onMount } from "svelte";
    import type { Component } from 'svelte';
  import Login from "./pages/Login.svelte";
  import Status404 from "./pages/Status404.svelte";

  const regex = /^http.\/\/.+?\/(.*)$/

  let currentURL = $state('')

  let slug = $derived(
    regex.exec(currentURL)?.[1] ?? ""
  )

  onMount(()=>{
    currentURL = window.location.href
  })

  const pages = new Map<string, Component>([
    ["login", Login],
  ])
</script>

{#if pages.has(slug)}
  {@const ActivePage = pages.get(slug)}
  <ActivePage />
{:else}
  <Status404 />
{/if}
