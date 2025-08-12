package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type respToken struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println(err)
		return
	}

	tokenInfo, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		errorHandler(w, r, http.StatusUnauthorized, "Could not find refresh token", err)
		return
	}

	if time.Now().After(tokenInfo.ExpiresAt) || tokenInfo.RevokedAt.Valid {
		errorHandler(w, r, http.StatusUnauthorized, "Token has been expired or revoked", nil)
		return
	}

	newToken, err := auth.MakeJWT(tokenInfo.UserID, cfg.secretKey, time.Hour)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "Could not create JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, respToken{
		Token: newToken,
	})
}
