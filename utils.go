package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v4"
)

func HashPassword(pw string) (string, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(bs), err
}

func CheckPasswordHash(pw, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}

func CreateToken(userID int, email string, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"is_admin": isAdmin,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
		"iat":     time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtSecret)
}
