package services

import (
	"context"
	"log"

	pb "crowd-backend-grpc/proto"
)

// Capitalize the struct name so it can be exported and accessed by main.go
type AddServer struct {
	pb.UnimplementedAddServiceServer
}

func (s *AddServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	log.Printf("[AddService] Processing: %d + %d", req.GetA(), req.GetB())
	result := req.GetA() + req.GetB()
	return &pb.AddResponse{Result: result}, nil
}
