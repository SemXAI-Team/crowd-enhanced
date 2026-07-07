import { defineEnvVars } from '@sveltejs/kit/hooks';

export const variables = defineEnvVars({
	DATABASE_URL: { description: 'The database connection string.' },
	ORIGIN: {
		description:
			'The app origin (base URL), e.g. `http://localhost:5173`. Used by adapter-node at runtime.'
	},
	AUTH_SECRET: {
		description:
			'Auth.js secret used to sign/encrypt session tokens. Minimum 32 chars, high entropy (`openssl rand -hex 32`). See [Auth.js](https://authjs.dev).'
	}
});
