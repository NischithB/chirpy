package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NischithB/chirpy/config"
	"github.com/NischithB/chirpy/controllers"
	"github.com/NischithB/chirpy/middlewares"
	"github.com/go-chi/chi"
)

func main() {
	port := 8080
	addr := fmt.Sprintf("localhost:%d", port)

	config.LoadEnvVars()
	config.ConfigureDB()

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	rootRouter := chi.NewRouter()
	rootRouter.Use(CorsMiddleware)
	rootRouter.Handle("/app", middlewares.FileServerHitCounter(fileServerHandler))
	rootRouter.Handle("/app/*", middlewares.FileServerHitCounter(fileServerHandler))

	rootRouter.Mount("/admin", controllers.GetAdminRouter())
	rootRouter.Mount("/api", controllers.GetAPIRouter())

	server := &http.Server{Addr: addr, Handler: rootRouter}

	log.Printf("Server listening on port: 8080")
	log.Fatal(server.ListenAndServe())
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
