package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/AnasK239/Chirpy/internal/auth"
	"github.com/AnasK239/Chirpy/internal/database"
)

const (
	JwtAccessTokenExpiration = time.Hour
	RefreshTokenExpiration = time.Hour * 24 * 60
)

type LoginResponse struct {
	User
	Token 		 string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type request struct{
	Password 	string `json:"password"`
	Email    	string `json:"email"`
}
	
func (cfg *apiConfig) handleLogin(w http.ResponseWriter , r *http.Request){

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

	
	accessToken , err := auth.MakeJWT(dbUser.ID , cfg.jwtSecret , JwtAccessTokenExpiration)
	if err != nil {
		log.Printf("Error creating user jwt")
		w.WriteHeader(500)
		return
	}

	
	
	values := database.CreateRefreshTokenParams{
		Token: auth.MakeRefreshToken(),
		ExpiresAt: time.Now().UTC().Add(RefreshTokenExpiration),
		UserID: dbUser.ID,
	}

	dbRefreshToken , err := cfg.dbQueries.CreateRefreshToken(r.Context() , values)
	if err != nil {
		log.Printf("Error creating user refreshToken")
		w.WriteHeader(500)
		return
	}



	response := LoginResponse {
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
		Token: accessToken,
		RefreshToken: dbRefreshToken.Token,
	}
	
	respondWithJSON(w , 200 , response)
}