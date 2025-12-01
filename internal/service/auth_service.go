package service

import (
	"context"
	"log"

	"github.com/ariefzainuri96/go-api-ecommerce-auth-service/internal/store"
	authpb "github.com/ariefzainuri96/go-api-ecommerce-auth-service/proto"
)

// AuthServiceImpl implements the RPC methods defined in auth.proto
type AuthServiceImpl struct {
	store store.Storage
	authpb.UnimplementedAuthServiceServer
}

func NewAuthService(store store.Storage) *AuthServiceImpl {	
	return &AuthServiceImpl{
		store: store,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	log.Println("Register called:", req.Email)

	// TODO: save user to DB	

	return &authpb.RegisterResponse{
		UserId: 123, // fake ID for now
	}, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	log.Println("Login called:", req.Email)

	// TODO: validate user + generate JWT

	return &authpb.LoginResponse{
		Id:    123,
		Token: "fake.jwt.token",
		Name:  "John Doe",
		Email: req.Email,
	}, nil
}
