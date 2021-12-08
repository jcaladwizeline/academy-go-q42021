package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	animeU "github.com/jcaladwizeline/academy-go-q42021/api/service"

	"github.com/gorilla/mux"
)

type ControllerInterfaceGet interface {
	Get(w http.ResponseWriter, r *http.Request)
}

func GetAllAnimes(w http.ResponseWriter, r *http.Request) {
	response, status, err := animeU.GetAllAnimes()
	if err != nil {
		http.Error(w, err.Error(), status)

		return
	}

	responseHandle(w, response, status)
}

func GetAnimeById(w http.ResponseWriter, r *http.Request) {
	// check params
	id := mux.Vars(r)["id"]
	animeID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	response, err := animeU.GetAnimeById(animeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	responseHandle(w, response, http.StatusOK)
}

func PostAnimeById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	response, err := animeU.PostAnimeById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	responseHandle(w, response, http.StatusOK)
}

func GetAnimesWorker(w http.ResponseWriter, r *http.Request) {
	itemsQuery := r.FormValue("items")
	items, err := strconv.Atoi(itemsQuery)
	if err != nil && itemsQuery != "" {
		log.Println("Unable convert string into int", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	itemsW := r.FormValue("items_per_workers")
	itemsPerWorker, err := strconv.Atoi(itemsW)
	if err != nil && itemsW != "" {
		log.Println("Unable convert string into int", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	jobType := r.FormValue("type")
	if (jobType != "odd" && jobType != "even") && jobType != "" {
		log.Println("Unable convert string into int", err)
		http.Error(w, errors.New("type value ").Error(), http.StatusBadRequest)

		return
	}

	animes, err := animeU.WorkerPool(items, itemsPerWorker, jobType)
	if err != nil {
		log.Println("Unable convert string into int", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	responseHandle(w, animes, http.StatusInternalServerError)
}
