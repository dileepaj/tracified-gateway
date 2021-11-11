package commons

import (
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

func GetHorizonClient() *horizon.Client {
	client := GoDotEnvVariable("HORIZONCLIENT")
	if client == "public" {
		return horizon.DefaultPublicNetClient
	} else {
		return horizon.DefaultTestNetClient
	}
}

func GetHorizonNetwork() build.Network {
	client := GoDotEnvVariable("HORIZONCLIENT")
	if client == "public" {
		return build.PublicNetwork
	} else {
		return build.TestNetwork
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
