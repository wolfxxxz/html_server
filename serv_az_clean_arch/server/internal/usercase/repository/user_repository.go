package repository

import (
	"context"
	"server/internal/domain/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req *models.User) (string, error)
	AddWordsToUser(ctx context.Context, user *models.User, words []*models.Word) error
	//GetUsersByPageAndPerPage(ctx context.Context, page, perPage int) ([]*models.User, error)
	//GetUserByID(ctx context.Context, userUUID *uuid.UUID) (*models.User, error)
	//UpdateUserByID(ctx context.Context, req *models.User) (string, error)
	//UpdateUserPasswordById(ctx context.Context, req *models.User) (string, error)
}
