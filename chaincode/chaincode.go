package chaincode

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract ...
type SmartContract struct {
	contractapi.Contract
}

// CarAsset ...
type CarAsset struct {
	ID             string `json:"id"`
	Brand          string `json:"brand"`
	Owner          string `json:"owner"`
	TransfersCount int    `json:"transfersCount"`
}

// InitLedger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	cars := []CarAsset{
		{Brand: "Toyota", ID: "12", Owner: "Juan", TransfersCount: 0},
		{Brand: "Honda", ID: "22", Owner: "Marcos", TransfersCount: 0},
	}

	for _, car := range cars {
		carJSON, err := json.Marshal(car)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(car.ID, carJSON)
		if err != nil {
			return fmt.Errorf("failed operation, %v", err)
		}
	}

	return nil
}

// GetCars ...
func (s *SmartContract) GetCars(ctx contractapi.TransactionContextInterface) ([]*CarAsset, error) {
	res, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var cars []*CarAsset

	for res.HasNext() {
		car, err := res.Next()
		if err != nil {
			return nil, err
		}

		var carAsset CarAsset
		err = json.Unmarshal(car.Value, &carAsset)
		if err != nil {
			return nil, err
		}

		cars = append(cars, &carAsset)
	}

	return cars, nil
}

// GetCarsByOwner
func (s *SmartContract) GetCarsByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]*CarAsset, error) {
	cars, err := s.GetCars(ctx)
	if err != nil {
		return nil, err
	}

	if len(cars) == 0 {
		return []*CarAsset{}, nil
	}

	var ownerCars []*CarAsset
	for _, car := range cars {
		if car.Owner == owner {
			ownerCars = append(ownerCars, car)
		}
	}

	return ownerCars, nil
}

// TransferCart ...
func (s *SmartContract) TransferCart(ctx contractapi.TransactionContextInterface, id, newOwner string) error {
	car, err := s.GetCar(ctx, id)
	if err != nil {
		return err
	}

	if ok, errTran := s.IsAbleToTransfer(car, newOwner); !ok {
		return errTran
	}

	car.Owner = newOwner
	car.TransfersCount++

	carJSON, err := json.Marshal(car)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, carJSON)
}

// IsAbleToTransfer ...
func (s *SmartContract) IsAbleToTransfer(car *CarAsset, newOwner string) (bool, error) {
	if car == nil {
		return false, fmt.Errorf("unable to process transaction, car does not exit")
	}

	if car.Owner == newOwner {
		return false, fmt.Errorf("unable to process transaction, car owner %s is equal to %s", car.Owner, newOwner)
	}

	totalTransAllwed := 3
	if car.TransfersCount >= totalTransAllwed {
		return false, fmt.Errorf("unable to process, total car transaction %d exeed the limit", car.TransfersCount)
	}

	return true, nil
}

// GetCar ...
func (s *SmartContract) GetCar(ctx contractapi.TransactionContextInterface, id string) (*CarAsset, error) {
	carJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("error getting car, %v", err)
	}

	if carJSON == nil {
		return nil, fmt.Errorf("car does not exist ID: %s", id)
	}

	var car CarAsset
	err = json.Unmarshal(carJSON, &car)
	if err != nil {
		return nil, err
	}

	return &car, nil
}

// CreateCar ...
func (s *SmartContract) CreateCar(ctx contractapi.TransactionContextInterface, id, brand, owner string) error {
	exist, err := s.ExistCar(ctx, id)
	if err != nil {
		return err
	}

	if exist {
		return fmt.Errorf("the car with id %s already exist", id)
	}

	if strings.TrimSpace(brand) == "" ||
		strings.TrimSpace(id) == "" ||
		strings.TrimSpace(owner) == "" {
		return fmt.Errorf("All fields are required")
	}

	newCar := CarAsset{
		Brand: brand,
		ID:    id,
		Owner: owner,
	}

	carJSON, err := json.Marshal(newCar)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, carJSON)
}

// ExistCar ...
func (s *SmartContract) ExistCar(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	carJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, err
	}

	return carJSON != nil, nil
}
