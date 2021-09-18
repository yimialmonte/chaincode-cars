package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract ...
type SmartContract struct {
	contractapi.Contract
}

// CarAsset ...
type CarAsset struct {
	Brand          string `json:"brand"`
	ID             string `json:"id"`
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
