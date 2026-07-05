import { redirect } from '@sveltejs/kit';

export const load = (event) => {
	if (!event.locals.user) {
		return redirect(302, '/login');
	}
	return { user: event.locals.user };
};

export const actions = {
	signOut: async (event) => {
		const { auth } = await import('$lib/server/auth');
		await auth.api.signOut({ headers: event.request.headers });
		return redirect(302, '/login');
	}
};
