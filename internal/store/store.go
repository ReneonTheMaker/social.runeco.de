package store

// sqlite gorm db for users and user info and posts

import (
	"app/internal/model"

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

// post crud

func (s *Store) CreatePost(post *model.Post) error {
	return s.DB.Create(post).Error
}

func (s *Store) GetPostByID(postID uint) (*model.Post, error) {
	var post model.Post
	err := s.DB.Preload("User").Preload("Parent").Where("id = ?", postID).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *Store) DeletePost(postID uint) error {
	return s.DB.Delete(&model.Post{}, postID).Error
}

// auth
func (s *Store) AuthenticateUser(username, passwordHash string) (*model.User, error) {
	var user model.User
	err := s.DB.Where("username = ? AND password_hash = ?", username, passwordHash).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
