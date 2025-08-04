package main

import (
	"encoding/json"
	"log"
	"net/http"
)


func errorHandler(w http.ResponseWriter, r *http.Request, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error : %s", err)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Typo", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error Marshaling json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}