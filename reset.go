package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) resetHandler(writer http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		writer.WriteHeader(403)
		writer.Write([]byte("Forbidden"))
		return
	}

	err := cfg.dbQueries.DeleteAllUsers(req.Context())
	if err != nil {
		log.Printf("Error deleting all users")
		writer.WriteHeader(500)
		return
	}
	
	cfg.fileServerHits.Store(0)
	writer.WriteHeader(200)
	writer.Write([]byte("Hits reset to 0 and all users deleted"))
}
