package controller

import (
	"net/http"
	"server/internal/apperrors"
	"server/internal/config"
	"server/internal/domain/models"
	"server/internal/domain/requests"
	"server/internal/infrastructure/webtemplate.go"
	"server/internal/usercase/interactor"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type handleController struct {
	libraryInteractor interactor.LibraryInteractor
	userInteractor    interactor.UserInteractor
	log               *logrus.Logger
	config            *config.Config
	tmpls             *webtemplate.WebTemplates
}

type HandleController interface {
	HomeHandler(c echo.Context) error
	GetTranslationHandler(c echo.Context) error
	QuickAnswerHandler(c echo.Context) error
	CreateUserHandler(c echo.Context) error
}

func NewUserController(ui interactor.UserInteractor, li interactor.LibraryInteractor, log *logrus.Logger, confg *config.Config, tmpls *webtemplate.WebTemplates) HandleController {
	return &handleController{li, ui, log, confg, tmpls}
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

//--------------Responds---------------------------

func (srv *handleController) respondErr(w http.ResponseWriter, appErr *apperrors.AppError) {
	err := srv.tmpls.Templates[errMes].ExecuteTemplate(w, errMes, appErr)
	if err != nil {
		appErr := apperrors.RespondErr.AppendMessage(err)
		srv.log.Error(appErr)
	}
}

/*
func (srv *userController) respondRegistrateErr(w http.ResponseWriter, appErr *apperrors.AppError) {
	err := srv.tmpls[registrate].ExecuteTemplate(w, registrate, appErr)
	if err != nil {
		appErr := apperrors.RespondErr.AppendMessage(err)
		srv.log.Error(appErr)
	}
}

func (srv *userController) respondAuthorizateErr(w http.ResponseWriter, appErr *apperrors.AppError) {
	err := srv.tmpls[registrate].ExecuteTemplate(w, registrate, appErr)
	if err != nil {
		appErr := apperrors.RespondErr.AppendMessage(err)
		srv.log.Error(appErr)
	}
}
*/

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

/*
func (srv *userController) getUserByIdHandler(c echo.Context) error {
	userID, _, ok := srv.getIdANdRoleFromRequest(c.r)
	if !ok {
		appErr := apperrors.GetUserByIdHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondErr(c.Response().Writer, appErr)
		return appErr
	}

	userService := services.NewUserService(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.sender, srv.logger)
	user, err := userService.GetUserById(c.r.Context(), userID)
	if err != nil {
		appErr := err.(*apperrors.AppError)
		srv.log.Error(appErr)
		srv.respondErr(c.Response().Writer, appErr)
		return appErr
	}

	hashTableUsers[userID] = user
	err = srv.tmpls[userInfo].ExecuteTemplate(c.Response().Writer, userInfo, user)
	if err != nil {
		appErr := apperrors.GetUserByIdHandlerErr.AppendMessage(err)
		srv.log.Error(appErr)
		srv.respondErr(c.Response().Writer, appErr)
		return appErr
	}

}

func (srv *userController) logoutHandler(c echo.Context) error {
	if c.r.Method == http.MethodGet {
		err := srv.tmpls[logout].ExecuteTemplate(w, logout, nil)
		if err != nil {
			appErr := apperrors.LogoutHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return appErr
		}
	}

	if c.r.Method == http.MethodPost {
		cookies := c.r.Cookies()
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

		srv.blacklist.AddToken(token)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	}

}

func (srv *userController) loginHandler(c echo.Context) error {
	if c.r.Method == http.MethodGet {
		err := srv.tmpls[authentification].ExecuteTemplate(w, authentification, nil)
		if err != nil {
			appErr := apperrors.LoginHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return appErr
		}
	}

	if r.Method == http.MethodPost {
		email := c.r.FormValue("email")
		password := c.r.FormValue("password")

		userService := services.NewUserService(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.sender, srv.logger)
		loginRequest := &requests.LoginRequest{
			Password: password,
			Email:    email,
		}

		getUserResp, err := userService.SignInUserWithJWT(r.Context(), loginRequest, srv.config.Server.SecretKey, srv.config.Server.ExpirationJWTInSeconds)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr) //--------------------------------------------------------------
			return appErr
		}

		if err := srv.tmpls[authentification].ExecuteTemplate(w, authentification, getUserResp); err != nil {
			appErr := apperrors.LoginHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return appErr
		}

	}
}


func (srv *userController) restoreUserPasswordHandler(c echo.Context) error {
	userService := services.NewUserService(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.sender, srv.logger)

	if c.r.Method == http.MethodGet {
		err := srv.tmpls[restoreUserPassword].ExecuteTemplate(c.Response().Writer, restoreUserPassword, nil)
		if err != nil {
			appErr := apperrors.UpdateUserHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return appErr
		}
	}

	if c.r.Method == http.MethodPost {
		email := c.r.FormValue("email")

		err := userService.RestoreUserPassword(c.r.Context(), email)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return appErr
		}

		http.Redirect(c.Response().Writer, c.r, "/login", http.StatusSeeOther)
	}

}

func (srv *userController) updateUserHandler(c echo.Context) error {
	userService := services.NewUserService(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.sender, srv.logger)
	userID, _, ok := srv.getIdANdRoleFromRequest(r)
	if !ok {
		appErr := apperrors.UpdateUserHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondErr(c.Response().Writer, appErr)
		return appErr
	}

	if c.r.Method == http.MethodGet {
		user := hashTableUsers[userID]
		err := srv.tmpls[updateUser].ExecuteTemplate(w, updateUser, user)
		if err != nil {
			appErr := apperrors.UpdateUserHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(c.Response().Writer, appErr)
			return appErr
		}
	}

	if r.Method == http.MethodPost {
		email := c.r.FormValue("email")
		name := c.r.FormValue("name")
		lastName := c.r.FormValue("last-name")
		password := c.r.FormValue("password")
		role := c.r.FormValue("role")
		//about := r.FormValue("about")

		createUserRequest := &requests.CreateUserRequest{
			Email:    email,
			Name:     name,
			LastName: lastName,
			Password: password,
			Role:     role,
		}

		user := hashTableUsers[userID]
		err := userService.UpdateUserById(r.Context(), user, createUserRequest)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		if err := srv.tmpls[registration].ExecuteTemplate(w, registration, createUserRequest); err != nil {
			appErr := apperrors.UpdateUserHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}
	}

}

func (srv *userController) updateUserPasswordHandler(c echo.Context) error {
	userService := services.NewUserService(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.sender, srv.logger)
	userID, _, ok := srv.getIdANdRoleFromRequest(r)
	if !ok {
		appErr := apperrors.UpdateUserPasswordHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondErr(w, appErr)
		return
	}

	if r.Method == http.MethodGet {
		user := hashTableUsers[userID]
		err := srv.tmpls[updateUserPassword].ExecuteTemplate(w, updateUserPassword, user)
		if err != nil {
			appErr := apperrors.UpdateUserPasswordHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}
	}

	if r.Method == http.MethodPost {
		oldPass := r.FormValue("old_password")
		newPass := r.FormValue("new_password")
		newPassSecond := r.FormValue("new_password_second")

		user := hashTableUsers[userID]
		err := userService.UpdateUserPasswordById(r.Context(), user, oldPass, newPass, newPassSecond)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		createUserRequest := &requests.CreateUserRequest{
			Email:    user.Email,
			Name:     user.Name,
			LastName: user.LastName,
			Password: newPass,
			Role:     user.Role,
		}

		if err := srv.tmpls[registration].ExecuteTemplate(w, registration, createUserRequest); err != nil {
			appErr := apperrors.UpdateUserPasswordHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}
	}

}

// ----------------tests------------------------
func (srv *userController) testHandler(c echo.Context) error {
	userID, _, ok := srv.getIdANdRoleFromRequest(r)
	if !ok {
		appErr := apperrors.TestHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondErr(w, appErr)
		return
	}

	userService := services.NewUserService(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.sender, srv.logger)
	getWordsByUsIdAndLimitRequest := &requests.GetWordsByUsIdAndLimitRequest{ID: userID, Limit: "5"}

	if r.Method == http.MethodGet {
		words, err := userService.GetWordsByUsIdAndLimit(r.Context(), getWordsByUsIdAndLimitRequest)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		pageData := &models.TestPageData{
			Words:      words,
			TestPassed: false,
		}

		comparer.HashTableWords[userID] = pageData

		err = srv.tmpls[test].ExecuteTemplate(w, test, pageData)
		if err != nil {
			appErr := apperrors.TestHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			appErr := apperrors.TestHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		compareServ := comparer.NewComparer(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.sender, srv.logger)
		err = compareServ.CompareTestWords(r, userID)
		if err != nil {
			appErr := apperrors.TestHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		err = srv.tmpls[test].ExecuteTemplate(w, test, comparer.HashTableWords[userID])
		if err != nil {
			appErr := apperrors.TestHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}
	}

}

func (srv *userController) learnHandler(c echo.Context) error {
	userID, _, ok := srv.getIdANdRoleFromRequest(r)
	if !ok {
		appErr := apperrors.LearnHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondErr(w, appErr)
		return
	}

	userService := services.NewUserService(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.sender, srv.logger)
	getWordsByUsIdAndLimitRequest := &requests.GetWordsByUsIdAndLimitRequest{ID: userID, Limit: "5"}

	if r.Method == http.MethodGet {
		words, err := userService.GetLearnByUsIdAndLimit(r.Context(), getWordsByUsIdAndLimitRequest)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		pageData := &models.TestPageData{
			Words: words,
		}

		comparer.HashTableWordsLearn[userID] = pageData
		err = srv.tmpls[learn].ExecuteTemplate(w, learn, pageData)
		if err != nil {
			appErr := apperrors.LearnHandlerErr.AppendMessage("User ID Err")
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			appErr := apperrors.LearnHandlerErr.AppendMessage("User ID Err")
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		compareServ := comparer.NewComparer(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.sender, srv.logger)
		err = compareServ.CompareLearnWords(r, userID)
		if err != nil {
			appErr := apperrors.LearnHandlerErr.AppendMessage("User ID Err")
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		err = srv.tmpls[learn].ExecuteTemplate(w, learn, comparer.HashTableWordsLearn[userID])
		if err != nil {
			appErr := apperrors.LearnHandlerErr
			srv.log.Error(appErr)
			srv.respondErr(w, &appErr)
			return
		}
	}

}

//------------------thematic tests--------------------

func (srv *userController) themesHandler(c echo.Context) error {
	serviceLibrary := services.NewLibraryService(srv.repoLibrary, srv.repoWords, srv.repoBackUp, srv.logger)
	if r.Method == http.MethodGet {
		topics, err := serviceLibrary.GetAllTopics()
		if err != nil {
			appErr := apperrors.ThemesHandlerErr.AppendMessage(err)
			srv.log.Error(err)
			srv.respondErr(w, appErr)
			return
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

		err = srv.tmpls[testThematic].ExecuteTemplate(w, testThematic, topicsAndBytes)
		if err != nil {
			appErr := apperrors.ThemesHandlerErr.AppendMessage(err)
			srv.log.Error(err)
			srv.respondErr(w, appErr)
			return
		}
	}

}

func (srv *userController) testUniversalHandler(c echo.Context) error {
	userID, _, ok := srv.getIdANdRoleFromRequest(r)
	if !ok {
		appErr := apperrors.TestUniversalHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondErr(w, appErr)
		return
	}

	userService := services.NewUserService(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.sender, srv.logger)
	getWordsByUsIdAndLimitRequest := &requests.GetWordsByUsIdAndLimitRequest{ID: userID, Limit: "5"}
	vars := mux.Vars(r)
	topic := vars["theme"]
	topicGet := strings.ReplaceAll(topic, "_", " ")

	if r.Method == http.MethodGet {
		words, err := userService.GetWordsByUserIdAndLimitAndTopic(r.Context(), getWordsByUsIdAndLimitRequest, topicGet)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		pageData := &models.TestPageData{
			Topic: topic,
			Words: words,
			//Result: results,
			TestPassed: false,
		}

		comparer.HashTableWords[userID] = pageData

		err = srv.tmpls[testThematicHandler].ExecuteTemplate(w, testThematicHandler, pageData)
		if err != nil {
			appErr := apperrors.TestUniversalHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			appErr := apperrors.TestUniversalHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		compareServ := comparer.NewComparer(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.sender, srv.logger)
		err = compareServ.CompareTestWords(r, userID)
		if err != nil {
			appErr := apperrors.TestUniversalHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		err = srv.tmpls[testThematicHandler].ExecuteTemplate(w, testThematicHandler, comparer.HashTableWords[userID])
		if err != nil {
			appErr := apperrors.TestUniversalHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}
	}

}

//------------Update Library role admin----------------------

func (srv *userController) updateLibraryHandler(c echo.Context) error {
	_, role, ok := srv.getIdANdRoleFromRequest(r)
	if !ok {
		appErr := apperrors.UpdateLibraryHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondErr(w, appErr)
		return
	}

	if role != "admin" {
		appErr := apperrors.UpdateLibraryHandlerErr.AppendMessage("ask for help in contacts")
		srv.log.Error(appErr)
		srv.respondErr(w, appErr)
		return
	}

	if r.Method == http.MethodGet {
		err := srv.tmpls[updateLib].ExecuteTemplate(w, updateLib, nil)
		if err != nil {
			appErr := apperrors.UpdateLibraryHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}
	}

	if r.Method == http.MethodPost {
		r.ParseMultipartForm(10 << 20) // Размер файла не должен превышать 10 MB
		file, _, err := r.FormFile("fileToUpload")
		if err != nil {
			appErr := apperrors.UpdateLibraryHandlerErr.AppendMessage(err)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		defer file.Close()

		libraryService := services.NewLibraryService(srv.repoLibrary, srv.repoWords, srv.repoBackUp, srv.logger)
		err = libraryService.UpdateLibraryOldAndNewWordsByMultyFile(r.Context(), &file)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.log.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		http.Redirect(w, r, "/user-info", http.StatusSeeOther)
	}

}

func (srv *userController) downloadHandler(c echo.Context) error {
	_, role, ok := srv.getIdANdRoleFromRequest(r)
	if !ok {
		appErr := apperrors.DownloadHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondErr(w, appErr)
		return
	}

	if role != "admin" {
		appErr := apperrors.DownloadHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondAuthorizateErr(w, appErr)
		return
	}
	libraryService := services.NewLibraryService(srv.repoLibrary, srv.repoWords, srv.repoBackUp, srv.logger)
	file, err := libraryService.DownloadXLXFromDb()
	if err != nil {
		appErr := err.(*apperrors.AppError)
		srv.log.Error(appErr)
		return
	}

	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename=library.xlsx")
	w.Header().Set("Content-Type", "application/octet-stream")

	if _, err := io.Copy(w, file); err != nil {
		srv.log.Error(err)
		return
	}

}


func (srv *userController) getAllUsersHandler(c echo.Context) error {
	srv.log.Info("getAllUsersHandler started")
	_, role, ok := srv.getIdANdRoleFromRequest(r)
	if !ok {
		appErr := apperrors.DownloadHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondErr(w, appErr)
		return
	}

	if role != "admin" {
		appErr := apperrors.DownloadHandlerErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		srv.respondAuthorizateErr(w, appErr)
		return
	}

	userService := services.NewUserService(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.sender, srv.logger)
	users, err := userService.GetAllUsers(r.Context())
	if err != nil {
		appErr := err.(*apperrors.AppError)
		srv.log.Error(appErr)
		return
	}

	err = srv.tmpls[usersInfo].ExecuteTemplate(w, usersInfo, users)
	if err != nil {
		appErr := apperrors.RespondErr.AppendMessage(err)
		srv.log.Error(appErr)
	}

	srv.log.Info("getAllUsersHandler sent tamplates")

}

func (srv *userController) getIdANdRoleFromRequest(r *http.Request) (string, string, bool) {
	id := r.Context().Value(contextKeyID)
	if id == nil {
		appErr := apperrors.GetIdANdRoleFromRequestErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		return "", "", false
	}

	userID, ok := id.(string)
	if !ok {
		appErr := apperrors.GetIdANdRoleFromRequestErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		return "", "", false
	}

	role := r.Context().Value(contextKeyRole)
	if role == nil {
		appErr := apperrors.GetIdANdRoleFromRequestErr.AppendMessage("UserIdErr")
		srv.log.Error(appErr)
		return "", "", false
	}

	userRole, ok := role.(string)
	if !ok {
		appErr := apperrors.GetIdANdRoleFromRequestErr.AppendMessage("UserRoleErr")
		srv.log.Error(appErr)
		return "", "", false
	}

	return userID, userRole, true
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
