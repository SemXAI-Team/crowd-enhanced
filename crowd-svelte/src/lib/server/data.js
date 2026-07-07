import { fetchCameras } from '$lib/api';

/**
 * Data adapter for the dashboard.
 *
 * All dashboard pages read through this module, never a transport directly —
 * so when the backend/runtime changes (new gRPC schema, REST, NATS, …), only
 * this file needs rewiring. Today: camera inventory comes from the Go gRPC
 * backend; live density/count metrics are mocked until the inference backend
 * publishes them.
 *
 * @typedef {'low' | 'medium' | 'high' | 'critical'} Density
 * @typedef {{
 *   id: string, name: string, location: string, online: boolean,
 *   density: Density, currentCount: number, maxCapacity: number
 * }} DashboardCamera
 */

/** @type {Density[]} */
const DENSITY_LEVELS = ['low', 'medium', 'high', 'critical'];

/** Deterministic per-id mock so values are stable across a session. */
function mockMetrics(/** @type {string} */ id) {
	let hash = 0;
	for (const ch of id) hash = (hash * 31 + ch.charCodeAt(0)) >>> 0;
	const density = DENSITY_LEVELS[hash % DENSITY_LEVELS.length];
	const maxCapacity = 200 + (hash % 6) * 100;
	const fill = { low: 0.2, medium: 0.5, high: 0.75, critical: 0.93 }[density];
	return { density, maxCapacity, currentCount: Math.round(maxCapacity * fill) };
}

/**
 * Camera inventory + live metrics for the dashboard. Degrades to an empty
 * list if the streaming backend is unreachable so pages still render.
 * @returns {Promise<DashboardCamera[]>}
 */
export async function getCameras() {
	try {
		const cameras = await fetchCameras();
		return cameras.map((c) => ({ ...c, ...mockMetrics(c.id) }));
	} catch {
		return [];
	}
}

/**
 * Lightweight backend health probe (used by the sidebar status indicator).
 * @returns {Promise<{ online: boolean, latencyMs: number }>}
 */
export async function getBackendHealth() {
	const start = Date.now();
	try {
		await fetchCameras();
		return { online: true, latencyMs: Date.now() - start };
	} catch {
		return { online: false, latencyMs: Date.now() - start };
	}
}
