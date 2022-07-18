package commons

import (
	"github.com/stellar/go/clients/horizonclient"
)

func GetHorizonClient() *horizonclient.Client {
	client := GoDotEnvVariable("HORIZONCLIENT")
	if client == "public" {
		return horizonclient.DefaultPublicNetClient
	} else {
		return horizonclient.DefaultTestNetClient
	}
}

func GetHorizonNetwork() *horizonclient.Client {
	client := GoDotEnvVariable("HORIZONCLIENT")
	if client == "public" {
		return horizonclient.DefaultPublicNetClient
	} else {
		return horizonclient.DefaultTestNetClient
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
