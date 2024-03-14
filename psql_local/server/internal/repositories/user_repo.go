package repositories

import (
	"context"
	"server/internal/apperrors"
	"server/internal/domain/models"
	"server/internal/domain/requests"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RepoUsers interface {
	UpdateUser(ctx context.Context, user *models.User) error
	UpdateUserPasswordById(ctx context.Context, userID, newPass string) error
	UpdateUserById(ctx context.Context, userReq *requests.CreateUserRequest) error
	CreateUser(ctx context.Context, user *models.User) (string, error)
	GetWordsByIDAndLimit(ctx context.Context, id *uuid.UUID, limit int) ([]*models.Word, error)
	GetLearnByIDAndLimit(ctx context.Context, id *uuid.UUID, limit int) ([]*models.Word, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserById(ctx context.Context, id *uuid.UUID) (*models.User, error)
	MoveWordToLearned(ctx context.Context, user *models.User, word *models.Word) error
	AddWordToLearn(ctx context.Context, user *models.User, word *models.Word) error
	DeleteLearnWordFromUserByWordID(ctx context.Context, user *models.User, word *models.Word) error
	GetWordsByUserIdAndLimitAndTopic(ctx context.Context, id *uuid.UUID, limit int, topic string) ([]*models.Word, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	AddWordsToUser(ctx context.Context, user *models.User, words []*models.Word) error
}

type repoUsers struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewRepoUsers(db *gorm.DB, log *logrus.Logger) RepoUsers {
	return &repoUsers{db: db, log: log}
}

func (usr *repoUsers) GetUserById(ctx context.Context, id *uuid.UUID) (*models.User, error) {
	var user *models.User
	err := usr.db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		appErr := apperrors.GetUserByIdErr.AppendMessage(err)
		usr.log.Error(appErr)
		return nil, appErr
	}

	return user, nil
}

func (usr *repoUsers) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user *models.User
	err := usr.db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		appErr := apperrors.GetUserByEmailErr.AppendMessage(err)
		usr.log.Error(appErr)
		return nil, appErr
	}

	return user, nil
}

func (usr *repoUsers) MoveWordToLearned(ctx context.Context, user *models.User, word *models.Word) error {
	tx := usr.db.Begin()
	if tx.Error != nil {
		appErr := apperrors.MoveWordToLearnedErr.AppendMessage(tx.Error)
		usr.log.Error(appErr)
		return appErr
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	association := tx.Model(user).Association("Words")
	if association.Error != nil {
		tx.Rollback()
		appErr := apperrors.MoveWordToLearnedErr.AppendMessage(association.Error)
		usr.log.Error(appErr)
		return appErr
	}

	if err := association.Delete(word); err != nil {
		tx.Rollback()
		appErr := apperrors.MoveWordToLearnedErr.AppendMessage(err)
		usr.log.Error(appErr)
		return appErr
	}

	err := tx.Model(user).Association("Learned").Append(word)
	if err != nil {
		tx.Rollback()
		appErr := apperrors.MoveWordToLearnedErr.AppendMessage(err)
		usr.log.Error(appErr)
		return appErr
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		appErr := apperrors.MoveWordToLearnedErr.AppendMessage(err)
		usr.log.Error(appErr)
		return appErr
	}

	return nil
}

func (usr *repoUsers) AddWordsToUser(ctx context.Context, user *models.User, words []*models.Word) error {
	for _, word := range words {
		err := usr.db.Model(user).Association("Words").Append(word)
		if err != nil {
			appErr := apperrors.AddWordsToUserErr.AppendMessage(word.English, "-word-", err)
			usr.log.Error(appErr)
			return appErr
		}
	}

	return nil
}

func (usr *repoUsers) AddWordToLearn(ctx context.Context, user *models.User, word *models.Word) error {
	err := usr.db.Model(user).Association("Learn").Append(word)
	if err != nil {
		appErr := apperrors.AddWordToLearnRepoErr.AppendMessage(err)
		usr.log.Error(appErr)
		return appErr
	}

	return nil
}

func (usr *repoUsers) UpdateUser(ctx context.Context, user *models.User) error {
	tx := usr.db.Begin()
	if tx.Error != nil {
		appErr := apperrors.UpdateUserErr.AppendMessage(tx.Error)
		usr.log.Error(appErr)
		return appErr
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		appErr := apperrors.UpdateUserErr.AppendMessage(err)
		usr.log.Error(appErr)
		return appErr
	}

	return tx.Commit().Error
}

func (usr *repoUsers) UpdateUserById(ctx context.Context, userReq *requests.CreateUserRequest) error {
	result := usr.db.Model(&models.User{}).Where("id = ?", userReq.Id).
		Updates(map[string]interface{}{
			"email":     userReq.Email,
			"name":      userReq.Name,
			"last_name": userReq.LastName,
			"role":      userReq.Role,
		})
	if result.Error != nil {
		appErr := apperrors.UpdateUserByIdErr.AppendMessage(result.Error)
		usr.log.Error(appErr)
		return appErr
	}

	if result.RowsAffected == 0 {
		appErr := &apperrors.UpdateWordRowAffectedErr
		usr.log.Info(appErr)
		return appErr
	}

	return nil
}

func (usr *repoUsers) UpdateUserPasswordById(ctx context.Context, userID, newPass string) error {
	result := usr.db.Model(&models.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"password": newPass,
		})
	if result.Error != nil {
		appErr := apperrors.UpdateUserByIdErr.AppendMessage(result.Error)
		usr.log.Error(appErr)
		return appErr
	}

	if result.RowsAffected == 0 {
		appErr := &apperrors.UpdateWordRowAffectedErr
		usr.log.Info(appErr)
		return appErr
	}

	return nil
}

