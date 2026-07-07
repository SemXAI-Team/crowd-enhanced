package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"

	pb "crowd-backend-grpc/proto"
	"crowd-backend-grpc/server"
	"crowd-backend-grpc/services"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var proto_method = "tcp"
var port = ":5005"

func main() {
	listener, err := net.Listen(proto_method, port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Connect to Postgres
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://crowd:crowd@localhost:5432/crowd?sslmode=disable"
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	// Start gRPC server with the Auth Interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(server.AuthInterceptor()),
	)

	camServer := services.NewCameraServer(db)
	iamServer := services.NewIAMServer(db)
	analyticsServer := services.NewAnalyticsServer()
	modelServer := services.NewModelServer()

	pb.RegisterGreeterServiceServer(grpcServer, &services.GreeterServer{})
	pb.RegisterCalculatorServiceServer(grpcServer, &services.CalculatorServer{})
	pb.RegisterCameraStreamServer(grpcServer, camServer)
	pb.RegisterIAMServiceServer(grpcServer, iamServer)
	pb.RegisterAnalyticsServiceServer(grpcServer, analyticsServer)
	pb.RegisterModelManagerServer(grpcServer, modelServer)

	// Enable reflection for easy debugging (Postman / grpcurl)
	reflection.Register(grpcServer)

	// Setup and run grpc-gateway proxy layer
	go func() {
		log.Println("grpc-gateway server successfully running on port :8080")
		
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		
		err := pb.RegisterCameraStreamHandlerFromEndpoint(ctx, mux, ":5005", opts)
		if err != nil {
			log.Fatalf("Failed to register CameraStream handler: %v", err)
		}
		
		err = pb.RegisterIAMServiceHandlerFromEndpoint(ctx, mux, ":5005", opts)
		if err != nil {
			log.Fatalf("Failed to register IAMService handler: %v", err)
		}

		err = pb.RegisterAnalyticsServiceHandlerFromEndpoint(ctx, mux, ":5005", opts)
		if err != nil {
			log.Fatalf("Failed to register AnalyticsService handler: %v", err)
		}
		
		err = pb.RegisterModelManagerHandlerFromEndpoint(ctx, mux, ":5005", opts)
		if err != nil {
			log.Fatalf("Failed to register ModelManager handler: %v", err)
		}

		// Simple CORS middleware
		corsMiddleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				if r.Method == "OPTIONS" {
					w.WriteHeader(http.StatusOK)
					return
				}
				next.ServeHTTP(w, r)
			})
		}

		if err := http.ListenAndServe(":8080", corsMiddleware(mux)); err != nil {
			log.Fatalf("Failed to start grpc-gateway server: %v", err)
		}
	}()

	log.Println("Server successfully running on port ", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
