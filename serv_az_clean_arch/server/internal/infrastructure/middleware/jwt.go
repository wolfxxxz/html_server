package middleware

import (
	"encoding/json"
	"net/http"
	"server/internal/apperrors"
	"server/internal/infrastructure/webtemplate.go"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

const (
	registrate = "registrate"
)

func JWTAuthentication(jc *JWTMiddlewareConfig, blacklist *Blacklist, tmpls *webtemplate.WebTemplates) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookies := c.Request().Cookies()
			var tokenGet string
			for _, cookie := range cookies {
				if cookie.Name == "user_token_translator" {
					tokenGet = cookie.Value
				}
			}

			if tokenGet == "" {
				appErr := apperrors.JWTMiddleware.AppendMessage("Vars Authorization")
				log.Error(appErr)
				return tmpls.Templates[registrate].ExecuteTemplate(c.Response().Writer, registrate, appErr.Message)
			}

			log.Info(tokenGet)

			if blacklist.IsTokenBlacklisted(tokenGet) {
				appErr := apperrors.JWTMiddleware.AppendMessage("Token is blacklisted")
				log.Error(appErr)
				return tmpls.Templates[registrate].ExecuteTemplate(c.Response().Writer, registrate, appErr.Message)
			}

			token, err := jwt.Parse(tokenGet, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					appErr := apperrors.JWTMiddleware.AppendMessage("invalid signature method")
					log.Error(appErr)
					return nil, appErr
				}

				return []byte(jc.SecretKey), nil
			})

			if err != nil {
				return tmpls.Templates[registrate].ExecuteTemplate(c.Response().Writer, registrate, err)
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				role, ok := claims["role"].(string)
				if !ok {
					appErr := apperrors.JWTMiddleware.AppendMessage("Role not found in token")
					c.Response().Writer.WriteHeader(http.StatusBadRequest)
					err := json.NewEncoder(c.Response().Writer).Encode(appErr.Message)
					if err != nil {
						log.Error(appErr)
					}
					//appErr := apperrors.JWTMiddleware.AppendMessage("Role not found in token")
					//return appErr
				}

				id, ok := claims["id"].(string)
				if !ok {
					appErr := apperrors.JWTMiddleware.AppendMessage("Email not found in token")
					return tmpls.Templates[registrate].ExecuteTemplate(c.Response().Writer, registrate, appErr.Message)
				}

				c.Set("role", role)
				c.Set("id", id)
				return next(c)
			}

			appErr := apperrors.JWTMiddleware.AppendMessage("The token has expired or is invalid")
			return tmpls.Templates[registrate].ExecuteTemplate(c.Response().Writer, registrate, appErr.Message)
		}
	}
}
