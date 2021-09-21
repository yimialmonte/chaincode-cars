package chaincode

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/yimialmonte/chaincode-cars/asset"
)

// SmartContract ...
type SmartContract struct {
	contractapi.Contract
}

// InitLedger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	cars := []asset.Car{
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
func (s *SmartContract) GetCars(ctx contractapi.TransactionContextInterface) ([]*asset.Car, error) {
	res, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var cars []*asset.Car

	for res.HasNext() {
		car, err := res.Next()
		if err != nil {
			return nil, err
		}

		var carAsset asset.Car
		err = json.Unmarshal(car.Value, &carAsset)
		if err != nil {
			return nil, err
		}

		cars = append(cars, &carAsset)
	}

	return cars, nil
}

// GetCarsByOwner
func (s *SmartContract) GetCarsByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]*asset.Car, error) {
	cars, err := s.GetCars(ctx)
	if err != nil {
		return nil, err
	}

	if len(cars) == 0 {
		return []*asset.Car{}, nil
	}

	var ownerCars []*asset.Car
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
func (s *SmartContract) IsAbleToTransfer(car *asset.Car, newOwner string) (bool, error) {
	if car == nil {
		return false, fmt.Errorf("unable to process transaction, car does not exist")
	}

	if car.Owner == newOwner {
		return false, fmt.Errorf("unable to process transaction, car owner %s is equal to %s", car.Owner, newOwner)
	}

	maxTransfert := 3
	if car.TransfersCount >= maxTransfert {
		return false, fmt.Errorf("unable to process, total car transaction %d exceed the limit", car.TransfersCount)
	}

	return true, nil
}

// GetCar ...
func (s *SmartContract) GetCar(ctx contractapi.TransactionContextInterface, id string) (*asset.Car, error) {
	carJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("error getting car, %v", err)
	}

	if carJSON == nil {
		return nil, fmt.Errorf("car does not exist ID: %s", id)
	}

	var car asset.Car
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

	newCar := asset.Car{
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
