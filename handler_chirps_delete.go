package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/mcriq/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		log.Printf("invalid chirp id: %s", err)
        w.WriteHeader(http.StatusBadRequest)
        return
	}

	if chirpID == uuid.Nil {
		log.Printf("chirp id must not be blank")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("chirp does not exist: %s", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Printf("error getting chirp", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if chirp.UserID != userID {
		log.Printf("not authorized to delete chirp")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = cfg.dbQueries.DeleteChirp(r.Context(), chirp.ID)
	if err != nil {
		log.Printf("unable to delete chirp: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("Successfully deleted chirp")
}