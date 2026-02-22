package store

// sqlite gorm db for users and user info and posts

import (
	"time"

	"app/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	DB *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{DB: db}
}

// user crud

func (s *Store) CreateUser(user *model.User) error {
	return s.DB.Create(user).Error
}

func (s *Store) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := s.DB.Preload("UserInfo").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	err := s.DB.Preload("UserInfo").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) UpdateUserInfo(userInfo *model.UserInfo) error {
	return s.DB.Save(userInfo).Error
}

func (s *Store) GetUserInfoByUserID(userID uint) (*model.UserInfo, error) {
	var userInfo model.UserInfo
	err := s.DB.Where("user_id = ?", userID).First(&userInfo).Error
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func (s *Store) DeleteUser(userID uint) error {
	return s.DB.Delete(&model.User{}, userID).Error
}

func (s *Store) GetPostByID(postID uint) (*model.Post, error) {
	var post model.Post
	err := s.DB.Preload("User").Preload("Parent").Where("id = ?", postID).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *Store) AuthenticateUser(username, passwordHash string) (*model.User, error) {
	var user model.User
	err := s.DB.Where("username = ? AND password_hash = ?", username, passwordHash).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) SetUserSession(loginID uint) string {
	sessionId := uuid.New().String()
	userLogin := &model.UserLogin{
		UserID:    loginID,
		TokenHash: sessionId,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	s.DB.Create(userLogin)
	return sessionId
}

func (s *Store) GetUserBySession(sessionId string) (*model.User, error) {
	var userLogin model.UserLogin
	err := s.DB.Where("token_hash = ?", sessionId).First(&userLogin).Error
	if err != nil {
		return nil, err
	}
	if userLogin.ExpiresAt.Before(time.Now()) {
		s.DB.Delete(&userLogin)
		return nil, nil
	}
	s.DB.Model(&userLogin).Update("last_seen_at", time.Now())
	var user model.User
	err = s.DB.Where("id = ?", userLogin.UserID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) GetUserIDFromSession(sessionId string) (uint, error) {
	var userLogin model.UserLogin
	err := s.DB.Where("token_hash = ?", sessionId).First(&userLogin).Error
	if err != nil {
		return 0, err
	}
	if userLogin.ExpiresAt.Before(time.Now()) {
		s.DB.Delete(&userLogin)
		return 0, nil
	}
	s.DB.Model(&userLogin).Update("last_seen_at", time.Now())
	var user model.User
	err = s.DB.Where("id = ?", userLogin.UserID).First(&user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (s *Store) GetUserFromSession(sessionId string) (*model.User, error) {
	var userLogin model.UserLogin
	err := s.DB.Where("token_hash = ?", sessionId).First(&userLogin).Error
	if err != nil {
		return nil, err
	}
	if userLogin.ExpiresAt.Before(time.Now()) {
		s.DB.Delete(&userLogin)
		return nil, nil
	}
	s.DB.Model(&userLogin).Update("last_seen_at", time.Now())
	var user model.User
	err = s.DB.Where("id = ?", userLogin.UserID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
