package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/JStephens72/chirpy/internal/auth"
	"github.com/google/uuid"
)

const EventUserStatusUpgrade string = "user.upgraded"

func (cfg *apiConfig) handlerPolkaWH(w http.ResponseWriter, r *http.Request) {
	//authenticate the event by validating the ApiKey
	token, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find ApiKey", err)
		return
	}

	if token != cfg.apiKey {
		respondWithError(w, http.StatusUnauthorized, "invalid APIKey", err)
		return
	}

	type UpgradeData struct {
		UserID uuid.UUID `json:"user_id"`
	}

	type UpgradeEvent struct {
		Event string      `json:"event"`
		Data  UpgradeData `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	event := UpgradeEvent{}
	err = decoder.Decode(&event)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if event.Event != EventUserStatusUpgrade {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	_, err = cfg.db.UpgradeUser(r.Context(), event.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "unknown user", err)
			return
		}
		respondWithError(w, http.StatusBadRequest, "error: ", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
