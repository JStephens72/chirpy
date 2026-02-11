package main

import (
	"database/sql"
	"errors"
	"net/http"
	"sort"

	"github.com/JStephens72/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	authorIDString := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	ascending := sortOrder != "desc"

	var dbChirps []database.Chirp
	var err error
	var author_id uuid.UUID

	if authorIDString == "" {
		dbChirps, err = cfg.db.GetChirps(r.Context())
	} else {
		author_id, err = uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid author id", err)
			return
		}
		dbChirps, err = cfg.db.GetChirpsByUser(r.Context(), author_id)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "no records", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	sortChirpsByCreatedAt(dbChirps, ascending)

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func sortChirpsByCreatedAt(chirps []database.Chirp, ascending bool) {
	sort.Slice(chirps, func(i, j int) bool {
		if ascending {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		}
		return chirps[j].CreatedAt.Before(chirps[i].CreatedAt)
	})
}
