package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type respParams struct {
	Body   string    `json:"body"`
	UserID uuid.UUID 
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) *respParams {

	decoder := json.NewDecoder(r.Body)
	params := respParams{}
	err := decoder.Decode(&params)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return nil
	}

	if len(params.Body) > 140 {
		errorHandler(w, r, http.StatusBadRequest, "Chirp too long", nil)
		return nil
	}

	params.Body = profanityHandler(&params)
	return &params
}
