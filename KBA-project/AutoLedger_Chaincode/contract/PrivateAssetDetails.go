package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type PrivateAssetDetailsContract struct {
	contractapi.Contract
}
type RCPrivateData struct {
	RCId            string `json:"rcId"`
	AssetType       string `json:"assetType"`
	Make            string `json:"make"`
	Model           string `json:"model"`
	Color           string `json:"color"`
	OwnerName       string `json:"ownerName"`
	OwnerAadhar     string `json:"ownerAadhar"`
	EngineNumber    string `json:"engineNumber"`
	InsuranceCert   string `json:"insuranceCert"`
	PollutionCert   string `json:"pollutionCert"`
	Status          string `json:"status"`
	RegisteredState string `json:"registeredState"`
}
type TransferRCData struct {
	RCId            string `json:"rcId"`
	AssetType       string `json:"assetType"`
	Make            string `json:"make"`
	Model           string `json:"model"`
	Color           string `json:"color"`
	OwnerName       string `json:"ownerName"`
	OwnerAadhar     string `json:"ownerAadhar"`
	EngineNumber    string `json:"engineNumber"`
	InsuranceCert   string `json:"insuranceCert"`
	PollutionCert   string `json:"pollutionCert"`
	Status          string `json:"status"`
	RegisteredState string `json:"registeredState"`
	TransferredFrom string `json:"transferredFrom"`
	TransferredTo   string `json:"transferredTo"`
}

// CreatePrivateRC function will add the Registration Certificate (RC)details on the state's PDC
func (p *PrivateAssetDetailsContract) CreatePrivateRC(ctx contractapi.TransactionContextInterface, state string) error {
	fmt.Println("Inside CreateRC details")
	// Authorization check: Ensure the client is from the correct organization
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("error getting client MSP ID: %v", err)
	}

	// Retrieve transient data map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("error getting transient data: %v", err)
	}

	// Retrieve the asset properties (including both on-chain and private data) from transient field
	transientRCJSON, ok := transientMap["rc_properties"]
	if !ok {
		return fmt.Errorf("RC properties not found in the transient map input")
	}

	// Define a struct for transient data (both on-chain and private information)
	type rcTransientInput struct {
		AssetType         string `json:"assetType"`
		RCId              string `json:"rcId"`
		Make              string `json:"make"`
		Model             string `json:"model"`
		Color             string `json:"color"`
		OwnerName         string `json:"ownerName"`
		OwnerAadhar       string `json:"ownerAadhar"`
		EngineNumber      string `json:"engineNumber"`
		InsuranceCert     string `json:"insuranceCert"`
		PollutionCert     string `json:"pollutionCert"`
		RegistereredState string `json:"registeredState"`
	}

	// Unmarshal the transient JSON into the struct
	var rcInput rcTransientInput
	err = json.Unmarshal(transientRCJSON, &rcInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal RC data JSON: %v", err)
	}

	// Validate the input fields (on-chain and private data)
	if len(rcInput.RCId) == 0 {
		return fmt.Errorf("rCId field must be a non-empty string")
	}
	if len(rcInput.AssetType) == 0 {
		return fmt.Errorf("rCId field must be a non-empty string")
	}
	if len(rcInput.Make) == 0 {
		return fmt.Errorf("make field must be a non-empty string")
	}
	if len(rcInput.Model) == 0 {
		return fmt.Errorf("code field must be a non-empty string")
	}
	if len(rcInput.Color) == 0 {
		return fmt.Errorf("color field must be a non-empty string")
	}
	if len(rcInput.OwnerName) == 0 {
		return fmt.Errorf("ownerName field must be a non-empty string")
	}
	if len(rcInput.OwnerAadhar) == 0 {
		return fmt.Errorf("ownerAadhar field must be a non-empty string")
	}
	if len(rcInput.EngineNumber) == 0 {
		return fmt.Errorf("engineNumber field must be a non-empty string")
	}
	if len(rcInput.InsuranceCert) == 0 {
		return fmt.Errorf("insuranceCert field must be a non-empty string")
	}
	if len(rcInput.PollutionCert) == 0 {
		return fmt.Errorf("pollutionCert field must be a non-empty string")
	}
	if len(rcInput.RegistereredState) == 0 {
		return fmt.Errorf("RegisteredState field must be a non-empty string")
	}

	// Check if the RC already exists
	rcAsBytes, err := ctx.GetStub().GetState(rcInput.RCId)
	if err != nil {
		return fmt.Errorf("failed to get RC: %v", err)
	}
	if rcAsBytes != nil {
		return fmt.Errorf("rc with ID %s already exists", rcInput.RCId)
	}
	// Create a private data (off-chain) collection for sensitive information
	rcPrivate := RCPrivateData{
		RCId:            rcInput.RCId,
		AssetType:       rcInput.AssetType,
		Make:            rcInput.Make,
		Model:           rcInput.Model,
		Color:           rcInput.Color,
		OwnerName:       rcInput.OwnerName,
		OwnerAadhar:     rcInput.OwnerAadhar,
		EngineNumber:    rcInput.EngineNumber,
		InsuranceCert:   rcInput.InsuranceCert,
		PollutionCert:   rcInput.PollutionCert,
		Status:          "Active",
		RegisteredState: rcInput.RegistereredState,
	}

	// Marshal the private RC data to bytes for storing in Private Data Collection (PDC)
	rcPrivateJSON, err := json.Marshal(rcPrivate)
	if err != nil {
		return fmt.Errorf("error marshalling private RC data: %v", err)
	}
	// Authorization check: Ensure the client is from the correct organization
	if clientOrgID == state+"MSP" {
		// Store the private data in the respective collection
		collectionName := state + "PDC"
		fmt.Println("updating PDC")
		err = ctx.GetStub().PutPrivateData(collectionName, rcInput.RCId, rcPrivateJSON)
		if err != nil {
			return fmt.Errorf("error storing private RC data: %v", err)
		}
	} else {
		return fmt.Errorf("user under MSP ID: %v can't perform this action", clientOrgID)
	}

	return nil
}

