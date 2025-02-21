package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

var profaneWords = []string{"kerfuffle", "sharbert", "fornax"}


func handlerValidateChirp(w http.ResponseWriter, r *http.Request){
    type parameters struct {
        Body string `json:"body"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }
    
	type invalid struct {
        Error string `json:"error"`
    }

	type valid struct {
		CleanedBody string `json:"cleaned_body"`
	}

	if len(params.Body) > 140 {
		err := writeJSON(w, http.StatusBadRequest, invalid{Error: "Chirp is too long"})
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
			}
		return
	}
	
	newBodyString := replaceProfanity(params.Body)

	err = writeJSON(w, http.StatusOK, valid{CleanedBody: newBodyString})
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
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