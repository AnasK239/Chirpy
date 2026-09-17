package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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
		Email string `json:"email"`
	}
	
	var reqStruct request
	
	err := json.NewDecoder(req.Body).Decode(&reqStruct)
	if err != nil {
		log.Printf("Error Unmarshalling JSON: %v", err)
		writer.WriteHeader(500)
		return
	}
	
	dbUser , err := cfg.dbQueries.CreateUser(req.Context() , reqStruct.Email)
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

	jsonData , err := json.Marshal(userResponse)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		writer.WriteHeader(500)
		return
	}

	writer.Header().Set("Content-Type" , "application/json")
	writer.WriteHeader(201)
	writer.Write(jsonData)
}