package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func ComparePassword(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err24 := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err24 != nil {
		return false
	} else {
		return true
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
