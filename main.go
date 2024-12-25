package main

import (
	"billy/api"
	"billy/handlers"
	"billy/storage/postgres"
	"log"
)

func main() {
	storage, err := postgres.NewPostgresStorage("user=postgres password=faxa host=localhost port=5432 database=tempi sslmode=disable")
	if err != nil {
		log.Fatal("unable to initialize storage")
	}

	handler := handlers.NewHandler(storage)
	api := api.NewApi(*handler)
	api.SetupRoutes()

	if err := api.Run(); err != nil {
		log.Fatal(err)
	}
}
