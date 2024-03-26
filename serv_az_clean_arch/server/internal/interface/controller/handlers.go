package controller

import (
	"io"
	"net/http"
	"server/internal/apperrors"
	"server/internal/config"
	"server/internal/domain/models"
	"server/internal/domain/requests"
	"server/internal/infrastructure/datastore"
	"server/internal/infrastructure/middleware"
	"server/internal/infrastructure/webtemplate.go"
	"server/internal/usercase/comparer"
	"server/internal/usercase/interactor"
	"strings"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

const (
	contextKeyRole string = "role"
	contextKeyID   string = "id"
)

type handleController struct {
	comparer          comparer.Comparer
	libraryInteractor interactor.LibraryInteractor
	userInteractor    interactor.UserInteractor
	hashDB            *datastore.HashDB
	log               *logrus.Logger
	config            *config.Config
	tmpls             *webtemplate.WebTemplates
}

type HandleController interface {
	HomeHandler(c echo.Context) error
	GetTranslationHandler(c echo.Context) error
	QuickAnswerHandler(c echo.Context) error
	CreateUserHandler(c echo.Context) error
	LoginHandler(c echo.Context) error
	LogoutHandler(blacklist *middleware.Blacklist) echo.HandlerFunc
	GetUserByIdHandler(c echo.Context) error
	RestoreUserPasswordHandler(c echo.Context) error
	UpdateUserHandler(c echo.Context) error
	UpdateUserPasswordHandler(c echo.Context) error
	UpdateLibraryHandler(c echo.Context) error
	DownloadHandler(c echo.Context) error
	GetAllUsersHandler(c echo.Context) error
	TestHandler(c echo.Context) error
	LearnHandler(c echo.Context) error
	ThemesHandler(c echo.Context) error
	TestUniversalHandler(c echo.Context) error
}

func NewHandlersController(comparer comparer.Comparer, ui interactor.UserInteractor, li interactor.LibraryInteractor, hashDB *datastore.HashDB, log *logrus.Logger, confg *config.Config, tmpls *webtemplate.WebTemplates) HandleController {
	return &handleController{comparer, li, ui, hashDB, log, confg, tmpls}
}

func (srv *handleController) HomeHandler(c echo.Context) error {
	err := srv.tmpls.Templates[home].ExecuteTemplate(c.Response().Writer, home, nil)
	if err != nil {
		appErr := apperrors.HomeHandlerErr.AppendMessage(err)
		srv.log.Error(err)
		srv.respondErr(c.Response().Writer, appErr)
		return appErr
	}

	return nil
}

func (srv *handleController) GetTranslationHandler(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		err := srv.tmpls.Templates[translate].ExecuteTemplate(c.Response().Writer, translate, nil)
		if err != nil {
			appErr := apperrors.GetTranslationHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return appErr
		}
	}

	if c.Request().Method == http.MethodPost {
		wordToTranslate := c.FormValue("word")
		if len(wordToTranslate) == 0 {
			http.Redirect(c.Response().Writer, c.Request(), "/translate", http.StatusSeeOther)
			return nil
		}

		//libService := services.NewLibraryService(srv.repoLibrary, srv.repoWords, srv.repoBackUp, srv.logger)
		words, err := srv.libraryInteractor.GetTranslationByWord(c.Request().Context(), wordToTranslate)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return appErr
		}

		if len(words) == 0 {
			word := &models.Library{
				English: "There is not such world in the library",
				Russian: "Такого слова нет в библиотеке",
			}

			words = append(words, word)
		}

		responseData := Rsvp{
			Words:   words,
			Word:    wordToTranslate,
			WordRus: words[0].English,
		}

		if err := srv.tmpls.Templates[translate].ExecuteTemplate(c.Response().Writer, translate, responseData); err != nil {
			appErr := apperrors.GetTranslationHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return appErr
		}
	}

	return nil
}

//--------------Respons---------------------------

