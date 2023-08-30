package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	port := 8080
	addr := fmt.Sprintf("localhost:%d", port)
	fileSysRoot := "."
	apiCfg := apiConfig{}

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(fileSysRoot)))

	r := chi.NewRouter()
	r.Use(CorsMiddleware)
	r.Handle("/app", apiCfg.metricsMiddleware(fileServerHandler))
	r.Handle("/app/*", apiCfg.metricsMiddleware(fileServerHandler))
	r.Get("/healthz", handleReadiness)
	r.Get("/metrics", apiCfg.handleMetrics)

	server := &http.Server{Addr: addr, Handler: r}

	log.Printf("Server listening on port: 8080")
	log.Fatal(server.ListenAndServe())
}

func handleReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
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
