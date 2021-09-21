package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yimialmonte/chaincode-cars/asset"
)

// CarStore ...
type CarStore interface {
	GetCars() ([]*asset.Car, error)
	GetCarsByOwner(owner string) ([]*asset.Car, error)
	TransferCart(id, owner string) error
}

// GetAllCars
type GetAllCars struct {
	Store CarStore
}

func (g *GetAllCars) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cars, err := g.Store.GetCars()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	carsJSON, err := json.Marshal(cars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(carsJSON)
}

// GetCarsOwner
type GetCarsOwner struct {
	Store CarStore
}

func (g *GetCarsOwner) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	cars, err := g.Store.GetCarsByOwner(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	carsJSON, err := json.Marshal(cars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(carsJSON)
}

// TransferCarOwner ...
type TransferCarOwner struct {
	Store CarStore
}

func (g *TransferCarOwner) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	car, err := getCarFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(car.ID) == "" || strings.TrimSpace(car.Owner) == "" {
		http.Error(w, errors.New("Supply Car ID and Owner").Error(), http.StatusBadRequest)
		return
	}

	err = g.Store.TransferCart(car.ID, car.Owner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getCarFromRequest(r *http.Request) (asset.Car, error) {
	var car asset.Car
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&car)
	if err != nil {
		return car, fmt.Errorf("error decoding json")
	}

	return car, nil
}
