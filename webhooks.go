package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AnasK239/Chirpy/internal/auth"
	"github.com/google/uuid"
)


type UpgradeRequest struct{
	Event 	string 				`json:"event"`
	Data 	map[string]string	`json:"data"`
}
const upgradeEvent = "user.upgraded"

func (cfg *apiConfig) handleUpgradeEvent(w http.ResponseWriter , r *http.Request){

	var request UpgradeRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Error unmarshalling json")
		w.WriteHeader(500)
		return
	}
	defer r.Body.Close()

	APIKey , err := auth.GetAPIKey(r.Header)
	if err != nil || APIKey != cfg.polkaKey {
		respondWithError(w , http.StatusUnauthorized , "Invalid api key " + err.Error())
		return
	}

	if request.Event != upgradeEvent {
		w.WriteHeader(204)
		return
	}

	userId , ok := request.Data["user_id"]
	if !ok {
		respondWithError(w , http.StatusBadRequest , "UserId not present in data")
		return
	}

	id , err := uuid.Parse(userId)
	if err != nil {
		respondWithError(w , http.StatusBadRequest , "Malformed userID")
		return
	}
	
	err = cfg.dbQueries.UpgradeToChirpyRed(r.Context() , id)
	if err != nil {
		respondWithError(w , http.StatusNotFound , "User doesn't exist")
		return
	}

	w.WriteHeader(204)
}