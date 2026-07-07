<script>
	import { onMount } from 'svelte';
	import Icon from '$lib/components/Icon.svelte';

	let isDark = $state(false);

	onMount(() => {
		// Check local storage or system preference
		const theme = localStorage.getItem('theme');
		if (theme === 'dark' || (!theme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
			isDark = true;
			document.documentElement.classList.add('dark');
		} else {
			isDark = false;
			document.documentElement.classList.remove('dark');
		}
	});

	function toggleTheme() {
		isDark = !isDark;
		if (isDark) {
			document.documentElement.classList.add('dark');
			localStorage.setItem('theme', 'dark');
		} else {
			document.documentElement.classList.remove('dark');
			localStorage.setItem('theme', 'light');
		}
	}
</script>

<button
	onclick={toggleTheme}
	class="relative inline-flex h-9 w-9 items-center justify-center rounded-full border border-border bg-card text-muted-foreground shadow-sm hover:bg-accent hover:text-accent-foreground transition-all duration-200"
	aria-label="Toggle theme"
>
	{#if isDark}
		<Icon name="sun" class="h-4 w-4" />
	{:else}
		<Icon name="moon" class="h-4 w-4" />
	{/if}
</button>