func (srv *handleController) respondErr(w http.ResponseWriter, appErr *apperrors.AppError) {
	err := srv.tmpls.Templates[errMes].ExecuteTemplate(w, errMes, appErr)
	if err != nil {
		appErr := apperrors.RespondErr.AppendMessage(err)
		srv.log.Error(appErr)
	}
}

func (srv *handleController) respondRegistrateErr(w http.ResponseWriter, appErr *apperrors.AppError) {
	err := srv.tmpls.Templates[registrate].ExecuteTemplate(w, registrate, appErr)
	if err != nil {
		appErr := apperrors.RespondErr.AppendMessage(err)
		srv.log.Error(appErr)
	}
}

func (srv *handleController) respondAuthorizateErr(w http.ResponseWriter, appErr *apperrors.AppError) {
	err := srv.tmpls.Templates[registrate].ExecuteTemplate(w, registrate, appErr)
	if err != nil {
		appErr := apperrors.RespondErr.AppendMessage(err)
		srv.log.Error(appErr)
	}
}

func (srv *handleController) QuickAnswerHandler(c echo.Context) error {
	key := c.QueryParam("key")
	words, err := srv.libraryInteractor.GetTranslationByPieceOfWord(c.Request().Context(), key)
	if err != nil {
		srv.log.Error()
		return err
	}

	c.String(http.StatusOK, words)
	return nil
}

//-------------CRUD USER------------------

func (srv *handleController) CreateUserHandler(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		err := srv.tmpls.Templates[registration].ExecuteTemplate(c.Response().Writer, registration, nil)
		if err != nil {
			appErr := apperrors.CreateUserHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return appErr
		}
	}

	if c.Request().Method == http.MethodPost {
		email := c.FormValue("email")
		name := c.FormValue("name")
		lastName := c.FormValue("last-name")
		password := c.FormValue("password")
		//role := r.FormValue("role")
		//about := r.FormValue("about")

		createUserRequest := &requests.CreateUserRequest{
			Email:    email,
			Name:     name,
			LastName: lastName,
			Password: password,
			Role:     "user",
		}
		getUserResp, err := srv.userInteractor.CreateUser(c.Request().Context(), createUserRequest)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return appErr
		}

		createUserRequest.Id = getUserResp.UserId
		if err := srv.tmpls.Templates[registration].ExecuteTemplate(c.Response().Writer, registration, createUserRequest); err != nil {
			appErr := apperrors.CreateUserHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return appErr
		}
	}

	return nil
}

func (srv *handleController) LoginHandler(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		err := srv.tmpls.Templates[authentification].ExecuteTemplate(c.Response().Writer, authentification, nil)
		if err != nil {
			appErr := apperrors.LoginHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return appErr
		}
	}

	if c.Request().Method == http.MethodPost {
		email := c.Request().FormValue("email")
		password := c.Request().FormValue("password")

		loginRequest := &requests.LoginRequest{
			Password: password,
			Email:    email,
		}

		getUserResp, err := srv.userInteractor.SignInUserWithJWT(c.Request().Context(), loginRequest, srv.config.Server.SecretKey, srv.config.Server.ExpirationJWTInSeconds)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return appErr
		}

		if err := srv.tmpls.Templates[authentification].ExecuteTemplate(c.Response().Writer, authentification, getUserResp); err != nil {
			appErr := apperrors.LoginHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return appErr
		}

	}

	return nil
}

func (srv *handleController) LogoutHandler(blacklist *middleware.Blacklist) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method == http.MethodGet {
			err := srv.tmpls.Templates[logout].ExecuteTemplate(c.Response().Writer, logout, nil)
			if err != nil {
				appErr := apperrors.LogoutHandlerErr.AppendMessage(err)
				srv.log.Error(appErr)
				srv.respondErr(c.Response().Writer, appErr)
				return appErr
			}
		}

		if c.Request().Method == http.MethodPost {
			cookies := c.Request().Cookies()
			var token string
			for _, cookie := range cookies {
				if cookie.Name == "user_token_translator" {
					token = cookie.Value
				}
			}

			if token == "" {
				appErr := apperrors.LogoutHandlerErr.AppendMessage("there is nothing to blacklist")
				srv.log.Error(appErr)
				srv.respondErr(c.Response().Writer, appErr)
				return appErr
			}

			blacklist.AddToken(token)
			http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusSeeOther)
			return nil
		}

		return nil
	}
}

