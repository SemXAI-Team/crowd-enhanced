import '@auth/sveltekit';

// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// `locals.auth()` is added by the Auth.js handle (@auth/sveltekit).
		// interface Locals {}

		// interface Error {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

// Extend the Auth.js session user with our Casbin role.
declare module '@auth/sveltekit' {
	interface User {
		role?: string;
	}
	interface Session {
		user: {
			id?: string;
			name?: string | null;
			email?: string | null;
			image?: string | null;
			role?: string;
		};
	}
}

export {};