//Add RC to TransferRCCollection will move the RC entry from "fromState" private data collection to "TransferRCCollection"

func (p *PrivateAssetDetailsContract) MovetoTransferRCCollection(ctx contractapi.TransactionContextInterface, rcId string, fromState string, toState string) error {
	// Authorization: Ensure only clients from the fromState can delete the RC
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("error getting client MSP ID: %v", err)
	}
	if clientOrgID != fromState+"MSP" {
		return fmt.Errorf("only a client from %s can delete the RC", fromState)
	}

	// // Fetch the RC from the ledger (on-chain)

	// if err != nil {
	// 	return fmt.Errorf("failed to read RC from the ledger: %v", err)
	// }
	// fmt.Println("Fetched RC details from ledger", rc)
	// // Check if the RC status is "Transfer Approved"
	// if rc.Status != "Transfer Approved" {
	// 	return fmt.Errorf("RC with ID %s is not ready for transfer, current status: %s", rcId, rc.Status)
	// }

	// Fetch the RC private data from the fromState PDC
	collectionName := fromState + "PDC"
	privateRC, err := p.GetAssetPrivateDetails(ctx, rcId, fromState)
	if err != nil {
		return fmt.Errorf("failed to get private RC data from %sPDC: %v", fromState, err)
	}
	fmt.Println("Fetched RC details from PDC", privateRC)

	// Add the Transfer approved private RC data to the TransferRCCollection, including TransferredFrom and TransferredTo fields
	transferRCCollection := "TransferRCCollection"

	// Fill the transfer data
	transferRC := TransferRCData{
		RCId:            privateRC.RCId,
		AssetType:       privateRC.AssetType,
		Make:            privateRC.Make,
		Model:           privateRC.Model,
		Color:           privateRC.Color,
		OwnerName:       privateRC.OwnerName,
		OwnerAadhar:     privateRC.OwnerAadhar,
		EngineNumber:    privateRC.EngineNumber,
		PollutionCert:   privateRC.PollutionCert,
		InsuranceCert:   privateRC.InsuranceCert,
		RegisteredState: privateRC.RegisteredState,
		Status:          "Dis-owned by " + fromState,
		TransferredFrom: fromState,
		TransferredTo:   toState,
	}
	fmt.Println("TransferRC details before adding in TransferRC Collection", transferRC)
	// Marshal the transfer RC data to bytes
	transferRCAsBytes, err := json.Marshal(transferRC)
	if err != nil {
		return fmt.Errorf("failed to marshal transfer RC: %v", err)
	}

	// Store the updated RC data in the TransferRCCollection
	err = ctx.GetStub().PutPrivateData(transferRCCollection, rcId, transferRCAsBytes)
	if err != nil {
		return fmt.Errorf("failed to store RC in %s: %v", transferRCCollection, err)
	}

	// Delete the RC from the private data collection of the fromState
	err = ctx.GetStub().DelPrivateData(collectionName, rcId)
	if err != nil {
		return fmt.Errorf("failed to delete RC from %sPDC: %v", fromState, err)
	}

	return nil
}
func (p *PrivateAssetDetailsContract) AddTransferredRC(ctx contractapi.TransactionContextInterface, rcID string, state string) error {
	// Authorization: Ensure client belongs to the state querying the private data
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("error getting client MSP ID: %v", err)
	}

	// Retrieve transient data map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("error getting transient data: %v", err)
	}

	// Retrieve the asset properties (including both on-chain and private data) from transient field
	transientRCJSON, ok := transientMap["add_rc_properties"]
	if !ok {
		return fmt.Errorf("RC properties not found in the transient map input")
	}

	// Define a struct for transient data (both on-chain and private information)
	type addrcTransientInput struct {
		AssetType     string `json:"assetType"`
		RCId          string `json:"rcId"`
		Make          string `json:"make"`
		Model         string `json:"model"`
		Color         string `json:"color"`
		OwnerName     string `json:"ownerName"`
		OwnerAadhar   string `json:"ownerAadhar"`
		EngineNumber  string `json:"engineNumber"`
		InsuranceCert string `json:"insuranceCert"`
		PollutionCert string `json:"pollutionCert"`
		//RegistereredState string `json:"registeredState"`
	}

	// Unmarshal the transient JSON into the struct
	var addrcInput addrcTransientInput
	err = json.Unmarshal(transientRCJSON, &addrcInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal RC data JSON: %v", err)
	}

	// Validate the input fields (on-chain and private data)
	if len(addrcInput.RCId) == 0 {
		return fmt.Errorf("rCId field must be a non-empty string")
	}
	if len(addrcInput.AssetType) == 0 {
		return fmt.Errorf("rCId field must be a non-empty string")
	}
	if len(addrcInput.Make) == 0 {
		return fmt.Errorf("make field must be a non-empty string")
	}
	if len(addrcInput.Model) == 0 {
		return fmt.Errorf("code field must be a non-empty string")
	}
	if len(addrcInput.Color) == 0 {
		return fmt.Errorf("color field must be a non-empty string")
	}
	if len(addrcInput.OwnerName) == 0 {
		return fmt.Errorf("ownerName field must be a non-empty string")
	}
	if len(addrcInput.OwnerAadhar) == 0 {
		return fmt.Errorf("ownerAadhar field must be a non-empty string")
	}
	if len(addrcInput.EngineNumber) == 0 {
		return fmt.Errorf("engineNumber field must be a non-empty string")
	}
	if len(addrcInput.InsuranceCert) == 0 {
		return fmt.Errorf("insuranceCert field must be a non-empty string")
	}
	if len(addrcInput.PollutionCert) == 0 {
		return fmt.Errorf("pollutionCert field must be a non-empty string")
	}
	// if len(addrcInput.RegistereredState) == 0 {
	// 	return fmt.Errorf("RegisteredState field must be a non-empty string")
	// }

	// Check if the RC already exists
	rcAsBytes, err := ctx.GetStub().GetState(rcID)
	if err != nil {
		return fmt.Errorf("failed to get RC: %v", err)
	}
	if rcAsBytes != nil {
		return fmt.Errorf("rc with ID %s already exists", rcID)
	}
	// Create a private data (off-chain) collection for sensitive information
	rcPrivate := RCPrivateData{
		RCId:            rcID,
		AssetType:       addrcInput.AssetType,
		Make:            addrcInput.Make,
		Model:           addrcInput.Model,
		Color:           addrcInput.Color,
		OwnerName:       addrcInput.OwnerName,
		OwnerAadhar:     addrcInput.OwnerAadhar,
		EngineNumber:    addrcInput.EngineNumber,
		InsuranceCert:   addrcInput.InsuranceCert,
		PollutionCert:   addrcInput.PollutionCert,
		Status:          "Active",
		RegisteredState: state,
	}

	fmt.Println("RC details before adding to approving state PDC", rcPrivate)
	// Marshal the private RC data to bytes for storing in Private Data Collection (PDC)
	rcPrivateJSON, err := json.Marshal(rcPrivate)
	if err != nil {
		return fmt.Errorf("error marshalling private RC data: %v", err)
	}
	// Authorization check: Ensure the client is from the correct organization
	if clientOrgID == state+"MSP" {
		// Store the private data in the respective collection
		collectionName := state + "PDC"
		fmt.Println("Before Adding the transferred rc to PDC")
		err = ctx.GetStub().PutPrivateData(collectionName, addrcInput.RCId, rcPrivateJSON)
		if err != nil {
			return fmt.Errorf("error storing private RC data: %v", err)
		}
	} else {
		return fmt.Errorf("user under MSP ID: %v can't perform this action", clientOrgID)
	}

	return nil
}

