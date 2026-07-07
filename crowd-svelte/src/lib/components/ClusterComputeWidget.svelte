<script>
	import Icon from '$lib/components/Icon.svelte';

	// Mock hardware data for the Threadripper + 8x H100 inference node
	let cpuUsage = $state(28);
	let systemRam = $state(42); // percentage
	
	let gpus = $state(
		Array.from({ length: 8 }).map((_, i) => ({
			id: i,
			name: 'H100 80GB',
			utilization: Math.floor(Math.random() * 40) + 40, // 40-80%
			vram: Math.floor(Math.random() * 30) + 20, // 20-50%
			temp: Math.floor(Math.random() * 20) + 50 // 50-70C
		}))
	);

	let totalFps = $state(1250);

	// Simulate live fluctuations
	$effect(() => {
		const interval = setInterval(() => {
			cpuUsage = Math.min(100, Math.max(10, cpuUsage + (Math.random() * 10 - 5)));
			totalFps = Math.floor(1250 + (Math.random() * 100 - 50));
			gpus = gpus.map(g => ({
				...g,
				utilization: Math.min(100, Math.max(0, g.utilization + (Math.random() * 10 - 5))),
				temp: Math.min(90, Math.max(30, g.temp + (Math.random() * 4 - 2)))
			}));
		}, 3000);
		return () => clearInterval(interval);
	});
</script>

<div class="rounded-2xl border border-border/50 bg-card p-6 shadow-sm mb-6">
	<div class="flex flex-col md:flex-row items-start md:items-center justify-between mb-6 gap-4">
		<div>
			<h2 class="text-xl font-bold tracking-tight flex items-center gap-2">
				<Icon name="activity" class="h-5 w-5 text-primary" />
				Compute Cluster Health
			</h2>
			<p class="text-sm text-muted-foreground mt-1">Central Inference Server (Threadripper PRO + 8x H100)</p>
		</div>
		<div class="flex items-center gap-6 bg-background border border-border/50 px-4 py-2 rounded-xl">
			<div class="flex flex-col">
				<span class="text-[10px] uppercase tracking-wider text-muted-foreground font-semibold">Global Inference Rate</span>
				<span class="font-mono text-xl font-bold text-primary">{totalFps} FPS</span>
			</div>
			<div class="h-8 w-px bg-border/50"></div>
			<div class="flex flex-col">
				<span class="text-[10px] uppercase tracking-wider text-muted-foreground font-semibold">Active Pipelines</span>
				<span class="font-mono text-xl font-bold">12</span>
			</div>
		</div>
	</div>

	<div class="grid grid-cols-1 xl:grid-cols-[1fr_2fr] gap-6">
		<!-- CPU & System -->
		<div class="flex flex-col gap-4 border border-border/50 bg-background/50 rounded-xl p-4">
			<h3 class="text-sm font-semibold flex items-center gap-1.5">
				<Icon name="database" class="h-4 w-4 text-muted-foreground" /> System Node
			</h3>
			
			<div class="space-y-1">
				<div class="flex justify-between text-xs font-medium">
					<span class="text-muted-foreground">AMD Threadripper CPU</span>
					<span class="font-mono">{Math.round(cpuUsage)}%</span>
				</div>
				<div class="h-2 w-full bg-muted rounded-full overflow-hidden">
					<div class="h-full bg-blue-500 rounded-full transition-all duration-500" style="width: {cpuUsage}%"></div>
				</div>
			</div>

			<div class="space-y-1">
				<div class="flex justify-between text-xs font-medium">
					<span class="text-muted-foreground">System RAM (512GB)</span>
					<span class="font-mono">{systemRam}%</span>
				</div>
				<div class="h-2 w-full bg-muted rounded-full overflow-hidden">
					<div class="h-full bg-purple-500 rounded-full transition-all duration-500" style="width: {systemRam}%"></div>
				</div>
			</div>
		</div>

		<!-- GPUs -->
		<div class="grid grid-cols-2 md:grid-cols-4 gap-3">
			{#each gpus as gpu}
				<div class="border border-border/50 bg-background/50 rounded-xl p-3 flex flex-col justify-between">
					<div class="flex items-center justify-between mb-2">
						<span class="text-[10px] font-bold text-muted-foreground bg-accent px-1.5 py-0.5 rounded uppercase">GPU {gpu.id}</span>
						<span class="text-[10px] font-mono {gpu.temp > 75 ? 'text-red-500' : 'text-green-500'}">{Math.round(gpu.temp)}°C</span>
					</div>
					
					<div class="space-y-2 mt-auto">
						<div>
							<div class="flex justify-between text-[10px] mb-0.5">
								<span class="text-muted-foreground">Util</span>
								<span class="font-mono">{Math.round(gpu.utilization)}%</span>
							</div>
							<div class="h-1.5 w-full bg-muted rounded-full overflow-hidden">
								<div class="h-full bg-primary rounded-full transition-all duration-500" style="width: {gpu.utilization}%"></div>
							</div>
						</div>
						<div>
							<div class="flex justify-between text-[10px] mb-0.5">
								<span class="text-muted-foreground">VRAM</span>
								<span class="font-mono">{Math.round(gpu.vram)}%</span>
							</div>
							<div class="h-1.5 w-full bg-muted rounded-full overflow-hidden">
								<div class="h-full bg-yellow-500 rounded-full transition-all duration-500" style="width: {gpu.vram}%"></div>
							</div>
						</div>
					</div>
				</div>
			{/each}
		</div>
	</div>
</div>
