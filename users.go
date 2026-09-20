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
	ID        	uuid.UUID `json:"id"`
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
	Email     	string    `json:"email"`
	IsChirpyRed bool	  `json:"is_chirpy_red"`
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
		log.Printf("Error creating user %v" , err)
		writer.WriteHeader(500)
		return
	}

	userResponse := User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}

	respondWithJSON(writer , http.StatusCreated , userResponse)
}


func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter , r *http.Request){
	type request struct {
		NewPassword  string `json:"password"`
		NewEmail 	 string	`json:"email"`
	}

	var reqStruct request

	err := json.NewDecoder(r.Body).Decode(&reqStruct)
	if err != nil {
		log.Printf("Error Unmarshalling JSON: %v", err)
		w.WriteHeader(500)
		return
	}
	defer r.Body.Close()

	authToken , err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w , 401 , err.Error())
		return
	}
	
	userID , err := auth.ValidateJWT(authToken , cfg.jwtSecret)
	if err != nil {
		respondWithError(w , 401 , "Invalid token")
		return
	}

	newHashedPass , err:= auth.HashPassword(reqStruct.NewPassword)
	if err != nil {
		log.Print("Couldn't hash password")
		w.WriteHeader(500)
		return
	}

	values := database.UpdatePasswordAndEmailParams{
		HashedPassword: newHashedPass,
		Email: reqStruct.NewEmail,
		ID: userID,
	}

	updatedDBUser , err := cfg.dbQueries.UpdatePasswordAndEmail(r.Context() , values)
	if err != nil {
		respondWithError(w , 500 , "Couldn't update user detail")
	}

	userResponse := User{
		ID: updatedDBUser.ID,
		CreatedAt: updatedDBUser.CreatedAt,
		UpdatedAt: updatedDBUser.UpdatedAt,
		Email: updatedDBUser.Email,
		IsChirpyRed: updatedDBUser.IsChirpyRed,
	}

	respondWithJSON(w , 200 , userResponse)
}


