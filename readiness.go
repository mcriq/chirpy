package main

import (
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	bodyStr := "OK"
	bodyBytes := []byte(bodyStr)
	w.Write(bodyBytes)
}