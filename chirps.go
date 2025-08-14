package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/chirpy/internal/auth"
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

func toChirp(DBChirp database.Chirp) Chirp {
	return Chirp{
		ID:        DBChirp.ID,
		CreatedAt: DBChirp.CreatedAt,
		UpdatedAt: DBChirp.UpdatedAt,
		Body:      DBChirp.Body,
		UserID:    DBChirp.UserID,
	}
}

func toChirps(DBChirps []database.Chirp) []Chirp {
	chirps := make([]Chirp, len(DBChirps))
	for i, DBChirp := range DBChirps {
		chirps[i] = toChirp(DBChirp)
	}
	return chirps
}

func (cfg *apiConfig) chirpHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Chirp
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println(err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		errorHandler(w, r, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	params := validateChirpHandler(w, r)
	if params == nil {
		return
	}
	args := database.CreateChirpParams{
		Body:   params.Body,
		UserID: userId,
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

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Chirp []Chirp
	}
	DBChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		log.Println("error getting chips")
		return
	}

	chirps := toChirps(DBChirps)
	respondWithJSON(w, http.StatusOK, response{
		Chirp: chirps,
	})
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Chirp
	}
	IDStr := r.PathValue("chirpID")
	id, err := uuid.Parse(IDStr)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "Couldn't parse UUID", err)
		return
	}

	DBChirp, err := cfg.db.GetChirpByID(r.Context(), id)
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "Not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Chirp: toChirp(DBChirp),
	})
}

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorHandler(w, r, http.StatusUnauthorized, "token malformed or missing", err)
	}

	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		errorHandler(w, r, http.StatusForbidden, "invalid token", err)
		return
	}

	IDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(IDStr)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "Couldn't parse uuid", err)
		return
	}
	chirpInfo, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "Chirp Not Found", err)
		return
	}
	if userID != chirpInfo.UserID {
		errorHandler(w, r, http.StatusForbidden, "Forbidden", err)
		return
	}

	if err := cfg.db.DeleteChirp(r.Context(), chirpID); err != nil {
		errorHandler(w, r, http.StatusNotFound, "chirp could not be deleted", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
