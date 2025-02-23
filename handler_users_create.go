package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type email struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	emailParams := email{}
	err := decoder.Decode(&emailParams)
	if err != nil {
		log.Printf("Error decoding user parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), emailParams.Email)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			log.Printf("unable to create user: %v", err)
			w.WriteHeader(http.StatusConflict)
			return
		}
		log.Printf("unable to create user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userParams := User{
		ID: user.ID, 
		CreatedAt: user.CreatedAt, 
		UpdatedAt: user.UpdatedAt, 
		Email: user.Email,
	}
	
	err = writeJSON(w, http.StatusCreated, userParams)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := cfg.dbQueries.DeleteUsers(r.Context())
	if err != nil {
		fmt.Printf("unable to delete users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}