<script>
	import { enhance } from '$app/forms';

	let { data } = $props();
</script>

<svelte:head>
	<title>Dashboard · CrowdGuard</title>
</svelte:head>

<div class="min-h-screen bg-slate-950 p-8 font-sans text-white">
	<div class="mx-auto max-w-3xl">
		<h1 class="text-2xl font-bold tracking-tight">CrowdGuard Dashboard</h1>
		<p class="mt-2 text-slate-400">
			Signed in as <span class="text-white">{data.user.email}</span>. Dashboard modules coming next.
		</p>
		<form method="post" action="?/signOut" use:enhance class="mt-6">
			<button
				type="submit"
				class="rounded-lg border border-slate-700 bg-slate-800/50 px-4 py-2 text-sm font-medium text-slate-200 transition-colors hover:bg-slate-800"
			>
				Sign out
			</button>
		</form>
	</div>
</div>
