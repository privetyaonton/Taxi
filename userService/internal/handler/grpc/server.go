package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/RipperAcskt/innotaxi/config"
	"github.com/RipperAcskt/innotaxi/internal/service"
	"github.com/RipperAcskt/innotaxi/pkg/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	listener   net.Listener
	grpcServer *grpc.Server
	log        *zap.Logger
	cfg        *config.Config
}

func New(log *zap.Logger, cfg *config.Config) *Server {
	return &Server{nil, nil, log, cfg}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.cfg.GRPC_HOST)

	if err != nil {
		return fmt.Errorf("listen failed: %w", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	s.listener = listener
	s.grpcServer = grpcServer

	proto.RegisterAuthServiceServer(grpcServer, s)
	err = grpcServer.Serve(listener)
	if err != nil {
		return fmt.Errorf("serve failed: %w", err)
	}

	return nil
}

func (s *Server) GetJWT(c context.Context, params *proto.Params) (*proto.Response, error) {
	tokenParams := service.TokenParams{
		ID:                params.DriverID,
		Type:              params.Type,
		HS256_SECRET:      s.cfg.HS256_SECRET,
		ACCESS_TOKEN_EXP:  s.cfg.ACCESS_TOKEN_EXP,
		REFRESH_TOKEN_EXP: s.cfg.REFRESH_TOKEN_EXP,
	}

	token, err := service.NewToken(tokenParams)
	if err != nil {
		return nil, fmt.Errorf("new token failed: %w", err)
	}

	response := &proto.Response{
		AccessToken:  token.Access,
		RefreshToken: token.RT,
	}
	return response, nil
}

func (s *Server) Stop() error {
	s.log.Info("Shuttig down grpc...")

	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("listener close failed: %w", err)
	}

	s.grpcServer.Stop()
	s.log.Info("Grpc server exiting.")
	return nil
}
