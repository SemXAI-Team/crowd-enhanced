<script>
	import Icon from '$lib/components/Icon.svelte';

	/**
	 * @type {{
	 *   cameras: import('$lib/server/data').DashboardCamera[],
	 *   onAlertClick?: (camera: any) => void
	 * }}
	 */
	let { cameras, onAlertClick } = $props();

	const alerts = $derived(cameras.filter((c) => c.density === 'critical' || c.density === 'high'));
</script>

<div class="flex h-full flex-col rounded-lg border border-destructive/20 bg-card shadow-sm">
	<div class="border-b border-border p-4">
		<h2 class="flex items-center gap-2 text-lg font-semibold">
			<Icon name="alert-triangle" class="h-5 w-5 text-destructive" />
			Active Alerts
			{#if alerts.length > 0}
				<span class="ml-auto inline-flex items-center rounded-full bg-destructive px-2.5 py-0.5 text-xs font-semibold text-destructive-foreground">
					{alerts.length}
				</span>
			{/if}
		</h2>
	</div>

	<div class="flex-1 overflow-y-auto">
		<div class="space-y-3 p-4">
			{#if alerts.length === 0}
				<div class="py-8 text-center text-sm text-muted-foreground">No active alerts at this time.</div>
			{:else}
				{#each alerts as camera (camera.id)}
					<button
						type="button"
						onclick={() => onAlertClick?.(camera)}
						class="group flex w-full cursor-pointer flex-col gap-2 rounded-lg border border-border bg-background/60 p-3 text-left shadow-sm transition-colors hover:bg-accent"
					>
						<div class="flex items-center justify-between">
							<span class="truncate font-medium">{camera.location}</span>
							<span
								class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-semibold
									{camera.density === 'critical'
									? 'bg-destructive text-destructive-foreground'
									: 'bg-orange-500 text-white'}"
							>
								{camera.density.toUpperCase()}
							</span>
						</div>

						<div class="flex items-center gap-4 text-xs text-muted-foreground">
							<span class="flex items-center gap-1"><Icon name="cctv" class="h-3 w-3" />{camera.name}</span>
							<span>Just now</span>
						</div>

						<div class="mt-1 flex items-center justify-between text-xs">
							<span class="font-medium text-foreground">{camera.currentCount.toLocaleString()} people</span>
							<span class="text-muted-foreground">
								Capacity: {Math.round((camera.currentCount / camera.maxCapacity) * 100)}%
							</span>
						</div>
					</button>
				{/each}
			{/if}
		</div>
	</div>
</div>
