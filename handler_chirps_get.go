package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)


func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	type ChirpResponse struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    string    `json:"user_id"`
	}

	chirpID := r.PathValue("chirpID")

	if len(chirpID) != 0 {
		parsedChirpID, err := uuid.Parse(chirpID)
		if err != nil {
			log.Printf("unable to parse chirp id: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		chirp, err := cfg.dbQueries.GetChirp(r.Context(), parsedChirpID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("no chirp found for id: %s: %v", chirpID, err)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			log.Printf("error fetching chirp with id: %s: %v", chirpID, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = writeJSON(w, http.StatusOK, ChirpResponse{ID: chirp.ID.String(), CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID.String()})
		if err != nil {
			log.Printf("unable to parse json: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		return
	}

	chirps, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		log.Printf("unable to fetch chirps")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	chirpResponses := []ChirpResponse{}
	for _, chirp := range chirps {
		chirpResponse := ChirpResponse{
			ID: chirp.ID.String(),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID.String(),
		}
		chirpResponses = append(chirpResponses, chirpResponse)
	}

	err = writeJSON(w, http.StatusOK, chirpResponses)
	if err != nil {
		log.Printf("unable to parse json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}