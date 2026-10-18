<script lang="ts">
  import type { Component } from "svelte";
  import Status404 from "./pages/Status404.svelte";
  import { GetUser } from "./lib/bgdata/user.svelte";
  import _ from "lodash";
  import { Toast } from "@skeletonlabs/skeleton-svelte";
  import { Toaster } from "./lib/apiutils/toaster";

  type LazyPage = () => Promise<{ default: Component<any, any, any> }>;

  type Route = {
    path: RegExp;
    load: LazyPage;
    public?: boolean;
    teacher?: boolean;
  };

  const routes: Route[] = [
    {
      path: /^\/login\/?$/,
      load: () => import("./pages/Login.svelte"),
      public: true,
    },
    {
      path: /^\/register\/?$/,
      load: () => import("./pages/Register.svelte"),
      public: true,
    },
    { path: /^\/home\/?$/, load: () => import("./pages/Home.svelte") },
    {
      path: /^\/teachery\/?$/,
      load: () => import("./pages/Treachery.svelte"),
      teacher: true,
    },
    {
      path: /^\/(?<runnerKey>\d{1,20})\/?$/,
      load: () => import("./pages/Running.svelte"),
    },
  ];

  function matchRoute(pathname: string) {
    for (const route of routes) {
      const m = route.path.exec(pathname);
      if (m)
        return { route, params: { ...m.groups } as Record<string, string> };
    }
    return null;
  }

  const sendToLogin = () => {
    window.location.href = "/login";
  };
  const sendToHome = () => {
    window.location.href = "/home";
  };

  const match = $derived(matchRoute(window.location.pathname));
  const requestor = GetUser();
</script>

{#if !match}
  <Status404 />
{:else if requestor === null && !match.route.public}
  {sendToLogin()}
{:else if match.route.teacher && !requestor?.isTeacher}
  {sendToHome()}
{:else}
  {#await match.route.load() then { default: Page }}
    <Page {...match.params} />
  {/await}
{/if}

<Toast.Group toaster={Toaster}>
  {#snippet children(toast)}
    <Toast {toast}>
      <Toast.Message>
        <Toast.Title>{toast.title}</Toast.Title>
        <Toast.Description>{toast.description}</Toast.Description>
      </Toast.Message>
      <Toast.CloseTrigger />
    </Toast>
  {/snippet}
</Toast.Group>
