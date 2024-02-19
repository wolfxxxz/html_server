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
	us.repoUser.UpdateUser(ctx, user)
	respCreateUser := &responses.CreateUserResponse{UserId: user.ID.String()}
	return respCreateUser, nil
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

func (us *UserService) GetWordsByUsIdAndLimit(ctx context.Context, getWordsReq *requests.GetWordsByUsIdAndLimitRequest) ([]*responses.WordResp, error) {
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

	words, err := us.repoUser.GetWordsByIDAndLimit(ctx, &userId, quantity)
	if err != nil {
		us.log.Error(err)
		return nil, err
	}

	wordsResp := mappers.MapWordsToWordsResp(words)
	return wordsResp, nil
}

func (us *UserService) GetLearnByUsIdAndLimit(ctx context.Context, getWordsReq *requests.GetWordsByUsIdAndLimitRequest) ([]*responses.WordResp, error) {
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

	words, err := us.repoUser.GetLearnByIDAndLimit(ctx, &userId, quantity)
	if err != nil {
		us.log.Error(err)
		return nil, err
	}

	wordsResp := mappers.MapWordsToWordsResp(words)
	return wordsResp, nil
}

func (us *UserService) GetUserById(ctx context.Context, id string) (*models.User, error) {
	userId, err := uuid.Parse(id)
	if err != nil {
		appErr := apperrors.GetUserByIdErr.AppendMessage(err)
		us.log.Error(appErr)
		return nil, err
	}
	user, err := us.repoUser.GetUserById(ctx, &userId)
	if err != nil {
		us.log.Error(err)
		return nil, err
	}

	return user, nil
}

func (us *UserService) MoveWordToLearned(ctx context.Context, deleteWordReq *requests.DeleteWordFromUserByIDRequest) error {
	userId, err := uuid.Parse(deleteWordReq.UserID)
	if err != nil {
		appErr := apperrors.MoveWordToLearnedErr.AppendMessage(err)
		us.log.Error(appErr)
		return err
	}

	user := &models.User{ID: &userId}
	wordId, err := uuid.Parse(deleteWordReq.WordID)
	if err != nil {
		appErr := apperrors.MoveWordToLearnedErr.AppendMessage(err)
		us.log.Error(appErr)
		return err
	}

	word := &models.Word{ID: &wordId}

	err = us.repoUser.MoveWordToLearned(ctx, user, word)
	if err != nil {
		us.log.Error(err)
		return err
	}

	return nil
}

func (us *UserService) AddWordToLearn(ctx context.Context, deleteWordReq *requests.DeleteWordFromUserByIDRequest) error {
	userId, err := uuid.Parse(deleteWordReq.UserID)
	if err != nil {
		appErr := apperrors.AddWordToLearnErr.AppendMessage(err)
		us.log.Error(appErr)
		return appErr
	}

	user := &models.User{ID: &userId}
	wordId, err := uuid.Parse(deleteWordReq.WordID)
	if err != nil {
		appErr := apperrors.AddWordToLearnErr.AppendMessage(err)
		us.log.Error(appErr)
		return appErr
	}

	word := &models.Word{ID: &wordId}

	err = us.repoUser.AddWordToLearn(ctx, user, word)
	if err != nil {
		us.log.Error(err)
		return err
	}

	return nil
}

func (us *UserService) DeleteLearnFromUserById(ctx context.Context, deleteWordReq *requests.DeleteWordFromUserByIDRequest) error {
	userId, err := uuid.Parse(deleteWordReq.UserID)
	if err != nil {
		appErr := apperrors.DeleteLearnFromUserByIdErr.AppendMessage(err)
		us.log.Error(appErr)
		return err
	}

	user := &models.User{ID: &userId}
	wordId, err := uuid.Parse(deleteWordReq.WordID)
	if err != nil {
		appErr := apperrors.DeleteLearnFromUserByIdErr.AppendMessage(err)
		us.log.Error(appErr)
		return err
	}

	word := &models.Word{ID: &wordId}

	err = us.repoUser.DeleteLearnWordFromUserByWordID(ctx, user, word)
	if err != nil {
		us.log.Error(err)
		return err
	}

	return nil
}
