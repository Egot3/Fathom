<script lang="ts">
  import { GetUserOrRedirect } from "../lib/bgdata/user.svelte";
  import { Navigate } from "../lib/router.svelte";

  const user = $state(GetUserOrRedirect());
</script>

<div class="h-full flex flex-col">
  <h1 class="text-2xl">{user.nickname}</h1>
  {#if user.isTeacher}
    <div class="badge w-fit bg-surface-300-700 text-surface-950-50">
      Teacher
    </div>
  {/if}

  <p class="opacity-80 mt-10">
    your UUID: {user.UUID}
  </p>

  {#if user.isTeacher}
    <button
      class="mt-auto btn preset-filled-primary-500 flex w-fit"
      onclick={(e) => {
        e.preventDefault();
        Navigate("/teachery");
      }}>Go to teacher's control panel</button
    >
  {/if}
</div>
