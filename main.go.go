package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Holder struct {
	FullName string `json:"fullName"`
}

type NFT struct {
	ID    string `json:"id"`
	Owner string `json:"owner"`
}

func generateUniqueID() string {
	// Implement a unique ID generation logic
	return "unique-id"
}

func (s *SmartContract) CreateNFT(ctx contractapi.TransactionContextInterface, fullName string) (*NFT, error) {
	holderAsBytes, err := ctx.GetStub().GetState(fullName)
	if err != nil || holderAsBytes == nil {
		return nil, fmt.Errorf("holder %s does not exist", fullName)
	}

	nftID := generateUniqueID()
	nft := NFT{
		ID:    nftID,
		Owner: fullName,
	}

	nftAsBytes, _ := json.Marshal(nft)
	ctx.GetStub().PutState(nftID, nftAsBytes)

	return &nft, nil
}

func (s *SmartContract) TransferNFT(ctx contractapi.TransactionContextInterface, nftID, newOwnerFullName string) error {
	nftAsBytes, err := ctx.GetStub().GetState(nftID)
	if err != nil || nftAsBytes == nil {
		return fmt.Errorf("NFT %s does not exist", nftID)
	}

	holderAsBytes, err := ctx.GetStub().GetState(newOwnerFullName)
	if err != nil || holderAsBytes == nil {
		return fmt.Errorf("holder %s does not exist", newOwnerFullName)
	}

	var nft NFT
	json.Unmarshal(nftAsBytes, &nft)
	nft.Owner = newOwnerFullName

	nftAsBytes, _ = json.Marshal(nft)
	ctx.GetStub().PutState(nftID, nftAsBytes)

	return nil
}

func (s *SmartContract) ListNFTs(ctx contractapi.TransactionContextInterface, fullName string) ([]NFT, error) {
	queryString := fmt.Sprintf(`{"selector":{"owner":"%s"}}`, fullName)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var nfts []NFT
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var nft NFT
		json.Unmarshal(queryResponse.Value, &nft)
		nfts = append(nfts, nft)
	}

	return nfts, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error create nft chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting nft chaincode: %s", err.Error())
	}
}
