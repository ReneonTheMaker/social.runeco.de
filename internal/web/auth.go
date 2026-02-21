package web

import (
	"app/internal/model"
	"app/internal/store"
	"golang.org/x/crypto/bcrypt"
)

func Auth(id uint, store *store.Store) bool {
	user, err := store.GetUserByID(id)
	if err != nil {
		return false
	}

	return user != nil
}

func HashString(s string) (string, error) {
	// bcrypt
	bytes, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPasswordHash(user model.User, passwordHash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(passwordHash))
	if err != nil {
		return false, err
	}
	return true, nil
}
