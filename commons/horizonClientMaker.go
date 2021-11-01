package commons

import "github.com/stellar/go/clients/horizon"

func GetHorizonClient() *horizon.Client {
	client := GoDotEnvVariable("HORIZONCLIENT")
	if client == "public" {
		return horizon.DefaultPublicNetClient
	} else {
		return horizon.DefaultTestNetClient
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
