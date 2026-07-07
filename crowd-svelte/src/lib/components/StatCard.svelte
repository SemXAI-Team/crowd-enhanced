<script>
	import Icon from '$lib/components/Icon.svelte';

	/** @type {{ title: string, value: string | number, icon: string, iconClass?: string }} */
	let { title, value, icon, iconClass = 'text-muted-foreground' } = $props();
</script>

<div class="group rounded-2xl border border-border/50 bg-card p-6 shadow-sm transition-all hover:shadow-md hover:border-primary/30">
	<div class="flex flex-row items-center justify-between">
		<h3 class="text-sm font-semibold text-muted-foreground tracking-tight">{title}</h3>
		<div class="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10 text-primary group-hover:scale-110 transition-transform duration-300">
			<Icon name={icon} class="h-5 w-5 {iconClass}" />
		</div>
	</div>
	<div class="mt-4">
		<div class="text-3xl font-bold tracking-tighter text-foreground">{value}</div>
	</div>
</div>
