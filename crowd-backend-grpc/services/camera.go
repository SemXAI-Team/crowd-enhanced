package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	pb "crowd-backend-grpc/proto"

	"github.com/pion/webrtc/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// cameraConfig describes a camera and how DeepStream exposes its footage.
type cameraConfig struct {
	ID       string
	Name     string
	Location string
	Source   string // input URI (e.g. rtsp://…) fed to the DeepStream pipeline
	RTPPort  int    // local UDP port the pipeline sends encoded RTP/H264 to
}

// CameraServer implements the CameraStream gRPC service. gRPC only carries the
// WebRTC signalling; the H.264 video produced by DeepStream flows peer-to-peer
// over WebRTC straight to the browser.
type CameraServer struct {
	pb.UnimplementedCameraStreamServer

	db      *sql.DB
	order   []string
	cameras map[string]*cameraConfig

	mu      sync.Mutex
	sources map[string]*cameraSource // one DeepStream pipeline per camera, lazily started
	scaler  *AutoScaler              // manages GPU batching across all camera sources
}

func NewCameraServer(db *sql.DB) *CameraServer {
	s := &CameraServer{
		db:      db,
		cameras: make(map[string]*cameraConfig),
		sources: make(map[string]*cameraSource),
		scaler:  NewAutoScaler(),
	}
	s.initDB()
	s.loadCameras()
	return s
}

func (s *CameraServer) initDB() {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS cameras (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			location TEXT NOT NULL,
			source TEXT NOT NULL,
			rtp_port INTEGER NOT NULL
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create cameras table: %v", err)
	}
}

func (s *CameraServer) loadCameras() {
	rows, err := s.db.Query("SELECT id, name, location, source, rtp_port FROM cameras")
	if err != nil {
		log.Printf("Failed to load cameras: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var c cameraConfig
		if err := rows.Scan(&c.ID, &c.Name, &c.Location, &c.Source, &c.RTPPort); err == nil {
			s.cameras[c.ID] = &c
			s.order = append(s.order, c.ID)
		}
	}
}

func (s *CameraServer) ListCameras(ctx context.Context, _ *pb.ListCamerasRequest) (*pb.ListCamerasResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	resp := &pb.ListCamerasResponse{}
	for _, id := range s.order {
		c := s.cameras[id]
		resp.Cameras = append(resp.Cameras, &pb.Camera{
			Id: c.ID, Name: c.Name, Location: c.Location, Online: true,
		})
	}
	return resp, nil
}

func (s *CameraServer) GetCamera(ctx context.Context, req *pb.GetCameraRequest) (*pb.Camera, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.cameras[req.GetCameraId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "camera not found")
	}
	return &pb.Camera{
		Id: c.ID, Name: c.Name, Location: c.Location, Online: true,
	}, nil
}

func (s *CameraServer) RegisterCamera(ctx context.Context, req *pb.RegisterCameraRequest) (*pb.Camera, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Simple ID generation for MVP
	id := fmt.Sprintf("cam-%d", time.Now().UnixNano())
	// Find next available RTP port (start from 5100)
	port := 5100
	for {
		used := false
		for _, c := range s.cameras {
			if c.RTPPort == port {
				used = true
				break
			}
		}
		if !used {
			break
		}
		port++
	}

	c := &cameraConfig{
		ID:       id,
		Name:     req.GetName(),
		Location: req.GetLocation(),
		Source:   req.GetSource(),
		RTPPort:  port,
	}

	_, err := s.db.ExecContext(ctx, "INSERT INTO cameras (id, name, location, source, rtp_port) VALUES ($1, $2, $3, $4, $5)", c.ID, c.Name, c.Location, c.Source, c.RTPPort)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert camera: %v", err)
	}

	s.cameras[c.ID] = c
	s.order = append(s.order, c.ID)

	return &pb.Camera{
		Id: c.ID, Name: c.Name, Location: c.Location, Online: true,
	}, nil
}

func (s *CameraServer) UpdateCamera(ctx context.Context, req *pb.UpdateCameraRequest) (*pb.Camera, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.cameras[req.GetCameraId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "camera not found")
	}

	_, err := s.db.ExecContext(ctx, "UPDATE cameras SET name = $1, location = $2, source = $3 WHERE id = $4", req.GetName(), req.GetLocation(), req.GetSource(), c.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update camera: %v", err)
	}

	c.Name = req.GetName()
	c.Location = req.GetLocation()
	c.Source = req.GetSource()

	return &pb.Camera{
		Id: c.ID, Name: c.Name, Location: c.Location, Online: true,
	}, nil
}

