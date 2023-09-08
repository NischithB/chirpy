package controllers

import (
	"fmt"
	"net/http"

	"github.com/NischithB/chirpy/config"
	"github.com/go-chi/chi"
)

func GetAPIRouter() chi.Router {
	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handleReadiness)
	apiRouter.Mount("/chirps", getChirpRouter())
	apiRouter.Mount("/users", getUsersRouter())
	apiRouter.Mount("/", getAuthController())
	apiRouter.Mount("/polka", getPaymentController())
	return apiRouter
}

func GetAdminRouter() chi.Router {
	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", handleMetrics)
	return adminRouter
}

func handleReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	body := fmt.Sprintf(`
	<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
	</html>
	`, config.Config.FileServerHits)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}
