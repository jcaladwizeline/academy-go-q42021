package api

import (
	"github.com/jcaladwizeline/academy-go-q42021/api/controller"

	"github.com/gorilla/mux"
)

func New() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/animes", controller.GetAllAnimes).Methods("GET")
	router.HandleFunc("/animes/{id}", controller.GetAnimeById).Methods("GET")
	router.HandleFunc("/animes/{id}", controller.PostAnimeById).Methods("POST")
	router.HandleFunc("/animesWorker", controller.GetAnimesWorker).Methods("GET").Queries()

	return router
}