func (srv *handleController) GetUserByIdHandler(c echo.Context) error {
	userID, _, ok := srv.getIdANdRoleFromRequest(c)
	if !ok {
		appErr := apperrors.GetUserByIdHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondErr(c.Response().Writer, appErr)
		return appErr
	}

	user, err := srv.userInteractor.GetUserById(c.Request().Context(), userID)
	if err != nil {
		appErr := err.(*apperrors.AppError)
		srv.log.Error(appErr)
		srv.respondErr(c.Response().Writer, appErr)
		return appErr
	}

	srv.hashDB.DB[userID] = user
	err = srv.tmpls.Templates[userInfo].ExecuteTemplate(c.Response().Writer, userInfo, user)
	if err != nil {
		appErr := apperrors.GetUserByIdHandlerErr.AppendMessage(err)
		srv.log.Error(appErr)
		srv.respondErr(c.Response().Writer, appErr)
		return appErr
	}

	return nil
}

func (srv *handleController) RestoreUserPasswordHandler(c echo.Context) error {

	if c.Request().Method == http.MethodGet {
		err := srv.tmpls.Templates[restoreUserPassword].ExecuteTemplate(c.Response().Writer, restoreUserPassword, nil)
		if err != nil {
			appErr := apperrors.UpdateUserHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return appErr
		}
	}

	if c.Request().Method == http.MethodPost {
		email := c.Request().FormValue("email")

		err := srv.userInteractor.RestoreUserPassword(c.Request().Context(), email)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return appErr
		}

		http.Redirect(c.Response().Writer, c.Request(), "/login", http.StatusSeeOther)
	}

	return nil
}

func (srv *handleController) UpdateUserHandler(c echo.Context) error {
	userID, _, ok := srv.getIdANdRoleFromRequest(c)
	if !ok {
		appErr := apperrors.UpdateUserHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondErr(c.Response().Writer, appErr)
		return appErr
	}

	if c.Request().Method == http.MethodGet {
		user := srv.hashDB.DB[userID]
		err := srv.tmpls.Templates[updateUser].ExecuteTemplate(c.Response().Writer, updateUser, user)
		if err != nil {
			appErr := apperrors.UpdateUserHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return appErr
		}
	}

	if c.Request().Method == http.MethodPost {
		email := c.Request().FormValue("email")
		name := c.Request().FormValue("name")
		lastName := c.Request().FormValue("last-name")
		password := c.Request().FormValue("password")
		role := c.Request().FormValue("role")
		//about := r.FormValue("about")

		createUserRequest := &requests.CreateUserRequest{
			Email:    email,
			Name:     name,
			LastName: lastName,
			Password: password,
			Role:     role,
		}

		user := srv.hashDB.DB[userID]
		err := srv.userInteractor.UpdateUserById(c.Request().Context(), user, createUserRequest)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}

		if err := srv.tmpls.Templates[registration].ExecuteTemplate(c.Response().Writer, registration, createUserRequest); err != nil {
			appErr := apperrors.UpdateUserHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}
	}

	return nil
}

func (srv *handleController) UpdateUserPasswordHandler(c echo.Context) error {
	userID, _, ok := srv.getIdANdRoleFromRequest(c)
	if !ok {
		appErr := apperrors.UpdateUserPasswordHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondErr(c.Response().Writer, appErr)
		return nil
	}

	if c.Request().Method == http.MethodGet {
		user := srv.hashDB.DB[userID]
		err := srv.tmpls.Templates[updateUserPassword].ExecuteTemplate(c.Response().Writer, updateUserPassword, user)
		if err != nil {
			appErr := apperrors.UpdateUserPasswordHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}
	}

	if c.Request().Method == http.MethodPost {
		oldPass := c.Request().FormValue("old_password")
		newPass := c.Request().FormValue("new_password")
		newPassSecond := c.Request().FormValue("new_password_second")

		user := srv.hashDB.DB[userID]
		err := srv.userInteractor.UpdateUserPasswordById(c.Request().Context(), user, oldPass, newPass, newPassSecond)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}

		createUserRequest := &requests.CreateUserRequest{
			Email:    user.Email,
			Name:     user.Name,
			LastName: user.LastName,
			Password: newPass,
			Role:     user.Role,
		}

		if err := srv.tmpls.Templates[registration].ExecuteTemplate(c.Response().Writer, registration, createUserRequest); err != nil {
			appErr := apperrors.UpdateUserPasswordHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}
	}

	return nil

}

