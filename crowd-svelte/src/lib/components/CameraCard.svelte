<script>
	import Icon from '$lib/components/Icon.svelte';

	/**
	 * @type {{
	 *   camera: import('$lib/server/data').DashboardCamera,
	 *   onclick?: (camera: any) => void
	 * }}
	 */
	let { camera, onclick } = $props();

	const borderColors = {
		low: 'border-l-green-500',
		medium: 'border-l-yellow-500',
		high: 'border-l-orange-500',
		critical: 'border-l-red-500 animate-pulse'
	};
	const densityColor = {
		low: 'bg-green-500',
		medium: 'bg-yellow-500',
		high: 'bg-orange-500',
		critical: 'bg-red-500'
	};

	const fillPct = $derived(
		Math.min(100, Math.round((camera.currentCount / camera.maxCapacity) * 100))
	);
</script>

<button
	type="button"
	onclick={() => onclick?.(camera)}
	class="group w-full cursor-pointer overflow-hidden rounded-lg border border-border border-l-4 bg-card text-left text-card-foreground shadow-sm transition-all hover:shadow-md {borderColors[camera.density]}"
>
	<div class="space-y-3 p-4">
		<!-- Header -->
		<div class="flex items-start justify-between">
			<div>
				<h3 class="text-lg font-semibold leading-none">{camera.name}</h3>
				<p class="mt-1 text-sm text-muted-foreground">{camera.location}</p>
			</div>
			<span
				class="ml-2 inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold
					{camera.online
					? 'border-green-500/30 bg-green-500/10 text-green-500'
					: 'border-transparent bg-destructive text-destructive-foreground'}"
			>
				{camera.online ? 'LIVE' : 'OFFLINE'}
			</span>
		</div>

		<!-- Count + density -->
		<div class="flex items-center justify-between pt-2">
			<div class="flex items-center gap-2">
				<Icon name="users" class="h-4 w-4 text-muted-foreground" />
				<span class="text-sm font-medium">
					{camera.currentCount.toLocaleString()}
					<span class="text-muted-foreground">/ {camera.maxCapacity.toLocaleString()}</span>
				</span>
			</div>
			{#if camera.density === 'critical'}
				<span class="animate-pulse text-xs font-bold uppercase tracking-wide text-red-500">Critical</span>
			{/if}
		</div>

		<!-- Occupancy bar -->
		<div class="h-2 w-full overflow-hidden rounded-full bg-secondary">
			<div
				class="h-full transition-all duration-500 {densityColor[camera.density]}"
				style="width: {fillPct}%"
			></div>
		</div>
	</div>
</button>
