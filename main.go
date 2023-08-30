package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := 8080
	addr := fmt.Sprintf("localhost:%d", port)
	fileSysRoot := "."
	apiCfg := apiConfig{}

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(fileSysRoot)))

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.metricsMiddleware(fileServerHandler))
	mux.HandleFunc("/healthz", handleReadiness)
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hits: %d", apiCfg.fileServerHits)))
	})

	corsMux := CorsMiddleware(mux)

	server := &http.Server{Addr: addr, Handler: corsMux}

	log.Printf("Server listening on port: 8080")
	log.Fatal(server.ListenAndServe())
}

type apiConfig struct {
	fileServerHits int
}

func handleReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer next.ServeHTTP(w, r)

		cfg.fileServerHits++
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
