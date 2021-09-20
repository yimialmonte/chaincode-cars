package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yimialmonte/chaincode-cars/rest/handler"
	"github.com/yimialmonte/chaincode-cars/rest/repository"
)

func main() {
	route := mux.NewRouter()

	route.Handle("/cars", &handler.GetAllCars{Store: &repository.Car{}}).Methods(http.MethodGet)
	route.Handle("/cars/owner/{id}", &handler.GetCarsOwner{Store: &repository.Car{}}).Methods(http.MethodGet)
	route.Handle("/cars", &handler.TransferCarOwner{Store: &repository.Car{}}).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", route))
}
