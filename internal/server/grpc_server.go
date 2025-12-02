package server

import (
	"fmt"
	"log"
	"net"

	"github.com/ariefzainuri96/go-api-ecommerce-auth-service/internal/interceptor"
	"github.com/ariefzainuri96/go-api-ecommerce-auth-service/internal/service"
	"github.com/ariefzainuri96/go-api-ecommerce-auth-service/internal/store"
	authpb "github.com/ariefzainuri96/go-api-ecommerce-auth-service/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	port int
}

func NewGRPCServer(port int) *GRPCServer {
	return &GRPCServer{
		port: port,
	}
}

func (s *GRPCServer) Start(store store.Storage, logger *zap.Logger) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.CorrelationIDInterceptor(logger),
			interceptor.LoggingInterceptor(logger),
		),
	)	

	// register service implementation
	authpb.RegisterAuthServiceServer(grpcServer, service.NewAuthService(store, logger))

	log.Printf("gRPC Auth Service running on port %d", s.port)

	return grpcServer.Serve(lis)
}