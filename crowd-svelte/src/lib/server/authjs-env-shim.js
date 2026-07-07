/**
 * Compatibility shim for `$env/dynamic/private`.
 *
 * `@auth/sveltekit` imports `env` from `$env/dynamic/private`, but SvelteKit 3
 * replaced the `$env/*` virtual modules with `$app/env/*` (explicit env vars).
 * Auth.js only reads this for defaults (AUTH_SECRET, AUTH_TRUST_HOST, …) — we
 * pass `secret` and `trustHost` explicitly in `src/lib/server/auth.js`, so
 * exposing `process.env` here is sufficient and keeps Auth.js importable.
 *
 * Wired up via a Vite alias in `vite.config.js`.
 */
export const env = process.env;
