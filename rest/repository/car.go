package repository

import "github.com/yimialmonte/chaincode-cars/chaincode"

// Car ...
type Car struct {
}

// GetCars ...
func (c *Car) GetCars() ([]*chaincode.CarAsset, error) {
	return []*chaincode.CarAsset{
		{ID: "01", Brand: "Toyota", Owner: "Peter"},
		{ID: "02", Brand: "Homnda", Owner: "Max"},
	}, nil
}

// GetCarsByOwner ...
func (c *Car) GetCarsByOwner(owner string) ([]*chaincode.CarAsset, error) {
	return []*chaincode.CarAsset{
		{ID: "02", Brand: "Homnda", Owner: "Max"},
	}, nil
}

// TransferCart ...
func (c *Car) TransferCart(id, owner string) error {
	return nil
}
