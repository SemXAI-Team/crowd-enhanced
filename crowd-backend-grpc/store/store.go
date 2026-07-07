// Package store persists the camera registry and model catalog.
//
// It ports the legacy backend's Mongo `cameras` collection semantics
// (auto-assigned "cam-N" ids, source URIs, per-camera model binding) onto a
// simple JSON file with atomic writes — no external database required for a
// registry of this size. Swap this package if a real DB becomes necessary.
package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Camera is a registered video source.
type Camera struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	// Input URI for the DeepStream pipeline: rtsp://…, http(s)://…, file://…
	Source string `json:"source"`
	// Model package id used for inference ("" = passthrough, no inference).
	ModelID string `json:"model_id"`
	// Roles allowed to view this camera. Empty means every authenticated role.
	AllowedRoles []string `json:"allowed_roles,omitempty"`
	// Local UDP port the camera's pipeline sends RTP/H264 to (assigned once).
	RTPPort int `json:"rtp_port"`
}

// Allows reports whether a role may view this camera.
func (c *Camera) Allows(role string) bool {
	if len(c.AllowedRoles) == 0 || role == "admin" {
		return true
	}
	for _, r := range c.AllowedRoles {
		if strings.EqualFold(r, role) {
			return true
		}
	}
	return false
}

const rtpPortBase = 5100

// Store is a mutex-guarded, JSON-file-backed camera registry.
type Store struct {
	mu   sync.RWMutex
	path string
	cams map[string]*Camera
}

// Open loads (or initialises) the registry at path.
func Open(path string) (*Store, error) {
	s := &Store{path: path, cams: map[string]*Camera{}}

	data, err := os.ReadFile(path)
	switch {
	case err == nil:
		var cams []*Camera
		if err := json.Unmarshal(data, &cams); err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		for _, c := range cams {
			s.cams[c.ID] = c
		}
	case os.IsNotExist(err):
		// Seed defaults (parity with the legacy backend seeding when empty).
		s.cams["cam-1"] = &Camera{ID: "cam-1", Name: "North Gate", Location: "Entrance A",
			Source: "rtsp://127.0.0.1:8554/cam1", RTPPort: rtpPortBase}
		s.cams["cam-2"] = &Camera{ID: "cam-2", Name: "Main Hall", Location: "Concourse",
			Source: "rtsp://127.0.0.1:8554/cam2", RTPPort: rtpPortBase + 1}
		if err := s.flushLocked(); err != nil {
			return nil, err
		}
	default:
		return nil, err
	}
	return s, nil
}

// List returns cameras sorted by numeric id suffix (cam-1, cam-2, …).
func (s *Store) List() []*Camera {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Camera, 0, len(s.cams))
	for _, c := range s.cams {
		cp := *c
		out = append(out, &cp)
	}
	sort.Slice(out, func(i, j int) bool { return camNum(out[i].ID) < camNum(out[j].ID) })
	return out
}

// Get returns a copy of the camera, or nil.
func (s *Store) Get(id string) *Camera {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if c, ok := s.cams[id]; ok {
		cp := *c
		return &cp
	}
	return nil
}

// Register adds a camera, assigning the next "cam-N" id and a free RTP port
// (port of the legacy add_camera_logic).
func (s *Store) Register(c Camera) (*Camera, error) {
	if err := validate(c); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	maxID, maxPort := 0, rtpPortBase-1
	for _, existing := range s.cams {
		if n := camNum(existing.ID); n > maxID {
			maxID = n
		}
		if existing.RTPPort > maxPort {
			maxPort = existing.RTPPort
		}
	}
	c.ID = fmt.Sprintf("cam-%d", maxID+1)
	c.RTPPort = maxPort + 1

	s.cams[c.ID] = &c
	if err := s.flushLocked(); err != nil {
		delete(s.cams, c.ID)
		return nil, err
	}
	cp := c
	return &cp, nil
}

// Update mutates name/location/source/model/roles; id and RTP port are stable.
func (s *Store) Update(id string, c Camera) (*Camera, error) {
	if err := validate(c); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.cams[id]
	if !ok {
		return nil, os.ErrNotExist
	}
	prev := *existing
	existing.Name, existing.Location = c.Name, c.Location
	existing.Source, existing.ModelID = c.Source, c.ModelID
	existing.AllowedRoles = c.AllowedRoles
	if err := s.flushLocked(); err != nil {
		*existing = prev
		return nil, err
	}
	cp := *existing
	return &cp, nil
}

// Delete removes a camera. Returns false if it did not exist.
func (s *Store) Delete(id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.cams[id]
	if !ok {
		return false, nil
	}
	delete(s.cams, id)
	if err := s.flushLocked(); err != nil {
		s.cams[id] = existing
		return false, err
	}
	return true, nil
}

func validate(c Camera) error {
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("camera name is required")
	}
	src := strings.ToLower(strings.TrimSpace(c.Source))
	if src == "" {
		return fmt.Errorf("camera source is required")
	}
	for _, prefix := range []string{"rtsp://", "http://", "https://", "file://"} {
		if strings.HasPrefix(src, prefix) {
			return nil
		}
	}
	return fmt.Errorf("source must be an rtsp://, http(s):// or file:// URI")
}

func camNum(id string) int {
	n, _ := strconv.Atoi(strings.TrimPrefix(id, "cam-"))
	return n
}

// flushLocked atomically persists the registry. Callers hold s.mu.
func (s *Store) flushLocked() error {
	cams := make([]*Camera, 0, len(s.cams))
	for _, c := range s.cams {
		cams = append(cams, c)
	}
	sort.Slice(cams, func(i, j int) bool { return camNum(cams[i].ID) < camNum(cams[j].ID) })

	data, err := json.MarshalIndent(cams, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
