package userusecase

import (
	"context"
	"time"

	"antiscam-simulator/internal/user/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user userdomain.User) error
	GetHistory(ctx context.Context, userID string) ([]userdomain.TrainingHistoryItem, error)
}

type UsecaseUser struct {
	repo UserRepository
}

func NewUsecaseUser(repo UserRepository) *UsecaseUser {
	return &UsecaseUser{repo: repo}
}

func (u *UsecaseUser) Register(ctx context.Context, username string) (string, error) {
	userID := uuid.New().String()
	user := userdomain.User{
		ID:        userID,
		Username:  username,
		CreatedAt: time.Now(),
	}

	err := u.repo.Create(ctx, user)
	if err != nil {
		return "", err
	}

	return userID, nil
}



func (u *UsecaseUser) GetHistory(ctx context.Context, userID string) ([]userdomain.TrainingHistoryItem, error) {
	return u.repo.GetHistory(ctx, userID)
}
