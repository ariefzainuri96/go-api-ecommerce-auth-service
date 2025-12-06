package store

import (
	"context"
	"database/sql"

	"github.com/ariefzainuri96/go-api-ecommerce-auth-service/cmd/api/entity"
	db "github.com/ariefzainuri96/go-api-ecommerce-auth-service/internal/db"
	authpb "github.com/ariefzainuri96/go-api-ecommerce-auth-service/proto"
)

type Storage struct {
	IAuth interface {
		Login(context.Context, *authpb.LoginRequest) (entity.User, string, error)
		Register(context.Context, *authpb.RegisterRequest) (uint, error)
	}
}

func NewStorage(db *sql.DB, gorm *db.GormDB) Storage {
	return Storage{
		IAuth: &AuthStore{db, gorm},
	}
}
