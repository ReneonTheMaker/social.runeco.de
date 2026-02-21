package store

import "app/internal/model"

func (s *Store) GetTopPosts(limit int) ([]model.Post, error) {
	var posts []model.Post
	err := s.DB.Preload("User").Preload("Parent").Order("created_at desc").Limit(limit).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}
