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
	"github.com/yimialmonte/chaincode-cars/asset"
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

type stateReturn struct {
	state []byte
	err   error
}

func TestInitLedger(t *testing.T) {
	tests := []struct {
		state       stateReturn
		expectedErr error
	}{
		{
			stateReturn{nil, nil},
			nil,
		},
		{
			stateReturn{nil, errors.New("error")},
			errors.New("failed operation, error"),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			stu := &mocks.ChaincodeStub{}
			tctx := &mocks.TransactionContext{}
			tctx.GetStubReturns(stu)

			sc := SmartContract{}
			stu.PutStateReturns(test.state.err)

			err := sc.InitLedger(tctx)
			require.Equal(t, test.expectedErr, err)
		})
	}
}

func TestGetCars(t *testing.T) {
	car := asset.Car{ID: "123", Owner: "Peter", Brand: "Honda"}
	b, err := json.Marshal(car)
	assert.Nil(t, err)

	var cars []*asset.Car

	tests := []struct {
		state       stateReturn
		expectedErr error
		expectedCar []*asset.Car
	}{
		{
			stateReturn{b, nil},
			nil,
			[]*asset.Car{&car},
		},
		{
			stateReturn{nil, errors.New("connection failed")},
			errors.New("connection failed"),
			cars,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			it := &mocks.StateQueryIterator{}
			it.HasNextReturnsOnCall(0, true)
			it.HasNextReturnsOnCall(1, false)
			it.NextReturns(&queryresult.KV{Value: test.state.state}, nil)

			stub := &mocks.ChaincodeStub{}
			tctx := &mocks.TransactionContext{}
			tctx.GetStubReturns(stub)

			stub.GetStateByRangeReturns(it, test.state.err)
			sc := &SmartContract{}
			cars, err := sc.GetCars(tctx)
			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedCar, cars)
		})
	}
}

func TestGetCarsByOwner(t *testing.T) {
	car := asset.Car{ID: "000", Owner: "Max", Brand: "Toyota"}
	b, err := json.Marshal(car)
	assert.Nil(t, err)
	var cars []*asset.Car
	tests := []struct {
		owner       string
		state       stateReturn
		expectedErr error
		expectedCar []*asset.Car
	}{
		{
			"Max",
			stateReturn{b, nil},
			nil,
			[]*asset.Car{&car},
		},
		{
			"Juan",
			stateReturn{b, nil},
			nil,
			cars,
		},
		{
			"Juan",
			stateReturn{nil, errors.New("connection failed")},
			errors.New("connection failed"),
			cars,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			it := &mocks.StateQueryIterator{}
			it.HasNextReturnsOnCall(0, true)
			it.HasNextReturnsOnCall(1, false)
			it.NextReturns(&queryresult.KV{Value: test.state.state}, nil)

			stub := &mocks.ChaincodeStub{}
			tctx := &mocks.TransactionContext{}
			tctx.GetStubReturns(stub)

			stub.GetStateByRangeReturns(it, test.state.err)
			sc := &SmartContract{}
			cars, err := sc.GetCarsByOwner(tctx, test.owner)
			fmt.Println(err)
			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedCar, cars)
		})
	}
}

func TestGetCar(t *testing.T) {
	car := asset.Car{ID: "000", Owner: "Max", Brand: "Toyota"}
	b, err := json.Marshal(car)
	assert.Nil(t, err)

	tests := []struct {
		id          string
		state       stateReturn
		expectedErr error
		expectedCar *asset.Car
	}{
		{
			"000",
			stateReturn{b, nil},
			nil,
			&car,
		},
		{
			"000",
			stateReturn{nil, errors.New("connection failed")},
			errors.New("error getting car, connection failed"),
			nil,
		},
		{
			"000",
			stateReturn{nil, nil},
			errors.New("car does not exist ID: 000"),
			nil,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			stu := &mocks.ChaincodeStub{}
			tctx := &mocks.TransactionContext{}
			tctx.GetStubReturns(stu)

			stu.GetStateReturns(test.state.state, test.state.err)

			sc := SmartContract{}
			car, err := sc.GetCar(tctx, test.id)
			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedCar, car)
		})
	}
}

