package generalservices

import (
	"crypto/ecdsa"
	"errors"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

//Load the client and private and public keys
//1- Ethereum, 2-Polygon
func LoadClientAndKey(blockchainType int) (*ethclient.Client, *ecdsa.PrivateKey, common.Address, error) {
	var client *ethclient.Client
	var privateKey *ecdsa.PrivateKey
	var fromAddress common.Address
	var errWhenDialingEthClient, errWhenGettingECDSAKey error

	logrus.Info("Loading the Client and Keys")

	if blockchainType == 1 {
		//dial infura client
		client, errWhenDialingEthClient = ethclient.Dial(commons.GoDotEnvVariable("ETHEREUMTESTNETLINK"))
		if errWhenDialingEthClient != nil {
			logrus.Error("Error when dialing the Ethereum client : " + errWhenDialingEthClient.Error())
			return client, privateKey, fromAddress, errWhenDialingEthClient
		}
	} else if blockchainType == 2 {
		//dial alchemy client
		client, errWhenDialingEthClient = ethclient.Dial(commons.GoDotEnvVariable("POLYGONALCHEMYAPILINK") + commons.GoDotEnvVariable("POLYGONALCHEMYAPIKEY"))
		if errWhenDialingEthClient != nil {
			logrus.Error("Error when dialing the Polygon client : " + errWhenDialingEthClient.Error())
			return client, privateKey, fromAddress, errWhenDialingEthClient
		}
	}

	//load ECDSA private key
	privateKey, errWhenGettingECDSAKey = crypto.HexToECDSA(commons.GoDotEnvVariable("ETHEREUMSECKEY"))
	if errWhenGettingECDSAKey != nil {
		logrus.Error("Error when getting ECDSA key " + errWhenGettingECDSAKey.Error())
		return client, privateKey, fromAddress, errWhenGettingECDSAKey
	}

	//get the public key
	publicKey := privateKey.Public()
	//get public key ECDSA
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		logrus.Error("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return client, privateKey, fromAddress, errors.New("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fA := crypto.PubkeyToAddress(*publicKeyECDSA)

	fromAddress = common.Address(fA)

	return client, privateKey, fromAddress, nil
}
