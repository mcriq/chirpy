package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/mcriq/chirpy/internal/auth"
	"github.com/mcriq/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type userResp struct {
		Email     string `json:"email"`
		ID        uuid.UUID `json:"id"`
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

	type requestBody struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	requestParams := requestBody{}
	err = decoder.Decode(&requestParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if requestParams.Email == "" || requestParams.Password == "" {
		log.Printf("Email and password must not be blank")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// hash the password parameter
	hashedPW, err := auth.HashPassword(requestParams.Password)
	if err != nil {
		log.Printf("Unable to hash password: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// write the updated info to the database by userid
	err = cfg.dbQueries.UpdateUserByID(r.Context(), database.UpdateUserByIDParams{Email: requestParams.Email, HashedPassword: hashedPW, ID: userID})
	if err != nil {
		log.Printf("Unable to update user information: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	updateResp := userResp{
		Email: requestParams.Email,
		ID: userID,
	}

	writeJSON(w, http.StatusOK, updateResp)
}