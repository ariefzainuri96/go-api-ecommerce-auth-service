package server

import (
	"fmt"
	"log"
	"net"

	"github.com/ariefzainuri96/go-api-ecommerce-auth-service/internal/service"
	"github.com/ariefzainuri96/go-api-ecommerce-auth-service/internal/store"
	authpb "github.com/ariefzainuri96/go-api-ecommerce-auth-service/proto"
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

func (s *GRPCServer) Start(store store.Storage) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()	

	// register service implementation
	authpb.RegisterAuthServiceServer(grpcServer, service.NewAuthService(store))

	log.Printf("gRPC Auth Service running on port %d", s.port)

	return grpcServer.Serve(lis)
}