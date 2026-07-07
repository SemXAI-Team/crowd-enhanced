<script>
	import Icon from '$lib/components/Icon.svelte';
	import { listUsers, updateUserRole, deleteUser, listCameras, listModels, deleteCamera } from '$lib/api';
	import { onMount } from 'svelte';

	let users = $state([]);
	let cameras = $state([]);
	let models = $state([]);
	let loading = $state(true);
	let errorMsg = $state('');

	// Add User State
	let showAddUser = $state(false);
	let newName = $state('');
	let newEmail = $state('');
	let newPassword = $state('');
	let newRole = $state('operator');
	let addUserLoading = $state(false);

	// Change Password State
	let changePwUserId = $state(null);
	let changePwValue = $state('');
	let changePwLoading = $state(false);

	onMount(async () => {
		try {
			const [u, c, m] = await Promise.all([
				listUsers(),
				listCameras(),
				listModels()
			]);
			users = u;
			cameras = c;
			models = m;
		} catch (e) {
			errorMsg = e.message;
		} finally {
			loading = false;
		}
	});

	async function removeCam(id) {
		if (!confirm('Are you sure you want to delete this camera?')) return;
		try {
			await deleteCamera(id);
			cameras = cameras.filter(c => c.id !== id);
		} catch (error) {
			alert('Failed to delete camera: ' + error.message);
		}
	}

	async function updateRole(userId, newRole) {
		try {
			await updateUserRole(userId, newRole);
			const u = users.find(u => u.id === userId);
			if (u) u.role = newRole;
		} catch (e) {
			alert('Failed to update role: ' + e.message);
		}
	}

	async function removeUser(userId) {
		if (!confirm('Are you sure you want to completely delete this user?')) return;
		try {
			await deleteUser(userId);
			users = users.filter(u => u.id !== userId);
		} catch (e) {
			alert('Failed to delete user: ' + e.message);
		}
	}

	async function addUser(e) {
		e.preventDefault();
		addUserLoading = true;
		try {
			const res = await fetch('/api/users', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ name: newName, email: newEmail, password: newPassword, role: newRole })
			});
			if (!res.ok) throw new Error(await res.text());
			
			// Reload users
			users = await listUsers();
			showAddUser = false;
			newName = ''; newEmail = ''; newPassword = ''; newRole = 'operator';
		} catch (e) {
			alert('Failed to create user: ' + e.message);
		} finally {
			addUserLoading = false;
		}
	}

	async function updatePassword(e) {
		e.preventDefault();
		changePwLoading = true;
		try {
			const res = await fetch(`/api/users/${changePwUserId}/password`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ password: changePwValue })
			});
			if (!res.ok) throw new Error(await res.text());
			
			alert('Password updated successfully');
			changePwUserId = null;
			changePwValue = '';
		} catch (e) {
			alert('Failed to update password: ' + e.message);
		} finally {
			changePwLoading = false;
		}
	}
</script>

<svelte:head>
	<title>Settings · CrowdGuard</title>
</svelte:head>

