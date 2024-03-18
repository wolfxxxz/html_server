package server

import "server/internal/domain/models"

const (
	header              = "templates/header.html"
	footer              = "templates/footer.html"
	translate           = "translate"
	home                = "home"
	registration        = "registration"
	authentification    = "authentification"
	logout              = "logout"
	userInfo            = "user_info"
	test                = "test"
	learn               = "learn"
	errMes              = "err"
	registrate          = "registrate"
	testThematic        = "test_thematic"
	testThematicHandler = "test_thematic_handler"
	updateLib           = "update_lib"
	updateUser          = "update_user"
	updateUserPassword  = "update_user_password"
	usersInfo           = "users_info"
)

var hashTableUsers = make(map[string]*models.User)
