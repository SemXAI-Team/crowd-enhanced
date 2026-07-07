import { redirect, error } from '@sveltejs/kit';
import { sequence } from '@sveltejs/kit/hooks';

import { authHandle } from '$lib/server/auth';
import { authorize } from '$lib/server/authz';

/**
 * Route guard. Runs after Auth.js has populated `event.locals.auth()`.
 *
 * - Requests without a matched route (static assets, Auth.js `/auth/*`) pass through.
 * - No valid session → redirect to `/login`.
 * - Authenticated users hitting `/` or `/login` → redirect to `/dashboard`.
 * - Otherwise Casbin decides based on (role, path, method); denied → 403.
 *
 * @type {import('@sveltejs/kit').Handle}
 */
const guard = async ({ event, resolve }) => {
	// Only guard real routed pages/endpoints; let assets and unknowns pass.
	if (!event.route.id) return resolve(event);

	const session = await event.locals.auth();
	const { pathname } = event.url;

	if (!session?.user) {
		if (pathname === '/login') return resolve(event);
		throw redirect(303, '/login');
	}

	// Signed in: keep users out of the login screen / bare root.
	if (pathname === '/login' || pathname === '/') {
		throw redirect(303, '/dashboard');
	}

	const role = session.user.role ?? 'viewer';
	const allowed = await authorize(role, pathname, event.request.method);
	if (!allowed) throw error(403, 'You do not have permission to access this resource.');

	return resolve(event);
};

export const handle = sequence(authHandle, guard);
