package main

import (
	"fmt"
	"net/http"

	"github.com/FilipStam97/rssator/internal/auth"
	"github.com/FilipStam97/rssator/internal/database"
)

// problem func signature doesent match the handler fun signatrue that we have used, solution is the func down
type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {

	//closure function has accsess to everything in apiConfig
	return func(w http.ResponseWriter, r *http.Request) {
		//this is an authenticated endpoint
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			//403 permision error
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		//r.Context(), context package in std library, gives you a way to track something thatt is happening across multiple go routines(when we are using them)
		//most important thing is canceling the context wwhich will kill the http request, so it is important to use the current context
		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get a user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
