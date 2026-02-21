package web

import (
	"errors"

	"app/internal/model"
	"app/internal/store"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrInvalidCredentials = fiber.NewError(fiber.StatusUnauthorized, "Invalid username or password")

func Auth(sessionId string, store *store.Store) bool {
	user, err := store.GetUserBySession(sessionId)
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

func CreateUser(username, password string, store *store.Store) (*model.User, error) {
	passwordHash, err := HashString(password)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Username:     username,
		PasswordHash: passwordHash,
	}
	err = store.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func LoginOrSignUp(username, password string, store *store.Store) (*model.User, error) {
	user, err := store.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return CreateUser(username, password, store)
		}
		return nil, err
	}
	// check password
	match, err := CheckPasswordHash(*user, password)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}
