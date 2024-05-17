package repository

import (
	"gorm.io/gorm"
	"inventory-management-system/model/domain"
)

type SessionRepository interface {
	Add(session domain.Sessions) (domain.Sessions, error)
	Delete(username string) error
	GetByUsername(username string) (domain.Sessions, error)
}

type sessionRepository struct {
	*gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{DB: db}
}

func (s *sessionRepository) Add(session domain.Sessions) (domain.Sessions, error) {
	if err := s.DB.Create(&session).Error; err != nil {
		return domain.Sessions{}, err
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

func (s *sessionRepository) GetByUsername(username string) (domain.Sessions, error) {
	session := domain.Sessions{}
	if err := s.DB.Where("username = ?", username).First(&session).Error; err != nil {
		return domain.Sessions{}, err
	}

	return session, nil
}
