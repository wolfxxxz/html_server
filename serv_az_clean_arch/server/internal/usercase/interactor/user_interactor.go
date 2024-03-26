package interactor

import (
	"context"
	"math/rand"
	"server/internal/apperrors"
	"server/internal/domain/mappers"
	"server/internal/domain/models"
	"server/internal/domain/requests"
	"server/internal/domain/responses"
	"server/internal/infrastructure/email"
	"server/internal/usercase/repository"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type userInteractor struct {
	UserRepository  repository.UserRepository
	WordsRepository repository.WordsRepository
	Sender          email.Sender
}

type UserInteractor interface {
	CreateUser(ctx context.Context, userReq *requests.CreateUserRequest) (*responses.CreateUserResponse, error)
	SignInUserWithJWT(ctx context.Context, logReq *requests.LoginRequest, secretKey string, expiresAt string) (*responses.LoginResponse, error)
	RestoreUserPassword(ctx context.Context, email string) error
	UpdateUserById(ctx context.Context, user *models.User, userReq *requests.CreateUserRequest) error
	UpdateUserPasswordById(ctx context.Context, user *models.User, oldPass, newPass, newPassSec string) error
	GetWordsByUserIdAndLimitAndTopic(ctx context.Context, getWordsReq *requests.GetWordsByUsIdAndLimitRequest, topic string) ([]*models.Word, error)
	GetWordsByUsIdAndLimit(ctx context.Context, getWordsReq *requests.GetWordsByUsIdAndLimitRequest) ([]*models.Word, error)
	GetLearnByUsIdAndLimit(ctx context.Context, getWordsReq *requests.GetWordsByUsIdAndLimitRequest) ([]*models.Word, error)
	GetUserById(ctx context.Context, id string) (*models.User, error)
	MoveWordToLearned(ctx context.Context, userID, wordID string) error
	AddWordToLearn(ctx context.Context, userID, wordID string) error
	DeleteLearnFromUserById(ctx context.Context, userID, wordID string) error
	GetAllUsers(ctx context.Context) ([]*models.User, error)
}

func NewUserInteractor(u repository.UserRepository, w repository.WordsRepository, sender email.Sender) UserInteractor {
	return &userInteractor{UserRepository: u, WordsRepository: w, Sender: sender}
}

func (us *userInteractor) CreateUser(ctx context.Context, userReq *requests.CreateUserRequest) (*responses.CreateUserResponse, error) {
	user := mappers.MapReqCreateUsToUser(userReq)
	hashPass, err := hashPassword(userReq.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashPass
	createdUserID, err := us.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	userUUID, err := uuid.Parse(createdUserID)
	if err != nil {
		appErr := apperrors.CreateUserErr.AppendMessage(err)
		return nil, appErr
	}

	user.ID = &userUUID
	words, err := us.WordsRepository.GetAllWords()
	if err != nil {
		return nil, err
	}

	if len(words) == 0 {
		respCreateUser := &responses.CreateUserResponse{UserId: user.ID.String()}
		return respCreateUser, nil
	}

	err = us.UserRepository.AddWordsToUser(ctx, user, words)
	if err != nil {
		return nil, err
	}

	respCreateUser := &responses.CreateUserResponse{UserId: user.ID.String()}
	return respCreateUser, nil
}

func (us *userInteractor) SignInUserWithJWT(ctx context.Context, logReq *requests.LoginRequest, secretKey string, expiresAt string) (*responses.LoginResponse, error) {
	user, err := us.UserRepository.GetUserByEmail(ctx, logReq.Email)
	if err != nil {
		return nil, err
	}

	if !checkPasswordHash(logReq.Password, user.Password) {
		appErr := apperrors.SignInUserWithJWTErr.AppendMessage("check password err")
		return nil, appErr
	}

	token, err := claimJWTToken(user.Role, user.ID.String(), expiresAt, []byte(secretKey))
	if err != nil {
		return nil, err
	}

	return mappers.MapTokenToLoginResponse(token, expiresAt), nil
}

// --------------------------------------------
func (us *userInteractor) RestoreUserPassword(ctx context.Context, email string) error {
	user, err := us.UserRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	newPass := randomPassword()
	err = us.Sender.Send(email, "restore password translator", newPass)
	if err != nil {
		return err
	}

	userHashPassword, _ := hashPassword(newPass)

	err = us.UserRepository.UpdateUserPasswordById(ctx, user.ID.String(), userHashPassword)
	if err != nil {
		return err
	}

	return nil
}

func (us *userInteractor) UpdateUserById(ctx context.Context, user *models.User, userReq *requests.CreateUserRequest) error {
	if !checkPasswordHash(userReq.Password, user.Password) {
		appErr := apperrors.UpdateUserByIdErr.AppendMessage("WRONG Password")
		return appErr
	}

	userReq.Id = user.ID.String()
	return us.UserRepository.UpdateUserById(ctx, userReq)
}

func (us *userInteractor) UpdateUserPasswordById(ctx context.Context, user *models.User, oldPass, newPass, newPassSec string) error {
	if !checkPasswordHash(oldPass, user.Password) {
		appErr := apperrors.UpdateUserPasswordByIdErr.AppendMessage("WRONG Password")
		return appErr
	}

	if !strings.EqualFold(newPass, newPassSec) {
		appErr := apperrors.UpdateUserPasswordByIdErr.AppendMessage("WRONG New Password")
		return appErr
	}

	hashPass, err := hashPassword(newPass)
	if err != nil {
		return err
	}

	return us.UserRepository.UpdateUserPasswordById(ctx, user.ID.String(), hashPass)
}

func (us *userInteractor) GetWordsByUserIdAndLimitAndTopic(ctx context.Context, getWordsReq *requests.GetWordsByUsIdAndLimitRequest, topic string) ([]*models.Word, error) {
	quantity, err := strconv.Atoi(getWordsReq.Limit)
	if err != nil {
		appErr := apperrors.GetWordsByUserIdAndLimitAndTopicErr.AppendMessage(err)
		return nil, appErr
	}

	userId, err := uuid.Parse(getWordsReq.ID)
	if err != nil {
		appErr := apperrors.GetWordsByUserIdAndLimitAndTopicErr.AppendMessage(err)
		return nil, appErr
	}

	return us.UserRepository.GetWordsByUserIdAndLimitAndTopic(ctx, &userId, quantity, topic)
}

func (us *userInteractor) GetWordsByUsIdAndLimit(ctx context.Context, getWordsReq *requests.GetWordsByUsIdAndLimitRequest) ([]*models.Word, error) {
	quantity, err := strconv.Atoi(getWordsReq.Limit)
	if err != nil {
		appErr := apperrors.GetWordsByUsIdAndLimitServiceErr.AppendMessage(err)
		return nil, appErr
	}

	userId, err := uuid.Parse(getWordsReq.ID)
	if err != nil {
		appErr := apperrors.GetWordsByUsIdAndLimitServiceErr.AppendMessage(err)
		return nil, appErr
	}

	return us.UserRepository.GetWordsByIDAndLimit(ctx, &userId, quantity)
}

func (us *userInteractor) GetLearnByUsIdAndLimit(ctx context.Context, getWordsReq *requests.GetWordsByUsIdAndLimitRequest) ([]*models.Word, error) {
	quantity, err := strconv.Atoi(getWordsReq.Limit)
	if err != nil {
		appErr := apperrors.GetLearnByUsIdAndLimitErr.AppendMessage(err)
		return nil, appErr
	}

	userId, err := uuid.Parse(getWordsReq.ID)
	if err != nil {
		appErr := apperrors.GetLearnByUsIdAndLimitErr.AppendMessage(err)
		return nil, appErr
	}

	return us.UserRepository.GetLearnByIDAndLimit(ctx, &userId, quantity)
}

func (us *userInteractor) GetUserById(ctx context.Context, id string) (*models.User, error) {
	userId, err := uuid.Parse(id)
	if err != nil {
		appErr := apperrors.GetUserByIdErr.AppendMessage(err)
		return nil, appErr
	}

	return us.UserRepository.GetUserById(ctx, &userId)
}

func (us *userInteractor) MoveWordToLearned(ctx context.Context, userID, wordID string) error {
	userId, err := uuid.Parse(userID)
	if err != nil {
		appErr := apperrors.MoveWordToLearnedErr.AppendMessage(err)
		return appErr
	}

	user := &models.User{ID: &userId}
	wordId, err := strconv.Atoi(wordID)
	if err != nil {
		appErr := apperrors.MoveWordToLearnedErr.AppendMessage(err)
		return appErr
	}

	word := &models.Word{ID: wordId}
	return us.UserRepository.MoveWordToLearned(ctx, user, word)
}

func (us *userInteractor) AddWordToLearn(ctx context.Context, userID, wordID string) error {
	userId, err := uuid.Parse(userID)
	if err != nil {
		appErr := apperrors.AddWordToLearnErr.AppendMessage(err)
		return appErr
	}

	user := &models.User{ID: &userId}
	wordId, err := strconv.Atoi(wordID)
	if err != nil {
		appErr := apperrors.AddWordToLearnErr.AppendMessage(err)
		return appErr
	}

	word := &models.Word{ID: wordId}
	return us.UserRepository.AddWordToLearn(ctx, user, word)
}

func (us *userInteractor) DeleteLearnFromUserById(ctx context.Context, userID, wordID string) error {
	userId, err := uuid.Parse(userID)
	if err != nil {
		appErr := apperrors.DeleteLearnFromUserByIdErr.AppendMessage(err)
		return appErr
	}

	user := &models.User{ID: &userId}
	wordId, err := strconv.Atoi(wordID)
	if err != nil {
		appErr := apperrors.DeleteLearnFromUserByIdErr.AppendMessage(err)
		return appErr
	}

	word := &models.Word{ID: wordId}
	return us.UserRepository.DeleteLearnWordFromUserByWordID(ctx, user, word)
}

func (us *userInteractor) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return us.UserRepository.GetAllUsers(ctx)
}

func randomPassword() string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	const (
		chars          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		passwordLength = 7
	)

	password := make([]byte, passwordLength)
	for i := range password {
		password[i] = chars[random.Intn(len(chars))]
	}

	return string(password)
}
