<script>
	import { onMount } from 'svelte';
	import { registerCamera, listModels } from '$lib/api';

	let name = $state('');
	let location = $state('');
	let source = $state('');
	let modelId = $state('');
	let saving = $state(false);
	let error = $state('');
	let models = $state([]);

	onMount(async () => {
		try {
			models = await listModels();
			if (models.length > 0) modelId = models[0].id;
		} catch (e) {
			console.error("Failed to load models", e);
		}
	});

	async function save() {
		saving = true;
		error = '';
		try {
			await registerCamera(name, location, source, modelId);
			window.location.href = '/dashboard/cameras';
		} catch (e) {
			error = e.message;
			saving = false;
		}
	}
</script>

<div class="max-w-xl mx-auto space-y-6 mt-10">
	<h1 class="text-3xl font-bold tracking-tight">Add Camera</h1>
	<p class="text-muted-foreground">Register a new RTSP camera stream.</p>

	<form onsubmit={(e) => { e.preventDefault(); save(); }} class="space-y-4 rounded-lg border border-border bg-card p-6 shadow-sm">
		<div>
			<label class="block text-sm font-medium mb-1">Name</label>
			<input type="text" bind:value={name} required class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" placeholder="e.g. Lobby Cam" />
		</div>
		<div>
			<label class="block text-sm font-medium mb-1">Location</label>
			<input type="text" bind:value={location} required class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" placeholder="e.g. First Floor" />
		</div>
		<div>
			<label class="block text-sm font-medium mb-1">Source (RTSP URL)</label>
			<input type="url" bind:value={source} required class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" placeholder="rtsp://192.168.1.100:554/stream" />
		</div>
		<div>
			<label class="block text-sm font-medium mb-1">Inference Model</label>
			<select bind:value={modelId} class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm">
				{#each models as model}
					<option value={model.id}>{model.name} - {model.description}</option>
				{/each}
			</select>
		</div>

		{#if error}
			<p class="text-red-500 text-sm font-medium">{error}</p>
		{/if}

		<div class="flex justify-end gap-2 mt-6">
			<a href="/dashboard/cameras" class="px-4 py-2 text-sm rounded-md hover:bg-muted">Cancel</a>
			<button type="submit" disabled={saving} class="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50">
				{saving ? 'Saving...' : 'Add Camera'}
			</button>
		</div>
	</form>
</div>
