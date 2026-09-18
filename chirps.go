package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/AnasK239/Chirpy/internal/auth"
	"github.com/AnasK239/Chirpy/internal/database"
	"github.com/google/uuid"
)

var profaneWords = map[string]struct{} {
	"kerfuffle" : {},
	"sharbert" : {},
	"fornax" : {},
}

type Chirp struct{
	ID 		  uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body 	  string	`json:"body"`
	UserId	  uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter , r *http.Request){
	type request struct{
		Body   string `json:"body"`
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

	
	body , ok := validateAndCleanChirp(reqStruct.Body)
	if !ok {
		respondWithError(w , 400 , "Chirp is too long")
		return
	}

	values := database.CreateChirpParams{
		Body: body,
		UserID: userID,
	}
	
	dbChirp , err := cfg.dbQueries.CreateChirp(r.Context() , values)
	if err != nil {
		log.Printf("Error creating user")
		w.WriteHeader(500)
		return
	}
	
	chirp := Chirp{
		ID: dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body: dbChirp.Body,
		UserId: dbChirp.UserID,
	}

	respondWithJSON(w , http.StatusCreated , chirp)
}



func (cfg *apiConfig) handleGetAllChirps(w http.ResponseWriter , r *http.Request){

	dbChrips , err := cfg.dbQueries.GetChirpsOrderByCreatedAtAsc(r.Context())
	if err != nil {
		log.Printf("Error Unmarshalling JSON: %v", err)
		w.WriteHeader(500)
		return
	}

	responseChirps := make([]Chirp , 0 , len(dbChrips))

	for _ , dbChirp := range dbChrips {
		responseChirps = append(responseChirps , Chirp{
			ID: dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body: dbChirp.Body,
			UserId: dbChirp.UserID,
		})
	}
	
	respondWithJSON(w , http.StatusOK , responseChirps)
}


func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter , r *http.Request){


	id , err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w , 500 , "INVALID ID")
		return
	}
	
	dbChirp , err := cfg.dbQueries.GetChirp(r.Context() , id)
	if err != nil {
		respondWithError(w , http.StatusNotFound , "NOT FOUND")
		return
	}
	
	chirp := Chirp{
		ID: dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body: dbChirp.Body,
		UserId: dbChirp.UserID,
	}

	respondWithJSON(w , http.StatusOK , chirp)
}

func validateAndCleanChirp(requestData string) (string , bool) {
	
	if len(requestData) > 140 {
		return "" , false;
	}

	cleanedBody := cleanProfanity(requestData , profaneWords)

	return cleanedBody , true

}


func cleanProfanity(message string , words map[string]struct{}) string {

	slice := strings.Split(message, " ")

	for i := 0 ; i < len(slice) ; i++ {
		word := strings.ToLower(slice[i])
		if _ , ok := words[word] ; ok {
			slice[i] = "****"
		}
	}

	return strings.Join(slice , " ")
}