func (s *CameraServer) DeleteCamera(ctx context.Context, req *pb.DeleteCameraRequest) (*pb.DeleteCameraResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.cameras[req.GetCameraId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "camera not found")
	}

	_, err := s.db.ExecContext(ctx, "DELETE FROM cameras WHERE id = $1", c.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete camera: %v", err)
	}

	delete(s.cameras, c.ID)
	// Remove from order slice
	for i, id := range s.order {
		if id == c.ID {
			s.order = append(s.order[:i], s.order[i+1:]...)
			break
		}
	}

	// Stop source if running
	if src, ok := s.sources[c.ID]; ok {
		if src.cmd != nil {
			src.cmd.Process.Kill()
		}
		if src.conn != nil {
			src.conn.Close()
		}
		delete(s.sources, c.ID)
		
		// Detach from the shared AutoScaler model instance
		s.scaler.RemoveStream(c.ID, "default-model-id")
	}

	return &pb.DeleteCameraResponse{Success: true}, nil
}

func (s *CameraServer) Negotiate(ctx context.Context, req *pb.NegotiateRequest) (*pb.NegotiateResponse, error) {
	cfg, ok := s.cameras[req.GetCameraId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "unknown camera %q", req.GetCameraId())
	}
	if strings.TrimSpace(req.GetSdpOffer()) == "" {
		return nil, status.Error(codes.InvalidArgument, "sdp_offer is required")
	}

	source, err := s.getOrStartSource(cfg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "start camera source: %v", err)
	}

	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{{URLs: []string{"stun:stun.l.google.com:19302"}}},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create peer connection: %v", err)
	}

	// Broadcast the camera's shared track to this viewer.
	if _, err := pc.AddTrack(source.track); err != nil {
		_ = pc.Close()
		return nil, status.Errorf(codes.Internal, "add track: %v", err)
	}

	// Tear the connection down once it drops so we don't leak peer connections.
	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		log.Printf("[CameraStream] %s peer state: %s", cfg.ID, state)
		if state == webrtc.PeerConnectionStateFailed ||
			state == webrtc.PeerConnectionStateClosed ||
			state == webrtc.PeerConnectionStateDisconnected {
			_ = pc.Close()
		}
	})

	if err := pc.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer, SDP: req.GetSdpOffer(),
	}); err != nil {
		_ = pc.Close()
		return nil, status.Errorf(codes.InvalidArgument, "set remote description: %v", err)
	}

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		_ = pc.Close()
		return nil, status.Errorf(codes.Internal, "create answer: %v", err)
	}

	// Non-trickle ICE: gather all candidates before replying so the answer SDP
	// is self-contained and a single round-trip completes the handshake.
	gatherComplete := webrtc.GatheringCompletePromise(pc)
	if err := pc.SetLocalDescription(answer); err != nil {
		_ = pc.Close()
		return nil, status.Errorf(codes.Internal, "set local description: %v", err)
	}
	select {
	case <-gatherComplete:
	case <-ctx.Done():
		_ = pc.Close()
		return nil, status.FromContextError(ctx.Err()).Err()
	}

	return &pb.NegotiateResponse{SdpAnswer: pc.LocalDescription().SDP}, nil
}

func (s *CameraServer) StreamTelemetry(req *pb.TelemetryRequest, stream pb.CameraStream_StreamTelemetryServer) error {
	s.mu.Lock()
	cfg, ok := s.cameras[req.GetCameraId()]
	s.mu.Unlock()
	if !ok {
		return status.Errorf(codes.NotFound, "unknown camera %q", req.GetCameraId())
	}

	// This is a placeholder mock for the DeepStream Telemetry Connector.
	// In reality, this loop would block on an MQTT/Redis channel receiving JSON metrics for this cfg.ID.
	for i := 0; i < 100; i++ {
		sample := &pb.TelemetrySample{
			CameraId:  cfg.ID,
			Count:     int32(10 + (i % 5)),
			LatencyMs: 45.2,
			Metrics: []*pb.Metric{
				{Label: "CROWD COUNT", Value: fmt.Sprintf("%d", 10+(i%5)), Trend: "up", Color: "blue"},
				{Label: "DENSITY", Value: "High", Trend: "neutral", Color: "red"},
			},
		}

		if err := stream.Send(sample); err != nil {
			return err
		}

		// Simulate 1 Hz telemetry
		select {
		case <-stream.Context().Done():
			return nil
		case <-time.After(1 * time.Second):
		}
	}

	return nil
}

