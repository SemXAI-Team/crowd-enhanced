import { json } from '@sveltejs/kit';
import { getBackendHealth } from '$lib/server/data';

/**
 * Backend health probe for the sidebar status indicator. Auth + Casbin are
 * enforced by hooks.server.js before this runs.
 * @type {import('./$types').RequestHandler}
 */
export const GET = async () => json(await getBackendHealth());