func (usr *repoUsers) CreateUser(ctx context.Context, user *models.User) (string, error) {
	if user == nil {
		appErr := apperrors.CreateUserErr.AppendMessage("user is nil")
		usr.log.Error(appErr)
		return "", appErr
	}

	tx := usr.db.WithContext(ctx)
	if tx.Error != nil {
		appErr := apperrors.CreateUserErr.AppendMessage(tx.Error)
		usr.log.Error(appErr)
		return "", appErr
	}

	result := tx.Create(user)
	if result.Error != nil {
		appErr := apperrors.CreateUserErr.AppendMessage(result.Error)
		usr.log.Error(appErr)
		return "", appErr
	}

	if result.RowsAffected == 0 {
		appErr := apperrors.CreateUserErr.AppendMessage("no rows affected")
		usr.log.Error(appErr)
		return "", appErr
	}

	createdUser := &models.User{}
	if err := tx.First(createdUser, "id = ?", user.ID).Error; err != nil {
		appErr := apperrors.CreateUserErr.AppendMessage(err)
		usr.log.Error(appErr)
		return "", appErr
	}

	return createdUser.ID.String(), nil
}

func (usr *repoUsers) GetWordsByIDAndLimit(ctx context.Context, id *uuid.UUID, limit int) ([]*models.Word, error) {
	var user *models.User
	err := usr.db.Preload("Words", func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}).Where("id = ?", id).Find(&user).Error
	if err != nil {
		appErr := apperrors.GetWordsByIDAndLimitErr.AppendMessage(err)
		usr.log.Error(appErr)
		return nil, appErr
	}

	usr.log.Info(len(user.Words))
	return user.Words, nil
}

func (usr *repoUsers) GetWordsByUserIdAndLimitAndTopic(ctx context.Context, id *uuid.UUID, limit int, topic string) ([]*models.Word, error) {
	words := []*models.Word{}
	err := usr.db.
		Unscoped().
		Table("user_words").
		Select("id, english, russian, theme, parts_of_speech", "created_at", "updated_at").
		Joins("JOIN words ON words.id = user_words.word_id").
		Where("user_words.user_id = ? AND words.theme = ?", id, topic). //"Irregular verb"
		Limit(limit).
		Find(&words).Error

	if err != nil {
		appErr := apperrors.GetWordsByUserIdAndLimitAndTopicErr.AppendMessage(err)
		usr.log.Error(appErr)
		return nil, appErr
	}

	return words, nil
}

func (usr *repoUsers) GetLearnByIDAndLimit(ctx context.Context, id *uuid.UUID, limit int) ([]*models.Word, error) {
	var user *models.User
	err := usr.db.Preload("Learn", func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}).Where("id = ?", id).Find(&user).Error
	if err != nil {
		appErr := apperrors.GetLearnByIDAndLimitErr.AppendMessage(err)
		usr.log.Error(appErr)
		return nil, appErr
	}

	usr.log.Info(user)
	usr.log.Info(len(user.Learn))
	return user.Learn, nil
}

func (usr *repoUsers) DeleteLearnWordFromUserByWordID(ctx context.Context, user *models.User, word *models.Word) error {
	association := usr.db.Model(user).Association("Learn")
	if association.Error != nil {
		appErr := apperrors.DeleteLearnWordFromUserByWordIDErr.AppendMessage(association.Error)
		usr.log.Error(appErr)
		return appErr
	}

	if err := association.Delete(word); err != nil {
		appErr := apperrors.DeleteLearnWordFromUserByWordIDErr.AppendMessage(err)
		usr.log.Error(appErr)
		return appErr
	}

	return nil
}

func (usr *repoUsers) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	err := usr.db.Find(&users).Error
	if err != nil {
		appErr := apperrors.GetAllUsersErr.AppendMessage(err)
		usr.log.Error(appErr)
		return nil, appErr
	}

	return users, nil
}
