package main

import (
	"encoding/json"
	"net/http"

	"github.com/JStephens72/chirpy/internal/auth"
	"github.com/JStephens72/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type UserUpdateParams struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't extract token", err)
	}

	userID, err := auth.ValidateJWT(token, cfg.serverSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := UserUpdateParams{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't decode parameters", err)
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash new password", err)
		return
	}

	user, err := cfg.db.GetUserByUserID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find user", err)
		return
	}

	userUpdateParams := database.UpdateUserParams{
		ID:             user.ID,
		Email:          params.Email,
		HashedPassword: hashed_password,
	}

	user, err = cfg.db.UpdateUser(r.Context(), userUpdateParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	respBody := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusOK, respBody)
}
