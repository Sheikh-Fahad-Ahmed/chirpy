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
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type userReqParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserFromDB(dbUser database.CreateUserRow) User {
	return User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
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

func (cfg *apiConfig) userPUTHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorHandler(w, r, http.StatusUnauthorized, "token malformed or missing", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		errorHandler(w, r, http.StatusUnauthorized, "invalid token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := userReqParams{}
	err = decoder.Decode(&params)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "couldn't Decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "error hashing password", err)
		return
	}

	args := database.UpdateUserEmailAndPassParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	}

	DBUser, err := cfg.db.UpdateUserEmailAndPass(r.Context(), args)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "Couldn't update email and password", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:          DBUser.ID,
		CreatedAt:   DBUser.CreatedAt,
		UpdatedAt:   DBUser.UpdatedAt,
		Email:       DBUser.Email,
		IsChirpyRed: DBUser.IsChirpyRed,
	})
}
