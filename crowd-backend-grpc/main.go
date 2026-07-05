package main

import (
	"log"
	"net"

	pb "crowd-backend-grpc/proto"
	"crowd-backend-grpc/services" // Import your new modularized package

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var proto_method = "tcp"
var port = ":5005"

func main() {
	listener, err := net.Listen(proto_method, port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register services by instantiating the structs exported from the services package
	pb.RegisterGreeterServiceServer(grpcServer, &services.GreeterServer{})
	pb.RegisterAddServiceServer(grpcServer, &services.AddServer{})

	// Enable reflection for easy debugging (Postman / grpcurl)
	reflection.Register(grpcServer)

	log.Println("Server successfully running on port ", port	)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
