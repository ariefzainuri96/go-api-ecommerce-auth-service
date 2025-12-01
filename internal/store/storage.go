package store

import (
	"context"
	"database/sql"

	"github.com/ariefzainuri96/go-api-ecommerce-auth-service/cmd/api/entity"
	authpb "github.com/ariefzainuri96/go-api-ecommerce-auth-service/proto"
	"gorm.io/gorm"
)

type Storage struct {
	IAuth interface {
		Login(context.Context, *authpb.LoginRequest) (entity.User, string, error)
		Register(context.Context, *authpb.RegisterRequest) error
	}
}

func NewStorage(db *sql.DB, gormDb *gorm.DB) Storage {
	return Storage{
		IAuth: &AuthStore{db, gormDb},
	}
}