func TestCreateCar(t *testing.T) {
	tests := []struct {
		state       stateReturn
		expectedErr error
		car         asset.Car
	}{
		{
			stateReturn{nil, nil},
			nil,
			asset.Car{ID: "11", Brand: "Toyota", Owner: "Peter"},
		},
		{
			stateReturn{[]byte{}, nil},
			errors.New("the car with id 11 already exist"),
			asset.Car{ID: "11", Brand: "Toyota", Owner: "Peter"},
		},
		{
			stateReturn{nil, errors.New("connection failed")},
			errors.New("connection failed"),
			asset.Car{ID: "11", Brand: "Toyota", Owner: "Peter"},
		},
		{
			stateReturn{nil, nil},
			errors.New("All fields are required"),
			asset.Car{ID: "11", Brand: "Toyota"},
		},
		{
			stateReturn{nil, nil},
			errors.New("All fields are required"),
			asset.Car{ID: "11", Owner: "Max"},
		},
		{
			stateReturn{nil, nil},
			errors.New("All fields are required"),
			asset.Car{Brand: "Honda", Owner: "Max"},
		},
		{
			stateReturn{nil, nil},
			errors.New("All fields are required"),
			asset.Car{ID: " ", Brand: " ", Owner: "Max"},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			stu := &mocks.ChaincodeStub{}
			tctx := &mocks.TransactionContext{}
			tctx.GetStubReturns(stu)

			sc := SmartContract{}
			stu.GetStateReturns(test.state.state, test.state.err)
			err := sc.CreateCar(tctx, test.car.ID, test.car.Brand, test.car.Owner)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestExistcar(t *testing.T) {
	tests := []struct {
		state         stateReturn
		expectedFound bool
		expectedErr   error
	}{
		{
			stateReturn{[]byte{}, nil},
			true,
			nil,
		},
		{
			stateReturn{nil, nil},
			false,
			nil,
		},
		{
			stateReturn{nil, errors.New("Error occurs")},
			false,
			errors.New("Error occurs"),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			stu := &mocks.ChaincodeStub{}
			tctx := &mocks.TransactionContext{}
			tctx.GetStubReturns(stu)

			sc := SmartContract{}
			stu.GetStateReturns(test.state.state, test.state.err)
			res, err := sc.ExistCar(tctx, "")
			assert.Equal(t, test.expectedFound, res)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestIsAbleToTransfer(t *testing.T) {
	tests := []struct {
		car            *asset.Car
		newOwner       string
		expectedErr    error
		expectedResult bool
	}{
		{
			&asset.Car{Owner: "Juan", TransfersCount: 3},
			"Max",
			errors.New(fmt.Sprintf("unable to process, total car transaction %d exceed the limit", 3)),
			false,
		},
		{
			&asset.Car{Owner: "Juan", TransfersCount: 1},
			"Juan",
			errors.New(fmt.Sprintf("unable to process transaction, car owner %s is equal to %s", "Juan", "Juan")),
			false,
		},
		{
			&asset.Car{Owner: "Peter", TransfersCount: 3},
			"Peter",
			errors.New(fmt.Sprintf("unable to process transaction, car owner %s is equal to %s", "Peter", "Peter")),
			false,
		},
		{
			nil,
			"Peter",
			errors.New("unable to process transaction, car does not exist"),
			false,
		},
		{
			&asset.Car{Owner: "Ana", TransfersCount: 1},
			"Peter",
			nil,
			true,
		},
		{
			&asset.Car{Owner: "Mar", TransfersCount: 0},
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

func TestTransferCart(t *testing.T) {
	tests := []struct {
		car         asset.Car
		expectedErr error
		newOwner    string
	}{
		{
			asset.Car{
				Brand:          "Toyota",
				ID:             "123",
				Owner:          "Max",
				TransfersCount: 1,
			},
			nil,
			"Peter",
		},
		{
			asset.Car{
				Brand:          "Toyota",
				ID:             "123",
				Owner:          "Max",
				TransfersCount: 5,
			},
			errors.New("unable to process, total car transaction 5 exceed the limit"),
			"Peter",
		},
		{
			asset.Car{
				Brand:          "Toyota",
				ID:             "123",
				Owner:          "Max",
				TransfersCount: 1,
			},
			errors.New("unable to process transaction, car owner Max is equal to Max"),
			"Max",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			stub := &mocks.ChaincodeStub{}
			tctx := &mocks.TransactionContext{}
			tctx.GetStubReturns(stub)

			bytes, err := json.Marshal(test.car)
			require.NoError(t, err)

			stub.GetStateReturns(bytes, nil)
			sc := &SmartContract{}
			err = sc.TransferCart(tctx, "", test.newOwner)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
