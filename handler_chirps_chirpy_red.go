package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/mcriq/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpyRed(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserID uuid.UUID `json:"user_id"`
	}
	type requestBody struct {
		Event string `json:"event"`
		Data data `json:"data"`
	}

	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Unable to retreive api key: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if key != cfg.polkaKey {
		log.Printf("Unable to verify apikey")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	requestParams := requestBody{}
	err = decoder.Decode(&requestParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if requestParams.Event != "user.upgraded" {
		log.Printf("unsupported event: %s", requestParams.Event)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userExists, err := cfg.dbQueries.UserExists(r.Context(), requestParams.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !userExists {
		log.Printf("user not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = cfg.dbQueries.UpgradeToChirpyRedById(r.Context(), requestParams.Data.UserID)
	if err != nil {
		log.Printf("error upgrading user: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return

}