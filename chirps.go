package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"
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

	userID ,  err :=  cfg.checkCredentials(r.Header)
	if err != nil {
		respondWithError(w , http.StatusUnauthorized , err.Error())
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

	authorIDParam := r.URL.Query().Get("author_id")
	sortParam := r.URL.Query().Get("sort") 


	dbChirps := []database.Chirp{}
	var err error

	if authorIDParam != "" {
		
		id , err := uuid.Parse(authorIDParam)
		if err != nil {
			respondWithError(w , http.StatusBadRequest , "Malformed ID")
			return
		}
		
		dbChirps , err = cfg.dbQueries.GetAuthorChirps(r.Context() , id)
	} else {
		dbChirps , err = cfg.dbQueries.GetAllChirps(r.Context())
	}
	
	if err != nil {
		log.Printf("Error while fetching chirps: %v", err)
		w.WriteHeader(500)
		return
	}	

	if sortParam == "desc" {
		sort.Slice(dbChirps , func(i, j int) bool {
			return dbChirps[j].CreatedAt.Before(dbChirps[i].CreatedAt)  
		})
	}else {
		sort.Slice(dbChirps , func(i, j int) bool {
			return dbChirps[i].CreatedAt.Before(dbChirps[j].CreatedAt)  
		})
	}
	
	responseChirps := make([]Chirp , 0 , len(dbChirps))

	for _ , dbChirp := range dbChirps {
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

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter , r *http.Request) {

	userID ,  err :=  cfg.checkCredentials(r.Header)
	if err != nil {
		respondWithError(w , http.StatusUnauthorized , err.Error())
		return
	}
	
	chirpId , err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w , 500 , "INVALID ID")
		return
	}
	
	chirp , err := cfg.dbQueries.GetChirp(r.Context() , chirpId)
	if err != nil {
		respondWithError(w , http.StatusNotFound , "Chirp doesn't exist")
		return
	}

	if chirp.UserID != userID {
		respondWithError(w , http.StatusForbidden , "You don't own this resource")
		return
	}

	err = cfg.dbQueries.DeleteChirp(r.Context() , chirp.ID)
	if err != nil {
		respondWithError(w , 500 , "Couldn't delete chirp")
		return
	}

	w.WriteHeader(204)
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

func (cfg *apiConfig) checkCredentials(headers http.Header) (uuid.UUID, error) {

	authToken , err := auth.GetBearerToken(headers)
	if err != nil {
		return  uuid.UUID{},  err
	}
	
	userID , err := auth.ValidateJWT(authToken , cfg.jwtSecret)
	if err != nil {
		return  uuid.UUID{}, errors.New("Invalid token")
	}

	return userID,  nil
}
