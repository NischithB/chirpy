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

	rootRouter := chi.NewRouter()
	rootRouter.Use(CorsMiddleware)
	rootRouter.Handle("/app", apiCfg.metricsMiddleware(fileServerHandler))
	rootRouter.Handle("/app/*", apiCfg.metricsMiddleware(fileServerHandler))

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handleMetrics)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handleReadiness)

	rootRouter.Mount("/admin", adminRouter)
	rootRouter.Mount("/api", apiRouter)

	server := &http.Server{Addr: addr, Handler: rootRouter}

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
