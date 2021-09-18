package chaincode

import (
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
