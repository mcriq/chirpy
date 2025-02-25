package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mcriq/chirpy/internal/auth"
	"github.com/mcriq/chirpy/internal/database"
)

var profaneWords = []string{"kerfuffle", "sharbert", "fornax"}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body string `json:"body"`
	}

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

	decoder := json.NewDecoder(r.Body)
	requestParams := requestBody{}
	err = decoder.Decode(&requestParams)
	if err != nil {
		log.Printf("Error decoding chirp parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(requestParams.Body) == 0 {
		log.Printf("Cannot add empty chirp")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(requestParams.Body) > 140 {
		log.Printf("Chirp is too long")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	newBodyString := replaceProfanity(requestParams.Body)
	requestParams.Body = newBodyString

	
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: requestParams.Body,
		UserID: userID,
	})
	if err != nil {
		log.Printf("unable to create chirp: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	responseChirp := Chirp{
		ID:        chirp.ID,         
		CreatedAt: chirp.CreatedAt,  
		UpdatedAt: chirp.UpdatedAt,  
		Body:      chirp.Body,       
		UserID:    chirp.UserID,     
	}

	err = writeJSON(w, http.StatusCreated, responseChirp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func writeJSON[T any](w http.ResponseWriter, status int, v T) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    return json.NewEncoder(w).Encode(v)
}

func replaceProfanity(text string) string {
	bodySlice := strings.Split(text, " ")
	for i, word := range bodySlice {
		if slices.Contains(profaneWords, strings.ToLower(word)) {
			bodySlice[i] = "****"
		}
	}
	return strings.Join(bodySlice, " ")
}

func safeWriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
    err := writeJSON(w, statusCode, data)
    if err != nil {
        log.Printf("unable to parse json: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
    }
    return err
}