import os from 'node:os';
import fs from 'node:fs';
import path from 'node:path';
import grpc from '@grpc/grpc-js';
import protoLoader from '@grpc/proto-loader';

// Inline the proto text (via Vite `?raw`) and materialise it to a temp file so
// @grpc/proto-loader can read it — works in dev and the bundled prod server.
import protoSource from './proto/camera_stream.proto?raw';

const BACKEND = process.env.GRPC_BACKEND_URL || '127.0.0.1:5005';

const protoFile = path.join(os.tmpdir(), 'crowd-camera_stream.proto');
fs.writeFileSync(protoFile, protoSource);

const packageDef = protoLoader.loadSync(protoFile, {
	keepCase: true,
	longs: String,
	enums: String,
	defaults: true,
	oneofs: true
});
const CameraStream = grpc.loadPackageDefinition(packageDef).camera.CameraStream;

/** @type {import('@grpc/grpc-js').Client | undefined} */
let client;
function getClient() {
	if (!client) {
		client = new CameraStream(BACKEND, grpc.credentials.createInsecure());
	}
	return client;
}

/**
 * @returns {Promise<Array<{id:string,name:string,location:string,online:boolean}>>}
 */
export function listCameras() {
	return new Promise((resolve, reject) => {
		getClient().ListCameras({}, (/** @type {any} */ err, /** @type {any} */ res) => {
			if (err) reject(err);
			else resolve(res.cameras ?? []);
		});
	});
}

/**
 * Relay the browser's SDP offer to the backend and return its SDP answer.
 * @param {string} cameraId
 * @param {string} sdpOffer
 * @returns {Promise<string>}
 */
export function negotiate(cameraId, sdpOffer) {
	return new Promise((resolve, reject) => {
		getClient().Negotiate(
			{ camera_id: cameraId, sdp_offer: sdpOffer },
			(/** @type {any} */ err, /** @type {any} */ res) => {
				if (err) reject(err);
				else resolve(res.sdp_answer);
			}
		);
	});
}
