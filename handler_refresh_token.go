package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mcriq/chirpy/internal/auth"
)


func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refrToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	refreshToken, err := cfg.dbQueries.GetRefreshTokenByToken(r.Context(), refrToken)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid or expired refresh token",
			})
			return
		}
		log.Printf("Error retrieving refresh token: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := auth.MakeJWT(refreshToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		log.Printf("Error creating JWT: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = writeJSON(w, int(http.StatusOK), map[string]string{"token": token})
	if err != nil {
		log.Printf("Error encoding response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}