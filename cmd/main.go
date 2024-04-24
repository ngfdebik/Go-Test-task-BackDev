package main

import (
	"log"

	"example.com/m/pkg/api"
	"github.com/gorilla/mux"
)

func main() {
	api := api.New(mux.NewRouter())
	api.FillEndpoints()
	log.Fatal(api.ListenAndServe("localhost:8090"))
}
