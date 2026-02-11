package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering handlerReset...")
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))

	// delete users from database
	err := cfg.db.Reset(r.Context())
	if err != nil {
		log.Printf("error removing users from database: %v", err)
	}
}
