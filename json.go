package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		//errors in the 400 range are client errs and we dont really care about them but we do care about 500 rande since tose are bugs on our end
		log.Println("Responding with 5XX error: ", msg)
	}

	type errResponse struct {
		Error string `json:"error"` //in go we take a struct add these json tags to itt to specify how we want json marshal func to convert struct to json obj
		//here we are saying i have an Error field and i want the key for this field to be "error"
		//the final result will be a json like { "error": "something went wrong"}
		//with this we will have a consistent error format for our api that we can put into docs for users
	}

	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}
