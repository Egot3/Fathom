<script lang="ts">
	import { GetUserOrRedirect } from "../lib/bgdata/user.svelte";

	const user = $state(GetUserOrRedirect());
</script>

<hr class="hr mt-5 mb-5" />
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
		<a class="mt-auto btn preset-filled-primary-500 flex w-fit" href="/teachery"
			>Go to teacher's control panel</a
		>
	{/if}
</div>
