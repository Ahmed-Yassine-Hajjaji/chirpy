package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}
func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Swap(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}



func (cfg *apiConfig) middlewareMetricsInc(next http.Handler)	http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cfg.fileserverHits.Add(1)
			next.ServeHTTP(w, r)
		})
	} 


	
func main() {
	apiConfiguration := apiConfig{}
	// Creates a new server Mux(the router: looks for the request's
	// incoming path and decides which handler gets to deal with it)
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	})
	
	mux.Handle("/app/", apiConfiguration.middlewareMetricsInc(http.StripPrefix("/app",http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/metrics", apiConfiguration.metricsHandler)
	mux.HandleFunc("/reset", apiConfiguration.resetHandler)

	// mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	// mux.Handle("/assets/", http.FileServer(http.Dir(".")))

	// Creates a new server and listnens for incoming requests
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()

}
