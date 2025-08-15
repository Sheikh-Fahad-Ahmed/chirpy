package main

import (
	"encoding/json"
	"net/http"

	"github.com/Sheikh-Fahad-Ahmed/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	type UserDataParams struct {
		UserID string `json:"user_id"`
	}

	type webhookReqParams struct {
		Event string         `json:"event"`
		Data  UserDataParams `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := webhookReqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "couldn't decode params", err)
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		errorHandler(w, r, http.StatusUnauthorized, "couldn't get api key", err)
		return
	}

	if apiKey != cfg.polkaKey {
		errorHandler(w, r, http.StatusUnauthorized, "request Unauthorized", nil)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "could not parse uuid", err)
		return
	}

	if err := cfg.db.UpgradeUserToChirpyRed(r.Context(), userID); err != nil {
		errorHandler(w, r, http.StatusNotFound, "Could not find user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
