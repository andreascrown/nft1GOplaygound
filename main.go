package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Holder represents a physical entity uniquely identified by a full name
type Holder struct {
	FullName string `json:"fullName"`
}

// NFT represents a non-fungible token with a unique 10-symbol alphanumeric sequence
type NFT struct {
	ID    string `json:"id"`
	Owner string `json:"owner"`
}

// Database to simulate ledger
var holders = make(map[string]Holder)
var nfts = make(map[string]NFT)

// CreateNFT assigns a new NFT to a holder
func CreateNFT(fullName string) (*NFT, error) {
	holder, exists := holders[fullName]
	if !exists {
		return nil, fmt.Errorf("holder %s does not exist", fullName)
	}

	nftID, err := generateUniqueID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate NFT ID: %v", err)
	}

	nft := NFT{
		ID:    nftID,
		Owner: holder.FullName,
	}

	nfts[nftID] = nft

	return &nft, nil
}

// TransferNFT transfers an NFT from one holder to another
func TransferNFT(nftID, newOwnerFullName string) error {
	nft, exists := nfts[nftID]
	if !exists {
		return fmt.Errorf("NFT %s does not exist", nftID)
	}

	_, holderExists := holders[newOwnerFullName]
	if !holderExists {
		return fmt.Errorf("holder %s does not exist", newOwnerFullName)
	}

	nft.Owner = newOwnerFullName
	nfts[nftID] = nft
	return nil
}

// ListNFTs lists all NFTs owned by a holder
func ListNFTs(fullName string) ([]NFT, error) {
	var holderNFTs []NFT
	for _, nft := range nfts {
		if nft.Owner == fullName {
			holderNFTs = append(holderNFTs, nft)
		}
	}
	return holderNFTs, nil
}

// generateUniqueID generates a unique 10-symbol alphanumeric ID
func generateUniqueID() (string, error) {
	bytes := make([]byte, 5)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func main() {
	// Initialize holders
	holders["Andreas Koronias"] = Holder{FullName: "Andreas Koronias"}
	holders["Elias Iosif"] = Holder{FullName: "Elias Iosif"}
	holders["George Georgiou"] = Holder{FullName: "George Georgiou"} // New holder added

	// Create NFTs
	nft1, err := CreateNFT("Andreas Koronias")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		nft1JSON, _ := json.Marshal(nft1)
		fmt.Println("NFT Created:", string(nft1JSON))
	}

	nft2, err := CreateNFT("Elias Iosif")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		nft2JSON, _ := json.Marshal(nft2)
		fmt.Println("NFT Created:", string(nft2JSON))
	}

	// Transfer an NFT from Elias Iosif to George Georgiou
	err = TransferNFT(nft2.ID, "George Georgiou")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("NFT Transferred to George Georgiou")
	}

	// List NFTs for George Georgiou
	georgeNFTs, err := ListNFTs("George Georgiou")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("George Georgiou's NFTs:")
		for _, nft := range georgeNFTs {
			nftJSON, _ := json.Marshal(nft)
			fmt.Println(string(nftJSON))
		}
	}
}
