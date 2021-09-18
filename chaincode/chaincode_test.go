package chaincode

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/require"
	"github.com/yimialmonte/chaincode-cars/chaincode/mocks"
)

//go:generate counterfeiter -o mocks/transaction.go -fake-name TransactionContext . transactionContext
type transactionContext interface {
	contractapi.TransactionContextInterface
}

//go:generate counterfeiter -o mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//go:generate counterfeiter -o mocks/statequeryiterator.go -fake-name StateQueryIterator . stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

func TestInitLedger(t *testing.T) {
	stu := &mocks.ChaincodeStub{}
	tctx := &mocks.TransactionContext{}
	tctx.GetStubReturns(stu)

	sc := SmartContract{}
	err := sc.InitLedger(tctx)
	require.NoError(t, err)

	stu.PutStateReturns(fmt.Errorf("failed initledger"))
	err = sc.InitLedger(tctx)
	require.EqualError(t, err, "failed operation, failed initledger")
}
