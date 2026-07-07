package services

import (
	"context"
	pb "crowd-backend-grpc/proto"
)

type AnalyticsServer struct {
	pb.UnimplementedAnalyticsServiceServer
}

func NewAnalyticsServer() *AnalyticsServer {
	return &AnalyticsServer{}
}

func (s *AnalyticsServer) GetGlobalStats(ctx context.Context, req *pb.GetGlobalStatsRequest) (*pb.GlobalStatsResponse, error) {
	// Mock implementation
	return &pb.GlobalStatsResponse{
		ActiveCameras:   12,
		TotalCrowdCount: 450,
		AverageLatency:  35.5,
		ActiveAlerts:    2,
	}, nil
}
