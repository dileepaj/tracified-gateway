package commons

import (
	rp "github.com/gagliardetto/solana-go/rpc"
	"github.com/portto/solana-go-sdk/rpc"
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

func GetSolanaNetwork() string {
	client := GoDotEnvVariable("SOLANACLIENT")
	if client == "public" {
		return rpc.MainnetRPCEndpoint
	} else {
		return rpc.TestnetRPCEndpoint
	}
}

//get the connected solana networ(public or test)

func GetSolanaRPC() string {
	client := GoDotEnvVariable("SOLANACLIENT")
	if client == "public" {
		return rp.MainNetBeta_WS
	} else {
		return rp.TestNet_WS
	}
}

func GetHorizonClientNetworkName() string {
	clientNetworkName := GoDotEnvVariable("HORIZONCLIENT")
	return clientNetworkName
}

func GetStellarLaboratoryClient() string {
	clientUrl := GoDotEnvVariable("STELLARLABORATORYCLIENT")
	return clientUrl
}