// -------------------------doesn't work---------------------------

// ----------------tests------------------------
func (srv *handleController) TestHandler(c echo.Context) error {
	userID, _, ok := srv.getIdANdRoleFromRequest(c)
	if !ok {
		appErr := apperrors.TestHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondErr(c.Response().Writer, appErr)
		return nil
	}

	getWordsByUsIdAndLimitRequest := &requests.GetWordsByUsIdAndLimitRequest{ID: userID, Limit: "5"}

	if c.Request().Method == http.MethodGet {
		words, err := srv.userInteractor.GetWordsByUsIdAndLimit(c.Request().Context(), getWordsByUsIdAndLimitRequest)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}

		pageData := &models.TestPageData{
			Words:      words,
			TestPassed: false,
		}

		comparer.HashTableWords[userID] = pageData

		err = srv.tmpls.Templates[test].ExecuteTemplate(c.Response().Writer, test, pageData)
		if err != nil {
			appErr := apperrors.TestHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}
	}

	if c.Request().Method == http.MethodPost {
		err := c.Request().ParseForm()
		if err != nil {
			appErr := apperrors.TestHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}

		err = srv.comparer.CompareTestWords(c.Request(), userID)
		if err != nil {
			appErr := apperrors.TestHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}

		err = srv.tmpls.Templates[test].ExecuteTemplate(c.Response().Writer, test, comparer.HashTableWords[userID])
		if err != nil {
			appErr := apperrors.TestHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}
	}

	return nil
}

func (srv *handleController) LearnHandler(c echo.Context) error {
	userID, _, ok := srv.getIdANdRoleFromRequest(c)
	if !ok {
		appErr := apperrors.LearnHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondErr(c.Response().Writer, appErr)
		return nil
	}

	getWordsByUsIdAndLimitRequest := &requests.GetWordsByUsIdAndLimitRequest{ID: userID, Limit: "5"}

	if c.Request().Method == http.MethodGet {
		words, err := srv.userInteractor.GetLearnByUsIdAndLimit(c.Request().Context(), getWordsByUsIdAndLimitRequest)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}

		pageData := &models.TestPageData{
			Words: words,
		}

		comparer.HashTableWordsLearn[userID] = pageData
		err = srv.tmpls.Templates[learn].ExecuteTemplate(c.Response().Writer, learn, pageData)
		if err != nil {
			appErr := apperrors.LearnHandlerErr.AppendMessage("User ID Err")
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}
	}

	if c.Request().Method == http.MethodPost {
		err := c.Request().ParseForm()
		if err != nil {
			appErr := apperrors.LearnHandlerErr.AppendMessage("User ID Err")
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}

		err = srv.comparer.CompareLearnWords(c.Request(), userID)
		if err != nil {
			appErr := apperrors.LearnHandlerErr.AppendMessage("User ID Err")
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}

		err = srv.tmpls.Templates[learn].ExecuteTemplate(c.Response().Writer, learn, comparer.HashTableWordsLearn[userID])
		if err != nil {
			appErr := apperrors.LearnHandlerErr
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, &appErr)
			return nil
		}
	}

	return nil
}

//------------------thematic tests--------------------

