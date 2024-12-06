package contracts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type RegistrationCertificateContract struct {
	contractapi.Contract
}
type RegistrationCertificate struct {
	AssetType       string `json:"assetType"`
	RCId            string `json:"RCId"`
	Make            string `json:"make"`
	Model           string `json:"model"`
	Color           string `json:"color"`
	IsNOCapproved   bool   `json:"isNOCapproved"`
	RegisteredState string `json:"registeredState"`
	Status          string `json:"status"`
}
type HistoryQueryResult struct {
	RegCertificate *RegistrationCertificate `json:"record"`
	TxID           string                   `json:"txId"`
	Timestamp      string                   `json:"timestamp"`
	IsDelete       bool                     `json:"isDelete"`
}
type PaginatedQueryResult struct {
	Records             []*RegistrationCertificate `json:"records"`
	FetchedRecordsCount int32                      `json:"fetchedRecordsCount"`
	Bookmark            string                     `json:"bookmark"`
}

type EventData struct {
	Type  string
	State string
}

// CarExists returns true when asset with given ID exists in world state
func (r *RegistrationCertificateContract) RCExists(ctx contractapi.TransactionContextInterface, rcID string) (bool, error) {
	data, err := ctx.GetStub().GetState(rcID)

	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)

	}
	return data != nil, nil
}

// CreateRC creates RC entry in Ledger
func (r *RegistrationCertificateContract) CreateRC(ctx contractapi.TransactionContextInterface, rcID string, assetType string, make string, model string, color string, registeredState string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	if clientOrgID == registeredState+"MSP" {

		exists, err := r.RCExists(ctx, rcID)
		if err != nil {
			return "", fmt.Errorf("%s", err)
		} else if exists {
			return "", fmt.Errorf("the rc, %s already exists", rcID)
		}

		rc := RegistrationCertificate{
			AssetType:       assetType,
			RCId:            rcID,
			Color:           color,
			Make:            make,
			Model:           model,
			IsNOCapproved:   false,
			RegisteredState: registeredState,
			Status:          "Active",
		}

		bytes, _ := json.Marshal(rc)

		err = ctx.GetStub().PutState(rcID, bytes)
		fmt.Println("Create RC data ======= ", rc)

		if err != nil {
			return "", err
		} else {
			addRCEventData := EventData{
				Type:  "RC creation",
				State: registeredState,
			}
			eventDataByte, _ := json.Marshal(addRCEventData)
			ctx.GetStub().SetEvent("CreateRC", eventDataByte)

			return fmt.Sprintf("successfully added RC %v", rcID), nil
		}

	} else {
		return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgID)
	}
}

// Initiate RC Transfer from Onestate to another
func (r *RegistrationCertificateContract) InitiateTransferRC(ctx contractapi.TransactionContextInterface, rcId string, fromState string, toState string) error {
	// Authorization: Ensure only clients from the fromState can initiate transfer
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("error getting client MSP ID: %v", err)
	}
	if clientOrgID != fromState+"MSP" {
		return fmt.Errorf("only a client from %s can initiate transfer", fromState)
	}
	fmt.Println("Inside Initiate Transfer Function")
	// Fetch the RC from the ledger (on-chain)
	rcAsBytes, err := ctx.GetStub().GetState(rcId)
	if err != nil {
		return fmt.Errorf("failed to read RC from the ledger: %v", err)
	}
	if rcAsBytes == nil {
		return fmt.Errorf("RC with ID %s does not exist", rcId)
	}

	// Unmarshal the RC data
	var rc RegistrationCertificate
	err = json.Unmarshal(rcAsBytes, &rc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal RC data: %v", err)
	}

	// Update the status to "Transfer Initiated" and approve the NOC
	rc.Status = "Transfer Initiated"
	rc.IsNOCapproved = true

	// Marshal the updated RC back to bytes for saving to the ledger
	updatedRCAsBytes, err := json.Marshal(rc)
	if err != nil {
		return fmt.Errorf("failed to marshal updated RC: %v", err)
	}

	// Put the updated RC back to the ledger (on-chain)
	err = ctx.GetStub().PutState(rcId, updatedRCAsBytes)
	if err != nil {
		return fmt.Errorf("failed to update RC on ledger: %v", err)
	}

	return nil
}

// Approve transfer will be done by "toState"
func (r *RegistrationCertificateContract) ApproveTransferRC(ctx contractapi.TransactionContextInterface, rcId string, fromState string, toState string) error {
	// Authorization: Ensure only clients from the toState can approve the transfer
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("error getting client MSP ID: %v", err)
	}
	if clientOrgID != toState+"MSP" {
		return fmt.Errorf("only a client from %s can approve the transfer", toState)
	}

	// Fetch the RC from the ledger (on-chain)
	rcAsBytes, err := ctx.GetStub().GetState(rcId)
	if err != nil {
		return fmt.Errorf("failed to read RC from the ledger: %v", err)
	}
	if rcAsBytes == nil {
		return fmt.Errorf("RC with ID %s does not exist", rcId)
	}

	// Unmarshal the RC data
	var rc RegistrationCertificate
	err = json.Unmarshal(rcAsBytes, &rc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal RC data: %v", err)
	}

	// Ensure that NOC is approved before proceeding
	if !rc.IsNOCapproved {
		return fmt.Errorf("NOC has not been approved for RC ID %s", rcId)
	}

	// Update the status to "Transfer Approved"
	rc.Status = "Transfer Approved"

	// Marshal the updated RC data
	updatedRCAsBytes, err := json.Marshal(rc)
	if err != nil {
		return fmt.Errorf("failed to marshal updated RC: %v", err)
	}

	// Update the RC in the ledger (on-chain)
	err = ctx.GetStub().PutState(rcId, updatedRCAsBytes)
	if err != nil {
		return fmt.Errorf("failed to update RC on ledger: %v", err)
	}

	return nil
}

