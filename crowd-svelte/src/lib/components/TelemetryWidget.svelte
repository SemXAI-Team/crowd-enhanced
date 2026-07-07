<script>
	import { onMount, onDestroy, untrack } from 'svelte';
	import { streamTelemetry } from '$lib/api';
	import StatCard from '$lib/components/StatCard.svelte';
	import Icon from '$lib/components/Icon.svelte';

	let { cameraId } = $props();

	let data = $state({});
	let alertMsg = $state('');
	let activeCall = null;

	$effect(() => {
		const id = cameraId;
		
		untrack(() => {
			if (activeCall) {
				activeCall.abort();
				activeCall = null;
			}
			
			if (id) {
				const controller = new AbortController();
				activeCall = controller;
				
				const connectTelemetry = async () => {
					try {
						await streamTelemetry(id, (sample) => {
							untrack(() => {
								const newData = { ...data };
								if (sample.metrics) {
									for (const m of sample.metrics) {
										newData[m.label] = m.value;
									}
								}
								data = newData;
								if (sample.alert) {
									alertMsg = sample.alert_reason || 'Crowd threshold exceeded';
								} else {
									alertMsg = '';
								}
							});
						}, controller.signal);
					} catch (e) {
						if (e.name !== 'AbortError') {
							console.error('Telemetry stream error', e);
						}
					}
				};
				
				connectTelemetry();
			}
		});

		return () => untrack(() => {
			if (activeCall) {
				activeCall.abort();
				activeCall = null;
			}
		});
	});
</script>

{#if alertMsg}
	<div class="mb-4 flex items-center gap-3 rounded-lg border border-red-500 bg-red-500/10 px-4 py-3 font-semibold text-red-500 animate-pulse">
		<Icon name="alert-triangle" class="h-5 w-5" />
		<span>ALERT: {alertMsg}</span>
	</div>
{/if}

<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
	<!-- Dynamically render generic metrics from the telemetry JSON payload -->
	{#if Object.keys(data).length === 0}
		<div class="col-span-full text-center text-muted-foreground p-8 border border-dashed border-border rounded-xl">
			Waiting for inference telemetry...
		</div>
	{:else}
		{#each Object.entries(data) as [key, value]}
			<StatCard 
				title={key} 
				value={value} 
				icon="activity" 
				trend={0} 
			/>
		{/each}
	{/if}
</div>