// GetAssetPrivateDetails query the private data from the give state's PDC
func (p *PrivateAssetDetailsContract) GetAssetPrivateDetails(ctx contractapi.TransactionContextInterface, rcId string, state string) (*RCPrivateData, error) {
	// Authorization: Ensure client belongs to the state querying the private data
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("error getting client MSP ID: %v", err)
	}
	if clientOrgID != state+"MSP" {
		return nil, fmt.Errorf("only clients from %s can query private data", state)
	}

	// Fetch the RC private data from the state PDC
	collectionName := state + "PDC"
	rcPrivateAsBytes, err := ctx.GetStub().GetPrivateData(collectionName, rcId)
	if err != nil {
		return nil, fmt.Errorf("failed to read private RC data from %sPDC: %v", state, err)
	}
	if rcPrivateAsBytes == nil {
		return nil, fmt.Errorf("private RC data with ID %s does not exist in %sPDC", rcId, state)
	}

	// Unmarshal the private RC data
	var rcPrivate RCPrivateData
	err = json.Unmarshal(rcPrivateAsBytes, &rcPrivate)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal private RC data: %v", err)
	}

	return &rcPrivate, nil
}

// GetAssetPrivateDetails query the private data from the give state's PDC
func (p *PrivateAssetDetailsContract) GetTransferredRCDetails(ctx contractapi.TransactionContextInterface, rcId string) (*TransferRCData, error) {
	// Authorization: Ensure client belongs to the state querying the private data
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("error getting client MSP ID: %v", err)
	}
	if clientOrgID == "tnrtoMSP" || clientOrgID == "klrtoMSP" || clientOrgID == "knrtoMSP" {
		return nil, fmt.Errorf("only clients from %s can query private data", clientOrgID)
	}

	// Fetch the RC private data from the state PDC
	collectionName := "TransferRCCollection"
	transferrcAsBytes, err := ctx.GetStub().GetPrivateData(collectionName, rcId)
	if err != nil {
		return nil, fmt.Errorf("failed to read private RC data from %sPDC: %v", collectionName, err)
	}
	if transferrcAsBytes == nil {
		return nil, fmt.Errorf("private RC data with ID %s does not exist in %sPDC", rcId, collectionName)
	}

	// Unmarshal the private RC data
	var transferrc TransferRCData
	err = json.Unmarshal(transferrcAsBytes, &transferrc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal private RC data: %v", err)
	}

	return &transferrc, nil
}
func (p *PrivateAssetDetailsContract) GetAssetCollection(ctx contractapi.TransactionContextInterface, state string) ([]*RCPrivateData, error) {
	fmt.Println("InsideGetAssetCollection")
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("error getting client MSP ID: %v", err)
	}
	if clientOrgID != state+"MSP" {
		return nil, fmt.Errorf("only clients from %s can query private data", state)
	}
	fmt.Println("myclient", clientOrgID)
	collectionName := state + "PDC"
	queryString := `{"selector":{"assetType":"RegistrationCertificate"}}`
	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(collectionName, queryString)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the query result. %s", err)
	}
	defer resultsIterator.Close()
	return AssetPrivateIteratorFunction(resultsIterator)
}
func AssetPrivateIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*RCPrivateData, error) {
	fmt.Println("insidePrivateIteratorfunction", resultsIterator)
	var rcs []*RCPrivateData
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not fetch the details of the result iterator. %s", err)
		}
		var rc RCPrivateData
		err = json.Unmarshal(queryResult.Value, &rc)
		fmt.Println("Assets", &rc)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal the data. %s", err)
		}
		rcs = append(rcs, &rc)
	}

	return rcs, nil
}
func (p *PrivateAssetDetailsContract) GetTransferredRCCollection(ctx contractapi.TransactionContextInterface, state string) ([]*TransferRCData, error) {
	fmt.Println("TransferredRCCollection")
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("error getting client MSP ID: %v", err)
	}
	if clientOrgID != state+"MSP" {
		return nil, fmt.Errorf("only clients from %s can query private data", state)
	}
	//fmt.Println("myclient", clientOrgID)
	collectionName := "TransferRCCollection"
	queryString := `{"selector":{"assetType":"RegistrationCertificate"}}`
	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(collectionName, queryString)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the query result. %s", err)
	}
	defer resultsIterator.Close()
	return TransferRCIteratorFunction(resultsIterator)
}
func TransferRCIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*TransferRCData, error) {
	fmt.Println("insideTransferRCIteratorfunction", resultsIterator)
	var rcs []*TransferRCData
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not fetch the details of the result iterator. %s", err)
		}
		var rc TransferRCData
		err = json.Unmarshal(queryResult.Value, &rc)
		fmt.Println("Assets", &rc)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal the data. %s", err)
		}
		rcs = append(rcs, &rc)
	}

	return rcs, nil
}
