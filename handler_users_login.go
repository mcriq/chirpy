package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mcriq/chirpy/internal/auth"
	"github.com/mcriq/chirpy/internal/database"
)


func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type userLoginParams struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
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

	expiresIn := time.Hour
	
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Invalid password for user %s: %s", user.Email, err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "incorrect email or password",
		})
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, expiresIn)
	if err != nil {
		log.Printf("Error creating JWT: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refrToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error creating refresh token: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{Token: refrToken, UserID: user.ID})
	if err != nil {
		log.Printf("Error creating refresh token record: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userResp := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
		RefreshToken: refrToken,
		IsChirpyRed: user.IsChirpyRed,
	}

	err = writeJSON(w, http.StatusOK, userResp)
	if err != nil {
		log.Printf("Error encoding response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}