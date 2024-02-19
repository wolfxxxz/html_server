package apperrors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Message  string
	Code     string
	HTTPCode int
}

func NewAppError() *AppError {
	return &AppError{}
}

var (
	SetupDatabaseErr = AppError{
		Message: "Failed SetupDatabaseErr",
		Code:    database,
	}
	EnvConfigLoadError = AppError{
		Message: "Failed to load env file",
		Code:    envInit,
	}
	EnvConfigParseError = AppError{
		Message: "Failed to parse env file",
		Code:    envParse,
	}
	InitPostgressErr = AppError{
		Message: "Failed to InitPostgress",
		Code:    envParse,
	}
	NewLoggerErr = AppError{
		Message: "Failed to NewLog",
		Code:    log,
	}
	SetLevelErr = AppError{
		Message: "Failed to SetLevelErr",
		Code:    log,
	}
	JWTMiddleware = AppError{
		Message:  "Failed to JWTMiddlewareErr",
		Code:     middleware,
		HTTPCode: http.StatusUnauthorized,
	}
	GetAllFromBackUpErr = AppError{
		Message: "Failed to GetAllFromBackUp",
		Code:    backUpRepo,
	}
	SaveAllAsJsonErr = AppError{
		Message: "Failed to SaveAllAsJson",
		Code:    backUpRepo,
	}
	InsertWordsLibraryErr = AppError{
		Message: "Failed to InsertWordsLibraryErr",
		Code:    repoLibrary,
	}
	GetTranslationEnglLikeErr = AppError{
		Message: "Failed to GetTranslationEnglLikeErr",
		Code:    repoLibrary,
	}
	GetTranslationEnglErr = AppError{
		Message: "Failed to GetTranslationEnglErr",
		Code:    repoLibrary,
	}
	GetTranslationRusLikeErr = AppError{
		Message: "Failed to GetTranslationRusLikeErr",
		Code:    repoLibrary,
	}
	GetTranslationRusErr = AppError{
		Message: "Failed to GetTranslationRusErr",
		Code:    repoLibrary,
	}
	GetAllWordsLibErr = AppError{
		Message: "Failed to GetAllWords",
		Code:    repoLibrary,
	}
	UpdateWordErr = AppError{
		Message: "Failed to UpdateWord",
		Code:    repoLibrary,
	}
	GetAllWordsErr = AppError{
		Message: "Failed to GetAllWords",
		Code:    repoLibrary,
	}
	GetWordsWhereRAErr = AppError{
		Message: "Failed to GetWordsWhereRA",
		Code:    repoLibrary,
	}
	UpdateUserErr = AppError{
		Message: "Failed to UpdateUserErr",
		Code:    repoUsers,
	}
	CreateUserErr = AppError{
		Message: "Failed to CreateUser",
		Code:    repoUsers,
	}
	GetUserByEmailErr = AppError{
		Message: "Failed to GetUserByEmailErr",
		Code:    repoUsers,
	}
	GetWordsByIDAndLimitErr = AppError{
		Message: "Failed to GetWordsByIDAndLimitErr",
		Code:    repoUsers,
	}
	GetLearnByIDAndLimitErr = AppError{
		Message: "Failed to GetLearnByIDAndLimitErr",
		Code:    repoUsers,
	}
	DeleteLearnWordFromUserByWordErr = AppError{
		Message: "Failed to DeleteLearnWordFromUserByWordErr",
		Code:    repoUsers,
	}
	AddWordToLearnRepoErr = AppError{
		Message: "Failed to AddWordToLearnRepoErr",
		Code:    repoUsers,
	}
	DeleteLearnByUserIDAndLearnIDHandlerErr = AppError{
		Message: "Failed to deleteLearnByUserIDAndLearnIDHandlerErr",
		Code:    handlers,
	}
	GetWordsByUserIDAndLimitHandlerErr = AppError{
		Message: "Failed to GetWordsByUserIDAndLimitHandlerErr",
		Code:    handlers,
	}
	GetTranslationHandlerErr = AppError{
		Message: "Failed to GetTranslationHandlerErr",
		Code:    handlers,
	}
	LoginHandlerErr = AppError{
		Message: "Failed to LoginHandlerErr",
		Code:    handlers,
	}
	LogoutHandlerErr = AppError{
		Message: "Failed to LogoutHandlerErr",
		Code:    handlers,
	}
	GetUserByIdHandlerErr = AppError{
		Message: "Failed to GetUserByIdHandlerErr",
		Code:    handlers,
	}
	CreateUserHandlerErr = AppError{
		Message: "Failed to CreateUserHandlerErr",
		Code:    handlers,
	}
	GetLearnByUserIDAndLimitHandlerErr = AppError{
		Message: "Failed to GetLearnByUserIDAndLimitHandlerErr",
		Code:    handlers,
	}
	MoveWordToLearnedHandlerErr = AppError{
		Message: "Failed to MoveWordToLearnedHandlerErr",
		Code:    handlers,
	}
	AddWordToLearnHandlerErr = AppError{
		Message: "Failed to AddWordToLearnHandlerErr",
		Code:    handlers,
	}
	DeleteLearnFromUserByIdErr = AppError{
		Message: "Failed to DeleteLearnFromUserByIdErr",
		Code:    services,
	}
	GetTranslationByWordErr = AppError{
		Message: "Failed to GetTranslationByWordErr",
		Code:    services,
	}
	SignInUserWithJWTErr = AppError{
		Message: "Failed to SignInUserWithJWTErr",
		Code:    services,
	}
	GetUserByIdErr = AppError{
		Message: "Failed to GetUserByIdErr",
		Code:    services,
	}
	MoveWordToLearnedErr = AppError{
		Message: "Failed to MoveWordToLearnedErr",
		Code:    services,
	}
	AddWordToLearnErr = AppError{
		Message: "Failed to AddWordToLearnErr",
		Code:    services,
	}
	ClaimJWTTokenErr = AppError{
		Message: "Failed to ClaimJWTTokenErr",
		Code:    services,
	}
	HashPasswordErr = AppError{
		Message: "Failed to HashPasswordErr",
		Code:    services,
	}
	GetWordsByUsIdAndLimitServiceErr = AppError{
		Message: "Failed to GetWordsByUsIdAndLimitServiceErr",
		Code:    services,
	}
)

func (appError *AppError) Error() string {
	return appError.Code + ": " + appError.Message
}

func (appError *AppError) AppendMessage(anyErrs ...interface{}) *AppError {
	return &AppError{
		Message: fmt.Sprintf("%v : %v", appError.Message, anyErrs),
		Code:    appError.Code,
	}
}

func IsAppError(err1 error, err2 *AppError) bool {
	err, ok := err1.(*AppError)
	if !ok {
		return false
	}

	return err.Code == err2.Code
}
