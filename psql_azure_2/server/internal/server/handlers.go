package server

import (
	"io"
	"net/http"
	"server/internal/apperrors"
	"server/internal/domain/models"
	"server/internal/domain/requests"
	comparer "server/internal/s_comparer"
	"server/internal/services"
	"strings"

	"github.com/gorilla/mux"
)

func (srv *server) homeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := srv.tmpls[home].ExecuteTemplate(w, home, nil)
		if err != nil {
			appErr := apperrors.HomeHandlerErr.AppendMessage(err)
			srv.logger.Error(err)
			srv.respondErr(w, appErr)
			return
		}
	}
}

func (srv *server) getTranslationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/*
			text := r.URL.Query().Get("text")
			libServiceFastTr := services.NewLibraryService(srv.repoLibrary, srv.logger)
			words, err := libServiceFastTr.GetTranslationByWord(r.Context(), &requests.TranslationRequest{Word: text})
			if err != nil {
				appErr := err.(*apperrors.AppError)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}

			response := map[string]string{"translation": words[0].English}
			jsonResponse, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)
		*/

		if r.Method == http.MethodGet {
			err := srv.tmpls[translate].ExecuteTemplate(w, translate, nil)
			if err != nil {
				appErr := apperrors.GetTranslationHandlerErr.AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}
		}

		if r.Method == http.MethodPost {
			wordToTranslate := r.FormValue("word")
			if len(wordToTranslate) == 0 {
				http.Redirect(w, r, "/translate", http.StatusSeeOther)
				return
			}

			libService := services.NewLibraryService(srv.repoLibrary, srv.repoWords, srv.repoBackUp, srv.logger)
			words, err := libService.GetTranslationByWord(r.Context(), wordToTranslate)
			if err != nil {
				appErr := err.(*apperrors.AppError)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
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

			if err := srv.tmpls[translate].ExecuteTemplate(w, translate, responseData); err != nil {
				appErr := apperrors.GetTranslationHandlerErr.AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}
		}
	}
}

//-------------CRUD USER------------------

func (srv *server) getUserByIdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _, ok := srv.getIdANdRoleFromRequest(r)
		if !ok {
			appErr := apperrors.GetUserByIdHandlerErr.AppendMessage("UserIdErr")
			srv.logger.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		userService := services.NewUserService(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.logger)
		user, err := userService.GetUserById(r.Context(), userID)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.logger.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		hashTableUsers[userID] = user
		err = srv.tmpls[userInfo].ExecuteTemplate(w, userInfo, user)
		if err != nil {
			appErr := apperrors.GetUserByIdHandlerErr.AppendMessage(err)
			srv.logger.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}
	}
}

func (srv *server) logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := srv.tmpls[logout].ExecuteTemplate(w, logout, nil)
			if err != nil {
				appErr := apperrors.LogoutHandlerErr.AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}
		}

		if r.Method == http.MethodPost {
			cookies := r.Cookies()
			var token string
			for _, cookie := range cookies {
				if cookie.Name == "user_token_translator" {
					token = cookie.Value
				}
			}

			if token == "" {
				appErr := apperrors.LogoutHandlerErr.AppendMessage("there is nothing to blacklist")
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}

			srv.blacklist.AddToken(token)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
}

func (srv *server) loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := srv.tmpls[authentification].ExecuteTemplate(w, authentification, nil)
			if err != nil {
				appErr := apperrors.LoginHandlerErr.AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}
		}

		if r.Method == http.MethodPost {
			email := r.FormValue("email")
			password := r.FormValue("password")

			userService := services.NewUserService(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.logger)
			loginRequest := &requests.LoginRequest{
				Password: password,
				Email:    email,
			}

			getUserResp, err := userService.SignInUserWithJWT(r.Context(), loginRequest, srv.config.Server.SecretKey, srv.config.Server.ExpirationJWTInSeconds)
			if err != nil {
				appErr := err.(*apperrors.AppError)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr) //--------------------------------------------------------------
				return
			}

			if err := srv.tmpls[authentification].ExecuteTemplate(w, authentification, getUserResp); err != nil {
				appErr := apperrors.LoginHandlerErr.AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}
		}
	}
}

