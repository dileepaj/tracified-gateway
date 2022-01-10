package mintnft

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/nft/polygon/TracifiedNFT"
	"github.com/dileepaj/tracified-gateway/nft/polygon/polygonconnection"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
)

func MintNftWithPolygon() {
	println("-------inside mintNFT----------------------")
	//owener of the contract(issuer)
	nftIssuerPrivateKey:=commons.GoDotEnvVariable("NFTISSUESECRETKEYPOLYGON")
	privateKey, err := crypto.HexToECDSA(nftIssuerPrivateKey)
	reciverAddress := common.HexToAddress("0xb0f63fbB5c24dbAe40bc5cDF8E823226aC66E717")
    if err != nil {
        fmt.Println(err)
    }
	
	publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        logrus.Error("error casting public key to ECDSA")
    }

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
    nonce, err :=polygonconnection.GetPolygonClient().PendingNonceAt(context.Background(), fromAddress)
    if err != nil {
        logrus.Error(err)
    }

	gasPrice, err := polygonconnection.GetPolygonClient().SuggestGasPrice(context.Background())
    if err != nil {
        log.Fatal(err)
    }

	chainID, err := polygonconnection.GetPolygonClient().NetworkID(context.Background())
	if err != nil {
		log.Fatal("Error : ", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal("Error : ", err)
	}

    auth.Nonce = big.NewInt(int64(nonce))
    auth.Value = big.NewInt(0)     // in wei
    auth.GasLimit = uint64(300000) // in units
    auth.GasPrice = gasPrice

	client := polygonconnection.GetPolygonClient()
	contractAddress := common.HexToAddress("0x1b43258Ceb9e9cB40aa8c9cbc158205C22a6F45a")
    instance, err := TracifiedNFT.NewTracifiedNFT(contractAddress, client)
    if err != nil {
        log.Fatal(err)
    }

	//authenticate and write to smart contract
    tx, err := instance.MintNFT(auth,reciverAddress , "NFTURL")
    if err != nil {
        log.Fatal(err)
    }
	fmt.Printf("-------------------------------------------------trasactionHash sent: %s", tx.Hash().Hex())
}