func (srv *handleController) ThemesHandler(c echo.Context) error {

	topics, err := srv.libraryInteractor.GetAllTopics()
	if err != nil {
		appErr := apperrors.ThemesHandlerErr.AppendMessage(err)
		srv.log.Error(err)
		srv.respondErr(c.Response().Writer, appErr)
		return nil
	}

	topicsAndBytes := []struct {
		TopStr       string
		TopUnderLine string
	}{}

	for _, top := range topics {
		underLine := strings.ReplaceAll(top, " ", "_")

		biStr := struct {
			TopStr       string
			TopUnderLine string
		}{
			TopStr:       top,
			TopUnderLine: underLine,
		}

		topicsAndBytes = append(topicsAndBytes, biStr)
	}

	err = srv.tmpls.Templates[testThematic].ExecuteTemplate(c.Response().Writer, testThematic, topicsAndBytes)
	if err != nil {
		appErr := apperrors.ThemesHandlerErr.AppendMessage(err)
		srv.log.Error(err)
		srv.respondErr(c.Response().Writer, appErr)
		return nil
	}

	return nil
}

func (srv *handleController) TestUniversalHandler(c echo.Context) error {
	userID, _, ok := srv.getIdANdRoleFromRequest(c)
	if !ok {
		appErr := apperrors.TestUniversalHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondErr(c.Response().Writer, appErr)
		return nil
	}

	getWordsByUsIdAndLimitRequest := &requests.GetWordsByUsIdAndLimitRequest{ID: userID, Limit: "5"}
	//vars := mux.Vars(c.Request())
	topic := c.Param("theme")
	topicGet := strings.ReplaceAll(topic, "_", " ")
	//-------------------------------------------------------------------------------------------------------------------------
	if c.Request().Method == http.MethodGet {
		words, err := srv.userInteractor.GetWordsByUserIdAndLimitAndTopic(c.Request().Context(), getWordsByUsIdAndLimitRequest, topicGet)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}

		pageData := &models.TestPageData{
			Topic: topic,
			Words: words,
			//Result: results,
			TestPassed: false,
		}

		comparer.HashTableWords[userID] = pageData

		err = srv.tmpls.Templates[testThematicHandler].ExecuteTemplate(c.Response().Writer, testThematicHandler, pageData)
		if err != nil {
			appErr := apperrors.TestUniversalHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}
	}

	if c.Request().Method == http.MethodPost {
		err := c.Request().ParseForm()
		if err != nil {
			appErr := apperrors.TestUniversalHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}

		err = srv.comparer.CompareTestWords(c.Request(), userID)
		if err != nil {
			appErr := apperrors.TestUniversalHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}

		err = srv.tmpls.Templates[testThematicHandler].ExecuteTemplate(c.Response().Writer, testThematicHandler, comparer.HashTableWords[userID])
		if err != nil {
			appErr := apperrors.TestUniversalHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}
	}

	return nil
}

//------------Update Library role admin----------------------

func (srv *handleController) UpdateLibraryHandler(c echo.Context) error {
	_, role, ok := srv.getIdANdRoleFromRequest(c)
	if !ok {
		appErr := apperrors.UpdateLibraryHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondErr(c.Response().Writer, appErr)
		return nil
	}

	if role != "admin" {
		appErr := apperrors.UpdateLibraryHandlerErr.AppendMessage("ask for help in contacts")
		srv.log.Error(appErr)
		srv.respondErr(c.Response().Writer, appErr)
		return nil
	}

	if c.Request().Method == http.MethodGet {
		err := srv.tmpls.Templates[updateLib].ExecuteTemplate(c.Response().Writer, updateLib, nil)
		if err != nil {
			appErr := apperrors.UpdateLibraryHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}
	}

	if c.Request().Method == http.MethodPost {
		c.Request().ParseMultipartForm(10 << 20) // Размер файла не должен превышать 10 MB
		file, _, err := c.Request().FormFile("fileToUpload")
		if err != nil {
			appErr := apperrors.UpdateLibraryHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return nil
		}

		defer file.Close()

		err = srv.libraryInteractor.UpdateLibraryOldAndNewWordsByMultyFile(c.Request().Context(), &file)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.log.Error(appErr)
			srv.respondErr(c.Response(), appErr)
			return nil
		}

		http.Redirect(c.Response().Writer, c.Request(), "/user-info", http.StatusSeeOther)
	}

	return nil
}

