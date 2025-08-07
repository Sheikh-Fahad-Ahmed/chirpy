package main

import (
	"net/http"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) chirpHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Chirp
	}

	params := validateChirpHandler(w, r)
	if params == nil {
		return
	}
	args := database.CreateChirpParams{
		Body:   params.Body,
		UserID: params.UserID,
	}
	chirp, err := cfg.db.CreateChirp(r.Context(), args)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Chirp: Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		},
	})
}
