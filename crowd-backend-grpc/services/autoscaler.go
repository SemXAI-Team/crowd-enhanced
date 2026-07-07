package services

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// MaxStreamsPerInstance defines how many camera streams a single GPU model instance
// (e.g., DeepStream nvstreammux) can handle before we spawn a new instance.
const MaxStreamsPerInstance = 8

// ModelInstance represents a running AI pipeline (e.g., a TensorRT engine context / DeepStream pipeline).
type ModelInstance struct {
	ID            string
	ModelID       string
	ActiveStreams map[string]bool // Set of Camera IDs currently assigned
	// In a real implementation, this would hold the exec.Cmd or DeepStream process handle,
	// as well as the RTSP restreamer endpoints for legacy VMS integration.
}

// AutoScaler manages the allocation of camera streams to ModelInstances.
type AutoScaler struct {
	mu        sync.Mutex
	instances map[string]*ModelInstance
}

func NewAutoScaler() *AutoScaler {
	return &AutoScaler{
		instances: make(map[string]*ModelInstance),
	}
}

// AssignStream automatically assigns a camera to an available model instance.
// If all instances are at MaxCapacity, it scales up by spawning a new instance.
func (a *AutoScaler) AssignStream(cameraID string, modelID string) (*ModelInstance, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// 1. Check if the camera is already assigned to an instance
	for _, inst := range a.instances {
		if inst.ModelID == modelID && inst.ActiveStreams[cameraID] {
			return inst, nil // Already assigned
		}
	}

	// 2. Find an existing instance for this model that has capacity
	for _, inst := range a.instances {
		if inst.ModelID == modelID && len(inst.ActiveStreams) < MaxStreamsPerInstance {
			inst.ActiveStreams[cameraID] = true
			log.Printf("[AutoScaler] Assigned camera %s to existing instance %s (Load: %d/%d)", cameraID, inst.ID, len(inst.ActiveStreams), MaxStreamsPerInstance)
			return inst, nil
		}
	}

	// 3. No capacity found. Scale up!
	newInst := a.spawnInstance(modelID)
	newInst.ActiveStreams[cameraID] = true
	log.Printf("[AutoScaler] Scaled UP: Spawning new instance %s for model %s (Load: 1/%d)", newInst.ID, modelID, MaxStreamsPerInstance)
	return newInst, nil
}

// RemoveStream detaches a camera from its model instance.
// If the instance becomes empty, it triggers a graceful scale-down timeout.
func (a *AutoScaler) RemoveStream(cameraID string, modelID string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, inst := range a.instances {
		if inst.ModelID == modelID && inst.ActiveStreams[cameraID] {
			delete(inst.ActiveStreams, cameraID)
			log.Printf("[AutoScaler] Removed camera %s from instance %s (Load: %d/%d)", cameraID, inst.ID, len(inst.ActiveStreams), MaxStreamsPerInstance)
			
			if len(inst.ActiveStreams) == 0 {
				log.Printf("[AutoScaler] Instance %s is idle. Triggering scale down...", inst.ID)
				// In a real system, we'd start a timer. If no streams attach in 5 mins, kill the process.
				// For the prototype, we immediately scale down to free resources.
				a.killInstance(inst.ID)
			}
			return
		}
	}
}

// spawnInstance simulates launching a new DeepStream pipeline batch.
func (a *AutoScaler) spawnInstance(modelID string) *ModelInstance {
	id := fmt.Sprintf("inst-%s-%d", modelID, time.Now().UnixMilli())
	inst := &ModelInstance{
		ID:            id,
		ModelID:       modelID,
		ActiveStreams: make(map[string]bool),
	}
	a.instances[id] = inst
	return inst
}

// killInstance simulates terminating a DeepStream pipeline.
func (a *AutoScaler) killInstance(instanceID string) {
	delete(a.instances, instanceID)
	log.Printf("[AutoScaler] Scaled DOWN: Killed instance %s", instanceID)
}
