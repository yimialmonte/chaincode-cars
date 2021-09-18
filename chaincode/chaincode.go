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
