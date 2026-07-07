<script>
	import CameraStream from '$lib/components/CameraStream.svelte';
	import TelemetryWidget from '$lib/components/TelemetryWidget.svelte';
	import Icon from '$lib/components/Icon.svelte';
	import { deleteCamera } from '$lib/api';

	let { data } = $props();

	const cameras = $derived(data.cameras);

	/** Only stream the selected camera; the rest stay as cards to keep the page light. */
	let selected = $state(/** @type {string | null} */ (null));
	let gridMode = $state(1); // 1=1x1, 4=2x2, 9=3x3

	$effect(() => {
		if (!selected && cameras.length > 0) selected = cameras[0].id;
	});
	const selectedCamera = $derived(cameras.find((c) => c.id === selected));

	const gridCameras = $derived((() => {
		if (gridMode === 1) return [selected];
		// Auto-populate the grid with up to `gridMode` cameras
		return Array.from({ length: gridMode }, (_, i) => cameras[i]?.id || null);
	})());

	async function removeCam(id) {
		if (!confirm('Delete this camera?')) return;
		try {
			await deleteCamera(id);
			window.location.reload();
		} catch (e) {
			alert(e.message);
		}
	}
</script>

<svelte:head>
	<title>Live Feeds · CrowdGuard</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<h1 class="text-3xl font-bold tracking-tight">Live Feeds</h1>
		<a href="/dashboard/cameras/add" class="rounded-full bg-primary px-5 py-2.5 text-sm font-semibold text-primary-foreground hover:bg-primary/90 inline-flex items-center justify-center gap-2 shadow-sm transition-all hover:shadow-md hover:scale-[1.02]">
			<Icon name="plus" class="h-4 w-4" />
			<span>Add Camera</span>
		</a>
	</div>

	{#if cameras.length === 0}
		<p class="rounded-lg border border-border bg-card px-4 py-6 text-muted-foreground">
			No cameras available. Is the streaming backend running?
		</p>
	{:else}
		<div class="grid gap-6 lg:grid-cols-[240px_1fr]">
			<!-- Camera picker -->
			<nav class="flex flex-col gap-3">
				{#each cameras as cam (cam.id)}
					<button
						type="button"
						onclick={() => (selected = cam.id)}
						class="group rounded-2xl border px-5 py-4 text-left transition-all duration-300
							{selected === cam.id
							? 'border-primary bg-primary/5 shadow-sm'
							: 'border-border/50 bg-card hover:border-primary/30 hover:shadow-sm'}"
					>
						<div class="flex items-center gap-2">
							<span class="h-2 w-2 rounded-full {cam.online ? 'bg-green-500' : 'bg-muted-foreground'}"></span>
							<span class="font-medium">{cam.name}</span>
						</div>
						<span class="mt-0.5 block text-xs text-muted-foreground">{cam.location}</span>
						{#if selected === cam.id}
							<div class="mt-3 flex gap-3">
								<a href="/dashboard/cameras/{cam.id}/edit" class="text-xs text-primary hover:underline" onclick={(e) => e.stopPropagation()}>
									Edit Camera
								</a>
								<button type="button" onclick={(e) => { e.stopPropagation(); removeCam(cam.id); }} class="text-xs text-red-500 hover:underline">
									Delete Camera
								</button>
							</div>
						{/if}
					</button>
				{/each}
			</nav>

			<!-- Player Area -->
			<div class="flex flex-col gap-4">
				<!-- Grid Controls -->
				<div class="flex items-center justify-between rounded-2xl border border-border/50 bg-card p-3 shadow-sm">
					<div class="flex items-center gap-2">
						<button onclick={() => gridMode = 1} class="inline-flex h-8 items-center justify-center rounded-lg px-3 text-xs font-semibold transition-all {gridMode === 1 ? 'bg-primary text-primary-foreground' : 'bg-background border hover:bg-accent'}">
							<Icon name="square" class="mr-1.5 h-3.5 w-3.5" /> 1x1
						</button>
						<button onclick={() => gridMode = 4} class="inline-flex h-8 items-center justify-center rounded-lg px-3 text-xs font-semibold transition-all {gridMode === 4 ? 'bg-primary text-primary-foreground' : 'bg-background border hover:bg-accent'}">
							<Icon name="grid" class="mr-1.5 h-3.5 w-3.5" /> 2x2
						</button>
						<button onclick={() => gridMode = 9} class="inline-flex h-8 items-center justify-center rounded-lg px-3 text-xs font-semibold transition-all {gridMode === 9 ? 'bg-primary text-primary-foreground' : 'bg-background border hover:bg-accent'}">
							<Icon name="layout-grid" class="mr-1.5 h-3.5 w-3.5" /> 3x3
						</button>
					</div>
					<div class="text-xs text-muted-foreground font-medium pr-2">
						NVR View Mode
					</div>
				</div>

				<!-- Video Grid -->
				<div class="grid gap-3 transition-all {gridMode === 4 ? 'grid-cols-2' : gridMode === 9 ? 'grid-cols-3' : 'grid-cols-1'}">
					{#each gridCameras as camId, index}
						{#if camId}
							<div class="relative w-full overflow-hidden rounded-xl border border-border/50 bg-black {selected === camId && gridMode > 1 ? 'ring-2 ring-primary' : ''}">
								<CameraStream cameraId={camId} />
								{#if gridMode > 1}
									<!-- Overlay name in grid mode -->
									<div class="absolute bottom-2 left-2 rounded bg-black/60 px-2 py-1 text-[10px] font-medium text-white backdrop-blur">
										{cameras.find(c => c.id === camId)?.name || 'Unknown'}
									</div>
								{/if}
							</div>
						{:else}
							<div class="flex aspect-video w-full items-center justify-center rounded-xl border border-dashed border-border/50 bg-card/30 text-muted-foreground text-xs">
								Empty Slot
							</div>
						{/if}
					{/each}
				</div>

				<!-- Focus / Telemetry Section -->
				{#if selectedCamera}
					<div class="mt-4 flex items-center justify-between rounded-2xl border border-border/50 bg-card p-4 shadow-sm">
						<div>
							<h3 class="text-sm font-semibold tracking-tight text-foreground flex items-center gap-2">
								<Icon name="cctv" class="h-4 w-4 text-primary" />
								{selectedCamera.name}
							</h3>
							<p class="text-xs text-muted-foreground mt-0.5">{selectedCamera.location}</p>
						</div>
						<div class="flex items-center gap-3">
							<button class="inline-flex items-center justify-center gap-2 rounded-full border border-border/50 bg-card px-4 py-2 text-xs font-semibold hover:bg-accent hover:text-accent-foreground shadow-sm transition-all hover:border-border">
								<Icon name="download" class="h-3.5 w-3.5" />
								<span>Export Data</span>
							</button>
							<button class="inline-flex items-center justify-center gap-2 rounded-full bg-destructive/10 border border-destructive/20 px-4 py-2 text-xs font-semibold text-destructive hover:bg-destructive hover:text-destructive-foreground shadow-sm transition-all">
								<Icon name="alert-triangle" class="h-3.5 w-3.5" />
								<span>Trigger Alarm</span>
							</button>
						</div>
					</div>
					
					<div class="mt-6">
						<div class="mb-4 flex items-center gap-2">
							<Icon name="activity" class="h-5 w-5 text-primary" />
							<h3 class="text-lg font-semibold tracking-tight">Live Telemetry</h3>
						</div>
						<TelemetryWidget cameraId={selected} />
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>
