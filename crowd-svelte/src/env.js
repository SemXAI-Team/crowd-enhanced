import { defineEnvVars } from '@sveltejs/kit/hooks';

/**
 * A Standard Schema validator that marks a string variable as optional,
 * so an empty/unset value is allowed instead of being reported as missing.
 * @type {import('@standard-schema/spec').StandardSchemaV1<string | undefined, string | undefined>}
 */
const optionalString = {
	'~standard': {
		version: 1,
		vendor: 'crowd-svelte',
		validate: (value) => ({ value: value || undefined })
	}
};

export const variables = defineEnvVars({
	DATABASE_URL: { description: 'The database connection string.' },
	ORIGIN: {
		description: 'The app origin (base URL), e.g. `http://localhost:5173`.'
	},
	BETTER_AUTH_SECRET: {
		description:
			'Secret used to sign tokens. For production use 32 characters generated with high entropy. See [Better Auth installation](https://www.better-auth.com/docs/installation).'
	},
	GITHUB_CLIENT_ID: {
		schema: optionalString,
		description:
			'GitHub OAuth client ID (optional). See [Better Auth GitHub provider](https://www.better-auth.com/docs/authentication/github).'
	},
	GITHUB_CLIENT_SECRET: {
		schema: optionalString,
		description:
			'GitHub OAuth client secret (optional). See [Better Auth GitHub provider](https://www.better-auth.com/docs/authentication/github).'
	}
});
