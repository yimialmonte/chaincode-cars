package repository

import (
	"github.com/yimialmonte/chaincode-cars/asset"
)

// Car ...
type Car struct {
}

// GetCars ...
func (c *Car) GetCars() ([]*asset.Car, error) {
	return []*asset.Car{
		{ID: "01", Brand: "Toyota", Owner: "Peter"},
		{ID: "02", Brand: "Homnda", Owner: "Max"},
	}, nil
}

// GetCarsByOwner ...
func (c *Car) GetCarsByOwner(owner string) ([]*asset.Car, error) {
	return []*asset.Car{
		{ID: "02", Brand: "Homnda", Owner: "Max"},
	}, nil
}

// TransferCart ...
func (c *Car) TransferCart(id, owner string) error {
	return nil
}
