package api

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func (a *API) verifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if len(auth) < 7 {
			errorResponse(w, http.StatusUnauthorized, errors.New("no token"))
			return
		}
		tokenStr := auth[7:]

		token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
			return a.runtime.PrivateKey, nil
		})
		if err != nil {
			w.Header().Set("X-Show-Login", "1")
			errorResponse(w, http.StatusUnauthorized, err)
			return
		}

		if !token.Valid {
			w.Header().Set("X-Show-Login", "1")
			errorResponse(w, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
