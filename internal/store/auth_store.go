package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ariefzainuri96/go-api-ecommerce-auth-service/cmd/api/entity"
	authpb "github.com/ariefzainuri96/go-api-ecommerce-auth-service/proto"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthStore struct {
	db     *sql.DB
	gormDb *gorm.DB
}

func (store *AuthStore) Login(ctx context.Context, body *authpb.LoginRequest) (entity.User, string, error) {
	user := entity.User{
		Email: body.Email,
	}

	err := store.gormDb.
		WithContext(ctx).
		// get data by condition from user instance, which is by email
		Where(user).
		// insert data to [user] address
		First(&user).Error

	if err != nil {
		return user, "", err
	}

	// query := `SELECT id, name, email, password, is_admin FROM users WHERE email = $1;`

	// row := store.db.QueryRowContext(ctx, query, body.Email)

	// err := row.Scan(&login.ID, &login.Name, &login.Email, &password, &isAdmin)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		return user, "", errors.New("invalid email or password")
	}

	token, err := generateToken(body.Email, user.IsAdmin, int(user.ID))

	if err != nil {
		return user, "", err
	}

	return user, token, nil
}

func generateToken(email string, isAdmin bool, id int) (string, error) {
	jwtSecret := strings.TrimSpace(os.Getenv("SECRET_KEY"))

	claims := jwt.MapClaims{
		"user_id":  id,
		"email":    email,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(), // Token valid for 30 day
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (store *AuthStore) Register(ctx context.Context, body *authpb.RegisterRequest) error {
	var emailExists bool

	user := entity.User{
		Email: body.Email,
	}

	err := store.gormDb.
		WithContext(ctx).
		Model(&user).
		Where(user).
		Scan(&emailExists).Error

	if err != nil {
		return err
	} else if emailExists {
		return fmt.Errorf("email sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Name = body.FullName
	user.Password = string(hashedPassword)
	user.IsAdmin = false

	result := store.gormDb.WithContext(ctx).Create(&user)

	// query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3)`

	// result, err := store.db.ExecContext(ctx, query, body.Name, body.Email, string(hashedPassword))

	if result.Error != nil {
		return err
	}

	return nil
}
