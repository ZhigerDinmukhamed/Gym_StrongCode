package utils

import (
	"context"
	"errors"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	UserID  int    `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(getSecret())

func getSecret() string {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		return "CHANGE_ME_SECRET"
	}
	return s
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password empty")
	}
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateToken(userID int, email string, isAdmin bool) (string, error) {
	claims := JWTClaims{
		UserID:  userID,
		Email:   email,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(str string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(str, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token invalid")
	}
	return claims, nil
}

// Store user into context
type userKeyType string

var userKey userKeyType = "user"

func ContextWithUser(ctx context.Context, claims *JWTClaims) context.Context {
	return context.WithValue(ctx, userKey, claims)
}

func GetUserFromContext(ctx context.Context) *JWTClaims {
	if v := ctx.Value(userKey); v != nil {
		return v.(*JWTClaims)
	}
	return nil
}
