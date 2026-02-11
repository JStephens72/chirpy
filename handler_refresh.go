package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JStephens72/chirpy/internal/auth"
	"github.com/JStephens72/chirpy/internal/database"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type returnVal struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't extract refresh token", err)
		return
	}

	storedRefreshToken, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token", err)
		return
	}

	if !IsValid(storedRefreshToken) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired or revoked", fmt.Errorf("Refresh token expired or revoked"))
		return
	}

	authToken, err := auth.MakeJWT(storedRefreshToken.UserID, cfg.serverSecret, authTokenLifetime)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't generate auth token", err)
		return
	}

	respBody := returnVal{
		Token: authToken,
	}

	respondWithJSON(w, http.StatusOK, respBody)
}

func IsValid(token database.RefreshToken) bool {
	expired := token.ExpiresAt.Before(time.Now())
	revoked := token.RevokedAt.Valid == true && token.RevokedAt.Time.Before(time.Now())
	if expired || revoked {
		return false
	}
	return true
}
