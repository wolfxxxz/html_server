package server

import (
	"encoding/json"
	"net/http"
	"server/internal/apperrors"
	"server/internal/domain/mappers"
	"server/internal/domain/requests"
	"server/internal/domain/responses"
	"server/internal/services"

	"github.com/gorilla/mux"
)

func (srv *server) createUserHandler() http.HandlerFunc {
	srv.logger.Info("createUserHandler has been initiated.")
	return func(w http.ResponseWriter, r *http.Request) {
		createUserRequest := &requests.CreateUserRequest{}
		err := srv.decode(r, createUserRequest)
		if err != nil {
			appErr := apperrors.CreateUserHandlerErr.AppendMessage(err)
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusBadRequest)
			return
		}

		srv.logger.Infof("createUserHandler has been invoked. Name %v, Email %v", createUserRequest.Name, createUserRequest.Email)
		userService := services.NewUserService(srv.repoUsers, srv.repoLibrary, srv.logger)
		getUserResp, err := userService.CreateUser(r.Context(), createUserRequest)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusInternalServerError)
			return
		}

		srv.logger.Infof("createUserHandler has been processed. Response: %+v", getUserResp)
		srv.respond(w, getUserResp, http.StatusCreated)
	}
}

func (srv *server) loginHandler() http.HandlerFunc {
	srv.logger.Info("loginHandler has been initiated.")
	return func(w http.ResponseWriter, r *http.Request) {
		loginRequest := &requests.LoginRequest{}
		err := srv.decode(r, loginRequest)
		if err != nil {
			appErr := apperrors.LoginHandlerErr.AppendMessage(err)
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusBadRequest)
			return
		}

		srv.logger.Infof("loginHandler has been invoked.  Email %v", loginRequest.Email)
		userService := services.NewUserService(srv.repoUsers, srv.repoLibrary, srv.logger)
		getUserResp, err := userService.SignInUserWithJWT(r.Context(), loginRequest, srv.config.Server.SecretKey, srv.config.Server.ExpirationJWTInSeconds)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusInternalServerError)
			return
		}

		srv.logger.Infof("loginHandler has been processed. Response: %+v", getUserResp)
		srv.respond(w, getUserResp, http.StatusOK)
	}
}

func (srv *server) logoutHandler() http.HandlerFunc {
	srv.logger.Info("logoutHandler has been initiated.")
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			appErr := apperrors.LogoutHandlerErr.AppendMessage("HEADER GET Authorization")
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusBadRequest)
			return
		}

		srv.blacklist.AddToken(token)
		srv.logger.Info("Token has been blacklisted")
		srv.respond(w, "token deleted", http.StatusOK)
	}
}

func (srv *server) getUserByIdHandler() http.HandlerFunc {
	srv.logger.Info("getUserByIdHandler has been initiated.")
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := mux.Vars(r)["user_id"]
		if !ok {
			appErr := apperrors.GetUserByIdHandlerErr.AppendMessage("Vars User ID")
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusBadRequest)
			return
		}

		srv.logger.Infof("getUserByIdHandler has been invoked. Id %v", userID)
		userService := services.NewUserService(srv.repoUsers, srv.repoLibrary, srv.logger)
		words, err := userService.GetUserById(r.Context(), userID)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusInternalServerError)
			return
		}

		srv.logger.Infof("getUserByIdHandler has been processed. Response: %+v", words)
		srv.respond(w, words, http.StatusOK)
	}
}

func (srv *server) getWordsByUserIDAndLimitHandler() http.HandlerFunc {
	srv.logger.Info("getWordsByUserIDAndLimitHandler has been initiated.")
	return func(w http.ResponseWriter, r *http.Request) {
		getWordsByUsIdAndLimitRequest := &requests.GetWordsByUsIdAndLimitRequest{}
		err := srv.decode(r, getWordsByUsIdAndLimitRequest)
		if err != nil {
			appErr := apperrors.GetWordsByUserIDAndLimitHandlerErr.AppendMessage(err)
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusBadRequest)
			return
		}

		srv.logger.Infof("getWordsByUserIDAndLimitHandler has been invoked. Id %v, Limit %v", getWordsByUsIdAndLimitRequest.ID, getWordsByUsIdAndLimitRequest.Limit)
		userService := services.NewUserService(srv.repoUsers, srv.repoLibrary, srv.logger)
		words, err := userService.GetWordsByUsIdAndLimit(r.Context(), getWordsByUsIdAndLimitRequest)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusInternalServerError)
			return
		}

		srv.logger.Infof("getWordsByUserIDAndLimitHandler has been processed. Response: %+v", words)
		srv.respond(w, words, http.StatusOK)
	}
}

func (srv *server) getLearnByUserIDAndLimitHandler() http.HandlerFunc {
	srv.logger.Info("getLearnByUserIDAndLimitHandler has been initiated.")
	return func(w http.ResponseWriter, r *http.Request) {
		getWordsByUsIdAndLimitRequest := &requests.GetWordsByUsIdAndLimitRequest{}
		err := srv.decode(r, getWordsByUsIdAndLimitRequest)
		if err != nil {
			appErr := apperrors.GetLearnByUserIDAndLimitHandlerErr.AppendMessage(err)
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusBadRequest)
			return
		}

		srv.logger.Infof("getLearnByUserIDAndLimitHandler has been invoked. Id %v, Limit %v", getWordsByUsIdAndLimitRequest.ID, getWordsByUsIdAndLimitRequest.Limit)
		userService := services.NewUserService(srv.repoUsers, srv.repoLibrary, srv.logger)
		words, err := userService.GetLearnByUsIdAndLimit(r.Context(), getWordsByUsIdAndLimitRequest)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusInternalServerError)
			return
		}

		srv.logger.Infof("getLearnByUserIDAndLimitHandler has been processed. Response: %+v", words)
		srv.respond(w, words, http.StatusOK)
	}
}

