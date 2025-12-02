package service

import (
	"context"
	
	"github.com/ariefzainuri96/go-api-ecommerce-auth-service/internal/store"	
	authpb "github.com/ariefzainuri96/go-api-ecommerce-auth-service/proto"
	"go.uber.org/zap"
)

// AuthServiceImpl implements the RPC methods defined in auth.proto
type AuthServiceImpl struct {
	logger *zap.Logger
	store  store.Storage
	authpb.UnimplementedAuthServiceServer
}

func NewAuthService(store store.Storage, logger *zap.Logger) *AuthServiceImpl {
	return &AuthServiceImpl{
		logger: logger,
		store: store,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	id, err := s.store.IAuth.Register(ctx, req)

	if err != nil {
		return nil, err
	}

	return &authpb.RegisterResponse{
		UserId: int64(id), // fake ID for now
	}, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	// utils.LogMethod(s.logger, ctx, req)
	
	user, token, err := s.store.IAuth.Login(ctx, req)

	if err != nil {
		return nil, err
	}

	return &authpb.LoginResponse{
		Id:    int64(user.ID),
		Token: token,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
