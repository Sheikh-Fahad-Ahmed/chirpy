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

func (cfg *apiConfig) authenticateUser(w http.ResponseWriter, r *http.Request) {
	type loginReqParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type User struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
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

	expiresIn := time.Hour

	token, err := auth.MakeJWT(DBUser.ID, cfg.secretKey, expiresIn)
	if err != nil {
		log.Fatal("error creating JWT:", err)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Println("error making refresh token")
		return
	}

	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    DBUser.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	}

	refreshTokenData, err := cfg.db.CreateRefreshToken(r.Context(), refreshTokenParams)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "Could not create refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:           DBUser.ID,
			CreatedAt:    DBUser.CreatedAt,
			UpdatedAt:    DBUser.UpdatedAt,
			Email:        DBUser.Email,
			Token:        token,
			RefreshToken: refreshTokenData.Token,
			IsChirpyRed: DBUser.IsChirpyRed,
		},
	})
}
