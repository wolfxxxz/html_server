package webtemplate

import (
	"html/template"
	"server/internal/apperrors"

	"github.com/sirupsen/logrus"
)

type WebTemplates struct {
	Templates map[string]*template.Template
}

func InitializeTemplates(logger *logrus.Logger) (*WebTemplates, error) {
	tmplsList := make(map[string]*template.Template)

	tmpl, err := template.ParseFiles("templates/translate.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		logger.Error(appErr)
		return nil, appErr
	}
	tmplsList[translate] = tmpl

	tmpl, err = template.ParseFiles("templates/home.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		logger.Error(appErr)
		return nil, appErr
	}
	tmplsList[home] = tmpl

	tmpl, err = template.ParseFiles("templates/registration.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		logger.Error(appErr)
		return nil, appErr
	}
	tmplsList[registration] = tmpl

	tmpl, err = template.ParseFiles("templates/authentification.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		logger.Error(appErr)
		return nil, appErr
	}
	tmplsList[authentification] = tmpl

	tmpl, err = template.ParseFiles("templates/logout.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		logger.Error(appErr)
		return nil, appErr
	}
	tmplsList[logout] = tmpl

	tmpl, err = template.ParseFiles("templates/user_info.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		logger.Error(appErr)
		return nil, appErr
	}
	tmplsList[userInfo] = tmpl

	tmpl, err = template.ParseFiles("templates/test.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		logger.Error(appErr)
		return nil, appErr
	}
	tmplsList[test] = tmpl

	tmpl, err = template.ParseFiles("templates/test_thematic_handler.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		logger.Error(appErr)
		return nil, appErr
	}
	tmplsList[testThematicHandler] = tmpl

	tmpl, err = template.ParseFiles("templates/learn.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		logger.Error(appErr)
		return nil, appErr
	}
	tmplsList[learn] = tmpl

	tmpl, err = template.ParseFiles("templates/err.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		logger.Error(appErr)
		return nil, appErr
	}
	tmplsList[errMes] = tmpl

	tmpl, err = template.ParseFiles("templates/registrate.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		logger.Error(appErr)
		return nil, appErr
	}
	tmplsList[registrate] = tmpl

	tmpl, err = template.ParseFiles("templates/test_thematic.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		logger.Error(appErr)
		return nil, appErr
	}
	tmplsList[testThematic] = tmpl

	tmpl, err = template.ParseFiles("templates/update_lib.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		logger.Error(appErr)
		return nil, appErr
	}
	tmplsList[updateLib] = tmpl

	tmpl, err = template.ParseFiles("templates/update_user.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		logger.Error(appErr)
		return nil, appErr
	}
	tmplsList[updateUser] = tmpl

	tmpl, err = template.ParseFiles("templates/update_user_password.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		logger.Error(appErr)
		return nil, appErr
	}
	tmplsList[updateUserPassword] = tmpl

	tmpl, err = template.ParseFiles("templates/users_info.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		logger.Error(appErr)
		return nil, appErr
	}
	tmplsList[usersInfo] = tmpl
	//resore_user_password.html
	tmpl, err = template.ParseFiles("templates/resore_user_password.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		logger.Error(appErr)
		return nil, appErr
	}
	tmplsList[restoreUserPassword] = tmpl

	logger.Info("Templates have been registered")
	tmpls := &WebTemplates{Templates: tmplsList}
	return tmpls, nil
}
