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
		tokenGet := r.Header.Get("Authorization")
		if tokenGet == "" {
			appErr := apperrors.JWTMiddleware.AppendMessage("Vars Authorization")
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusUnauthorized)
			return
		}

		if srv.blacklist.IsTokenBlacklisted(tokenGet) {
			appErr := apperrors.JWTMiddleware.AppendMessage("Token is blacklisted")
			srv.logger.Error(appErr)
			srv.respond(w, appErr.Message, http.StatusUnauthorized)
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
			srv.respond(w, appErr.Message, http.StatusUnauthorized)
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
				srv.respond(w, appErr.Message, http.StatusUnauthorized)
				return
			}

			srv.logger.Infof("TimeOUT CONFIG %v", srv.config.Server.TimeoutContext)
			timeoutDuration, err := time.ParseDuration(srv.config.Server.TimeoutContext + "s")
			if err != nil {
				appErr := apperrors.JWTMiddleware.AppendMessage("Parse duration err").AppendMessage(err)
				srv.logger.Error(appErr)
				srv.respond(w, appErr.Message, http.StatusUnauthorized)
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), timeoutDuration)
			defer cancel()
			ctx = context.WithValue(ctx, contextKeyRole, role)
			ctx = context.WithValue(ctx, contextKeyID, id)
			r = r.WithContext(ctx)
			h(w, r)

			srv.logger.Info("jwtAuthentication success")
			return
		}

		appErr := apperrors.JWTMiddleware.AppendMessage("The token has expired or is invalid")
		srv.respond(w, appErr.Message, http.StatusUnauthorized)
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
