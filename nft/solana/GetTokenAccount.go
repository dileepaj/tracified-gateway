package solana

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/rpc"
)

type GetTokenAccountsByOwner struct {
	ID      int64         `json:"id"`
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type TokenAccount struct {
	ID      int64  `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Context struct {
			Slot int64 `json:"slot"`
		} `json:"context"`
		Value []struct {
			Account struct {
				Data struct {
					Parsed struct {
						Info struct {
							IsNative    bool   `json:"isNative"`
							Mint        string `json:"mint"`
							Owner       string `json:"owner"`
							State       string `json:"state"`
							TokenAmount struct {
								Amount         string  `json:"amount"`
								Decimals       int64   `json:"decimals"`
								UIAmount       float64 `json:"uiAmount"`
								UIAmountString string  `json:"uiAmountString"`
							} `json:"tokenAmount"`
						} `json:"info"`
						Type string `json:"type"`
					} `json:"parsed"`
					Program string  `json:"program"`
					Space   float64 `json:"space"`
				} `json:"data"`
				Executable bool    `json:"executable"`
				Lamports   float64 `json:"lamports"`
				Owner      string  `json:"owner"`
				RentEpoch  float64 `json:"rentEpoch"`
			} `json:"account"`
			Pubkey string `json:"pubkey"`
		} `json:"value"`
	} `json:"result"`
}

func GetTokenAccount(mintPK common.PublicKey, accountPK common.PublicKey) (*common.PublicKey, error) {
	payload := GetTokenAccountsByOwner{
		Jsonrpc: "2.0",
		Method:  "getTokenAccountsByOwner",
		Params: []interface{}{
			accountPK.String(),
			map[string]interface{}{
				"mint": mintPK.String(),
			},
			map[string]interface{}{
				"encoding": "jsonParsed",
			},
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return nil, err
	}

	result, err := http.Post(rpc.TestnetRPCEndpoint, "application/json", &buf)
	if err != nil {
		return nil, err
	}
	data, _ := ioutil.ReadAll(result.Body)
	if result.StatusCode == 200 {
		var tokenAccount TokenAccount
		err := json.Unmarshal(data, &tokenAccount)
		if err != nil {
			return nil, err
		}
		if len(tokenAccount.Result.Value) == 0 {
			return nil, errors.New("No token accounts")
		}
		pk := common.PublicKeyFromString(tokenAccount.Result.Value[0].Pubkey)
		return &pk, nil
	}
	return nil, errors.New("Unexpected Error")
}
