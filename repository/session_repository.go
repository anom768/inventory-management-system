package repository

import (
	"gorm.io/gorm"
	"inventory-management-system/model"
)

type SessionRepository interface {
	Add(session model.Sessions) (model.Sessions, error)
	Delete(username string) error
	GetByUsername(username string) (model.Sessions, error)
}

type sessionRepository struct {
	*gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{DB: db}
}

func (s *sessionRepository) Add(session model.Sessions) (model.Sessions, error) {
	if err := s.DB.Create(&session).Error; err != nil {
		return model.Sessions{}, err
	}

	return session, nil
}

func (s *sessionRepository) Delete(username string) error {
	session, err := s.GetByUsername(username)
	if err != nil {
		return err
	}

	if err := s.DB.Delete(&session).Error; err != nil {
		return err
	}

	return nil
}

func (s *sessionRepository) GetByUsername(username string) (model.Sessions, error) {
	session := model.Sessions{}
	if err := s.DB.Where("username = ?", username).First(&session).Error; err != nil {
		return model.Sessions{}, err
	}

	return session, nil
}
