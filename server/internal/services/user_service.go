package services

import (
	"context"
	"server/internal/apperrors"
	"server/internal/domain/mappers"
	"server/internal/domain/models"
	"server/internal/domain/requests"
	"server/internal/domain/responses"
	"server/internal/repositories"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	repoUser    repositories.RepoUsers
	repoLibrary repositories.RepoLibrary
	log         *logrus.Logger
}

func NewUserService(userRepo repositories.RepoUsers, repoLibrary repositories.RepoLibrary, log *logrus.Logger) *UserService {
	return &UserService{repoUser: userRepo, repoLibrary: repoLibrary, log: log}
}

func (us *UserService) CreateUser(ctx context.Context, userReq *requests.CreateUserRequest) (*responses.CreateUserResponse, error) {
	user := mappers.MapReqCreateUsToUser(userReq)
	hashPass, err := hashPassword(userReq.Password)
	if err != nil {
		us.log.Error(err)
		return nil, err
	}

	user.Password = hashPass
	createdUserID, err := us.repoUser.CreateUser(ctx, user)
	if err != nil {
		us.log.Error(err)
		return nil, err
	}

	userUUID, err := uuid.Parse(createdUserID)
	if err != nil {
		appErr := apperrors.CreateUserErr.AppendMessage(err)
		us.log.Error(appErr)
		return nil, appErr
	}

	user.ID = &userUUID

	library, err := us.repoLibrary.GetAllWords()
	if err != nil {
		return nil, err
	}

	words := mappers.MapLibraryToWords(library)
	user.Words = words
	err = us.repoUser.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	respCreateUser := &responses.CreateUserResponse{UserId: user.ID.String()}
	return respCreateUser, nil
}

func (us *UserService) UpdateUserById(ctx context.Context, user *models.User, userReq *requests.CreateUserRequest) error {
	if !checkPasswordHash(userReq.Password, user.Password) {
		appErr := apperrors.UpdateUserByIdErr.AppendMessage("WRONG Password")
		us.log.Error(appErr)
		return appErr
	}

	userReq.Id = user.ID.String()
	return us.repoUser.UpdateUserById(ctx, userReq)
}

func (us *UserService) UpdateUserPasswordById(ctx context.Context, user *models.User, oldPass, newPass, newPassSec string) error {
	if !checkPasswordHash(oldPass, user.Password) {
		appErr := apperrors.UpdateUserPasswordByIdErr.AppendMessage("WRONG Password")
		us.log.Error(appErr)
		return appErr
	}

	if !strings.EqualFold(newPass, newPassSec) {
		appErr := apperrors.UpdateUserPasswordByIdErr.AppendMessage("WRONG New Password")
		us.log.Error(appErr)
		return appErr
	}

	hashPass, err := hashPassword(newPass)
	if err != nil {
		us.log.Error(err)
		return err
	}

	return us.repoUser.UpdateUserPasswordById(ctx, user.ID.String(), hashPass)
}

func (us *UserService) SignInUserWithJWT(ctx context.Context, logReq *requests.LoginRequest, secretKey string, expiresAt string) (*responses.LoginResponse, error) {
	user, err := us.repoUser.GetUserByEmail(ctx, logReq.Email)
	if err != nil {
		us.log.Error(err)
		return nil, err
	}

	if !checkPasswordHash(logReq.Password, user.Password) {
		appErr := apperrors.SignInUserWithJWTErr.AppendMessage("check password err")
		us.log.Error(appErr)
		return nil, appErr
	}

	token, err := claimJWTToken(user.Role, user.ID.String(), expiresAt, []byte(secretKey))
	if err != nil {
		us.log.Error(err)
		return nil, err
	}

	return mappers.MapTokenToLoginResponse(token, expiresAt), nil
}

func (us *UserService) GetWordsByUserIdAndLimitAndTopic(ctx context.Context, getWordsReq *requests.GetWordsByUsIdAndLimitRequest, topic string) ([]*models.Word, error) {
	quantity, err := strconv.Atoi(getWordsReq.Limit)
	if err != nil {
		appErr := apperrors.GetWordsByUserIdAndLimitAndTopicErr.AppendMessage(err)
		us.log.Error(appErr)
		return nil, appErr
	}

	userId, err := uuid.Parse(getWordsReq.ID)
	if err != nil {
		appErr := apperrors.GetWordsByUserIdAndLimitAndTopicErr.AppendMessage(err)
		us.log.Error(appErr)
		return nil, appErr
	}

	return us.repoUser.GetWordsByUserIdAndLimitAndTopic(ctx, &userId, quantity, topic)
}

