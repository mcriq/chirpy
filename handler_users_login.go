package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/mcriq/chirpy/internal/auth"
)


func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type userLoginParams struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	
	decoder := json.NewDecoder(r.Body)
	params := userLoginParams{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding login parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user with email: %s found: %v", params.Email, err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "incorrect email or password",
			})
			return
		}
		log.Printf("Unable to get user: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Invalid password for user %s: %s", user.Email, err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "incorrect email or password",
		})
		return
	}

	userResp := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	err = writeJSON(w, http.StatusOK, userResp)
	if err != nil {
		log.Printf("Error encoding response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}