func (srv *server) createUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := srv.tmpls[registration].ExecuteTemplate(w, registration, nil)
			if err != nil {
				appErr := apperrors.CreateUserHandlerErr.AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}
		}

		if r.Method == http.MethodPost {
			email := r.FormValue("email")
			name := r.FormValue("name")
			lastName := r.FormValue("last-name")
			password := r.FormValue("password")
			//role := r.FormValue("role")
			//about := r.FormValue("about")

			userService := services.NewUserService(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.logger)
			createUserRequest := &requests.CreateUserRequest{
				Email:    email,
				Name:     name,
				LastName: lastName,
				Password: password,
				Role:     "user",
			}
			getUserResp, err := userService.CreateUser(r.Context(), createUserRequest)
			if err != nil {
				appErr := err.(*apperrors.AppError)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}

			createUserRequest.Id = getUserResp.UserId
			if err := srv.tmpls[registration].ExecuteTemplate(w, registration, createUserRequest); err != nil {
				appErr := apperrors.CreateUserHandlerErr.AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}
		}
	}
}

func (srv *server) updateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userService := services.NewUserService(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.logger)
		userID, _, ok := srv.getIdANdRoleFromRequest(r)
		if !ok {
			appErr := apperrors.UpdateUserHandlerErr.AppendMessage("UserIdErr")
			srv.logger.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		if r.Method == http.MethodGet {
			user := hashTableUsers[userID]
			err := srv.tmpls[updateUser].ExecuteTemplate(w, updateUser, user)
			if err != nil {
				appErr := apperrors.UpdateUserHandlerErr.AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}
		}

		if r.Method == http.MethodPost {
			email := r.FormValue("email")
			name := r.FormValue("name")
			lastName := r.FormValue("last-name")
			password := r.FormValue("password")
			role := r.FormValue("role")
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
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}

			if err := srv.tmpls[registration].ExecuteTemplate(w, registration, createUserRequest); err != nil {
				appErr := apperrors.UpdateUserHandlerErr.AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}
		}
	}
}

func (srv *server) updateUserPasswordHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userService := services.NewUserService(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.logger)
		userID, _, ok := srv.getIdANdRoleFromRequest(r)
		if !ok {
			appErr := apperrors.UpdateUserPasswordHandlerErr.AppendMessage("UserIdErr")
			srv.logger.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		if r.Method == http.MethodGet {
			user := hashTableUsers[userID]
			err := srv.tmpls[updateUserPassword].ExecuteTemplate(w, updateUserPassword, user)
			if err != nil {
				appErr := apperrors.UpdateUserPasswordHandlerErr.AppendMessage(err)
				srv.logger.Error(appErr)
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
				srv.logger.Error(appErr)
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
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}
		}
	}
}

// ----------------tests------------------------
func (srv *server) testHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _, ok := srv.getIdANdRoleFromRequest(r)
		if !ok {
			appErr := apperrors.TestHandlerErr.AppendMessage("UserIdErr")
			srv.logger.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		userService := services.NewUserService(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.logger)
		getWordsByUsIdAndLimitRequest := &requests.GetWordsByUsIdAndLimitRequest{ID: userID, Limit: "5"}

		if r.Method == http.MethodGet {
			words, err := userService.GetWordsByUsIdAndLimit(r.Context(), getWordsByUsIdAndLimitRequest)
			if err != nil {
				appErr := err.(*apperrors.AppError)
				srv.logger.Error(appErr)
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
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}
		}

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				appErr := apperrors.TestHandlerErr.AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}

			compareServ := comparer.NewComparer(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.logger)
			err = compareServ.CompareTestWords(r, userID)
			if err != nil {
				appErr := apperrors.TestHandlerErr.AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}

			err = srv.tmpls[test].ExecuteTemplate(w, test, comparer.HashTableWords[userID])
			if err != nil {
				appErr := apperrors.TestHandlerErr.AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}
		}
	}
}

