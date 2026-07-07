<script>
	import { goto } from '$app/navigation';
	import { signOut } from '@auth/sveltekit/client';
	import Icon from '$lib/components/Icon.svelte';
	import ThemeToggle from '$lib/components/ThemeToggle.svelte';

	/** @type {'none' | 'notifications' | 'user'} */
	let open = $state('none');

	function toggle(/** @type {'notifications' | 'user'} */ menu) {
		open = open === menu ? 'none' : menu;
	}

	// Close any open menu when clicking outside the topbar menus.
	function onWindowClick(/** @type {MouseEvent} */ e) {
		if (open !== 'none' && e.target instanceof Element && !e.target.closest('[data-menu]')) {
			open = 'none';
		}
	}

	const notifications = [
		{
			color: 'bg-red-500',
			title: 'Crowd Surge Detected',
			body: 'North Gate reported high density.',
			time: '2 mins ago'
		},
		{
			color: 'bg-yellow-500',
			title: 'System Update',
			body: 'New model version available.',
			time: '1 hour ago'
		}
	];
</script>

<svelte:window onclick={onWindowClick} />

<header class="sticky top-0 z-30 flex h-16 items-center gap-4 border-b border-border bg-background/80 px-8 backdrop-blur-xl transition-colors duration-300">
	<div class="flex flex-1 items-center gap-4">
		<div class="relative w-full max-w-sm">
			<span class="absolute left-3 top-2.5 text-muted-foreground"><Icon name="search" class="h-4 w-4" /></span>
			<input
				type="search"
				placeholder="Search cameras, zones..."
				class="h-9 w-full rounded-full border border-input bg-card/50 pl-9 pr-4 text-sm outline-none placeholder:text-muted-foreground focus:ring-1 focus:ring-primary focus:bg-background transition-all"
			/>
		</div>
	</div>

	<div class="flex items-center gap-3">
		<ThemeToggle />
		<!-- Notifications -->
		<div class="relative" data-menu>
			<button
				type="button"
				onclick={() => toggle('notifications')}
				aria-label="Notifications"
				class="relative inline-flex h-9 w-9 items-center justify-center rounded-full transition-colors hover:bg-accent hover:text-accent-foreground"
			>
				<Icon name="bell" class="h-4 w-4" />
				<span class="absolute right-1.5 top-1.5 h-2 w-2 rounded-full bg-red-500"></span>
			</button>
			{#if open === 'notifications'}
				<div class="absolute right-0 mt-2 w-[300px] rounded-md border border-border bg-popover p-1 shadow-lg">
					<p class="px-2 py-1.5 text-sm font-semibold">Notifications</p>
					<div class="my-1 h-px bg-border"></div>
					<div class="flex flex-col gap-2 p-2">
						{#each notifications as n (n.title)}
							<div class="flex items-start gap-2 text-sm">
								<span class="mt-1.5 h-2 w-2 flex-shrink-0 rounded-full {n.color}"></span>
								<div>
									<span class="font-medium">{n.title}</span>
									<p class="text-xs text-muted-foreground">{n.body}</p>
									<span class="text-[10px] text-muted-foreground">{n.time}</span>
								</div>
							</div>
						{/each}
					</div>
					<div class="my-1 h-px bg-border"></div>
					<button class="w-full cursor-pointer rounded-sm px-2 py-1.5 text-center text-xs text-muted-foreground hover:bg-accent">
						View all notifications
					</button>
				</div>
			{/if}
		</div>

		<!-- User menu -->
		<div class="relative" data-menu>
			<button
				type="button"
				onclick={() => toggle('user')}
				aria-label="Account"
				class="inline-flex h-9 w-9 items-center justify-center rounded-full transition-colors hover:bg-accent hover:text-accent-foreground"
			>
				<Icon name="user" class="h-4 w-4" />
			</button>
			{#if open === 'user'}
				<div class="absolute right-0 mt-2 w-44 rounded-md border border-border bg-popover p-1 shadow-lg">
					<p class="px-2 py-1.5 text-sm font-semibold">My Account</p>
					<div class="my-1 h-px bg-border"></div>
					<button
						onclick={() => {
							open = 'none';
							goto('/dashboard/settings');
						}}
						class="w-full rounded-sm px-2 py-1.5 text-left text-sm hover:bg-accent"
					>
						Settings
					</button>
					<button class="w-full rounded-sm px-2 py-1.5 text-left text-sm hover:bg-accent">Support</button>
					<div class="my-1 h-px bg-border"></div>
					<button
						onclick={() => signOut({ redirectTo: '/login' })}
						class="w-full rounded-sm px-2 py-1.5 text-left text-sm text-red-500 hover:bg-accent"
					>
						Log out
					</button>
				</div>
			{/if}
		</div>
	</div>
</header>
