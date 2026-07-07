package server

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// In production, this should be an environment variable matching AUTH_SECRET
var jwtSecret = []byte("my-super-secret-auth-key-12345")

func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Define which methods require which roles
		requiresAdmin := map[string]bool{
			"/camera.ModelManager/UploadModel":    true,
			"/camera.ModelManager/DeleteModel":    true,
			"/camera.IAMService/UpdateUserRole":   true,
			"/camera.IAMService/RevokeAccess":     true,
			"/camera.CameraStream/RegisterCamera": true,
			"/camera.CameraStream/UpdateCamera":   true,
			"/camera.CameraStream/DeleteCamera":   true,
		}

		if requiresAdmin[info.FullMethod] {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
			}

			authHeader := md["authorization"]
			if len(authHeader) == 0 {
				return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
			}

			tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return nil, status.Errorf(codes.Unauthenticated, "invalid token claims")
			}

			role, ok := claims["role"].(string)
			if !ok || role != "admin" {
				return nil, status.Errorf(codes.PermissionDenied, "permission denied: requires admin role")
			}
		}

		// Proceed with the request
		return handler(ctx, req)
	}
}
