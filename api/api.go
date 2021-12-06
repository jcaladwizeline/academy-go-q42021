package api

import (
	controllers "bootCampApi/api/controllers"

	"github.com/gorilla/mux"
)

func New() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/health", controllers.HealthCheck).Methods("GET")
	router.HandleFunc("/animes", controllers.GetAllAnimes).Methods("GET")
	router.HandleFunc("/animes/{id:[0-9]+}", controllers.GetAnimeById).Methods("GET")
	router.HandleFunc("/animes/{id:[0-9]+}", controllers.PostAnimeById).Methods("POST")
	router.HandleFunc("/animesWorker/{type}/{items:[0-9]+}/{items_per_workers:[0-9]+}", controllers.GetAnimesWorker).Methods("POST")

	return router
}
