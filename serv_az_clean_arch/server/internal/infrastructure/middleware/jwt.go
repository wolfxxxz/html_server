package middleware

import (
	"server/internal/apperrors"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func JWTAuthentication(jc *JWTMiddlewareConfig, blacklist *Blacklist) echo.MiddlewareFunc {
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
				//srv.respondRegistrateErr(w, appErr)
				return c.JSON(appErr.HTTPCode, appErr.Message)
			}

			log.Info(tokenGet)

			if blacklist.IsTokenBlacklisted(tokenGet) {
				appErr := apperrors.JWTMiddleware.AppendMessage("Token is blacklisted")
				log.Error(appErr)
				return c.JSON(appErr.HTTPCode, appErr.Message)
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
				log.Error(err)
				appErr := apperrors.JWTMiddleware.AppendMessage("Token is invalidd")
				return c.JSON(appErr.HTTPCode, appErr.Message)
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				role, ok := claims["role"].(string)
				if !ok {
					appErr := apperrors.JWTMiddleware.AppendMessage("Role not found in token")
					return appErr
				}

				id, ok := claims["id"].(string)
				if !ok {
					appErr := apperrors.JWTMiddleware.AppendMessage("Email not found in token")
					return appErr
				}

				c.Set("role", role)
				c.Set("id", id)
				return next(c)
			}

			appErr := apperrors.JWTMiddleware.AppendMessage("The token has expired or is invalid")
			return c.JSON(appErr.HTTPCode, appErr.Message)
		}
	}
}
