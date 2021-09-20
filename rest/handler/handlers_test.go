package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yimialmonte/chaincode-cars/chaincode"
)

type testCartStore struct {
	called       int
	carsResponse []*chaincode.CarAsset
	errResponse  error
}

func (t *testCartStore) GetCars() ([]*chaincode.CarAsset, error) {
	t.called++
	return t.carsResponse, t.errResponse
}

func (t *testCartStore) GetCarsByOwner(owner string) ([]*chaincode.CarAsset, error) {
	t.called++
	return t.carsResponse, t.errResponse
}

func (t *testCartStore) TransferCart(id, owner string) error {
	t.called++
	return t.errResponse
}
func TestGetAllCars(t *testing.T) {
	tests := []struct {
		response     []*chaincode.CarAsset
		expectedErr  error
		expectedCode int
		expectedRes  string
	}{
		{
			[]*chaincode.CarAsset{
				{ID: "000", Brand: "Honda", Owner: "Juan"},
			},
			nil,
			http.StatusOK,
			string(`[{"id":"000","brand":"Honda","owner":"Juan","transfersCount":0}]`),
		},
		{
			[]*chaincode.CarAsset{},
			fmt.Errorf("internal server error"),
			http.StatusInternalServerError,
			fmt.Sprintf("internal server error\n"),
		},
		{
			[]*chaincode.CarAsset{},
			nil,
			http.StatusOK,
			"[]",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/cars", nil)
			record := httptest.NewRecorder()

			store := &testCartStore{carsResponse: test.response, errResponse: test.expectedErr}
			car := GetAllCars{Store: store}

			car.ServeHTTP(record, r)

			assert.Equal(t, test.expectedCode, record.Code)
			assert.Equal(t, store.called, 1)
			assert.Equal(t, test.expectedRes, record.Body.String())
		})
	}
}

func TestGetCarsOwner(t *testing.T) {
	tests := []struct {
		response     []*chaincode.CarAsset
		expectedErr  error
		expectedCode int
		expectedRes  string
	}{
		{
			[]*chaincode.CarAsset{
				{ID: "000", Brand: "Toyota", Owner: "Max", TransfersCount: 10},
			},
			nil,
			http.StatusOK,
			string(`[{"id":"000","brand":"Toyota","owner":"Max","transfersCount":10}]`),
		},
		{
			[]*chaincode.CarAsset{},
			nil,
			http.StatusOK,
			"[]",
		},
		{
			[]*chaincode.CarAsset{},
			fmt.Errorf("internal server error"),
			http.StatusInternalServerError,
			fmt.Sprintf("internal server error\n"),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/cars/owner/juan", nil)
			record := httptest.NewRecorder()

			store := &testCartStore{carsResponse: test.response, errResponse: test.expectedErr}
			car := GetCarsOwner{Store: store}

			car.ServeHTTP(record, r)

			assert.Equal(t, test.expectedCode, record.Code)
			assert.Equal(t, store.called, 1)
			assert.Equal(t, test.expectedRes, record.Body.String())
		})
	}
}

func TestTransferCarOwner(t *testing.T) {
	tests := []struct {
		requestBody  string
		expectedErr  error
		expectedCode int
		respond      string
	}{
		{
			`{"id":"000","owner":"Max"}`,
			nil,
			http.StatusCreated,
			"",
		},
		{
			`{"id":"000","owner":"Max"}`,
			fmt.Errorf("internal server error"),
			http.StatusInternalServerError,
			fmt.Sprintf("internal server error\n"),
		},
		{
			`{"id":"","owner":""}`,
			nil,
			http.StatusBadRequest,
			fmt.Sprintf("Supply Car ID and Owner\n"),
		},
		{
			`{"id":"","owner"}`,
			nil,
			http.StatusBadRequest,
			fmt.Sprintf("error decoding json\n"),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/cars", strings.NewReader(test.requestBody))
			record := httptest.NewRecorder()

			store := &testCartStore{errResponse: test.expectedErr}
			car := TransferCarOwner{Store: store}

			car.ServeHTTP(record, r)

			assert.Equal(t, test.expectedCode, record.Code)
			assert.Equal(t, test.respond, record.Body.String())
		})
	}
}
