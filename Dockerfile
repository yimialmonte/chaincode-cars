FROM golang:latest

WORKDIR /app

COPY go.mod .
COPY go.sum . 

RUN go mod download 

COPY . .

EXPOSE 8080

RUN go build

CMD ["./chaincode-cars"]

