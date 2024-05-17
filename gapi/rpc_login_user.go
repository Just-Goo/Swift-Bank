package gapi

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zde37/Swift_Bank/models"
	"github.com/zde37/Swift_Bank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := server.service.LoginUser(ctx, models.LoginUserRequest{
		UserName: req.GetUsername(),
		Password: req.GetPassword(),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, status.Errorf(codes.AlreadyExists, "a db error occurred: %s", err )

		}
		return nil, status.Errorf(codes.Internal, "failed to login user: %s", err )
	}


	accessToken, _, err := server.tokenMaker.CreateToken(req.GetUsername(), server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create token: %s", err )
	}

	return &pb.LoginUserResponse{
		User: convertUser(user),
		AccessToken: accessToken,
	}, nil
}
