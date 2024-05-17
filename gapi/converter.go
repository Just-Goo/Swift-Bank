package gapi

import (
	"github.com/zde37/Swift_Bank/models"
	"github.com/zde37/Swift_Bank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user models.User) *pb.User {
	return &pb.User{
		Username:          user.UserName,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
