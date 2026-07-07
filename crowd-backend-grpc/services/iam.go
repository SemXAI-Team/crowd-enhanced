package services

import (
	"context"
	"database/sql"

	pb "crowd-backend-grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IAMServer struct {
	pb.UnimplementedIAMServiceServer
	db *sql.DB
}

func NewIAMServer(db *sql.DB) *IAMServer {
	return &IAMServer{db: db}
}

func (s *IAMServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, email, name, role FROM users")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to query users: %v", err)
	}
	defer rows.Close()

	var users []*pb.UserInfo
	for rows.Next() {
		var u pb.UserInfo
		if err := rows.Scan(&u.Id, &u.Email, &u.Name, &u.Role); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan user: %v", err)
		}
		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "rows iteration error: %v", err)
	}

	return &pb.ListUsersResponse{Users: users}, nil
}

func (s *IAMServer) UpdateUserRole(ctx context.Context, req *pb.UpdateUserRoleRequest) (*pb.UserInfo, error) {
	var u pb.UserInfo
	err := s.db.QueryRowContext(ctx, 
		"UPDATE users SET role = $1 WHERE id = $2 RETURNING id, email, name, role", 
		req.GetRole(), req.GetUserId(),
	).Scan(&u.Id, &u.Email, &u.Name, &u.Role)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user role: %v", err)
	}

	return &u, nil
}

func (s *IAMServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	res, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return &pb.DeleteUserResponse{Success: true}, nil
}

func (s *IAMServer) RevokeAccess(ctx context.Context, req *pb.RevokeAccessRequest) (*pb.RevokeAccessResponse, error) {
	// Let's implement RevokeAccess by just setting their role to 'viewer' or some disabled role.
	_, err := s.db.ExecContext(ctx, "UPDATE users SET role = 'viewer' WHERE id = $1", req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to revoke access: %v", err)
	}
	return &pb.RevokeAccessResponse{Success: true}, nil
}
