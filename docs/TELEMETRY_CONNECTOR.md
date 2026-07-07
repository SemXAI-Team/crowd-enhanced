# Telemetry Connector Guide

This guide explains how AI Engineers can pipe custom inference data from their DeepStream models to the CrowdGuard dashboard.

## Architecture

The dashboard supports generic, real-time telemetry rendering. You are not limited to just "crowd count". You can send any metric or state, and the frontend will automatically render it as a StatCard or chart.

1. **DeepStream Parsing:** Your pipeline extracts metadata (`NvDsMeta`) into JSON using the `nvmsgconv` plugin.
2. **Broker Publishing:** You use `nvmsgbroker` to publish this JSON to the backend's ingestion endpoint (or MQTT broker).
3. **Frontend Rendering:** The SvelteKit dashboard consumes the JSON via WebSockets and renders it dynamically.

## Connecting to the REST Ingestion Endpoint

If you are not using a dedicated MQTT broker, you can send HTTP POST requests directly to the Go backend's REST connector:

**Endpoint:** `POST http://localhost:8080/api/v1/telemetry/publish`  
**Content-Type:** `application/json`

### JSON Schema

The payload must be a JSON object containing the `camera_id` and a `data` object with your arbitrary key-value pairs.

```json
{
  "camera_id": "cam-1",
  "data": {
    "crowd_count": 142,
    "anomalies_detected": 0,
    "avg_dwell_time_seconds": 45.2,
    "zone_density": "high"
  }
}
```

The SvelteKit `TelemetryWidget` will dynamically loop through the keys in the `data` object and render a card for each metric.

### DeepStream `nvmsgconv` Integration

To automate this within your GStreamer pipeline without writing a custom Python/C++ HTTP sender:

1. Configure `nvmsgconv` to use the standard DeepStream schema.
2. Configure `nvmsgbroker` with an HTTP/MQTT adapter pointing to our backend.

*Example `msgbroker` config:*
```txt
[message-broker]
enable=1
msg-broker-proto-lib=/opt/nvidia/deepstream/deepstream/lib/libnvds_kafka_proto.so
msg-broker-conn-str=localhost;9092;telemetry_topic
```
*(Note: If using Kafka/MQTT, ensure the Go backend is updated to subscribe to that specific broker. The REST endpoint is available for direct script injection).*
