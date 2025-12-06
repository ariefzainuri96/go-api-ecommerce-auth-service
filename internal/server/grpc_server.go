package server

import (
	"context"
	"fmt"
	"net"
	"time"

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

func (s *GRPCServer) Start(ctx context.Context, store store.Storage, logger *zap.Logger) error {
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

	// Register health server
	// healthServer := health.NewServer()
	// healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	// grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	logger.Info("gRPC Auth Service running on port %d", zap.Int("addr", s.port))

	// return grpcServer.Serve(lis)

	serveErrCh := make(chan error, 1)
	go func() {
		logger.Info("starting grpc server serve loop")
		if err := grpcServer.Serve(lis); err != nil {
			serveErrCh <- fmt.Errorf("grpc serve error: %w", err)
			return
		}
		serveErrCh <- nil
	}()

	// Wait for context cancel or serve error
	select {
	case <-ctx.Done():
		// Graceful stop: stop accepting new connections, finish RPCs
		logger.Info("shutting down grpc server (GracefulStop)")
		// give some time for in-flight RPCs to finish (optional)
		done := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			logger.Info("grpc server stopped gracefully")
		case <-time.After(15 * time.Second):
			logger.Warn("grpc graceful stop timeout — forcing Stop()")
			grpcServer.Stop()
		}
		return nil
	case err := <-serveErrCh:
		// grpc Serve returned an error (possibly fatal)
		return err
	}
}
