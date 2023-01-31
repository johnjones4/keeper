package api

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type newTokenRequest struct {
	Password string `json:"password"`
}

type newTokenResponse struct {
	Token string `json:"token"`
}

func (a *API) newToken(w http.ResponseWriter, r *http.Request) {
	var in newTokenRequest
	err := readJson(r, &in)
	if err != nil {
		errorResponse(w, http.StatusBadGateway, err)
		return
	}

	err = bcrypt.CompareHashAndPassword(a.runtime.PasswordHash, []byte(in.Password))
	if err != nil {
		errorResponse(w, http.StatusForbidden, err)
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
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(w, http.StatusOK, newTokenResponse{
		Token: ss,
	})
}
