package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AnasK239/Chirpy/internal/auth"
)


func (cfg *apiConfig) handleLogin(w http.ResponseWriter , r *http.Request){
	type request struct{
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	var reqData request

	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		log.Printf("Error Unmarshalling JSON: %v", err)
		w.WriteHeader(500)
		return
	}
	defer r.Body.Close()
	
	dbUser , err := cfg.dbQueries.FindUserByEmail(r.Context() , reqData.Email)
	if err != nil{
		respondWithError(w , http.StatusUnauthorized , "Incorrect email or password")
		return
	}

	
	ok , err := auth.CheckPassword(reqData.Password, dbUser.HashedPassword)
	if !ok || err != nil {
		respondWithError(w , http.StatusUnauthorized , "Incorrect email or password")
		return
	}

	responseUser := User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
	}
	
	respondWithJSON(w , 200 , responseUser)
}