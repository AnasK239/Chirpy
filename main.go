package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/AnasK239/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits 	atomic.Int32
	dbQueries 		*database.Queries
	platform 		string
	jwtSecret 		string
	polkaKey		string
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM") 
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")
	
	db , err := sql.Open("postgres" , dbURL)
	if err != nil{
		log.Fatalf("Error opening database conneciton")
	}
	
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
		dbQueries: database.New(db),
		platform: platform,
		jwtSecret: jwtSecret,
		polkaKey: polkaKey,
	}

	serveMux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	FileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	FileHandler = apiCfg.middleWareMetricsIncr(FileHandler)

	serveMux.Handle("/app/", FileHandler)

	serveMux.HandleFunc("POST /api/polka/webhooks" , apiCfg.handleUpgradeEvent)
	
	serveMux.HandleFunc("GET /api/healthz", HealthHandler)
	
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	serveMux.HandleFunc("POST /api/users" , apiCfg.handlerCreateUser)
	serveMux.HandleFunc("PUT /api/users" , apiCfg.handleUpdateUser)

	serveMux.HandleFunc("POST /api/chirps" , apiCfg.handleCreateChirp)
	serveMux.HandleFunc("GET /api/chirps" , apiCfg.handleGetAllChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}" , apiCfg.handleGetChirp)
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handleDeleteChirp)
	
	serveMux.HandleFunc("POST /api/login" , apiCfg.handleLogin)
	serveMux.HandleFunc("POST /api/refresh", apiCfg.HandleRefreshAccessToken)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.HandleRevokeRefreshToken)
	
	
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
