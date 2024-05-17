package gapi

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zde37/Swift_Bank/helpers"
	"github.com/zde37/Swift_Bank/models"
	"github.com/zde37/Swift_Bank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status" 
)


func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := helpers.HashPassword(req.GetPassword()) // use 'GetPassword' because this function checks if the request object is nil
	if err != nil { 
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err )
	}
	
 
	createdUser, err := server.service.CreateUser(ctx, models.CreateUserRequest{
		UserName: req.GetUsername(),
		Password: hashedPassword,
		FullName: req.GetFullName(),
		Email: req.GetEmail(),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, status.Errorf(codes.AlreadyExists, "a db error occurred: %s", err )

		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err )
	}

	return &pb.CreateUserResponse{
		User: convertUser(createdUser),
	}, nil
}

// func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	 
// }
