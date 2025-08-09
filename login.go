package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) authenticateUser(w http.ResponseWriter, r *http.Request) {
	type loginReqParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
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

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        DBUser.ID,
			CreatedAt: DBUser.CreatedAt,
			UpdatedAt: DBUser.UpdatedAt,
			Email:     DBUser.Email,
		},
	})
}
