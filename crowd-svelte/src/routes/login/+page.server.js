import { fail, redirect } from '@sveltejs/kit';

import { auth } from '$lib/server/auth';
import { APIError } from 'better-auth/api';

export const load = (event) => {
	// Already authenticated — skip the login screen.
	if (event.locals.user) {
		return redirect(302, '/dashboard');
	}
	return {};
};

export const actions = {
	signInEmail: async (event) => {
		const formData = await event.request.formData();
		const email = formData.get('email')?.toString().trim() ?? '';
		const password = formData.get('password')?.toString() ?? '';

		if (!email || !password) {
			return fail(400, { email, message: 'Email and password are required.' });
		}

		try {
			await auth.api.signInEmail({
				body: { email, password }
			});
		} catch (error) {
			if (error instanceof APIError) {
				return fail(400, { email, message: error.message || 'Invalid credentials.' });
			}
			return fail(500, { email, message: 'Something went wrong. Please try again.' });
		}

		return redirect(302, '/dashboard');
	}
};
