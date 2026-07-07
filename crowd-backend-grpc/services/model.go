package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "crowd-backend-grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ModelServer struct {
	pb.UnimplementedModelManagerServer
	mu     sync.Mutex
	models map[string]*pb.ModelInfo
}

func NewModelServer() *ModelServer {
	return &ModelServer{
		models: map[string]*pb.ModelInfo{
			"1": {Id: "1", Name: "YOLOv8", Description: "Default Object Detector", Builtin: true},
		},
	}
}

func (s *ModelServer) ListModels(ctx context.Context, req *pb.ListModelsRequest) (*pb.ListModelsResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	resp := &pb.ListModelsResponse{}
	for _, m := range s.models {
		resp.Models = append(resp.Models, m)
	}
	return resp, nil
}

func (s *ModelServer) UploadModel(ctx context.Context, req *pb.UploadModelRequest) (*pb.ModelInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("model-%d", time.Now().UnixNano())
	m := &pb.ModelInfo{
		Id:          id,
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Builtin:     false,
	}
	s.models[id] = m
	return m, nil
}

func (s *ModelServer) DeleteModel(ctx context.Context, req *pb.DeleteModelRequest) (*pb.DeleteModelResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	m, ok := s.models[req.GetModelId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "model not found")
	}
	if m.Builtin {
		return nil, status.Errorf(codes.InvalidArgument, "cannot delete builtin model")
	}
	delete(s.models, req.GetModelId())
	return &pb.DeleteModelResponse{Success: true}, nil
}
