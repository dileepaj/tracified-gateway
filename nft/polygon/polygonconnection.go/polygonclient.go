package polygonconnection

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

var testNetPRC = "https://polygon-mumbai.g.alchemy.com/v2/x6utaPxydF90g72SIUfhhn8T1UoW7eIy"

func CallPolygonClient() *ethclient.Client {

	Client, err := ethclient.Dial(testNetPRC)
	if err != nil {
		log.Fatal("Error when connecting to polygon network " + err.Error())
	}

	fmt.Println("we have a connection")
	defer Client.Close()
	return Client
}
