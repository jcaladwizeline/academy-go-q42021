package controller

import (
	"fmt"
	"net/http"
	"strconv"

	animeU "github.com/jcaladwizeline/academy-go-q42021/api/service"

	"github.com/gorilla/mux"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Ok...")
}

func GetAllAnimes(w http.ResponseWriter, r *http.Request) {
	response := animeU.GetAllAnimes()
	if response == nil {
		responseHandle(w, nil, http.StatusInternalServerError)

		return
	}

	responseHandle(w, response, http.StatusOK)
}

func GetAnimeById(w http.ResponseWriter, r *http.Request) {
	// check params
	id := mux.Vars(r)["id"]
	response := animeU.GetAnimeById(id)

	responseHandle(w, response, http.StatusOK)
}

func PostAnimeById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	httpStatus := animeU.PostAnimeById(id)
	if httpStatus == http.StatusInternalServerError {
		responseHandle(w, nil, http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Ok...")
}

func GetAnimesWorker(w http.ResponseWriter, r *http.Request) {
	itemsQuery := mux.Vars(r)["items"]
	items, _ := strconv.Atoi(itemsQuery)
	itemsW := mux.Vars(r)["items_per_workers"]
	itemsPerWorker, _ := strconv.Atoi(itemsW)
	jobType := mux.Vars(r)["type"]

	animes, _ := animeU.WorkerPool(items, itemsPerWorker, jobType)

	responseHandle(w, animes, http.StatusInternalServerError)
}
