<script>
	import Sidebar from '$lib/components/Sidebar.svelte';
	import Topbar from '$lib/components/Topbar.svelte';
	import { setAuthToken } from '$lib/api';

	let { data, children } = $props();
	
	// Inject the server-generated JWT into the API client for browser fetches
	$effect(() => {
		if (data?.backendToken) {
			setAuthToken(data.backendToken);
		}
	});
</script>

<div class="flex min-h-screen flex-col bg-background md:flex-row">
	<Sidebar user={data.user} />
	<div class="flex-1 transition-all duration-300 md:ml-64">
		<Topbar />
		<main class="p-6">
			{@render children()}
		</main>
	</div>
</div>
