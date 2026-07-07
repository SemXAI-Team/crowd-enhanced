import { getCameras } from '$lib/server/data';
import jwt from 'jsonwebtoken';

/**
 * Shared data for every dashboard page. Auth + Casbin are enforced by
 * hooks.server.js before any of this runs.
 * @type {import('./$types').LayoutServerLoad}
 */
export const load = async ({ locals }) => {
	const session = await locals.auth();

	let backendToken = '';
	if (session?.user) {
		// Generate a standard JWT for the Go backend grpc-gateway
		backendToken = jwt.sign(
			{ 
				id: session.user.id, 
				role: session.user.role ?? 'viewer' 
			}, 
			'my-super-secret-auth-key-12345', 
			{ expiresIn: '1d' }
		);
	}

	return {
		user: session?.user,
		backendToken,
		cameras: await getCameras()
	};
};
