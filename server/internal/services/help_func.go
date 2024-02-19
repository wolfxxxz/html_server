package services

import (
	"server/internal/apperrors"
	"strconv"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", apperrors.HashPasswordErr.AppendMessage(err)
	}

	return string(bytes), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func claimJWTToken(role string, id string, expiresAt string, sekretKey []byte) (string, error) {
	expiresAtNum, err := strconv.Atoi(expiresAt)
	if err != nil {
		appErr := apperrors.ClaimJWTTokenErr.AppendMessage(err)
		return "", appErr
	}

	t := time.Duration(expiresAtNum) * time.Second
	claims := jwt.MapClaims{
		"role": role,
		"id":   id,
		"exp":  time.Now().Add(t).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(sekretKey)
	if err != nil {
		appErr := apperrors.ClaimJWTTokenErr.AppendMessage(err)
		return "", appErr
	}

	return signedToken, nil
}

func capitalizeFirstRune(line string) string {
	runes := []rune(line)
	for i, r := range runes {
		if i == 0 {
			runes[i] = unicode.ToUpper(r)
		}
	}

	return string(runes)
}

func isCyrillic(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Cyrillic, r) {
			return true
		}
	}

	return false
}
