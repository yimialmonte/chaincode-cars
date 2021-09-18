package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/stretchr/testify/assert"
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

func TestGetCars(t *testing.T) {
	car := CarAsset{Brand: "Toyota", ID: "123", Owner: "Max"}
	bytes, err := json.Marshal(car)
	require.NoError(t, err)

	it := &mocks.StateQueryIterator{}
	it.HasNextReturnsOnCall(0, true)
	it.HasNextReturnsOnCall(1, false)
	it.NextReturns(&queryresult.KV{Value: bytes}, nil)

	stub := &mocks.ChaincodeStub{}
	tctx := &mocks.TransactionContext{}
	tctx.GetStubReturns(stub)

	stub.GetStateByRangeReturns(it, nil)
	sc := &SmartContract{}
	cars, err := sc.GetCars(tctx)
	assert.Nil(t, err)
	assert.Equal(t, []*CarAsset{&car}, cars)

	it.HasNextReturns(true)
	it.NextReturns(nil, fmt.Errorf("error getting next car"))
	cars, err = sc.GetCars(tctx)
	assert.EqualError(t, err, "error getting next car")
	assert.Nil(t, cars)

	stub.GetStateByRangeReturns(nil, fmt.Errorf("error getting cars"))
	cars, err = sc.GetCars(tctx)
	assert.EqualError(t, err, "error getting cars")
	assert.Nil(t, cars)
}

func TestGetCarsByOwner(t *testing.T) {
	car := CarAsset{Brand: "Toyota", ID: "123", Owner: "Max"}
	bytes, err := json.Marshal(car)
	require.NoError(t, err)

	it := &mocks.StateQueryIterator{}
	it.HasNextReturnsOnCall(0, true)
	it.HasNextReturnsOnCall(1, false)
	it.NextReturns(&queryresult.KV{Value: bytes}, nil)

	stub := &mocks.ChaincodeStub{}
	tctx := &mocks.TransactionContext{}
	tctx.GetStubReturns(stub)

	stub.GetStateByRangeReturns(it, nil)
	sc := &SmartContract{}
	cars, err := sc.GetCarsByOwner(tctx, "Max")
	assert.Nil(t, err)
	assert.Equal(t, []*CarAsset{&car}, cars)

	cars, err = sc.GetCarsByOwner(tctx, "Peter")
	assert.Nil(t, err)
	assert.Equal(t, []*CarAsset{}, cars)

	it.HasNextReturns(true)
	it.NextReturns(nil, fmt.Errorf("error getting next car"))
	cars, err = sc.GetCarsByOwner(tctx, "Juan")
	assert.EqualError(t, err, "error getting next car")
	assert.Nil(t, cars)

}

func TestGetCar(t *testing.T) {
	stu := &mocks.ChaincodeStub{}
	tctx := &mocks.TransactionContext{}
	tctx.GetStubReturns(stu)

	expectedCar := CarAsset{ID: "123"}
	b, err := json.Marshal(expectedCar)
	assert.Nil(t, err)

	stu.GetStateReturns(b, nil)
	sc := SmartContract{}
	car, err := sc.GetCar(tctx, "")
	assert.Nil(t, err)
	assert.Equal(t, &expectedCar, car)

	stu.GetStateReturns(nil, fmt.Errorf("connection error"))
	_, err = sc.GetCar(tctx, "")
	assert.EqualError(t, err, "error getting car, connection error")

	stu.GetStateReturns(nil, nil)
	car, err = sc.GetCar(tctx, "000")
	assert.EqualError(t, err, "car does not exist ID: 000")
	assert.Nil(t, car)
}

func TestCreateCar(t *testing.T) {
	stu := &mocks.ChaincodeStub{}
	tctx := &mocks.TransactionContext{}
	tctx.GetStubReturns(stu)

	sc := SmartContract{}
	err := sc.CreateCar(tctx, "", "", "")
	assert.Nil(t, err)

	stu.GetStateReturns([]byte{}, nil)
	err = sc.CreateCar(tctx, "00", "0", "")
	assert.EqualError(t, err, "the car with id 00 already exist")

	stu.GetStateReturns(nil, fmt.Errorf("unable to process"))
	err = sc.CreateCar(tctx, "00", "", "")
	assert.EqualError(t, err, "unable to process")

}

func TestExistcar(t *testing.T) {
	stu := &mocks.ChaincodeStub{}
	tctx := &mocks.TransactionContext{}
	tctx.GetStubReturns(stu)

	sc := SmartContract{}
	_, err := sc.ExistCar(tctx, "")
	assert.Nil(t, err)

	stu.GetStateReturns([]byte{}, nil)
	exist, _ := sc.ExistCar(tctx, "")
	assert.Equal(t, exist, true)

	stu.GetStateReturns(nil, nil)
	exist, _ = sc.ExistCar(tctx, "")
	assert.Equal(t, exist, false)

	stu.GetStateReturns(nil, fmt.Errorf("unable to process"))
	exist, err = sc.ExistCar(tctx, "")
	assert.EqualError(t, err, "unable to process")
	assert.Equal(t, exist, false)
}

func TestIsAbleToTransfer(t *testing.T) {
	tests := []struct {
		car            *CarAsset
		newOwner       string
		expectedErr    error
		expectedResult bool
	}{
		{
			&CarAsset{Owner: "Juan", TransfersCount: 3},
			"Max",
			errors.New(fmt.Sprintf("unable to process, total car transaction %d exeed the limit", 3)),
			false,
		},
		{
			&CarAsset{Owner: "Juan", TransfersCount: 1},
			"Juan",
			errors.New(fmt.Sprintf("unable to process transaction, car owner %s is equal to %s", "Juan", "Juan")),
			false,
		},
		{
			&CarAsset{Owner: "Peter", TransfersCount: 3},
			"Peter",
			errors.New(fmt.Sprintf("unable to process transaction, car owner %s is equal to %s", "Peter", "Peter")),
			false,
		},
		{
			nil,
			"Peter",
			errors.New("unable to process transaction, car does not exit"),
			false,
		},
		{
			&CarAsset{Owner: "Ana", TransfersCount: 1},
			"Peter",
			nil,
			true,
		},
		{
			&CarAsset{Owner: "Mar", TransfersCount: 0},
			"Jackson",
			nil,
			true,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			sc := SmartContract{}
			res, err := sc.IsAbleToTransfer(test.car, test.newOwner)
			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}
