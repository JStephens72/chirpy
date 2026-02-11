package main

import (
	"net/http"

	"github.com/JStephens72/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't extract refresh token", err)
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token in database", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
