package polygonconnection

import (
	"fmt"
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/ethereum/go-ethereum/ethclient"
)

var testNetRPC =commons.GoDotEnvVariable("POLYGONTESTNETRPC")
var mainNetRPC=commons.GoDotEnvVariable("POLYGONMAINNETRPC")

func GetPolygonClient() *ethclient.Client {
	horizonClient := commons.GoDotEnvVariable("HORIZONCLIENT")

	if horizonClient == "public" {
		Client, err := ethclient.Dial(mainNetRPC)
		if err != nil {
			log.Fatal("Error while connecting to polygon network " + err.Error())
		}
		fmt.Println("connection build with polygon mainnet")
		defer Client.Close()
		return Client
	} else {
		Client, err := ethclient.Dial(testNetRPC)
		if err != nil {
			log.Fatal("Error while connecting to polygon network " + err.Error())
		}
		fmt.Println("connection build with polygon testnet")
		defer Client.Close()
		return Client
	}
}
