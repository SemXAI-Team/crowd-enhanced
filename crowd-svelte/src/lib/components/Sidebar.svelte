<script>
	import { page } from '$app/state';
	import { signOut } from '@auth/sveltekit/client';
	import Icon from '$lib/components/Icon.svelte';

	/** @type {{ user: { email?: string|null, name?: string|null, role?: string } | undefined }} */
	let { user } = $props();

	const publicLinks = [
		{ href: '/dashboard', icon: 'layout-dashboard', label: 'Overview', exact: true },
		{ href: '/dashboard/cameras', icon: 'cctv', label: 'Live Feeds', exact: false },
		{ href: '/dashboard/analytics', icon: 'bar-chart', label: 'Analytics', exact: false }
	];
	const adminLinks = [
		{ href: '/dashboard/models', icon: 'sliders', label: 'Models', exact: false },
		{ href: '/dashboard/settings', icon: 'settings', label: 'Settings', exact: false }
	];
	const links = $derived(user?.role === 'admin' ? [...publicLinks, ...adminLinks] : publicLinks);

	const isActive = (/** @type {{href:string,exact:boolean}} */ link) =>
		link.exact
			? page.url.pathname === link.href
			: page.url.pathname.startsWith(link.href);

	// Backend health poll for the status indicator (5s keep-alive, like legacy).
	let health = $state({ online: false, latencyMs: 0, checking: true });
	$effect(() => {
		let alive = true;
		const check = async () => {
			try {
				const res = await fetch('/api/health');
				const data = await res.json();
				if (alive) health = { ...data, checking: false };
			} catch {
				if (alive) health = { online: false, latencyMs: 0, checking: false };
			}
		};
		check();
		const id = setInterval(check, 5000);
		return () => {
			alive = false;
			clearInterval(id);
		};
	});
</script>

<aside class="fixed left-0 top-0 z-40 hidden h-screen w-64 border-r border-border bg-card/50 backdrop-blur-xl md:block">
	<div class="flex h-16 items-center border-b border-border px-6">
		<Icon name="shield-alert" class="mr-2 h-6 w-6 text-primary" />
		<span class="text-lg font-bold tracking-tight">CrowdGuard <span class="text-muted-foreground font-medium">by SemXAI</span></span>
	</div>

	<nav class="flex flex-col gap-1 p-4">
		{#each links as link (link.href)}
			<a
				href={link.href}
				class="group relative flex items-center rounded-lg px-3 py-2.5 text-sm font-medium transition-all duration-200
					{isActive(link)
					? 'bg-primary/10 text-primary'
					: 'text-muted-foreground hover:bg-accent/50 hover:text-foreground'}"
			>
				{#if isActive(link)}
					<div class="absolute left-0 top-1/2 -mt-3 h-6 w-1 rounded-r-full bg-primary shadow-[0_0_8px_var(--primary)]"></div>
				{/if}
				<Icon name={link.icon} class="mr-3 h-5 w-5 {isActive(link) ? 'text-primary' : 'text-muted-foreground group-hover:text-foreground transition-colors'}" />
				{link.label}
			</a>
		{/each}
	</nav>

	<div class="absolute bottom-4 left-0 w-full px-4">
		<div class="rounded-xl border border-border bg-card shadow-sm p-4 text-xs text-muted-foreground">
			<p class="font-semibold text-foreground mb-1">System Status</p>

			<div class="flex items-center gap-2 py-1">
				<span
					class="h-2 w-2 rounded-full transition-colors duration-300
						{health.checking
						? 'bg-yellow-500'
						: health.online
							? 'bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.6)]'
							: 'bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.6)]'}"
				></span>
				<span class="flex flex-col">
					<span class="text-xs font-medium">
						{health.checking ? 'Connecting…' : health.online ? 'System Online' : 'Backend Offline'}
					</span>
					{#if health.online}
						<span class="font-mono text-[10px] text-muted-foreground/60">{health.latencyMs}ms latency</span>
					{/if}
				</span>
			</div>

			{#if user}
				<div class="mt-3 pt-3 border-t border-border/50">
					<p class="text-foreground/70">
						<span class="font-medium text-foreground">{user.name ?? user.email}</span>
						{#if user.role}
							<span class="ml-1 rounded text-[10px] uppercase tracking-wider text-primary font-semibold">{user.role}</span>
						{/if}
					</p>
					<button
						type="button"
						onclick={() => signOut({ redirectTo: '/login' })}
						class="mt-2 flex w-full items-center gap-2 rounded text-xs font-medium text-muted-foreground transition-colors hover:text-foreground"
					>
						<Icon name="log-out" class="h-3 w-3" />
						Sign Out
					</button>
				</div>
			{/if}
		</div>
	</div>
</aside>
