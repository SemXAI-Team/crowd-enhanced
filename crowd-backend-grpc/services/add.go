package services

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "crowd-backend-grpc/proto"
)

// Capitalize the struct name so it can be exported and accessed by main.go
type CalculatorServer struct {
	pb.UnimplementedCalculatorServiceServer
}

func (s *CalculatorServer) Operation(ctx context.Context, req *pb.BinaryRequest) (*pb.BinaryResponse, error) {
	switch req.Operation {
	case pb.BinaryRequest_ADD:
		return &pb.BinaryResponse{Result: req.A + req.B}, nil
	case pb.BinaryRequest_SUBTRACT:
		return &pb.BinaryResponse{Result: req.A - req.B}, nil
	case pb.BinaryRequest_MULTIPLY:
		return &pb.BinaryResponse{Result: req.A * req.B}, nil
	case pb.BinaryRequest_DIVIDE:
		return &pb.BinaryResponse{Result: req.A / req.B}, nil
	case pb.BinaryRequest_MODULO:
		return &pb.BinaryResponse{Result: req.A % req.B}, nil
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unknown operation: %v", req.Operation)
	}

}
