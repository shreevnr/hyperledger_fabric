package main

import (
	"log"

	contracts "autoledger/contract"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	rcContract := new(contracts.RegistrationCertificateContract)
	rcprivatConract := new(contracts.PrivateAssetDetailsContract)

	chaincode, err := contractapi.NewChaincode(rcContract, rcprivatConract)

	if err != nil {
		log.Panicf("Could not create chaincode : %v", err)
	}

	err = chaincode.Start()

	if err != nil {
		log.Panicf("Failed to start chaincode : %v", err)
	}
}
