<script>
	import { untrack } from 'svelte';
	import { negotiateWebRTC } from '$lib/api';
	import Icon from '$lib/components/Icon.svelte';

	/** @type {{ cameraId: string | null, children?: any }} */
	let { cameraId, children } = $props();

	/** @type {HTMLVideoElement | undefined} */
	let videoEl = $state();
	/** @type {MediaStream | null} */
	let stream = $state(null);
	let status = $state('idle');
	/** @type {RTCPeerConnection | null} */
	let pc = null;
	let isFullscreen = $state(false);
	/** @type {HTMLElement | undefined} */
	let containerEl = $state();

	/**
	 * Resolve once ICE gathering finishes, or after `timeoutMs` — whichever comes
	 * first — so a slow/unreachable STUN server can't stall the handshake. Host
	 * candidates gather almost instantly and are enough for local/LAN peers.
	 */
	function iceComplete(/** @type {RTCPeerConnection} */ conn, timeoutMs = 2500) {
		return new Promise((resolve) => {
			if (conn.iceGatheringState === 'complete') return resolve(undefined);
			const done = () => {
				conn.removeEventListener('icegatheringstatechange', check);
				clearTimeout(timer);
				resolve(undefined);
			};
			const check = () => {
				if (conn.iceGatheringState === 'complete') done();
			};
			const timer = setTimeout(done, timeoutMs);
			conn.addEventListener('icegatheringstatechange', check);
		});
	}

	function disconnect() {
		if (pc) {
			pc.close();
			pc = null;
		}
		stream = null;
	}

	async function connect(/** @type {string} */ id) {
		const conn = new RTCPeerConnection({
			iceServers: [{ urls: 'stun:stun.l.google.com:19302' }]
		});
		pc = conn;
		status = 'connecting';

		conn.addTransceiver('video', { direction: 'recvonly' });
		conn.ontrack = (e) => {
			console.log('[cam] ontrack fired, streams=', e.streams.length, 'current=', conn === pc);
			if (conn === pc) stream = e.streams[0] ?? new MediaStream([e.track]);
		};
		conn.onconnectionstatechange = () => {
			if (conn !== pc) return;
			if (conn.connectionState === 'connected') status = 'live';
			else if (['failed', 'disconnected', 'closed'].includes(conn.connectionState))
				status = 'offline';
		};

		try {
			console.log('[cam] connect start');
			await conn.setLocalDescription(await conn.createOffer());
			console.log('[cam] local set, gathering=', conn.iceGatheringState);
			await iceComplete(conn);
			console.log('[cam] ice complete');
			if (conn !== pc) return;

			const sdp = await negotiateWebRTC(id, conn.localDescription?.sdp || '');
			
			console.log('[cam] got answer, sdp len=', sdp?.length);
			if (conn === pc) await conn.setRemoteDescription({ type: 'answer', sdp });
		} catch (e) {
			console.log('[cam] connect error', e);
			if (conn === pc) status = 'error';
		}
	}

	// (Re)connect whenever the selected camera changes; untrack the rest so the
	// effect depends only on `cameraId`.
	$effect(() => {
		const id = cameraId;
		console.log('[cam] connect effect: cameraId=', JSON.stringify(id));
		untrack(() => {
			disconnect();
			if (id) connect(id);
		});
		return () => untrack(disconnect);
	});

	// Attach the remote stream once both the <video> and the stream exist —
	// avoids the race between the async ontrack callback and bind:this.
	$effect(() => {
		console.log('[cam] attach effect: videoEl=', !!videoEl, 'stream=', !!stream);
		if (videoEl) videoEl.srcObject = stream;
	});

	const statusColor = $derived(
		{
			live: 'bg-emerald-500',
			connecting: 'bg-amber-500 animate-pulse',
			idle: 'bg-slate-600',
			offline: 'bg-slate-500',
			error: 'bg-red-500'
		}[status] ?? 'bg-slate-600'
	);

	function toggleFullscreen() {
		if (!containerEl) return;
		if (!document.fullscreenElement) {
			containerEl.requestFullscreen().catch(err => {
				console.error(`Error attempting to enable fullscreen: ${err.message}`);
			});
		} else {
			document.exitFullscreen();
		}
	}

	$effect(() => {
		const handleFullscreenChange = () => {
			isFullscreen = !!document.fullscreenElement;
		};
		document.addEventListener('fullscreenchange', handleFullscreenChange);
		return () => document.removeEventListener('fullscreenchange', handleFullscreenChange);
	});
</script>

<div bind:this={containerEl} class="relative overflow-hidden rounded-xl border border-slate-800 bg-black {isFullscreen ? 'h-screen w-screen flex items-center justify-center' : ''}">
	<video bind:this={videoEl} autoplay muted playsinline class="{isFullscreen ? 'h-full w-full' : 'aspect-video w-full'} bg-black"></video>

	<div
		class="absolute left-3 top-3 flex items-center gap-2 rounded-full bg-black/60 px-3 py-1 text-xs font-medium text-white backdrop-blur"
	>
		<span class="h-2 w-2 rounded-full {statusColor}"></span>
		{status.toUpperCase()}
	</div>

	<button
		onclick={toggleFullscreen}
		class="absolute bottom-3 left-3 flex h-8 w-8 items-center justify-center rounded-full bg-black/40 text-white backdrop-blur transition-colors hover:bg-white/20"
	>
		<Icon name={isFullscreen ? 'minimize-2' : 'maximize-2'} class="h-4 w-4" />
	</button>

	{#if status !== 'live'}
		<div class="pointer-events-none absolute inset-0 flex items-center justify-center text-slate-500">
			{#if status === 'connecting'}
				Connecting to camera…
			{:else if status === 'error'}
				Stream unavailable
			{:else if !cameraId}
				Select a camera
			{:else}
				Offline
			{/if}
		</div>
	{/if}

	<!-- Fullscreen Overlay Slot -->
	{#if isFullscreen && children}
		<div class="absolute bottom-12 left-8 right-8 z-50 flex gap-4 overflow-x-auto pb-4 pointer-events-none">
			<div class="w-full max-w-4xl mx-auto pointer-events-auto">
				<div class="bg-black/60 backdrop-blur-md border border-white/20 rounded-xl p-6 shadow-2xl">
					{@render children()}
				</div>
			</div>
		</div>
	{/if}
</div>
