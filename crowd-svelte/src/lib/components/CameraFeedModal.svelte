<script>
	import Icon from '$lib/components/Icon.svelte';
	import CameraStream from '$lib/components/CameraStream.svelte';

	/**
	 * @type {{
	 *   camera: import('$lib/server/data').DashboardCamera | null,
	 *   onclose?: () => void
	 * }}
	 */
	let { camera, onclose } = $props();

	function onKeydown(/** @type {KeyboardEvent} */ e) {
		if (e.key === 'Escape') onclose?.();
	}
</script>

<svelte:window onkeydown={onKeydown} />

{#if camera}
	<!-- Backdrop -->
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 p-4 backdrop-blur-sm"
		onclick={(e) => e.target === e.currentTarget && onclose?.()}
		role="presentation"
	>
		<div
			class="w-full max-w-3xl rounded-xl border border-border bg-card shadow-2xl"
			role="dialog"
			aria-modal="true"
			aria-label="{camera.name} live feed"
		>
			<div class="flex items-center justify-between border-b border-border px-5 py-4">
				<div>
					<h2 class="text-lg font-semibold leading-none">{camera.name}</h2>
					<p class="mt-1 text-sm text-muted-foreground">{camera.location}</p>
				</div>
				<div class="flex items-center gap-3">
					<span class="text-sm text-muted-foreground">
						<Icon name="users" class="mr-1 inline h-4 w-4" />{camera.currentCount.toLocaleString()} people
					</span>
					<button
						type="button"
						onclick={() => onclose?.()}
						aria-label="Close"
						class="inline-flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
					>
						<Icon name="x" class="h-4 w-4" />
					</button>
				</div>
			</div>
			<div class="p-5">
				<!-- Mounted only while the modal is open, so the WebRTC peer connection
				     is created on demand and torn down on close. -->
				<CameraStream cameraId={camera.id} />
			</div>
		</div>
	</div>
{/if}
