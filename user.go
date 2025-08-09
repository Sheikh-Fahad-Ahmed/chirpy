package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/chirpy/internal/auth"
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
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserFromDB(dbUser database.CreateUserRow) User {
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

	hashPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Println("error hashing the password:", err)
		return
	}

	args := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashPassword,
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), args)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	user := NewUserFromDB(dbUser)
	respondWithJSON(w, http.StatusCreated, user)

}
