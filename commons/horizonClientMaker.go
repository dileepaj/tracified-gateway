package commons

import (
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
)

func GetHorizonClient() *horizonclient.Client {
	client := GoDotEnvVariable("HORIZONCLIENT")
	if client == "public" {
		return horizonclient.DefaultPublicNetClient
	} else {
		return horizonclient.DefaultTestNetClient
	}
}

func GetStellarNetwork() string {
	client := GoDotEnvVariable("HORIZONCLIENT")
	if client == "public" {
		return network.PublicNetworkPassphrase
	} else {
		return network.TestNetworkPassphrase
	}
}

//get the connected stellar networ(public or test)

func GetHorizonClientNetworkName() string {
	clientNetworkName := GoDotEnvVariable("HORIZONCLIENT")
	return clientNetworkName
}

func GetStellarLaboratoryClient() string {
	clientUrl := GoDotEnvVariable("STELLARLABORATORYCLIENT")
	return clientUrl
}
