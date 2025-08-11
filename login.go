package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) authenticateUser(w http.ResponseWriter, r *http.Request) {
	type loginReqParams struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}

	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := loginReqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "Couldn't Decode parameters", err)
		return
	}

	DBUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "User not found", err)
		return
	}

	if err := auth.CheckPasswordHash(params.Password, DBUser.HashedPassword); err != nil {
		errorHandler(w, r, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	var expiresIn time.Duration
	if params.ExpiresInSeconds == nil {
		expiresIn = time.Hour
	} else {
		expiresIn = time.Duration(*params.ExpiresInSeconds) * time.Second
		if expiresIn > time.Hour {
			expiresIn = time.Hour
		}
	}

	token, err := auth.MakeJWT(DBUser.ID, cfg.secretKey, expiresIn)
	if err != nil {
		log.Fatal("error creating JWT:", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        DBUser.ID,
			CreatedAt: DBUser.CreatedAt,
			UpdatedAt: DBUser.UpdatedAt,
			Email:     DBUser.Email,
			Token:     token,
		},
	})
}
