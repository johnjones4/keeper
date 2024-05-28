package api

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func (a *API) verifyToken(w http.ResponseWriter, r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if len(auth) < 7 {
		a.errorResponse(w, http.StatusUnauthorized, errors.New("no token"))
		return false
	}
	tokenStr := auth[7:]

	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return a.runtime.PrivateKey, nil
	})
	if err != nil {
		w.Header().Set("X-Show-Login", "1")
		a.errorResponse(w, http.StatusUnauthorized, err)
		return false
	}

	if !token.Valid {
		w.Header().Set("X-Show-Login", "1")
		a.errorResponse(w, http.StatusUnauthorized, errors.New("invalid token"))
		return false
	}

	return true
}
