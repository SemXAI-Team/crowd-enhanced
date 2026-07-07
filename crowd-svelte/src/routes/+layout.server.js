/**
 * Expose the Auth.js session to every page/layout via `data.session`.
 * @type {import('./$types').LayoutServerLoad}
 */
export const load = async ({ locals }) => {
	return { session: await locals.auth() };
};
