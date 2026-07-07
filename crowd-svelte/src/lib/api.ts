export const API_BASE = 'http://localhost:8080/api/v1';

let authToken = '';

export function setAuthToken(token: string) {
    authToken = token;
}

function getHeaders() {
    const headers: Record<string, string> = {
        'Content-Type': 'application/json'
    };
    if (authToken) {
        headers['Authorization'] = `Bearer ${authToken}`;
    }
    return headers;
}

export async function fetchCameras() {
    const res = await fetch(`${API_BASE}/cameras`, { headers: getHeaders() });
    if (!res.ok) throw new Error('Failed to fetch cameras');
    const data = await res.json();
    return data.cameras || [];
}

export async function negotiateWebRTC(cameraId: string, sdpOffer: string) {
    const res = await fetch(`${API_BASE}/webrtc/negotiate`, {
        method: 'POST',
        headers: getHeaders(),
        body: JSON.stringify({ camera_id: cameraId, sdp_offer: sdpOffer })
    });
    if (!res.ok) throw new Error('Failed to negotiate WebRTC');
    const data = await res.json();
    return data.sdp_answer;
}

export function getTelemetryStreamUrl(cameraId: string) {
    return `${API_BASE}/telemetry/stream/${cameraId}`;
}

// For use with the EventSource or ReadableStream
export async function streamTelemetry(cameraId: string, onSample: (sample: any) => void, signal?: AbortSignal) {
    const res = await fetch(getTelemetryStreamUrl(cameraId), {
        headers: getHeaders(),
        signal
    });
    
    if (!res.ok) {
        throw new Error('Failed to connect to telemetry stream');
    }
    
    if (!res.body) return;
    
    const reader = res.body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';
    
    try {
        while (true) {
            const { done, value } = await reader.read();
            if (done) break;
            
            buffer += decoder.decode(value, { stream: true });
            
            // grpc-gateway streams are newline delimited JSON
            const lines = buffer.split('\n');
            buffer = lines.pop() || '';
            
            for (const line of lines) {
                if (line.trim()) {
                    try {
                        const sample = JSON.parse(line);
                        // grpc-gateway wraps the message in {"result": {...}}
                        if (sample.result) {
                            onSample(sample.result);
                        } else {
                            onSample(sample);
                        }
                    } catch (e) {
                        console.error('Failed to parse streaming line:', line);
                    }
                }
            }
        }
    } finally {
        reader.releaseLock();
    }
}

export async function uploadModel(name: string, description: string, packageBase64: string) {
    const res = await fetch(`${API_BASE}/models/upload`, {
        method: 'POST',
        headers: getHeaders(),
        body: JSON.stringify({
            name,
            description,
            package_base64: packageBase64
        })
    });
    if (!res.ok) {
        throw new Error(await res.text());
    }
    return res.json();
}

export async function getGlobalStats() {
    const res = await fetch(`${API_BASE}/analytics`, { headers: getHeaders() });
    if (!res.ok) throw new Error('Failed to fetch stats');
    return res.json();
}

export async function listUsers() {
    const res = await fetch(`${API_BASE}/iam/users`, { headers: getHeaders() });
    if (!res.ok) throw new Error('Failed to fetch users');
    const data = await res.json();
    return data.users || [];
}

export async function updateUserRole(userId: string, role: string) {
    const res = await fetch(`${API_BASE}/iam/users/${userId}/role`, {
        method: 'PUT',
        headers: getHeaders(),
        body: JSON.stringify({ user_id: userId, role })
    });
    if (!res.ok) throw new Error('Failed to update role');
    return res.json();
}

export async function revokeAccess(userId: string) {
    const res = await fetch(`${API_BASE}/iam/users/${userId}/access`, {
        method: 'DELETE',
        headers: getHeaders()
    });
    if (!res.ok) throw new Error('Failed to revoke access');
    return res.json();
}

export async function deleteUser(userId: string) {
    const res = await fetch(`${API_BASE}/iam/users/${userId}`, {
        method: 'DELETE',
        headers: getHeaders()
    });
    if (!res.ok) throw new Error('Failed to delete user');
    return res.json();
}

export async function listCameras() {
    const res = await fetch(`${API_BASE}/cameras`, {
        headers: getHeaders()
    });
    if (!res.ok) throw new Error('Failed to fetch cameras');
    const data = await res.json();
    return data.cameras || [];
}

export async function registerCamera(name: string, location: string, source: string, modelId: string = '') {
    const res = await fetch(`${API_BASE}/cameras`, {
        method: 'POST',
        headers: getHeaders(),
        body: JSON.stringify({ name, location, source, model_id: modelId })
    });
    if (!res.ok) throw new Error('Failed to register camera');
    return res.json();
}

export async function updateCamera(cameraId: string, name: string, location: string, source: string, modelId: string = '') {
    const res = await fetch(`${API_BASE}/cameras/${cameraId}`, {
        method: 'PUT',
        headers: getHeaders(),
        body: JSON.stringify({ camera_id: cameraId, name, location, source, model_id: modelId })
    });
    if (!res.ok) throw new Error('Failed to update camera');
    return res.json();
}

export async function deleteCamera(cameraId: string) {
    const res = await fetch(`${API_BASE}/cameras/${cameraId}`, {
        method: 'DELETE',
        headers: getHeaders()
    });
    if (!res.ok) throw new Error('Failed to delete camera');
    return res.json();
}

export async function deleteModel(modelId: string) {
    const res = await fetch(`${API_BASE}/models/${modelId}`, {
        method: 'DELETE',
        headers: getHeaders()
    });
    if (!res.ok) throw new Error('Failed to delete model');
    return res.json();
}

export async function listModels() {
    const res = await fetch(`${API_BASE}/models`, {
        headers: getHeaders()
    });
    if (!res.ok) throw new Error('Failed to fetch models');
    const data = await res.json();
    return data.models || [];
}