func (srv *server) moveWordToLearnedHandler() http.HandlerFunc {
	srv.logger.Info("moveWordToLearnedHandler has been initiated.")
	return func(w http.ResponseWriter, r *http.Request) {
		deleteWordFromUserByIDRequest := &requests.DeleteWordFromUserByIDRequest{}
		err := srv.decode(r, deleteWordFromUserByIDRequest)
		if err != nil {
			appErr := apperrors.MoveWordToLearnedHandlerErr.AppendMessage(err)
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusBadRequest)
			return
		}

		srv.logger.Infof("moveWordToLearnedHandler has been invoked. User Id %v, Word Id %v", deleteWordFromUserByIDRequest.UserID, deleteWordFromUserByIDRequest.WordID)
		userService := services.NewUserService(srv.repoUsers, srv.repoLibrary, srv.logger)
		err = userService.MoveWordToLearned(r.Context(), deleteWordFromUserByIDRequest)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusInternalServerError)
			return
		}

		result := &responses.Result{Answer: "success"}
		srv.logger.Infof("moveWordToLearnedHandler has been processed. Response: %+v", result)
		srv.respond(w, result, http.StatusOK)
	}
}

func (srv *server) addWordToLearnHandler() http.HandlerFunc {
	srv.logger.Info("addWordToLearnHandler has been initiated.")
	return func(w http.ResponseWriter, r *http.Request) {
		deleteWordFromUserByIDRequest := &requests.DeleteWordFromUserByIDRequest{}
		err := srv.decode(r, deleteWordFromUserByIDRequest)
		if err != nil {
			appErr := apperrors.AddWordToLearnHandlerErr.AppendMessage(err)
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusBadRequest)
			return
		}

		srv.logger.Infof("addWordToLearnHandler has been invoked. User Id %v, Word Id %v", deleteWordFromUserByIDRequest.UserID, deleteWordFromUserByIDRequest.WordID)
		userService := services.NewUserService(srv.repoUsers, srv.repoLibrary, srv.logger)
		err = userService.AddWordToLearn(r.Context(), deleteWordFromUserByIDRequest)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusInternalServerError)
			return
		}

		result := &responses.Result{Answer: "success"}
		srv.logger.Infof("addWordToLearnHandler has been processed. Response: %+v", result)
		srv.respond(w, result, http.StatusOK)
	}
}

func (srv *server) deleteLearnByUserIDAndLearnIDHandler() http.HandlerFunc {
	srv.logger.Info("deleteLearnByUserIDAndLearnIDHandler has been initiated.")
	return func(w http.ResponseWriter, r *http.Request) {
		deleteWordFromUserByIDRequest := &requests.DeleteWordFromUserByIDRequest{}
		err := srv.decode(r, deleteWordFromUserByIDRequest)
		if err != nil {
			appErr := apperrors.DeleteLearnByUserIDAndLearnIDHandlerErr.AppendMessage(err)
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusBadRequest)
			return
		}

		srv.logger.Infof("deleteLearnByUserIDAndLearnIDHandler has been invoked. User Id %v, Word Id %v", deleteWordFromUserByIDRequest.UserID, deleteWordFromUserByIDRequest.WordID)
		userService := services.NewUserService(srv.repoUsers, srv.repoLibrary, srv.logger)
		err = userService.DeleteLearnFromUserById(r.Context(), deleteWordFromUserByIDRequest)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusInternalServerError)
			return
		}

		result := &responses.Result{Answer: "success"}
		srv.logger.Infof("deleteLearnByUserIDAndLearnIDHandler has been processed. Response: %+v", result)
		srv.respond(w, result, http.StatusOK)
	}
}

func (srv *server) getTranslationHandler() http.HandlerFunc {
	srv.logger.Info("getTranslationHandler has been initiated.")
	return func(w http.ResponseWriter, r *http.Request) {
		translationReq := &requests.TranslationRequest{}
		err := srv.decode(r, translationReq)
		if err != nil {
			appErr := apperrors.GetTranslationHandlerErr.AppendMessage(err)
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusBadRequest)
			return
		}

		srv.logger.Infof("getTranslationHandler has been invoked.  Word  %v", translationReq.Word)
		libService := services.NewLibraryService(srv.repoLibrary, srv.logger)
		words, err := libService.GetTranslationByWord(r.Context(), translationReq)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusInternalServerError)
			return
		}

		wordsResp := mappers.MapLibraryToWordsGetTranslResponse(words)
		srv.logger.Infof("getTranslationHandler has been processed. Response : %v words", len(wordsResp))
		srv.respond(w, wordsResp, http.StatusOK)
	}
}

func (srv *server) decode(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (srv *server) respond(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	if data == nil {
		return
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		srv.logger.Error(err)
	}
}
