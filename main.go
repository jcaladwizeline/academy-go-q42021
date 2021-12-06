package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jcaladwizeline/academy-go-q42021/api"
)

func main() {
	router := api.New()
	http.Handle("/", router)
	log.Println("Api running on port 8080")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}
}
