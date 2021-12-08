package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jcaladwizeline/academy-go-q42021/api"
	c "github.com/jcaladwizeline/academy-go-q42021/config"

	"github.com/spf13/viper"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}
	router := api.New()
	http.Handle("/", router)
	log.Println("Api running on port 8080")

	err = http.ListenAndServe(config.Server.Port, router)
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}
}

func LoadConfig() (config c.Config, err error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	if err = viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&config)

	return
}
