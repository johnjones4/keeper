package api

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func (a *API) PostToken(w http.ResponseWriter, r *http.Request) {
	var in TokenRequest
	err := readJson(r, &in)
	if err != nil {
		a.errorResponse(w, http.StatusBadGateway, err)
		return
	}

	err = bcrypt.CompareHashAndPassword(a.runtime.PasswordHash, []byte(in.Password))
	if err != nil {
		a.errorResponse(w, http.StatusForbidden, err)
		return
	}

	now := time.Now().UTC()

	claims := &jwt.StandardClaims{
		ExpiresAt: now.Add(expiration).Unix(),
		IssuedAt:  now.Unix(),
		Issuer:    issuer,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(a.runtime.PrivateKey)
	if err != nil {
		a.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	a.jsonResponse(w, http.StatusOK, TokenResponse{
		Token: ss,
	})
}
