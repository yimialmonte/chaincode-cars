# Chaincode Cars

Provide a chaincode (smartcontract) for Hyperledger Fabric platform \
REST API for making transaction inside the network

## Installation

Create a Hyperledger Fabric Network.\
Connect REST API with the network.

## Run the API
`docker build . -t api` \
`docker run -it -p 8080:8080 api`

## Get list of cars

### Request

`GET /cars/`

    curl -i -H 'Accept: application/json' http://localhost:8080/cars

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Tue, 21 Sep 2021 15:41:52 GMT
    Content-Length: 127

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

`GET /cars/owner/{name}`

    curl -i -H 'Accept: application/json' http://localhost:8080/cars/owner/Max

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Tue, 21 Sep 2021 15:40:58 GMT
    Content-Length: 63

    [
    {
        "id": "02",
        "brand": "Homnda",
        "owner": "Max",
        "transfersCount": 0
    }
    ]

`POST /thing/`

    curl -i -H 'Accept: application/json' -d '{"id":"002","owner":"max"}' http://localhost:8080/cars

### Response

    HTTP/1.1 201 Created
    Date: Tue, 21 Sep 2021 15:40:02 GMT
    Content-Length: 0

    
## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)