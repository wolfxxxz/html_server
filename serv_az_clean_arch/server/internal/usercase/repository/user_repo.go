package repository

import (
	"context"
	"server/internal/domain/models"
	"server/internal/domain/requests"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req *models.User) (string, error)
	AddWordsToUser(ctx context.Context, user *models.User, words []*models.Word) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	UpdateUserPasswordById(ctx context.Context, userID, newPass string) error
	UpdateUserById(ctx context.Context, userReq *requests.CreateUserRequest) error
	GetWordsByIDAndLimit(ctx context.Context, id *uuid.UUID, limit int) ([]*models.Word, error)
	GetLearnByIDAndLimit(ctx context.Context, id *uuid.UUID, limit int) ([]*models.Word, error)
	GetUserById(ctx context.Context, id *uuid.UUID) (*models.User, error)
	MoveWordToLearned(ctx context.Context, user *models.User, word *models.Word) error
	AddWordToLearn(ctx context.Context, user *models.User, word *models.Word) error
	DeleteLearnWordFromUserByWordID(ctx context.Context, user *models.User, word *models.Word) error
	GetWordsByUserIdAndLimitAndTopic(ctx context.Context, id *uuid.UUID, limit int, topic string) ([]*models.Word, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
}
