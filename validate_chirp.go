package main

import (
	"encoding/json"
	"net/http"
)

type respParams struct {
	Body string `json:"body"`
}

type errorRespParams struct {
	Error string `json:"error"`
}

type cleanRespParams struct {
	CleanedBody string `json:"cleaned_body"`
}

func validateHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := respParams{}
	err := decoder.Decode(&params)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(params.Body) > 140 {
		errorHandler(w, r, http.StatusBadRequest, "Chirp too long", nil)
		return
	}

	text := profanityHandler(&params)
	respondWithJSON(w, 200, cleanRespParams{
		CleanedBody: text,
	})
}
