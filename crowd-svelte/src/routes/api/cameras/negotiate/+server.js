import { json, error } from '@sveltejs/kit';
import { negotiate } from '$lib/server/camera-grpc';

/**
 * WebRTC signalling proxy: the browser POSTs its SDP offer here, we relay it to
 * the Go/DeepStream backend over gRPC and return the SDP answer. The request is
 * already authenticated + Casbin-authorized by hooks.server.js before it runs.
 * @type {import('./$types').RequestHandler}
 */
export const POST = async ({ request }) => {
	const { cameraId, sdp } = await request.json();
	if (!cameraId || !sdp) {
		throw error(400, 'cameraId and sdp are required');
	}
	try {
		const answer = await negotiate(cameraId, sdp);
		return json({ sdp: answer });
	} catch (/** @type {any} */ err) {
		throw error(502, `camera backend: ${err?.details ?? err?.message ?? 'unavailable'}`);
	}
};
