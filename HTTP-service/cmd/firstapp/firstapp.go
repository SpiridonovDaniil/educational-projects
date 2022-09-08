package main

import (
	"fmt"
	"net/http"

	"serverdb/repository/mongo"
	"serverdb/server"
	"serverdb/server/service"
)

func main() {
	mongoRepo := mongo.NewRepo("mongodb://mongo:27017", "userdb", "users")
	defer mongoRepo.Close()

	service := service.NewService(mongoRepo)

	mux := server.NewMux(service)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err)
	}
}
