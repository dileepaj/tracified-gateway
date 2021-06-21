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