<div class="space-y-6">
	<h1 class="text-3xl font-bold tracking-tight">Settings & IAM</h1>
	<p class="text-muted-foreground">Manage roles, permissions, and system access.</p>

	<!-- Camera Streams Management -->
	<div class="rounded-lg border border-border bg-card shadow-sm">
		<div class="p-6">
			<div class="flex items-center justify-between mb-4">
				<div>
					<h2 class="text-xl font-semibold">Camera Streams</h2>
					<p class="text-sm text-muted-foreground">Manage active video feeds and sources.</p>
				</div>
				<a href="/dashboard/cameras/add" class="rounded-full bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90 inline-flex items-center justify-center gap-2 shadow-sm transition-all hover:shadow-md">
					<Icon name="plus" class="h-4 w-4" /> 
					<span>Add Stream</span>
				</a>
			</div>
			
			{#if loading}
				<div class="flex items-center gap-2 text-muted-foreground">
					<Icon name="loader-2" class="animate-spin h-4 w-4" /> Loading cameras...
				</div>
			{:else}
				<div class="overflow-x-auto">
					<table class="w-full text-sm text-left">
						<thead class="bg-muted text-muted-foreground">
							<tr>
								<th class="px-4 py-3 font-medium rounded-tl-md">ID</th>
								<th class="px-4 py-3 font-medium">Name & Location</th>
								<th class="px-4 py-3 font-medium">Model</th>
								<th class="px-4 py-3 font-medium">Status</th>
								<th class="px-4 py-3 font-medium text-right rounded-tr-md">Actions</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-border">
							{#each cameras as camera}
								<tr class="hover:bg-muted/50 transition-colors">
									<td class="px-4 py-3 font-mono text-xs text-muted-foreground">{camera.id}</td>
									<td class="px-4 py-3">
										<div class="font-medium flex items-center gap-2">
											<Icon name="cctv" class="h-4 w-4 text-muted-foreground" />
											{camera.name}
										</div>
										<div class="text-xs text-muted-foreground pl-6">{camera.location}</div>
									</td>
									<td class="px-4 py-3">
										<span class="inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold font-mono">
											{models.find(m => m.id === camera.model_id)?.name || camera.model_id || 'simulation'}
										</span>
									</td>
									<td class="px-4 py-3">
										<span class="inline-flex items-center rounded-md px-2.5 py-0.5 text-xs font-semibold {camera.online ? 'bg-primary text-primary-foreground' : 'bg-secondary text-secondary-foreground'}">
											{camera.online ? 'active' : 'offline'}
										</span>
									</td>
									<td class="px-4 py-3 text-right">
										<a href="/dashboard/cameras/{camera.id}/edit" class="inline-flex items-center justify-center rounded-full bg-accent text-sm font-medium hover:bg-accent/80 hover:text-accent-foreground h-9 w-9 mr-2 transition-colors">
											<Icon name="settings-2" class="h-4 w-4" />
										</a>
										<button onclick={() => removeCam(camera.id)} class="inline-flex items-center justify-center rounded-full text-sm font-medium text-destructive hover:bg-destructive/10 h-9 w-9 transition-colors">
											<Icon name="trash-2" class="h-4 w-4" />
										</button>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		</div>
	</div>

	<!-- User Management -->
	<div class="rounded-lg border border-border bg-card shadow-sm">
		<div class="p-6">
			<div class="flex items-center justify-between mb-4">
				<h2 class="text-xl font-semibold">Identity & Access Management (RBAC)</h2>
				<button onclick={() => showAddUser = !showAddUser} class="rounded-full bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90 inline-flex items-center justify-center gap-2 shadow-sm transition-all hover:shadow-md">
					<Icon name="user-plus" class="h-4 w-4" /> 
					<span>Add User</span>
				</button>
			</div>
			
			{#if showAddUser}
				<div class="mb-6 p-4 rounded-lg border border-border bg-accent/30">
					<h3 class="font-medium mb-3">Create New User</h3>
					<form onsubmit={addUser} class="grid grid-cols-1 md:grid-cols-2 gap-4">
						<div>
							<label class="block text-xs font-medium mb-1">Name</label>
							<input type="text" bind:value={newName} required class="w-full rounded bg-background border px-2 py-1 text-sm" />
						</div>
						<div>
							<label class="block text-xs font-medium mb-1">Email / Username</label>
							<input type="email" bind:value={newEmail} required class="w-full rounded bg-background border px-2 py-1 text-sm" />
						</div>
						<div>
							<label class="block text-xs font-medium mb-1">Password</label>
							<input type="password" bind:value={newPassword} required minlength="4" class="w-full rounded bg-background border px-2 py-1 text-sm" />
						</div>
						<div>
							<label class="block text-xs font-medium mb-1">Role</label>
							<select bind:value={newRole} class="w-full rounded bg-background border px-2 py-1 text-sm">
								<option value="admin">Admin</option>
								<option value="operator">Operator</option>
								<option value="viewer">Viewer</option>
							</select>
						</div>
						<div class="col-span-1 md:col-span-2 flex justify-end items-center gap-3 mt-4">
							<button type="button" onclick={() => showAddUser = false} class="px-4 py-2 text-sm font-medium hover:bg-accent rounded-full transition-colors">Cancel</button>
							<button type="submit" disabled={addUserLoading} class="rounded-full bg-primary px-5 py-2 text-sm text-primary-foreground font-semibold inline-flex justify-center items-center disabled:opacity-50 shadow-sm hover:shadow-md transition-all">Create User</button>
						</div>
					</form>
				</div>
			{/if}

			{#if loading}
				<div class="flex items-center gap-2 text-muted-foreground">
					<Icon name="loader-2" class="animate-spin h-4 w-4" /> Loading users...
				</div>
			{:else if errorMsg}
				<div class="text-red-500">{errorMsg}</div>
			{:else}
				<div class="overflow-x-auto">
					<table class="w-full text-sm text-left">
						<thead class="bg-muted text-muted-foreground">
							<tr>
								<th class="px-4 py-3 font-medium rounded-tl-md">Name</th>
								<th class="px-4 py-3 font-medium">Email</th>
								<th class="px-4 py-3 font-medium">Role</th>
								<th class="px-4 py-3 font-medium rounded-tr-md">Actions</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-border">
							{#each users as user}
								<tr class="hover:bg-muted/50 transition-colors">
									<td class="px-4 py-3 font-medium">{user.name}</td>
									<td class="px-4 py-3 text-muted-foreground">{user.email}</td>
									<td class="px-4 py-3">
										<select 
											class="bg-background border border-border rounded px-2 py-1 text-sm focus:ring-1 focus:ring-primary"
											value={user.role}
											onchange={(e) => updateRole(user.id, e.target.value)}
										>
											<option value="admin">Admin</option>
											<option value="operator">Operator</option>
											<option value="viewer">Viewer</option>
										</select>
									</td>
									<td class="px-4 py-3">
										{#if changePwUserId === user.id}
											<form onsubmit={updatePassword} class="flex items-center gap-2">
												<input type="password" bind:value={changePwValue} placeholder="New password" required minlength="4" class="w-32 rounded bg-background border px-2 py-1 text-sm" />
												<button type="submit" disabled={changePwLoading} class="text-xs bg-primary text-primary-foreground px-2 py-1 rounded">Save</button>
												<button type="button" onclick={() => changePwUserId = null} class="text-xs hover:underline">Cancel</button>
											</form>
										{:else}
											<div class="flex items-center gap-3">
												<button class="text-muted-foreground hover:text-foreground font-medium transition-colors inline-flex items-center justify-center gap-1.5 text-xs" onclick={() => changePwUserId = user.id}>
													<Icon name="key" class="h-3.5 w-3.5" /> 
													<span>Password</span>
												</button>
												<button class="text-red-500 hover:text-red-400 font-medium transition-colors inline-flex items-center justify-center gap-1.5 text-xs" onclick={() => removeUser(user.id)}>
													<Icon name="user-minus" class="h-3.5 w-3.5" />
													<span>Delete</span>
												</button>
											</div>
										{/if}
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		</div>
	</div>

	<!-- System Information Card -->
	<div class="rounded-lg border border-border bg-card shadow-sm">
		<div class="p-6">
			<h2 class="text-xl font-semibold mb-2">System Information</h2>
			<p class="text-sm text-muted-foreground mb-4">Backend connection details.</p>
			
			<div class="grid grid-cols-2 gap-4 text-sm max-w-md">
				<div class="text-muted-foreground">Backend URL</div>
				<div class="font-mono">http://localhost:8000</div>
				<div class="text-muted-foreground">Status</div>
				<div class="text-green-500 font-medium">Connected</div>
			</div>
		</div>
	</div>

	<!-- Enterprise NVR Configurations -->
	<div class="mb-4 mt-8 flex items-center justify-between">
		<h2 class="text-xl font-semibold">Inference Cluster Node Settings</h2>
	</div>
	<div class="grid grid-cols-1 lg:grid-cols-2 gap-6 pb-12">
		<!-- Storage & Retention -->
		<div class="rounded-lg border border-border bg-card shadow-sm flex flex-col">
			<div class="p-6 flex-1">
				<div class="flex items-center gap-2 mb-4 text-primary">
					<Icon name="database" class="h-5 w-5" />
					<h2 class="text-lg font-semibold text-foreground">Local Secure Storage (NVMe RAID)</h2>
				</div>
				<div class="space-y-4">
					<div>
						<div class="flex justify-between text-sm mb-1">
							<span class="text-muted-foreground">Local Disk Array (Strict Security)</span>
							<span class="font-medium">2.1 TB / 16.0 TB (13%)</span>
						</div>
						<div class="h-2 w-full bg-muted rounded-full overflow-hidden">
							<div class="h-full bg-primary rounded-full" style="width: 13%"></div>
						</div>
					</div>
					<div class="pt-2">
						<label class="block text-xs font-medium mb-1">Local Recording Retention Policy</label>
						<select class="w-full rounded bg-background border px-3 py-2 text-sm focus:ring-1 focus:ring-primary">
							<option>Overwrite oldest data when disk is 90% full</option>
							<option>Keep logs/video for 7 Days</option>
							<option>Keep logs/video for 30 Days</option>
							<option>Stop Recording at 95% full (Strict Compliance)</option>
						</select>
					</div>
					<div>
						<p class="text-[10px] text-muted-foreground bg-accent/30 p-2 rounded border border-border/50">
							<Icon name="shield-alert" class="inline h-3 w-3 mr-1" /> All data is written exclusively to the local Threadripper NVMe arrays. No external Data Lake or NAS export is permitted per security policy.
						</p>
					</div>
				</div>
			</div>
			<div class="px-6 py-4 border-t border-border/50 bg-muted/10 flex justify-end">
				<button type="button" class="inline-flex items-center justify-center rounded-full bg-primary px-4 py-2 text-xs font-semibold text-primary-foreground hover:bg-primary/90 shadow-sm transition-all">Save Storage Policy</button>
			</div>
		</div>

		<!-- Network Configuration -->
		<div class="rounded-lg border border-border bg-card shadow-sm flex flex-col">
			<div class="p-6 flex-1">
				<div class="flex items-center gap-2 mb-4 text-primary">
					<Icon name="network" class="h-5 w-5" />
					<h2 class="text-lg font-semibold text-foreground">Direct Ingest Topology</h2>
				</div>
				<div class="space-y-4">
					<p class="text-xs text-muted-foreground">Manage the star topology IP cam network separating ingress traffic from standard management traffic.</p>
					
					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="block text-xs font-medium mb-1">Ingress Interface (eth0)</label>
							<input type="text" value="10.0.0.1 (100GbE)" class="w-full rounded bg-muted border border-transparent px-3 py-2 text-sm text-muted-foreground" disabled />
						</div>
						<div>
							<label class="block text-xs font-medium mb-1">Management Interface (eth1)</label>
							<input type="text" value="192.168.1.10 (1GbE)" class="w-full rounded bg-muted border border-transparent px-3 py-2 text-sm text-muted-foreground" disabled />
						</div>
						<div class="col-span-2">
							<label class="block text-xs font-medium mb-1">Star Topology Switch Trunk Port</label>
							<input type="text" value="VLAN 20 (Camera Net)" class="w-full rounded bg-background border px-3 py-2 text-sm focus:ring-1 focus:ring-primary" />
						</div>
					</div>
				</div>
			</div>
			<div class="px-6 py-4 border-t border-border/50 bg-muted/10 flex justify-end">
				<button type="button" class="inline-flex items-center justify-center rounded-full bg-primary px-4 py-2 text-xs font-semibold text-primary-foreground hover:bg-primary/90 shadow-sm transition-all">Apply Routing</button>
			</div>
		</div>

		<!-- Hardware Acceleration & GPU Pipeline -->
		<div class="rounded-lg border border-border bg-card shadow-sm flex flex-col">
			<div class="p-6 flex-1">
				<div class="flex items-center gap-2 mb-4 text-primary">
					<Icon name="activity" class="h-5 w-5" />
					<h2 class="text-lg font-semibold text-foreground">AI Pipeline & GPU Acceleration</h2>
				</div>
				<div class="space-y-4">
					<div>
						<label class="block text-xs font-medium mb-1">TensorRT Inference Precision</label>
						<select class="w-full rounded bg-background border px-3 py-2 text-sm focus:ring-1 focus:ring-primary">
							<option>FP16 (Recommended for H100)</option>
							<option>INT8 (Maximum Throughput, requires calibration)</option>
							<option>FP32 (Highest Accuracy, lowest speed)</option>
						</select>
					</div>
					<div>
						<label class="block text-xs font-medium mb-1">GPU Load Balancing Strategy</label>
						<select class="w-full rounded bg-background border px-3 py-2 text-sm focus:ring-1 focus:ring-primary">
							<option>Dynamic Auto-Scale across 8x H100 Cluster</option>
							<option>Manual Pinning (Advanced)</option>
							<option>Power Save (Consolidate pipelines)</option>
						</select>
					</div>
					<div class="flex items-center justify-between pt-2">
						<span class="text-sm font-medium">Enable NVDEC Hardware Decoding</span>
						<label class="relative inline-flex items-center cursor-pointer">
							<input type="checkbox" checked class="sr-only peer">
							<div class="w-9 h-5 bg-muted peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-primary rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-primary"></div>
						</label>
					</div>
				</div>
			</div>
			<div class="px-6 py-4 border-t border-border/50 bg-muted/10 flex justify-end">
				<button type="button" class="inline-flex items-center justify-center rounded-full bg-primary px-4 py-2 text-xs font-semibold text-primary-foreground hover:bg-primary/90 shadow-sm transition-all">Update TensorRT Config</button>
			</div>
		</div>
		<div class="rounded-lg border border-border bg-card shadow-sm flex flex-col">
			<div class="p-6 flex-1">
				<div class="flex items-center gap-2 mb-4 text-primary">
					<Icon name="bell-ring" class="h-5 w-5" />
					<h2 class="text-lg font-semibold text-foreground">Event & Alarm Routing</h2>
				</div>
				<div class="space-y-4">
					<div>
						<label class="block text-xs font-medium mb-1">Global Webhook URL (e.g., Slack, PagerDuty)</label>
						<input type="url" placeholder="https://hooks.slack.com/services/..." class="w-full rounded bg-background border px-3 py-2 text-sm focus:ring-1 focus:ring-primary" />
					</div>
					<div>
						<label class="block text-xs font-medium mb-1">SMTP Server (Email Alerts)</label>
						<input type="text" placeholder="smtp.mailgun.org:587" class="w-full rounded bg-background border px-3 py-2 text-sm focus:ring-1 focus:ring-primary" />
					</div>
					<div class="flex items-center justify-between pt-2">
						<span class="text-sm font-medium">Trigger Siren on Crowd Surge</span>
						<label class="relative inline-flex items-center cursor-pointer">
							<input type="checkbox" checked class="sr-only peer">
							<div class="w-9 h-5 bg-muted peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-primary rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-primary"></div>
						</label>
					</div>
				</div>
			</div>
			<div class="px-6 py-4 border-t border-border/50 bg-muted/10 flex justify-end">
				<button type="button" class="inline-flex items-center justify-center rounded-full bg-primary px-4 py-2 text-xs font-semibold text-primary-foreground hover:bg-primary/90 shadow-sm transition-all">Save Routing</button>
			</div>
		</div>

		<!-- System Maintenance -->
		<div class="rounded-lg border border-border bg-card shadow-sm flex flex-col">
			<div class="p-6 flex-1">
				<div class="flex items-center gap-2 mb-4 text-primary">
					<Icon name="wrench" class="h-5 w-5" />
					<h2 class="text-lg font-semibold text-foreground">System Maintenance</h2>
				</div>
				<div class="space-y-4">
					<p class="text-xs text-muted-foreground">Perform routine administrative tasks and system lifecycle management.</p>
					
					<div class="flex flex-col gap-3">
						<button class="inline-flex items-center justify-center gap-2 rounded-full border border-border/50 bg-background px-4 py-2.5 text-sm font-semibold hover:bg-accent transition-colors text-left w-full">
							<Icon name="download" class="h-4 w-4" /> 
							<span>Export Configuration Backup</span>
						</button>
						<button class="inline-flex items-center justify-center gap-2 rounded-full border border-border/50 bg-background px-4 py-2.5 text-sm font-semibold hover:bg-accent transition-colors text-left w-full">
							<Icon name="upload" class="h-4 w-4" /> 
							<span>Firmware Update</span>
						</button>
					</div>
				</div>
			</div>
			<div class="px-6 py-4 border-t border-border/50 bg-destructive/5 flex items-center justify-between">
				<span class="text-xs text-muted-foreground">Warning: Reboot drops active feeds.</span>
				<button type="button" class="inline-flex items-center justify-center gap-2 rounded-full bg-destructive/10 border border-destructive/20 text-destructive px-4 py-2 text-xs font-semibold hover:bg-destructive hover:text-destructive-foreground shadow-sm transition-all">
					<Icon name="power" class="h-3.5 w-3.5" />
					Reboot System
				</button>
			</div>
		</div>
	</div>
</div>
