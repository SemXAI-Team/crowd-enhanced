<script>
	import StatCard from '$lib/components/StatCard.svelte';
	import { getGlobalStats } from '$lib/api';
	import { onMount } from 'svelte';

	let stats = $state(null);
	let loading = $state(true);
	let errorMsg = $state('');

	onMount(async () => {
		try {
			stats = await getGlobalStats();
		} catch (e) {
			errorMsg = e.message;
		} finally {
			loading = false;
		}
	});
</script>

<svelte:head>
	<title>Global Analytics · CrowdGuard</title>
</svelte:head>

<div class="space-y-6">
	<h1 class="text-3xl font-bold tracking-tight">Global Analytics</h1>
	<p class="text-muted-foreground">Aggregated telemetry and insights across all registered camera streams.</p>

	{#if loading}
		<div class="p-12 text-center text-muted-foreground animate-pulse">Loading global statistics...</div>
	{:else if errorMsg}
		<div class="p-12 text-center text-red-500 bg-red-500/10 rounded-lg border border-red-500/20">{errorMsg}</div>
	{:else if stats}
		<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
			<StatCard 
				title="Active Cameras" 
				value={stats.activeCameras} 
				icon="video" 
				trend={2} 
			/>
			<StatCard 
				title="Total Crowd Count" 
				value={stats.totalCrowdCount} 
				icon="users" 
				trend={15} 
			/>
			<StatCard 
				title="Avg Latency (ms)" 
				value={stats.averageLatency?.toFixed(1) ?? '0.0'} 
				icon="activity" 
				trend={-4} 
			/>
			<StatCard 
				title="Active Alerts" 
				value={stats.activeAlerts} 
				icon="alert-triangle" 
				trend={0} 
			/>
		</div>
		
		<div class="mt-8 rounded-lg border border-border bg-card p-6 shadow-sm min-h-[400px] flex items-center justify-center">
			<p class="text-muted-foreground text-sm">Historical graphs placeholder (requires timeseries DB integration)</p>
		</div>
	{/if}
</div>
