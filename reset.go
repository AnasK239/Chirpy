package main

import(
	"net/http"
)

func (cfg *apiConfig) resetHandler(writer http.ResponseWriter, req *http.Request) {
	cfg.fileServerHits.Swap(0)
	writer.WriteHeader(200)
	writer.Write([]byte("Hits reset to 0"))
}