// getOrStartSource returns the camera's running source, starting its DeepStream
// pipeline + RTP reader on first use.
func (s *CameraServer) getOrStartSource(cfg *cameraConfig) (*cameraSource, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if src, ok := s.sources[cfg.ID]; ok {
		return src, nil
	}

	// 1. Assign this camera stream to an optimal Model Instance (GPU batching algorithm)
	modelID := "default-model-id" // In a real system, this comes from camera config
	instance, err := s.scaler.AssignStream(cfg.ID, modelID)
	if err != nil {
		return nil, err
	}
	log.Printf("[CameraServer] Camera %s attached to Model Instance %s", cfg.ID, instance.ID)

	// 2. Start the physical pipeline (simulated here per-camera for WebRTC viewing)
	src, err := newCameraSource(cfg)
	if err != nil {
		s.scaler.RemoveStream(cfg.ID, modelID)
		return nil, err
	}
	s.sources[cfg.ID] = src
	return src, nil
}

// cameraSource runs one DeepStream pipeline and fans its RTP packets out to
// every connected viewer via a single shared WebRTC track.
type cameraSource struct {
	cfg   *cameraConfig
	conn  *net.UDPConn
	cmd   *exec.Cmd
	track *webrtc.TrackLocalStaticRTP
}

func newCameraSource(cfg *cameraConfig) (*cameraSource, error) {
	track, err := webrtc.NewTrackLocalStaticRTP(
		webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264},
		"video", "camera-"+cfg.ID,
	)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: cfg.RTPPort})
	if err != nil {
		return nil, fmt.Errorf("listen udp :%d: %w", cfg.RTPPort, err)
	}
	_ = conn.SetReadBuffer(1 << 20)

	src := &cameraSource{cfg: cfg, conn: conn, track: track}
	if err := src.startPipeline(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	go src.readLoop()
	log.Printf("[CameraStream] source %q started (RTP on udp/%d)", cfg.ID, cfg.RTPPort)
	return src, nil
}

// readLoop copies RTP datagrams from the DeepStream udpsink onto the shared
// track, which forwards them to all connected peers.
func (s *cameraSource) readLoop() {
	buf := make([]byte, 1500)
	for {
		n, _, err := s.conn.ReadFrom(buf)
		if err != nil {
			log.Printf("[CameraStream] source %q read loop stopped: %v", s.cfg.ID, err)
			return
		}
		if _, err := s.track.Write(buf[:n]); err != nil && !errors.Is(err, io.ErrClosedPipe) {
			log.Printf("[CameraStream] source %q track write: %v", s.cfg.ID, err)
		}
	}
}

// startPipeline launches the GStreamer/DeepStream pipeline that decodes the
// camera, runs inference, re-encodes to H.264 and sends it as RTP to our UDP
// port. Set CAMERA_TEST_SOURCE=1 to use a synthetic source with no GPU.
func (s *cameraSource) startPipeline() error {
	var pipeline string
	if os.Getenv("CAMERA_TEST_SOURCE") == "1" {
		pipeline = fmt.Sprintf(
			"videotestsrc is-live=true pattern=ball ! video/x-raw,width=640,height=480,framerate=30/1 "+
				"! x264enc tune=zerolatency bitrate=1500 speed-preset=ultrafast key-int-max=30 "+
				"! rtph264pay config-interval=1 pt=96 ! udpsink host=127.0.0.1 port=%d",
			s.cfg.RTPPort)
	} else {
		// DeepStream integration point — adapt nvinfer's config-file-path to your
		// model and add nvstreammux/nvtracker as your deployment requires.
		config := os.Getenv("DEEPSTREAM_CONFIG")
		if config == "" {
			config = "config_infer_primary.txt"
		}

		// The pipeline uses a `tee` to split the processed stream:
		// 1. WebRTC Branch: sends encoded RTP to our local UDP port (for the CrowdGuard Dashboard)
		// 2. Legacy VMS Branch: sends standard RTSP out to a media server for existing NVRs/VMS integration.
		pipeline = fmt.Sprintf(
			"uridecodebin uri=%s ! nvvideoconvert ! nvinfer config-file-path=%s ! nvvideoconvert "+
				"! nvv4l2h264enc insert-sps-pps=1 idrinterval=30 ! rtph264pay config-interval=1 pt=96 "+
				"! tee name=t "+
				"t. ! queue ! udpsink host=127.0.0.1 port=%d "+
				"t. ! queue ! rtspclientsink location=rtsp://localhost:8554/processed/%s protocols=tcp",
			s.cfg.Source, config, s.cfg.RTPPort, s.cfg.ID)
	}

	cmd := exec.Command("gst-launch-1.0", append([]string{"-e"}, strings.Fields(pipeline)...)...)
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start gstreamer pipeline: %w", err)
	}
	s.cmd = cmd
	go func() {
		if err := cmd.Wait(); err != nil {
			log.Printf("[CameraStream] source %q pipeline exited: %v", s.cfg.ID, err)
		}
	}()
	return nil
}
