package main

import(
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middleWareMetricsIncr(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(writer http.ResponseWriter, req *http.Request) {

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	result := fmt.Sprintf(
		`
	<html>
  		<body>
    	<h1>Welcome, Chirpy Admin</h1>
     	<p>Chirpy has been visited %d times!</p>
    	</body>
	</html>
	`, cfg.fileServerHits.Load())

	writer.WriteHeader(200)
	writer.Write([]byte(result))
}
