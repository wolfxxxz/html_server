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
	PingEveryMinutsErr = AppError{
		Message: "Failed PingEveryMinutsErr",
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
	MapMultipartToXLSErr = AppError{
		Message: "Failed to GetAllFromBackUp",
		Code:    mapers,
	}
	InitializeTemplatesErr = AppError{
		Message: "Failed to InitializeTemplatesErr",
		Code:    server,
	}
	GetAllWordsFromBackUpXlsxErr = AppError{
		Message: "Failed to GetAllFromBackUp",
		Code:    backUpRepo,
	}
	GetAllFromBackUpErr = AppError{
		Message: "Failed to GetAllFromBackUp",
		Code:    backUpRepo,
	}
	SaveAllAsJsonErr = AppError{
		Message: "Failed to SaveAllAsJson",
		Code:    backUpRepo,
	}
	SaveWordsAsXLSXErr = AppError{
		Message: "Failed to SaveWordsAsXLSXErr",
		Code:    backUpRepo,
	}
	InsertWordsLibraryErr = AppError{
		Message: "Failed to InsertWordsLibraryErr",
		Code:    repoLibrary,
	}
	InsertWordLibraryErr = AppError{
		Message: "Failed to InsertWordLibraryErr",
		Code:    repoLibrary,
	}
	UpdateWordRowAffectedErr = AppError{
		Message: "Failed to UpdateWordRowAffectedErr",
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
	InitWordsMapErr = AppError{
		Message: "Failed to InitWordsMapErr",
		Code:    repoLibrary,
	}
	UpdateWordsMapErr = AppError{
		Message: "Failed to UpdateWordsMapErr",
		Code:    repoLibrary,
	}
	JWTMiddleware = AppError{
		Message:  "Failed to JWTMiddlewareErr",
		Code:     middleware,
		HTTPCode: http.StatusUnauthorized,
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
	GetAllTopicsErr = AppError{
		Message: "Failed to GetAllTopicsErr",
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
	UpdateUserByIdErr = AppError{
		Message: "Failed to UpdateUserByIdErr",
		Code:    repoUsers,
	}
	GetWordsByUserIdAndLimitAndTopicErr = AppError{
		Message: "Failed to GetWordsByUserIdAndLimitAndTopicErr",
		Code:    repoUsers,
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
	DeleteLearnWordFromUserByWordIDErr = AppError{
		Message: "Failed to DeleteLearnWordFromUserByWordErr",
		Code:    repoUsers,
	}
	UpdateLibraryHandlerErr = AppError{
		Message: "Failed to UpdateLibraryHandlerErr",
		Code:    repoUsers,
	}
	AddWordToLearnRepoErr = AppError{
		Message: "Failed to AddWordToLearnRepoErr",
		Code:    repoUsers,
	}
	GetAllUsersErr = AppError{
		Message: "Failed to GetAllUsersErr",
		Code:    handlers,
	}
	HomeHandlerErr = AppError{
		Message: "Failed to HomeHandlerErr",
		Code:    handlers,
	}
	DownloadHandlerErr = AppError{
		Message: "Failed to HomeHandlerErr",
		Code:    handlers,
	}
	TestUniversalHandlerErr = AppError{
		Message: "Failed to TestUniversalHandlerErr",
		Code:    handlers,
	}
	GetIdANdRoleFromRequestErr = AppError{
		Message: "Failed to GetIdANdRoleFromRequestErr",
		Code:    handlers,
	}
	ThemesHandlerErr = AppError{
		Message: "Failed to ThemesHandlerErr",
		Code:    handlers,
	}
	TestHandlerErr = AppError{
		Message: "Failed to testHandlerErr",
		Code:    handlers,
	}
	GetTranslationHandlerErr = AppError{
		Message: "Failed to GetTranslationHandlerErr",
		Code:    handlers,
	}
	UpdateUserPasswordHandlerErr = AppError{
		Message: "Failed to UpdateUserPasswordHandlerErr",
		Code:    handlers,
	}
	UpdateUserHandlerErr = AppError{
		Message: "Failed to UpdateUserHandlerErr",
		Code:    handlers,
	}
	LoginHandlerErr = AppError{
		Message: "Failed to LoginHandlerErr",
		Code:    handlers,
	}
	RespondErr = AppError{
		Message: "Failed to RespondErr",
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
	LearnHandlerErr = AppError{
		Message: "Failed to LearnHandlerErr",
		Code:    handlers,
	}
	DeleteLearnFromUserByIdErr = AppError{
		Message: "Failed to DeleteLearnFromUserByIdErr",
		Code:    services,
	}

	ComparerTestErr = AppError{
		Message: "Failed to ComparerTestErr",
		Code:    services,
	}
	ComparerLearnErr = AppError{
		Message: "Failed to ComparerLearnErr",
		Code:    services,
	}
	GetLearnByUsIdAndLimitErr = AppError{
		Message: "Failed to GetLearnByUsIdAndLimitErr",
		Code:    services,
	}
	GetTranslationByWordErr = AppError{
		Message: "Failed to GetTranslationByWordErr",
		Code:    services,
	}
	SignInUserWithJWTErr = AppError{
		Message:  "Failed to SignInUserWithJWTErr",
		Code:     services,
		HTTPCode: http.StatusUnauthorized,
	}
	ClaimJWTTokenErr = AppError{
		Message: "Failed to ClaimJWTTokenErr",
		Code:    services,
	}
	UpdateUserPasswordByIdErr = AppError{
		Message: "Failed to UpdateUserPasswordByIdErr",
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
	AddWordsToUserErr = AppError{
		Message: "Failed to AddWordsToUserErr",
		Code:    services,
	}
	AddWordToLearnErr = AppError{
		Message: "Failed to AddWordsToUserErr",
		Code:    services,
	}
	GetAllTopicsLibServErr = AppError{
		Message: "Failed to GetAllTopicsLibServErr",
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
		Message:  fmt.Sprintf("%v : %v", appError.Message, anyErrs),
		Code:     appError.Code,
		HTTPCode: appError.HTTPCode,
	}
}

func IsAppError(err1 error, err2 *AppError) bool {
	err, ok := err1.(*AppError)
	if !ok {
		return false
	}

	return err.Code == err2.Code
}
