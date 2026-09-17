package main

import(
	"net/http"
	"strings"
	"encoding/json"
	"log"
)

var profaneWords = map[string]struct{} {
	"kerfuffle" : {},
	"sharbert" : {},
	"fornax" : {},
}

func validationHandler(w http.ResponseWriter , r *http.Request){
	type data struct {
		Body string `json:"body"`
	}

	var requestData data

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		log.Printf("Error Decoding JSON: %v", err)
		w.WriteHeader(500)
		return
	}

	if len(requestData.Body) > 140 {
		respondWithError(w , 400 , "Chirp is too long")
		return
	}

	cleanedBody := cleanProfanity(requestData.Body , profaneWords)

	type result struct {
		CleanedBody string `json:"cleaned_body"`
	}

	respondWithJSON(w , 200 , result{CleanedBody: cleanedBody,})

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
