<script lang="ts">
  import { onMount } from "svelte";
  import type { Component } from "svelte";
  import Status404 from "./pages/Status404.svelte";
  import { GetUser } from "./lib/bgdata/user.svelte";
  import _ from "lodash";

  const regex = /^https?:\/\/[^\/]+\/([^?#]+)/;

  let currentURL = $derived(window.location.href);

  let hier_part = $derived(regex.exec(currentURL)?.[1] ?? "");

  $inspect(hier_part);
  type LazyPage = () => Promise<{ default: Component<any, any, any> }>;

  const pages = new Map<string, LazyPage>([
    ["login", () => import("./pages/Login.svelte")],
    ["register", () => import("./pages/Register.svelte")],
    ["home", () => import("./pages/Home.svelte")],
    ["teachery", () => import("./pages/Treachery.svelte")],
  ]);

  const publicRoutes = new Set(["login", "register"]);
  const teacherRoutes = new Set(["teachery"]);

  const sendToLogin = () => {
    window.location.href = "/login";
  };
  const sendToHome = () => {
    window.location.href = "/home";
  };
  async function resolvePage(key: string) {
    const loader = pages.get(key);
    if (!loader) return null;
    const { default: PageComponent } = await loader();
    return PageComponent;
  }
</script>

{const requestor = GetUser()}
{#if requestor === null && !publicRoutes.has(hier_part)}
  {sendToLogin()}
{:else}
  {#if !requestor?.isTeacher && teacherRoutes.has(hier_part)}
    {sendToHome()}
  {:else}
    {const ap = resolvePage(hier_part)}

    {#await ap then ActivePage}
      {#if ActivePage === null}
        {console.log("is null")}
        <Status404 />
      {:else}
        {console.log("is not null")}
        <ActivePage />
      {/if}
    {/await}
  {/if}
{/if}
