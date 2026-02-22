package store

import "app/internal/model"

func (s *Store) GetTopPosts(limit int) ([]model.Post, error) {
	var posts []model.Post
	err := s.DB.Preload("User").Preload("Parent").Where("parent_id is null").Order("created_at desc").Limit(limit).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *Store) GetReplyPosts(parentID uint) ([]model.Post, error) {
	var posts []model.Post
	err := s.DB.Preload("User").Preload("Parent").Where("parent_id = ?", parentID).Order("created_at asc").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *Store) GetNumberOfReplies(postID uint) (int64, error) {
	var count int64
	err := s.DB.Model(&model.Post{}).Where("parent_id = ?", postID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Store) CreateReply(userID uint, parentID uint, content string) (model.Post, error) {
	user, _ := s.GetUserByID(userID)
	reply := model.Post{
		UserID:   userID,
		User:     *user,
		ParentID: &parentID,
		Content:  content,
	}
	return reply, s.DB.Create(&reply).Error
}

func (s *Store) CreatePost(userID uint, content string) (model.Post, error) {
	user, _ := s.GetUserByID(userID)
	post := model.Post{
		UserID:  userID,
		User:    *user,
		Content: content,
	}
	return post, s.DB.Create(&post).Error
}

func (s *Store) EndSession(sessionId string) error {
	return s.DB.Where("session_id = ?", sessionId).Delete(&model.UserLogin{}).Error
}

func (s *Store) IsUserModerator(userID uint) (bool, error) {
	var user model.User
	err := s.DB.First(&user, userID).Error
	if err != nil {
		return false, err
	}
	return user.Mod, nil
}

func (s *Store) CanDeletePost(userID uint, postID uint) (bool, error) {
	var post model.Post
	err := s.DB.First(&post, postID).Error
	if err != nil {
		return false, err
	}
	if post.UserID == userID {
		return true, nil
	}
	return s.IsUserModerator(userID)
}

func (s *Store) DeletePost(postID uint) error {
	return s.DB.Delete(&model.Post{}, postID).Error
}
