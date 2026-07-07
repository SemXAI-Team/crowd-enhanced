<script>
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { updateCamera, listModels } from '$lib/api';
	import Icon from '$lib/components/Icon.svelte';
	import CameraStream from '$lib/components/CameraStream.svelte';
	import TelemetryWidget from '$lib/components/TelemetryWidget.svelte';

	let id = $page.params.id;
	
	// Real Data
	let name = $state('');
	let location = $state('');
	let source = $state('');
	let modelId = $state('');
	
	let saving = $state(false);
	let error = $state('');
	let models = $state([]);

	// AI Overlay Toggle State
	let showOverlay = $state(true);

	// Mocked Advanced Telemetry & Configs for Enterprise Dashboard
	let showCredentials = $state(false);
	const mockData = {
		mac: "00:1A:2B:3C:4D:5E",
		ip: "10.0.0.105",
		type: "PTZ Dome - 4K ANPR Capable",
		onvif: "Profile S, G, T Compliant",
		uptime: "45 days, 12 hours, 3 mins",
		interface: "eth0 (100GbE Trunk VLAN 20)",
		edgeUser: "admin",
		edgePass: "s3cur3_cam_2024!",
		token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjYW1lcmEiOiJjYW0tMSIsImV4cCI6MTcxMjM0NTY3OH0...",
		allowedUsers: ["Global Admins", "L2 Security Operations", "Traffic Analysts"]
	};

	onMount(async () => {
		try {
			models = await listModels();
			
			const res = await fetch('http://localhost:8000/api/v1/cameras');
			if (!res.ok) throw new Error('Failed to fetch cameras');
			const data = await res.json();
			const cam = (data.cameras || []).find((c) => c.id === id);
			
			if (cam) {
				name = cam.name;
				location = cam.location;
				source = cam.source || cam.source_path || '';
				modelId = cam.model_id || (models.length > 0 ? models[0].id : '');
			} else {
				error = "Camera not found";
			}
		} catch (e) {
			error = "Failed to load camera data: " + e.message;
		}
	});

	async function save() {
		saving = true;
		error = '';
		try {
			await updateCamera(id, name, location, source, modelId);
			window.location.href = '/dashboard/settings'; // route back to settings
		} catch (e) {
			error = e.message;
			saving = false;
		}
	}
</script>

<svelte:head>
	<title>Camera {name || id} · CrowdGuard</title>
</svelte:head>