func (srv *server) learnHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _, ok := srv.getIdANdRoleFromRequest(r)
		if !ok {
			appErr := apperrors.LearnHandlerErr.AppendMessage("UserIdErr")
			srv.logger.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		userService := services.NewUserService(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.logger)
		getWordsByUsIdAndLimitRequest := &requests.GetWordsByUsIdAndLimitRequest{ID: userID, Limit: "5"}

		if r.Method == http.MethodGet {
			words, err := userService.GetLearnByUsIdAndLimit(r.Context(), getWordsByUsIdAndLimitRequest)
			if err != nil {
				appErr := err.(*apperrors.AppError)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}

			pageData := &models.TestPageData{
				Words: words,
				//Result: results
				//TestPassed: false,
			}

			comparer.HashTableWordsLearn[userID] = pageData
			err = srv.tmpls[learn].ExecuteTemplate(w, learn, pageData)
			if err != nil {
				appErr := apperrors.LearnHandlerErr.AppendMessage("User ID Err")
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}
		}

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				appErr := apperrors.LearnHandlerErr.AppendMessage("User ID Err")
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}

			compareServ := comparer.NewComparer(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.logger)
			err = compareServ.CompareLearnWords(r, userID)
			if err != nil {
				appErr := apperrors.LearnHandlerErr.AppendMessage("User ID Err")
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}

			err = srv.tmpls[learn].ExecuteTemplate(w, learn, comparer.HashTableWordsLearn[userID])
			if err != nil {
				appErr := apperrors.LearnHandlerErr
				srv.logger.Error(appErr)
				srv.respondErr(w, &appErr)
				return
			}
		}
	}
}

//------------------thematic tests--------------------

func (srv *server) themesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceLibrary := services.NewLibraryService(srv.repoLibrary, srv.repoWords, srv.repoBackUp, srv.logger)
		if r.Method == http.MethodGet {
			topics, err := serviceLibrary.GetAllTopics()
			if err != nil {
				appErr := apperrors.ThemesHandlerErr.AppendMessage(err)
				srv.logger.Error(err)
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
				srv.logger.Error(err)
				srv.respondErr(w, appErr)
				return
			}
		}
	}
}

func (srv *server) testUniversalHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _, ok := srv.getIdANdRoleFromRequest(r)
		if !ok {
			appErr := apperrors.TestUniversalHandlerErr.AppendMessage("UserIdErr")
			srv.logger.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		userService := services.NewUserService(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.logger)
		getWordsByUsIdAndLimitRequest := &requests.GetWordsByUsIdAndLimitRequest{ID: userID, Limit: "5"}
		vars := mux.Vars(r)
		topic := vars["theme"]
		topicGet := strings.ReplaceAll(topic, "_", " ")

		if r.Method == http.MethodGet {
			words, err := userService.GetWordsByUserIdAndLimitAndTopic(r.Context(), getWordsByUsIdAndLimitRequest, topicGet)
			if err != nil {
				appErr := err.(*apperrors.AppError)
				srv.logger.Error(appErr)
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
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}
		}

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				appErr := apperrors.TestUniversalHandlerErr.AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}

			compareServ := comparer.NewComparer(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.logger)
			err = compareServ.CompareTestWords(r, userID)
			if err != nil {
				appErr := apperrors.TestUniversalHandlerErr.AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}

			err = srv.tmpls[testThematicHandler].ExecuteTemplate(w, testThematicHandler, comparer.HashTableWords[userID])
			if err != nil {
				appErr := apperrors.TestUniversalHandlerErr.AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}
		}
	}
}

//------------Update Library role admin----------------------

func (srv *server) updateLibraryHandler() http.HandlerFunc {
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

		if r.Method == http.MethodGet {
			err := srv.tmpls[updateLib].ExecuteTemplate(w, updateLib, nil)
			if err != nil {
				appErr := apperrors.UpdateLibraryHandlerErr.AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}
		}

		if r.Method == http.MethodPost {
			r.ParseMultipartForm(10 << 20) // Размер файла не должен превышать 10 MB
			file, _, err := r.FormFile("fileToUpload")
			if err != nil {
				appErr := apperrors.UpdateLibraryHandlerErr.AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}

			defer file.Close()

			libraryService := services.NewLibraryService(srv.repoLibrary, srv.repoWords, srv.repoBackUp, srv.logger)
			err = libraryService.UpdateLibraryOldAndNewWordsByMultyFile(r.Context(), &file)
			if err != nil {
				appErr := err.(*apperrors.AppError)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}

			http.Redirect(w, r, "/user-info", http.StatusSeeOther)
		}
	}
}

