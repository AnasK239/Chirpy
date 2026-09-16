package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

// http.FileServer(http.Dir(".")) is a handler that serves files from the exact URL resource path given
// http.NewServeMux creates an Interception point for requests
// calling .Handle on a serveMux requires giving the request path to be handled
// and the handler that handles it

func main() {

	var apiCfg apiConfig

	serveMux := http.NewServeMux()

	FileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	FileHandler = apiCfg.middleWareMetricsIncr(FileHandler)

	serveMux.Handle("/app/", FileHandler)
	serveMux.HandleFunc("GET /healthz", HealthHandler)
	serveMux.HandleFunc("GET /metrics", apiCfg.metricsHandler)
	serveMux.HandleFunc("POST /reset", apiCfg.resetHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	err := server.ListenAndServe()

	if err != nil {
		fmt.Println(fmt.Errorf("Error Serving requests, %v", err))
	}

}

func HealthHandler(writer http.ResponseWriter, req *http.Request) {

	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("OK"))

}

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middleWareMetricsIncr(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(writer http.ResponseWriter, req *http.Request) {

	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	result := fmt.Sprintf("Hits: %v", cfg.fileServerHits.Load())
	writer.WriteHeader(200)
	writer.Write([]byte(result))
}

func (cfg *apiConfig) resetHandler(writer http.ResponseWriter, req *http.Request) {
	cfg.fileServerHits.Swap(0)
	writer.WriteHeader(200)
}
