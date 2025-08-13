package main

import (
	"log"
	"net/http"

	"github.com/Sheikh-Fahad-Ahmed/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println(err)
		return
	}

	err = cfg.db.RevokeToken(r.Context(), refreshToken)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "could not revoke refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
