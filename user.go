package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type userReqParams struct {
	Email string `json:"email"`
}

func NewUserFromDB(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
}

func (cfg *apiConfig) userHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := userReqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "Couldn't Decode parameters", err)
		return
	}
 
	dbUser, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	user := NewUserFromDB(dbUser)
	respondWithJSON(w, http.StatusCreated, user)

}