func (srv *handleController) DownloadHandler(c echo.Context) error {
	_, role, ok := srv.getIdANdRoleFromRequest(c)
	if !ok {
		appErr := apperrors.DownloadHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondErr(c.Response().Writer, appErr)
		return nil
	}

	if role != "admin" {
		appErr := apperrors.DownloadHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondAuthorizateErr(c.Response().Writer, appErr)
		return nil
	}
	file, err := srv.libraryInteractor.DownloadXLXFromDb()
	if err != nil {
		appErr := err.(*apperrors.AppError)
		srv.log.Error(appErr)
		return nil
	}

	defer file.Close()

	c.Response().Writer.Header().Set("Content-Disposition", "attachment; filename=library.xlsx")
	c.Response().Writer.Header().Set("Content-Type", "application/octet-stream")

	if _, err := io.Copy(c.Response().Writer, file); err != nil {
		srv.log.Error(err)
		return nil
	}

	return nil
}

func (srv *handleController) GetAllUsersHandler(c echo.Context) error {
	srv.log.Info("getAllUsersHandler started")
	_, role, ok := srv.getIdANdRoleFromRequest(c)
	if !ok {
		appErr := apperrors.DownloadHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondErr(c.Response().Writer, appErr)
		return nil
	}

	if role != "admin" {
		appErr := apperrors.DownloadHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondAuthorizateErr(c.Response().Writer, appErr)
		return nil
	}

	users, err := srv.userInteractor.GetAllUsers(c.Request().Context())
	if err != nil {
		appErr := err.(*apperrors.AppError)
		srv.log.Error(appErr)
		return nil
	}

	err = srv.tmpls.Templates[usersInfo].ExecuteTemplate(c.Response().Writer, usersInfo, users)
	if err != nil {
		appErr := apperrors.RespondErr.AppendMessage(err)
		srv.log.Error(appErr)
	}

	srv.log.Info("getAllUsersHandler sent tamplates")
	return nil
}

func (srv *handleController) getIdANdRoleFromRequest(c echo.Context) (string, string, bool) {
	// err := c.Get("err").(string)
	// if err != "" {
	// 	appErr := apperrors.GetIdANdRoleFromRequestErr.AppendMessage(err)
	// 	srv.log.Error(appErr)
	// 	return "", "", false
	// }
	id := c.Get(contextKeyID).(string)
	if id == "" {
		appErr := apperrors.GetIdANdRoleFromRequestErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		return "", "", false
	}

	role := c.Get(contextKeyRole).(string)
	if role == "" {
		appErr := apperrors.GetIdANdRoleFromRequestErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		return "", "", false
	}

	return id, role, true
}

/* drop library and users
func (srv *server) dropLibrary() http.HandlerFunc { // hasn't been implicit yet

		return func(w http.ResponseWriter, r *http.Request) {
			_, role, ok := srv.getIdANdRoleFromRequest(r)
			if !ok {
				appErr := apperrors.UpdateLibraryHandlerErr.AppendMessage("UserIdErr")
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}

			if role != "admin" {
				appErr := apperrors.UpdateLibraryHandlerErr.AppendMessage("ask for help in contacts")
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}

			http.Redirect(w, r, "/library-drop", http.StatusSeeOther)
		}
	}

func (srv *server) dropUsers() http.HandlerFunc { // hasn't been implicit yet

		return func(w http.ResponseWriter, r *http.Request) {
			_, role, ok := srv.getIdANdRoleFromRequest(r)
			if !ok {
				appErr := apperrors.UpdateLibraryHandlerErr.AppendMessage("UserIdErr")
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}

			if role != "admin" {
				appErr := apperrors.UpdateLibraryHandlerErr.AppendMessage("ask for help in contacts")
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}

			http.Redirect(w, r, "/library-drop", http.StatusSeeOther)
		}
	}
*/
