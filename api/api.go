package api

import (
	"github.com/gorilla/mux"
	"github.com/jcaladwizeline/academy-go-q42021/api/controller"
)

func New() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/health", controller.HealthCheck).Methods("GET")
	router.HandleFunc("/animes", controller.GetAllAnimes).Methods("GET")
	router.HandleFunc("/animes/{id:[0-9]+}", controller.GetAnimeById).Methods("GET")
	router.HandleFunc("/animes/{id:[0-9]+}", controller.PostAnimeById).Methods("POST")
	router.HandleFunc("/animesWorker/{type}/{items:[0-9]+}/{items_per_workers:[0-9]+}", controller.GetAnimesWorker).Methods("POST")

	return router
}
