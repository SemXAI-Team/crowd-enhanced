import { newEnforcer, newModelFromString, StringAdapter } from 'casbin';

/**
 * Casbin RBAC model.
 *
 * Request:  (role, path, method)
 * Policy:   (role, path-pattern, method-pattern)
 * `g`:      role inheritance (e.g. admin inherits operator inherits viewer)
 * Matcher:  role match (with inheritance) + keyMatch2 path glob + regex method.
 */
const MODEL = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
`;

/**
 * Access policy. `keyMatch2` treats `:param` and `*` as wildcards, so
 * `/dashboard/*` matches `/dashboard/cameras`, etc.
 *
 * Role hierarchy: admin ⊃ operator ⊃ viewer.
 */
const POLICY = `
p, viewer, /dashboard, GET
p, viewer, /dashboard/cameras, GET
p, viewer, /dashboard/analytics, GET
p, viewer, /api/health, GET
p, viewer, /api/cameras, GET
p, viewer, /api/cameras/*, (GET|POST)
p, operator, /dashboard/*, (GET|POST)
p, admin, /*, (GET|POST|PUT|PATCH|DELETE)

g, operator, viewer
g, admin, operator
`;

/** @type {Promise<import('casbin').Enforcer> | null} */
let enforcerPromise = null;

function getEnforcer() {
	if (!enforcerPromise) {
		enforcerPromise = newEnforcer(newModelFromString(MODEL), new StringAdapter(POLICY));
	}
	return enforcerPromise;
}

/**
 * Returns true if `role` may perform `method` on `path`.
 * @param {string} role
 * @param {string} path
 * @param {string} method
 * @returns {Promise<boolean>}
 */
export async function authorize(role, path, method) {
	const enforcer = await getEnforcer();
	return enforcer.enforce(role, path, method.toUpperCase());
}