<div class="max-w-8xl mx-auto space-y-8 mt-6 pb-20">
	<div class="flex items-center justify-between">
		<div>
			<div class="flex items-center gap-3">
				<a href="/dashboard/settings" class="inline-flex h-8 w-8 items-center justify-center rounded-full bg-accent hover:bg-accent/80 transition-colors text-muted-foreground hover:text-foreground">
					<Icon name="x" class="h-4 w-4" />
				</a>
				<h1 class="text-3xl font-bold tracking-tight">Camera Configuration</h1>
			</div>
			<p class="text-muted-foreground mt-1 ml-11">Advanced hardware, topology, and security settings for <span class="font-mono text-xs">{id}</span></p>
		</div>
		<button onclick={save} disabled={saving || !!error} class="inline-flex items-center justify-center gap-2 rounded-full bg-primary px-6 py-2.5 text-sm font-semibold text-primary-foreground hover:bg-primary/90 shadow-sm transition-all disabled:opacity-50 hover:shadow-md hover:scale-[1.02]">
			<Icon name={saving ? 'loader-2' : 'settings-2'} class="h-4 w-4 {saving ? 'animate-spin' : ''}" />
			<span>{saving ? 'Applying...' : 'Apply Configurations'}</span>
		</button>
	</div>

	{#if error}
		<div class="rounded-lg bg-destructive/10 border border-destructive/20 p-4 text-destructive flex items-center gap-3">
			<Icon name="alert-triangle" class="h-5 w-5" />
			<span class="font-medium">{error}</span>
		</div>
	{/if}

	<!-- Top Section: Live Preview & Analytics -->
	<div class="grid grid-cols-1 lg:grid-cols-[1fr_350px] gap-6">
		<div class="space-y-4">
			<div class="relative overflow-hidden rounded-xl border border-border/50 bg-black shadow-sm aspect-video w-full">
				<!-- Passing TelemetryWidget as children so it renders on top of the fullscreen player -->
				<CameraStream cameraId={id}>
					<TelemetryWidget cameraId={id} />
				</CameraStream>
				
				<!-- Overlay Toggle Button -->
				<div class="absolute top-4 right-4 z-20">
					<button 
						type="button"
						onclick={() => showOverlay = !showOverlay} 
						class="inline-flex items-center gap-2 rounded-full bg-black/60 px-3 py-1.5 text-xs font-semibold text-white backdrop-blur border {showOverlay ? 'border-primary' : 'border-white/20'} transition-all hover:bg-black/80"
					>
						<Icon name={showOverlay ? 'eye' : 'eye-off'} class="h-3.5 w-3.5 {showOverlay ? 'text-primary' : 'text-white'}" />
						<span>AI Bounding Boxes: {showOverlay ? 'ON' : 'OFF'}</span>
					</button>
				</div>
			</div>
		</div>
		
		<div class="rounded-xl border border-border/50 bg-card shadow-sm overflow-hidden flex flex-col">
			<div class="bg-muted/30 px-5 py-3 border-b border-border/50 flex items-center gap-2">
				<Icon name="activity" class="h-4 w-4 text-primary" />
				<h2 class="text-sm font-semibold">Live Pipeline Statistics</h2>
			</div>
			<div class="p-5 flex-1 overflow-y-auto">
				<TelemetryWidget cameraId={id} />
			</div>
		</div>
	</div>

	<!-- Bottom Section: Configurations -->
	<div class="grid grid-cols-1 xl:grid-cols-[1fr_400px] gap-6">
		<!-- Left Column: Core Settings & Networking -->
		<div class="space-y-6">
			
			<!-- AI Pipeline & Core Settings -->
			<div class="rounded-xl border border-border/50 bg-card shadow-sm overflow-hidden">
				<div class="bg-muted/30 px-6 py-4 border-b border-border/50 flex items-center justify-between">
					<h2 class="text-lg font-semibold flex items-center gap-2">
						<Icon name="activity" class="h-5 w-5 text-primary" /> Core AI Pipeline
					</h2>
				</div>
				<div class="p-6 grid grid-cols-1 md:grid-cols-2 gap-6">
					<div>
						<label class="block text-xs font-semibold mb-1.5 uppercase tracking-wider text-muted-foreground">Camera Name</label>
						<input type="text" bind:value={name} class="w-full rounded-lg border border-input bg-background px-4 py-2.5 text-sm focus:ring-2 focus:ring-primary/20 transition-all" />
					</div>
					<div>
						<label class="block text-xs font-semibold mb-1.5 uppercase tracking-wider text-muted-foreground">Physical Location</label>
						<input type="text" bind:value={location} class="w-full rounded-lg border border-input bg-background px-4 py-2.5 text-sm focus:ring-2 focus:ring-primary/20 transition-all" />
					</div>
					<div class="md:col-span-2">
						<label class="block text-xs font-semibold mb-1.5 uppercase tracking-wider text-muted-foreground">Assigned Inference Model (DeepStream)</label>
						<select bind:value={modelId} class="w-full rounded-lg border border-input bg-background px-4 py-2.5 text-sm focus:ring-2 focus:ring-primary/20 transition-all">
							{#each models as model}
								<option value={model.id}>{model.name} - {model.description}</option>
							{/each}
						</select>
					</div>
				</div>
			</div>

			<!-- Topology & Hardware -->
			<div class="rounded-xl border border-border/50 bg-card shadow-sm overflow-hidden">
				<div class="bg-muted/30 px-6 py-4 border-b border-border/50 flex items-center justify-between">
					<h2 class="text-lg font-semibold flex items-center gap-2">
						<Icon name="network" class="h-5 w-5 text-primary" /> Topology & Hardware
					</h2>
					<span class="inline-flex items-center rounded-full bg-green-500/10 px-2 py-0.5 text-[10px] font-bold text-green-500 uppercase tracking-wider">
						Online
					</span>
				</div>
				<div class="p-6 space-y-6">
					<!-- Stream Routing -->
					<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
						<div class="col-span-1 md:col-span-2">
							<label class="block text-xs font-semibold mb-1.5 uppercase tracking-wider text-muted-foreground">Ingress RTSP (Edge Device -> Server)</label>
							<div class="flex">
								<span class="inline-flex items-center px-3 rounded-l-lg border border-r-0 border-input bg-muted text-muted-foreground text-xs font-mono">
									<Icon name="download" class="h-3 w-3" />
								</span>
								<input type="url" bind:value={source} class="flex-1 rounded-none rounded-r-lg border border-input bg-background px-4 py-2 text-sm font-mono focus:ring-2 focus:ring-primary/20" />
							</div>
						</div>
						<div class="col-span-1 md:col-span-2">
							<label class="block text-xs font-semibold mb-1.5 uppercase tracking-wider text-muted-foreground">Egress RTSP (Server -> VMS Restream)</label>
							<div class="flex">
								<span class="inline-flex items-center px-3 rounded-l-lg border border-r-0 border-input bg-muted text-muted-foreground text-xs font-mono">
									<Icon name="upload" class="h-3 w-3" />
								</span>
								<input type="text" readonly value="rtsp://localhost:8554/processed/{id}" class="flex-1 rounded-none rounded-r-lg border border-input bg-muted/50 px-4 py-2 text-sm font-mono text-muted-foreground cursor-not-allowed" />
							</div>
						</div>
					</div>

					<!-- Hardware Stats Grid -->
					<div class="grid grid-cols-2 md:grid-cols-4 gap-4 pt-4 border-t border-border/50">
						<div>
							<div class="text-[10px] uppercase tracking-wider text-muted-foreground font-semibold mb-1">Camera Type</div>
							<div class="text-sm font-medium">{mockData.type}</div>
						</div>
						<div>
							<div class="text-[10px] uppercase tracking-wider text-muted-foreground font-semibold mb-1">ONVIF Status</div>
							<div class="text-sm font-medium text-green-500">{mockData.onvif}</div>
						</div>
						<div>
							<div class="text-[10px] uppercase tracking-wider text-muted-foreground font-semibold mb-1">MAC Address</div>
							<div class="text-sm font-mono">{mockData.mac}</div>
						</div>
						<div>
							<div class="text-[10px] uppercase tracking-wider text-muted-foreground font-semibold mb-1">Switch Interface</div>
							<div class="text-sm font-mono">{mockData.interface}</div>
						</div>
					</div>
				</div>
			</div>
		</div>

		<!-- Right Column: Security & Edge Credentials -->
		<div class="space-y-6">
			<div class="rounded-xl border border-border/50 bg-card shadow-sm overflow-hidden">
				<div class="bg-muted/30 px-6 py-4 border-b border-border/50 flex items-center justify-between">
					<h2 class="text-lg font-semibold flex items-center gap-2">
						<Icon name="shield-alert" class="h-5 w-5 text-primary" /> Security & Access
					</h2>
				</div>
				<div class="p-6 space-y-6">
					
					<!-- Edge Credentials -->
					<div>
						<div class="flex items-center justify-between mb-2">
							<label class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Edge Credentials (HTTP/RTSP)</label>
							<button type="button" onclick={() => showCredentials = !showCredentials} class="text-xs text-primary hover:underline">
								{showCredentials ? 'Hide' : 'Reveal'}
							</button>
						</div>
						<div class="grid grid-cols-2 gap-3">
							<input type="text" readonly value={mockData.edgeUser} class="rounded-lg border border-input bg-muted/50 px-3 py-2 text-sm font-mono text-muted-foreground" />
							<input type={showCredentials ? "text" : "password"} readonly value={mockData.edgePass} class="rounded-lg border border-input bg-muted/50 px-3 py-2 text-sm font-mono text-muted-foreground" />
						</div>
					</div>

					<!-- Feed Access Token -->
					<div>
						<label class="block text-xs font-semibold mb-2 uppercase tracking-wider text-muted-foreground">Active JWT Feed Token</label>
						<textarea readonly rows="3" class="w-full rounded-lg border border-input bg-muted/50 px-3 py-2 text-xs font-mono text-muted-foreground resize-none">{mockData.token}</textarea>
						<div class="flex justify-end mt-2">
							<button type="button" class="text-xs text-destructive hover:underline flex items-center gap-1 font-medium">
								<Icon name="alert-triangle" class="h-3 w-3" /> Revoke & Rotate Token
							</button>
						</div>
					</div>

					<!-- Uptime -->
					<div class="pt-4 border-t border-border/50">
						<label class="block text-xs font-semibold mb-1 uppercase tracking-wider text-muted-foreground">Physical Uptime</label>
						<div class="text-sm font-medium">{mockData.uptime}</div>
					</div>

					<!-- IAM Authorization -->
					<div class="pt-4 border-t border-border/50">
						<label class="block text-xs font-semibold mb-3 uppercase tracking-wider text-muted-foreground">Authorized IAM Groups</label>
						<div class="flex flex-wrap gap-2">
							{#each mockData.allowedUsers as group}
								<span class="inline-flex items-center gap-1.5 rounded-full bg-accent px-3 py-1 text-xs font-medium text-foreground border border-border/50">
									<Icon name="users" class="h-3 w-3 text-muted-foreground" /> {group}
								</span>
							{/each}
						</div>
					</div>

					<!-- Danger Zone -->
					<div class="pt-6 border-t border-destructive/20 mt-6">
						<button type="button" class="w-full inline-flex items-center justify-center gap-2 rounded-lg bg-destructive/10 border border-destructive/20 px-4 py-2.5 text-sm font-semibold text-destructive hover:bg-destructive hover:text-destructive-foreground transition-all">
							<Icon name="power" class="h-4 w-4" /> Remote Reboot via ONVIF
						</button>
					</div>

				</div>
			</div>
		</div>
	</div>
</div>
