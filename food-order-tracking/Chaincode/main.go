package main

import (
	"foodtraker/orderTracker"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	chaincode, err := contractapi.NewChaincode(new(orderTracker.SmartContract))
	if err != nil {
		log.Panicf("Error creating food-order-tracking chaincode: %v", err)
	}
	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting food-order-tracking chaincode: %v", err)
	}
}
