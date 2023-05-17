package commons

import (
	rp "github.com/gagliardetto/solana-go/rpc"
	"github.com/portto/solana-go-sdk/rpc"
)

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
