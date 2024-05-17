package gapi

import (
	"fmt"

	"github.com/zde37/Swift_Bank/config"
	"github.com/zde37/Swift_Bank/pb"
	"github.com/zde37/Swift_Bank/service"
	"github.com/zde37/Swift_Bank/token"
)

// Server serves gRPC requests for the banking service
type Server struct {
	pb.UnimplementedSwiftBankServer
	service    service.ServiceProvider
	tokenMaker token.Maker
	config     config.Config
}

// NewServer creates a new gRPC server
func NewServer(c config.Config, s service.ServiceProvider) (*Server, error) {
	tokenMaker, err := token.NewPastoMaker(c.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		service:    s,
		config:     c,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
