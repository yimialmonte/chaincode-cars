# Chaincode Cars

Provide a chaincode (smartcontract) for Hyperledger Fabric platform \
REST API for making transaction inside the network

## Installation

Create a Hyperledger Fabric Network.\
Connect REST API with the network.

## REST API

## Get list of cars

### Request

`GET /cars/`

    curl -i -H 'Accept: application/json' http://localhost:8080/cars/

### Response

    HTTP/1.1 200 OK
    Date: Thu, 24 Feb 2011 12:36:30 GMT
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Content-Length: 2

    [
    {
        "id": "01",
        "brand": "Toyota",
        "owner": "Peter",
        "transfersCount": 0
    },
    {
        "id": "02",
        "brand": "Homnda",
        "owner": "Max",
        "transfersCount": 0
    }
    ]


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)