func (us *UserService) GetWordsByUsIdAndLimit(ctx context.Context, getWordsReq *requests.GetWordsByUsIdAndLimitRequest) ([]*models.Word, error) {
	quantity, err := strconv.Atoi(getWordsReq.Limit)
	if err != nil {
		appErr := apperrors.GetWordsByUsIdAndLimitServiceErr.AppendMessage(err)
		us.log.Error(appErr)
		return nil, appErr
	}

	userId, err := uuid.Parse(getWordsReq.ID)
	if err != nil {
		appErr := apperrors.GetWordsByUsIdAndLimitServiceErr.AppendMessage(err)
		us.log.Error(appErr)
		return nil, appErr
	}

	return us.repoUser.GetWordsByIDAndLimit(ctx, &userId, quantity)
}

func (us *UserService) GetLearnByUsIdAndLimit(ctx context.Context, getWordsReq *requests.GetWordsByUsIdAndLimitRequest) ([]*models.Word, error) {
	quantity, err := strconv.Atoi(getWordsReq.Limit)
	if err != nil {
		appErr := apperrors.GetLearnByUsIdAndLimitErr.AppendMessage(err)
		us.log.Error(appErr)
		return nil, appErr
	}

	userId, err := uuid.Parse(getWordsReq.ID)
	if err != nil {
		appErr := apperrors.GetLearnByUsIdAndLimitErr.AppendMessage(err)
		us.log.Error(appErr)
		return nil, appErr
	}

	return us.repoUser.GetLearnByIDAndLimit(ctx, &userId, quantity)
}

func (us *UserService) GetUserById(ctx context.Context, id string) (*models.User, error) {
	userId, err := uuid.Parse(id)
	if err != nil {
		appErr := apperrors.GetUserByIdErr.AppendMessage(err)
		us.log.Error(appErr)
		return nil, err
	}

	return us.repoUser.GetUserById(ctx, &userId)
}

func (us *UserService) MoveWordToLearned(ctx context.Context, userID, wordID string) error {
	userId, err := uuid.Parse(userID)
	if err != nil {
		appErr := apperrors.MoveWordToLearnedErr.AppendMessage(err)
		us.log.Error(appErr)
		return err
	}

	user := &models.User{ID: &userId}
	wordId, err := uuid.Parse(wordID)
	if err != nil {
		appErr := apperrors.MoveWordToLearnedErr.AppendMessage(err)
		us.log.Error(appErr)
		return err
	}

	word := &models.Word{ID: &wordId}
	return us.repoUser.MoveWordToLearned(ctx, user, word)
}

func (us *UserService) AddWordToLearn(ctx context.Context, userID, wordID string) error {
	userId, err := uuid.Parse(userID)
	if err != nil {
		appErr := apperrors.AddWordToLearnErr.AppendMessage(err)
		us.log.Error(appErr)
		return appErr
	}

	user := &models.User{ID: &userId}
	wordId, err := uuid.Parse(wordID)
	if err != nil {
		appErr := apperrors.AddWordToLearnErr.AppendMessage(err)
		us.log.Error(appErr)
		return appErr
	}

	word := &models.Word{ID: &wordId}
	return us.repoUser.AddWordToLearn(ctx, user, word)
}

func (us *UserService) DeleteLearnFromUserById(ctx context.Context, userID, wordID string) error {
	userId, err := uuid.Parse(userID)
	if err != nil {
		appErr := apperrors.DeleteLearnFromUserByIdErr.AppendMessage(err)
		us.log.Error(appErr)
		return err
	}

	user := &models.User{ID: &userId}
	wordId, err := uuid.Parse(wordID)
	if err != nil {
		appErr := apperrors.DeleteLearnFromUserByIdErr.AppendMessage(err)
		us.log.Error(appErr)
		return err
	}

	word := &models.Word{ID: &wordId}
	return us.repoUser.DeleteLearnWordFromUserByWordID(ctx, user, word)
}

func (us *UserService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return us.repoUser.GetAllUsers(ctx)
}
