package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism accross languages
// golang keeps the order when marshal to json but doesn't order automatically
type Vote struct {
	Candidate      string `json:"Candidate"`
	CPR            string `json:"CPR"`
	ID             string `json:"ID"`
	Name           string `json:"Name"`
	PoliticalParty string `json:"PoliticalParty"`
}

// InitLedger adds a base set of Vote to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Vote{
		{ID: "asset1", Candidate: "blue", CPR: "202000-1821", Name: "Tomoko", PoliticalParty: "Venstre"},
		{ID: "asset2", Candidate: "red", CPR: "202000-1821", Name: "Brad", PoliticalParty: "Socialdemokratiet"},
		{ID: "asset3", Candidate: "green", CPR: "202000-1821", Name: "Jin Soo", PoliticalParty: "Liberal alliance"},
		{ID: "asset4", Candidate: "yellow", CPR: "202000-1821", Name: "Max", PoliticalParty: "SF"},
		{ID: "asset5", Candidate: "black", CPR: "202000-1821", Name: "Adriana", PoliticalParty: "Dansk folkeparti"},
		{ID: "asset6", Candidate: "white", CPR: "202000-1821", Name: "Michel", PoliticalParty: "Enhedslisten"},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateVote(ctx contractapi.TransactionContextInterface, id string, candidate string, cpr string, name string, politicalParty string) error {
	exists, err := s.VoteExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Vote{
		ID:             id,
		Candidate:      candidate,
		CPR:            cpr,
		Name:           name,
		PoliticalParty: politicalParty,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// AssetExists returns true when asset with given CPR exists in world state
func (s *SmartContract) VoteExists(ctx contractapi.TransactionContextInterface, cpr string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(cpr)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllVotes(ctx contractapi.TransactionContextInterface) ([]*Vote, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Vote
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Vote
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	}
}
