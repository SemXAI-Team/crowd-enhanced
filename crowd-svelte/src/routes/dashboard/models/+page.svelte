<script>
	import Icon from '$lib/components/Icon.svelte';
	import { uploadModel, listModels, deleteModel } from '$lib/api';
	import { onMount } from 'svelte';

	let models = $state([]);
	let loadingModels = $state(true);

	let files = $state(/** @type {FileList | null} */ (null));
	let uploading = $state(false);
	let message = $state('');

	onMount(async () => {
		try {
			models = await listModels();
		} catch (e) {
			console.error(e);
		} finally {
			loadingModels = false;
		}
	});

	async function upload(/** @type {Event} */ e) {
		e.preventDefault();
		if (!files || files.length === 0) return;

		uploading = true;
		message = '';

		const file = files[0];
		
		try {
			const arrayBuffer = await file.arrayBuffer();
			const base64String = btoa(
				new Uint8Array(arrayBuffer)
					.reduce((data, byte) => data + String.fromCharCode(byte), '')
			);

			const m = await uploadModel(file.name, "Uploaded via dashboard", base64String);
			models = [...models, m];
			message = 'Model uploaded successfully. It will be reloaded in DeepStream shortly.';
			files = null;
		} catch (err) {
			message = 'Error uploading model: ' + err.message;
		} finally {
			uploading = false;
		}
	}

	async function removeModel(id) {
		if (!confirm('Are you sure you want to delete this model?')) return;
		try {
			await deleteModel(id);
			models = models.filter(m => m.id !== id);
		} catch (e) {
			alert('Failed to delete model: ' + e.message);
		}
	}
</script>

<svelte:head>
	<title>Models · CrowdGuard</title>
</svelte:head>

<div class="space-y-6">
	<h1 class="text-3xl font-bold tracking-tight">Models</h1>
	<p class="text-muted-foreground">
		Manage and upload inference models (e.g., .etlt, .engine, labels.txt) directly to the inference backend.
	</p>

	<div class="rounded-lg border border-border bg-card p-6 shadow-sm mb-6">
		<h2 class="text-xl font-semibold mb-4">Deployed Models</h2>
		{#if loadingModels}
			<div class="text-muted-foreground flex items-center gap-2">
				<Icon name="loader-2" class="h-4 w-4 animate-spin" /> Loading...
			</div>
		{:else if models.length === 0}
			<p class="text-muted-foreground text-sm">No models deployed.</p>
		{:else}
			<ul class="space-y-3">
				{#each models as model (model.id)}
					<li class="flex items-center justify-between p-3 border rounded-md bg-background">
						<div>
							<div class="font-medium flex items-center gap-2">
								{model.name}
								{#if model.builtin}
									<span class="text-[10px] uppercase bg-primary/10 text-primary px-1.5 py-0.5 rounded-sm">Built-in</span>
								{/if}
							</div>
							<div class="text-xs text-muted-foreground mt-0.5">{model.description}</div>
						</div>
						{#if !model.builtin}
							<button onclick={() => removeModel(model.id)} class="text-red-500 hover:text-red-400 font-medium transition-colors inline-flex items-center justify-center gap-1.5 text-xs p-2">
								<Icon name="trash-2" class="h-3.5 w-3.5" />
								<span>Delete</span>
							</button>
						{/if}
					</li>
				{/each}
			</ul>
		{/if}
	</div>

	<div class="rounded-lg border border-border bg-card p-6 shadow-sm">
		<h2 class="text-xl font-semibold mb-4">Upload New Model</h2>
		<form onsubmit={upload} class="space-y-4">
			<div>
				<label for="files" class="block text-sm font-medium">Select Model Package (ZIP file)</label>
				<input
					type="file"
					id="files"
					accept=".zip"
					class="mt-2 block w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
					bind:files
					required
				/>
				<p class="mt-2 text-xs text-muted-foreground">
					Check the <a href="/docs/MODEL_UPLOAD" class="text-primary hover:underline">Model Upload Documentation</a> for the required directory structures and formats.
				</p>
			</div>

			<button
				type="submit"
				disabled={uploading}
				class="inline-flex items-center justify-center gap-2 rounded-full bg-primary px-5 py-2.5 text-sm font-semibold text-primary-foreground transition-all shadow-sm hover:shadow-md hover:bg-primary/90 disabled:opacity-50"
			>
				<Icon name={uploading ? 'loader-2' : 'upload'} class="h-4 w-4 {uploading ? 'animate-spin' : ''}" />
				<span>{uploading ? 'Uploading...' : 'Upload Model'}</span>
			</button>
			
			{#if message}
				<p class="text-sm font-medium {message.includes('successfully') ? 'text-green-500' : 'text-red-500'}">
					{message}
				</p>
			{/if}
		</form>
	</div>
</div>
