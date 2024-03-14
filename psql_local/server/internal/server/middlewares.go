package server

import (
	"context"
	"encoding/json"
	"net/http"
	"server/internal/apperrors"
	"time"

	"github.com/golang-jwt/jwt"
)

type contextKey string

const (
	contextKeyRole contextKey = "role"
	contextKeyID   contextKey = "id"
)

func (srv *server) contextExpire(h http.HandlerFunc) http.HandlerFunc {
	srv.logger.Info("contextExpire")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
		defer cancel()

		r = r.WithContext(ctx)
		h(w, r)
	}
}

func (srv *server) jwtAuthentication(h http.HandlerFunc) http.HandlerFunc {
	srv.logger.Info("jwtAuthentication")
	return func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		var tokenGet string
		for _, cookie := range cookies {
			if cookie.Name == "user_token_translator" {
				tokenGet = cookie.Value
			}
		}

		if tokenGet == "" {
			appErr := apperrors.JWTMiddleware.AppendMessage("Vars Authorization")
			srv.logger.Error(appErr)
			srv.respondRegistrateErr(w, appErr)
			return
		}

		if srv.blacklist.IsTokenBlacklisted(tokenGet) {
			appErr := apperrors.JWTMiddleware.AppendMessage("Token is blacklisted")
			srv.logger.Error(appErr)
			srv.respondAuthorizateErr(w, appErr)
			return
		}

		token, err := jwt.Parse(tokenGet, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				appErr := apperrors.JWTMiddleware.AppendMessage("invalid signature method")
				srv.logger.Error(appErr)
				return nil, appErr
			}

			return []byte(srv.config.Server.SecretKey), nil
		})

		if err != nil {
			srv.logger.Error(err)
			appErr := apperrors.JWTMiddleware.AppendMessage("Token is invalid")
			srv.respondAuthorizateErr(w, appErr)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			role, ok := claims["role"].(string)
			if !ok {
				appErr := apperrors.JWTMiddleware.AppendMessage("Role not found in token")
				w.WriteHeader(http.StatusBadRequest)
				err := json.NewEncoder(w).Encode(appErr.Message)
				if err != nil {
					srv.logger.Error(appErr)
				}
			}

			id, ok := claims["id"].(string)
			if !ok {
				appErr := apperrors.JWTMiddleware.AppendMessage("Id not found in token")
				srv.respondErr(w, appErr)
				return
			}

			timeoutDuration, err := time.ParseDuration(srv.config.Server.TimeoutContext + "s")
			if err != nil {
				appErr := apperrors.JWTMiddleware.AppendMessage("Parse duration err").AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respondErr(w, appErr)
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), timeoutDuration)
			defer cancel()
			ctx = context.WithValue(ctx, contextKeyRole, role)
			ctx = context.WithValue(ctx, contextKeyID, id)
			r = r.WithContext(ctx)
			h(w, r)

			return
		}

		appErr := apperrors.JWTMiddleware.AppendMessage("The token has expired or is invalid")
		srv.respondErr(w, appErr)
	}
}

type blacklist struct {
	tokens map[string]bool
}

func newBlacklist() *blacklist {
	return &blacklist{
		tokens: make(map[string]bool),
	}
}

func (b *blacklist) AddToken(token string) {
	b.tokens[token] = true
}

func (b *blacklist) IsTokenBlacklisted(token string) bool {
	return b.tokens[token]
}
