package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/AnasK239/Chirpy/internal/auth"
	"github.com/AnasK239/Chirpy/internal/database"
	"github.com/google/uuid"
)



type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}


func (cfg *apiConfig) handlerCreateUser(writer http.ResponseWriter , req *http.Request){
	type request struct{
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	
	var reqStruct request
	
	err := json.NewDecoder(req.Body).Decode(&reqStruct)
	if err != nil {
		log.Printf("Error Unmarshalling JSON: %v", err)
		writer.WriteHeader(500)
		return
	}

	hashedPassword , err := auth.HashPassword(reqStruct.Password)
	if err != nil {
		log.Printf("Error hashing password: %v" , err)
		writer.WriteHeader(500)
		return 
	}
	
	values := database.CreateUserParams{
		Email: reqStruct.Email,
		HashedPassword: hashedPassword,
	}
	
	dbUser , err := cfg.dbQueries.CreateUser(req.Context() , values)
	if err != nil {
		log.Printf("Error creating user")
		writer.WriteHeader(500)
		return
	}

	userResponse := User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
	}

	respondWithJSON(writer , http.StatusCreated , userResponse)
}