<script>
	import { signIn } from '@auth/sveltekit/client';
	import { goto } from '$app/navigation';

	let email = $state('');
	let password = $state('');
	let showPassword = $state(false);
	let loading = $state(false);
	let errorMessage = $state('');

	/** @param {SubmitEvent} event */
	const submit = async (event) => {
		event.preventDefault();
		loading = true;
		errorMessage = '';

		const result = await signIn('credentials', { email, password, redirect: false });

		if (result?.error) {
			errorMessage = 'Invalid email or password.';
			loading = false;
			return;
		}

		await goto('/dashboard', { invalidateAll: true });
	};
</script>

<svelte:head>
	<title>Sign In · CrowdGuard</title>
</svelte:head>

<div
	class="relative flex min-h-screen items-center justify-center bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950 p-4 font-sans"
>
	<!-- Decorative background glow -->
	<div class="pointer-events-none absolute inset-0 overflow-hidden">
		<div class="absolute -left-40 -top-40 h-80 w-80 rounded-full bg-white/5 blur-3xl"></div>
		<div class="absolute -bottom-40 -right-40 h-80 w-80 rounded-full bg-white/5 blur-3xl"></div>
	</div>

	<div class="relative w-full max-w-md">
		<!-- Logo / header -->
		<div class="mb-8 text-center">
			<div
				class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-white/10 ring-1 ring-white/20"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="h-8 w-8 text-white"
					aria-hidden="true"
				>
					<path
						d="M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1.17 1.17 0 0 1 1.52 0C14.51 3.81 17 5 19 5a1 1 0 0 1 1 1z"
					/>
					<path d="M12 8v4" />
					<path d="M12 16h.01" />
				</svg>
			</div>
			<h1 class="text-3xl font-bold tracking-tight text-white">CrowdGuard</h1>
			<p class="mt-1 text-sm text-slate-400">by SemXAI</p>
		</div>

		<!-- Login card -->
		<div
			class="rounded-2xl border border-slate-800 bg-slate-900/80 p-8 shadow-2xl backdrop-blur-xl"
		>
			<h2 class="mb-6 text-xl font-semibold text-white">Sign In</h2>

			<form onsubmit={submit} class="space-y-5">
				<div>
					<label for="email" class="mb-2 block text-sm font-medium text-slate-300">Email</label>
					<!-- svelte-ignore a11y_autofocus -->
					<input
						id="email"
						name="email"
						type="email"
						bind:value={email}
						placeholder="admin@semxai.com"
						required
						autofocus
						autocomplete="username"
						class="w-full rounded-lg border border-slate-700 bg-slate-800/50 px-4 py-2.5 text-white placeholder-slate-500 outline-none transition-colors focus:border-white focus:ring-1 focus:ring-white"
					/>
				</div>

				<div>
					<label for="password" class="mb-2 block text-sm font-medium text-slate-300"
						>Password</label
					>
					<div class="relative">
						<input
							id="password"
							name="password"
							type={showPassword ? 'text' : 'password'}
							bind:value={password}
							placeholder="••••••••"
							required
							autocomplete="current-password"
							class="w-full rounded-lg border border-slate-700 bg-slate-800/50 px-4 py-2.5 pr-12 text-white placeholder-slate-500 outline-none transition-colors focus:border-white focus:ring-1 focus:ring-white"
						/>
						<button
							type="button"
							onclick={() => (showPassword = !showPassword)}
							aria-label={showPassword ? 'Hide password' : 'Show password'}
							class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 transition-colors hover:text-slate-200"
						>
							{#if showPassword}
								<svg
									xmlns="http://www.w3.org/2000/svg"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									stroke-linejoin="round"
									class="h-4 w-4"
									aria-hidden="true"
								>
									<path
										d="M10.733 5.076a10.744 10.744 0 0 1 11.205 6.575 1 1 0 0 1 0 .696 10.747 10.747 0 0 1-1.444 2.49"
									/>
									<path d="M14.084 14.158a3 3 0 0 1-4.242-4.242" />
									<path
										d="M17.479 17.499a10.75 10.75 0 0 1-15.417-5.151 1 1 0 0 1 0-.696 10.75 10.75 0 0 1 4.446-5.143"
									/>
									<path d="m2 2 20 20" />
								</svg>
							{:else}
								<svg
									xmlns="http://www.w3.org/2000/svg"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									stroke-linejoin="round"
									class="h-4 w-4"
									aria-hidden="true"
								>
									<path
										d="M2.062 12.348a1 1 0 0 1 0-.696 10.75 10.75 0 0 1 19.876 0 1 1 0 0 1 0 .696 10.75 10.75 0 0 1-19.876 0"
									/>
									<circle cx="12" cy="12" r="3" />
								</svg>
							{/if}
						</button>
					</div>
				</div>

				{#if errorMessage}
					<div
						class="rounded-lg border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400"
						role="alert"
					>
						{errorMessage}
					</div>
				{/if}

				<button
					type="submit"
					disabled={loading}
					class="flex w-full items-center justify-center gap-2 rounded-lg bg-white px-4 py-2.5 font-medium text-slate-900 transition-all hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-50"
				>
					{#if loading}
						<svg
							xmlns="http://www.w3.org/2000/svg"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
							class="h-4 w-4 animate-spin"
							aria-hidden="true"
						>
							<path d="M21 12a9 9 0 1 1-6.219-8.56" />
						</svg>
						Signing in…
					{:else}
						Sign In
					{/if}
				</button>
			</form>
		</div>

		<p class="mt-6 text-center text-xs text-slate-600">CrowdGuard Monitoring System v1.2.0</p>
	</div>
</div>
