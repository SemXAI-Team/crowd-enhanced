import { SvelteKitAuth } from '@auth/sveltekit';
import Credentials from '@auth/sveltekit/providers/credentials';
import { eq } from 'drizzle-orm';
import bcrypt from 'bcryptjs';

import { AUTH_SECRET } from '$app/env/private';
import { db } from '$lib/server/db';
import { users } from '$lib/server/db/schema';

/**
 * Auth.js (SvelteKit) with an email/password Credentials provider.
 * Sessions are JWT-based (required by Credentials), and the user's `role`
 * is threaded into the token/session for Casbin authorization.
 */
export const { handle: authHandle, signIn, signOut } = SvelteKitAuth({
	secret: AUTH_SECRET,
	trustHost: true,
	session: { strategy: 'jwt' },
	pages: { signIn: '/login' },
	providers: [
		Credentials({
			credentials: {
				email: { label: 'Email', type: 'email' },
				password: { label: 'Password', type: 'password' }
			},
			authorize: async (credentials) => {
				const email = String(credentials?.email ?? '')
					.trim()
					.toLowerCase();
				const password = String(credentials?.password ?? '');
				if (!email || !password) return null;

				const [user] = await db
					.select()
					.from(users)
					.where(eq(users.email, email))
					.limit(1);
				if (!user) return null;

				const valid = await bcrypt.compare(password, user.passwordHash);
				if (!valid) return null;

				return { id: user.id, email: user.email, name: user.name, role: user.role };
			}
		})
	],
	callbacks: {
		jwt: async ({ token, user }) => {
			if (user) token.role = user.role;
			return token;
		},
		session: async ({ session, token }) => {
			if (session.user) session.user.role = token.role;
			return session;
		}
	}
});
