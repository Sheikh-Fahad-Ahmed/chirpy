package main

import (
	"encoding/json"
	"log"
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

func toChirps(DBChirps []database.Chirp) []Chirp {
	chirps := make([]Chirp, len(DBChirps))
	for i, DBChirp := range DBChirps {
		chirps[i] = Chirp{
			ID: DBChirp.ID,
			CreatedAt: DBChirp.CreatedAt,
			UpdatedAt: DBChirp.UpdatedAt,
			Body: DBChirp.Body,
			UserID: DBChirp.UserID,
		} 
	}
	return chirps
}

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	
	DBChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		log.Println("error getting chips")
		return
	}

	chirps := toChirps(DBChirps)
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(chirps)
	if err != nil {
		log.Println("error json marshal")
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
