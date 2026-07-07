<script>
	import StatCard from '$lib/components/StatCard.svelte';
	import CameraCard from '$lib/components/CameraCard.svelte';
	import AlertsPanel from '$lib/components/AlertsPanel.svelte';
	import CameraFeedModal from '$lib/components/CameraFeedModal.svelte';
	import Icon from '$lib/components/Icon.svelte';
	import ClusterComputeWidget from '$lib/components/ClusterComputeWidget.svelte';

	let { data } = $props();

	/** @type {import('$lib/server/data').DashboardCamera | null} */
	let selectedCamera = $state(null);

	const cameras = $derived(data.cameras);
	const criticalCameras = $derived(
		cameras.filter((c) => c.density === 'critical' || c.density === 'high')
	);
	const calmCameras = $derived(
		cameras.filter((c) => c.density !== 'critical' && c.density !== 'high')
	);

	const totalCrowd = $derived(cameras.reduce((acc, c) => acc + c.currentCount, 0));
	const activeCameras = $derived(cameras.filter((c) => c.online).length);
	const avgOccupancy = $derived(
		cameras.length
			? Math.round(
					(cameras.reduce((acc, c) => acc + c.currentCount / c.maxCapacity, 0) / cameras.length) * 100
				)
			: 0
	);
</script>

<svelte:head>
	<title>Overview · CrowdGuard</title>
</svelte:head>

<CameraFeedModal camera={selectedCamera} onclose={() => (selectedCamera = null)} />

<div class="flex gap-6">
	<!-- Alerts panel (large screens) -->
	<div class="hidden w-80 flex-shrink-0 lg:block">
		<AlertsPanel {cameras} onAlertClick={(c) => (selectedCamera = c)} />
	</div>

	<!-- Main content -->
	<div class="flex flex-1 flex-col space-y-6">
		<div class="flex items-center justify-between">
			<h1 class="text-3xl font-bold tracking-tight">Overview</h1>
			<a
				href="/dashboard/cameras"
				class="rounded-full bg-primary px-5 py-2.5 text-sm font-semibold text-primary-foreground transition-all hover:bg-primary/90 shadow-sm hover:shadow-md"
			>
				Live Grid
			</a>
		</div>

		<ClusterComputeWidget />

		<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
			<StatCard title="Total Crowd Count" value={totalCrowd.toLocaleString()} icon="users" iconClass="text-blue-500" />
			<StatCard title="Active Cameras" value="{activeCameras}/{cameras.length}" icon="video" iconClass="text-green-500" />
			<StatCard title="Critical Zones" value={criticalCameras.length} icon="alert-triangle" iconClass="text-red-500" />
			<StatCard title="Avg Occupancy" value="{avgOccupancy}%" icon="activity" iconClass="text-yellow-500" />
		</div>

		{#if cameras.length === 0}
			<p class="rounded-lg border border-border bg-card px-4 py-6 text-muted-foreground">
				No cameras available. Is the streaming backend running?
			</p>
		{:else}
			{#if criticalCameras.length > 0}
				<div>
					<h3 class="mb-4 flex items-center gap-2 text-lg font-semibold text-destructive">
						<Icon name="alert-triangle" class="h-5 w-5" />
						Critical Attention Needed
					</h3>
					<div class="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
						{#each criticalCameras as cam (cam.id)}
							<CameraCard camera={cam} onclick={(c) => (selectedCamera = c)} />
						{/each}
					</div>
				</div>
			{/if}

			<div>
				<h3 class="mb-4 text-lg font-semibold text-muted-foreground">All Feeds</h3>
				<div class="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
					{#each calmCameras as cam (cam.id)}
						<CameraCard camera={cam} onclick={(c) => (selectedCamera = c)} />
					{/each}
				</div>
			</div>
		{/if}
	</div>
</div>