// Approve transfer will be done by "toState"
func (r *RegistrationCertificateContract) DeleteRC(ctx contractapi.TransactionContextInterface, rcId string, fromState string, toState string) error {
	// Authorization: Ensure only clients from the toState can approve the transfer
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("error getting client MSP ID: %v", err)
	}
	if clientOrgID != fromState+"MSP" {
		return fmt.Errorf("only a client from %s can do this operation", toState)
	}

	// Fetch the RC from the ledger (on-chain)
	rcAsBytes, err := ctx.GetStub().GetState(rcId)
	if err != nil {
		return fmt.Errorf("failed to read RC from the ledger: %v", err)
	}
	if rcAsBytes == nil {
		return fmt.Errorf("RC with ID %s does not exist", rcId)
	}

	// Unmarshal the RC data
	var rc RegistrationCertificate
	err = json.Unmarshal(rcAsBytes, &rc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal RC data: %v", err)
	}

	fmt.Printf("Updating the LedgerStatus")

	// Update the RC status on the ledger to "Transferred to TransferRCCollection"
	rc.Status = "Dis-owned by " + fromState

	fmt.Printf("Updating the TransferRCCollection")

	// Marshal the updated RC data
	updatedRCAsBytes, err := json.Marshal(rc)
	if err != nil {
		return fmt.Errorf("failed to marshal updated RC: %v", err)
	}

	// Update the RC in the ledger (on-chain)
	err = ctx.GetStub().PutState(rcId, updatedRCAsBytes)
	if err != nil {
		return fmt.Errorf("failed to update RC on ledger: %v", err)
	}

	return nil

}

func (r *RegistrationCertificateContract) GetAssetsbyState(ctx contractapi.TransactionContextInterface, state string) ([]*RegistrationCertificate, error) {
	// Define the range query
	queryString := fmt.Sprintf(`{"selector":{"registeredState":"%s"}}`, state)

	// Execute the query
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer resultsIterator.Close()

	// Process the results
	var results []*RegistrationCertificate
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next result: %v", err)
		}

		var rc RegistrationCertificate
		err = json.Unmarshal(queryResponse.Value, &rc)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal RC data: %v", err)
		}
		results = append(results, &rc)
	}

	return results, nil
}

// GetAssetDetails query the asset details from ledger
func (r *RegistrationCertificateContract) GetAssetDetails(ctx contractapi.TransactionContextInterface, rcId string) (*RegistrationCertificate, error) {
	// Fetch the RC from the ledger
	rcAsBytes, err := ctx.GetStub().GetState(rcId)
	if err != nil {
		return nil, fmt.Errorf("failed to read RC from the ledger: %v", err)
	}
	if rcAsBytes == nil {
		return nil, fmt.Errorf("RC with ID %s does not exist", rcId)
	}

	// Unmarshal the RC data
	var rc RegistrationCertificate
	err = json.Unmarshal(rcAsBytes, &rc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal RC data: %v", err)
	}

	return &rc, nil
}

// GetAssetHistory will fetch the history records of the given asset/RC
func (r *RegistrationCertificateContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, rcId string) ([]*HistoryQueryResult, error) {
	// Fetch the history of the RC
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(rcId)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for RC with ID %s: %v", rcId, err)
	}
	defer resultsIterator.Close()

	// Process the results
	var records []*HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next result from history: %v", err)
		}

		var rc RegistrationCertificate
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &rc)
			if err != nil {
				return nil, err
			}
		} else {
			rc = RegistrationCertificate{
				RCId: rcId,
			}
		}
		timestamp := response.Timestamp.AsTime()

		formattedTime := timestamp.Format(time.RFC1123)

		record := HistoryQueryResult{
			TxID:           response.TxId,
			Timestamp:      formattedTime,
			RegCertificate: &rc,
			IsDelete:       response.IsDelete,
		}
		records = append(records, &record)
	}

	return records, nil
}

// GetAssetsByRange will query the assets with given range
func (r *RegistrationCertificateContract) GetAssetsByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*RegistrationCertificate, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the  data by range. %s", err)
	}
	defer resultsIterator.Close()

	return AssetIteratorFunction(resultsIterator)
}
func (r *RegistrationCertificateContract) GetAllRCs(ctx contractapi.TransactionContextInterface) ([]*RegistrationCertificate, error) {

	queryString := `{"selector":{"assetType":"RegistrationCertificate"}, "sort":[{ "registeredState": "desc"}]}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the query result. %s", err)
	}
	defer resultsIterator.Close()
	return AssetIteratorFunction(resultsIterator)
}

func AssetIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*RegistrationCertificate, error) {
	var rcs []*RegistrationCertificate
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not fetch the details of the result iterator. %s", err)
		}
		var rc RegistrationCertificate
		err = json.Unmarshal(queryResult.Value, &rc)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal the data. %s", err)
		}
		rcs = append(rcs, &rc)
	}

	return rcs, nil
}
func (c *RegistrationCertificateContract) GetAssetsWithPagination(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {

	queryString := `{"selector":{"assetType":"RegistrationCertificate"}}`

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("could not get the rc records. %s", err)
	}
	defer resultsIterator.Close()

	assets, err := AssetIteratorFunction(resultsIterator)
	if err != nil {
		return nil, fmt.Errorf("could not return the rc records %s", err)
	}

	return &PaginatedQueryResult{
		Records:             assets,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}
