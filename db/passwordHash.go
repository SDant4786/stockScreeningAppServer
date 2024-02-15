package db

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(Password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(Password), 14)
	return string(bytes), err
}

func CheckPasswordHash(Password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(Password))
	return err == nil
}