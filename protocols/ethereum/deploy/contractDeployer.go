package deploy

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/dileepaj/tracified-gateway/commons"
	build "github.com/dileepaj/tracified-gateway/protocols/ethereum/contracts/build"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

/*
	Deploy smart contract in Ethereum
*/
func DeployContract() (string, string, error) {
	contractAddress := ""
	transactionHash := ""

	//Dial infura client
	client, errWhenDialingEthClinet := ethclient.Dial(commons.GoDotEnvVariable("SPOLIALINK"))
	if errWhenDialingEthClinet != nil {
		logrus.Error("Error when dialing the eth client " + errWhenDialingEthClinet.Error())
		return contractAddress, transactionHash, errors.New("Error when dialing eth client , ERROR : " + errWhenDialingEthClinet.Error())
	}

	//load ECDSA private key
	privateKey, errWhenGettingECDSAKey := crypto.HexToECDSA(commons.GoDotEnvVariable("ETHEREUMSECKEY"))
	if errWhenGettingECDSAKey != nil {
		logrus.Error("Error when getting ECDSA key " + errWhenGettingECDSAKey.Error())
		return contractAddress, transactionHash, errors.New("Error when getting ECDSA key , ERROR : " + errWhenGettingECDSAKey.Error())
	}

	//get the public key
	publicKey := privateKey.Public()
	//get public key ECDSA
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		logrus.Error("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return contractAddress, transactionHash, errors.New("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	//get nonce
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, errWhenGettingNonce := client.PendingNonceAt(context.Background(), fromAddress)
	if errWhenGettingNonce != nil {
		logrus.Error("Error when getting nonce " + errWhenGettingNonce.Error())
		return contractAddress, transactionHash, errors.New("Error when getting nonce , ERROR : " + errWhenGettingNonce.Error())
	}

	//get the gas price
	gasPrice, errWhenGettingGasPrice := client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if errWhenGettingGasPrice != nil {
		logrus.Error("Error when getting gas price " + errWhenGettingGasPrice.Error())
		return contractAddress, transactionHash, errors.New("Error when getting gas price , ERROR : " + errWhenGettingGasPrice.Error())
	}

	//create the keyed transactor
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(1000000) // in units
	auth.GasPrice = big.NewInt(int64(gasPrice))

	address, hash, instance, errInDeploy := build.DeployCalculations(auth, client)
	if errInDeploy != nil {
		logrus.Error("Error when deploying the contract " + errInDeploy.Error())
		return contractAddress, transactionHash, errors.New("Error when deploying the contract , ERROR : " + errInDeploy.Error())
	}

	logrus.Info("View contract at : https://sepolia.etherscan.io/address/", address)
	logrus.Info("View transaction at : https://sepolia.etherscan.io/tx/", hash.Hash().Hex())

	contractAddress = address.Hex()
	transactionHash = hash.Hash().Hex()
	_ = instance

	return contractAddress, transactionHash, nil
}
