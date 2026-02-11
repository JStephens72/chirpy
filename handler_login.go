package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/JStephens72/chirpy/internal/auth"
	"github.com/JStephens72/chirpy/internal/database"
	"github.com/google/uuid"
)

const authTokenLifetime time.Duration = 3600 * time.Second

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type returnVals struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("err retrieving user by email: %v", err)
		respondWithError(w, http.StatusBadRequest, "unknown user", err)
		return
	}

	ok, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "authentication failed", err)
		return
	}

	userJWT, err := auth.MakeJWT(user.ID, cfg.serverSecret, authTokenLifetime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error generating JWT", err)
		return
	}

	rt, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "err generating refresh token", err)
		return
	}

	createRefreshTokenParams := database.CreateRefreshTokenParams{
		Token:  rt,
		UserID: user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	}

	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), createRefreshTokenParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't insert refresh token", err)
		return
	}

	respBody := returnVals{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        userJWT,
		RefreshToken: refreshToken.Token,
	}

	respondWithJSON(w, http.StatusOK, respBody)
}