func (srv *server) downloadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, role, ok := srv.getIdANdRoleFromRequest(r)
		if !ok {
			appErr := apperrors.DownloadHandlerErr.AppendMessage("UserIdErr")
			srv.logger.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		if role != "admin" {
			appErr := apperrors.DownloadHandlerErr.AppendMessage("UserIdErr")
			srv.logger.Error(appErr)
			srv.respondAuthorizateErr(w, appErr)
			return
		}
		libraryService := services.NewLibraryService(srv.repoLibrary, srv.repoWords, srv.repoBackUp, srv.logger)
		file, err := libraryService.DownloadXLXFromDb()
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.logger.Error(appErr)
			return
		}

		defer file.Close()

		w.Header().Set("Content-Disposition", "attachment; filename=library.xlsx")
		w.Header().Set("Content-Type", "application/octet-stream")

		if _, err := io.Copy(w, file); err != nil {
			srv.logger.Error(err)
			return
		}
	}
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

func (srv *server) getAllUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srv.logger.Info("getAllUsersHandler started")
		_, role, ok := srv.getIdANdRoleFromRequest(r)
		if !ok {
			appErr := apperrors.DownloadHandlerErr.AppendMessage("UserIdErr")
			srv.logger.Error(appErr)
			srv.respondErr(w, appErr)
			return
		}

		if role != "admin" {
			appErr := apperrors.DownloadHandlerErr.AppendMessage("UserIdErr")
			srv.logger.Error(appErr)
			srv.respondAuthorizateErr(w, appErr)
			return
		}

		userService := services.NewUserService(srv.repoWords, srv.repoUsers, srv.repoLibrary, srv.logger)
		users, err := userService.GetAllUsers(r.Context())
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.logger.Error(appErr)
			return
		}

		err = srv.tmpls[usersInfo].ExecuteTemplate(w, usersInfo, users)
		if err != nil {
			appErr := apperrors.RespondErr.AppendMessage(err)
			srv.logger.Error(appErr)
		}

		srv.logger.Info("getAllUsersHandler sent tamplates")
	}
}

func (srv *server) getIdANdRoleFromRequest(r *http.Request) (string, string, bool) {
	id := r.Context().Value(contextKeyID)
	if id == nil {
		appErr := apperrors.GetIdANdRoleFromRequestErr.AppendMessage("UserIdErr")
		srv.logger.Error(appErr)
		return "", "", false
	}

	userID, ok := id.(string)
	if !ok {
		appErr := apperrors.GetIdANdRoleFromRequestErr.AppendMessage("UserIdErr")
		srv.logger.Error(appErr)
		return "", "", false
	}

	role := r.Context().Value(contextKeyRole)
	if role == nil {
		appErr := apperrors.GetIdANdRoleFromRequestErr.AppendMessage("UserIdErr")
		srv.logger.Error(appErr)
		return "", "", false
	}

	userRole, ok := role.(string)
	if !ok {
		appErr := apperrors.GetIdANdRoleFromRequestErr.AppendMessage("UserRoleErr")
		srv.logger.Error(appErr)
		return "", "", false
	}

	return userID, userRole, true
}

//--------------Responds---------------------------

func (srv *server) respondErr(w http.ResponseWriter, appErr *apperrors.AppError) {
	err := srv.tmpls[errMes].ExecuteTemplate(w, errMes, appErr)
	if err != nil {
		appErr := apperrors.RespondErr.AppendMessage(err)
		srv.logger.Error(appErr)
	}
}

func (srv *server) respondRegistrateErr(w http.ResponseWriter, appErr *apperrors.AppError) {
	err := srv.tmpls[registrate].ExecuteTemplate(w, registrate, appErr)
	if err != nil {
		appErr := apperrors.RespondErr.AppendMessage(err)
		srv.logger.Error(appErr)
	}
}

func (srv *server) respondAuthorizateErr(w http.ResponseWriter, appErr *apperrors.AppError) {
	err := srv.tmpls[registrate].ExecuteTemplate(w, registrate, appErr)
	if err != nil {
		appErr := apperrors.RespondErr.AppendMessage(err)
		srv.logger.Error(appErr)
	}
}
