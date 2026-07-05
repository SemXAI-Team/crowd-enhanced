package services

import (
	"context"
	"log"

	pb "crowd-backend-grpc/proto"
)

type GreeterServer struct {
	pb.UnimplementedGreeterServiceServer
}

func (s *GreeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("[GreeterService] Request from: %s", req.GetName())
	greeting := "Hello, " + req.GetName()
	return &pb.HelloResponse{Greeting: greeting}, nil
}
