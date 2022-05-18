package solana

import (
	"context"
	//"errors"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/gagliardetto/solana-go"
	rp "github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/tokenprog"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

func Transfer(mintPK common.PublicKey) *string {

	var fromWallet, _ = types.AccountFromBytes([]byte(commons.GoDotEnvVariable("FROMWALLETSECRET")))
	var toWallet, _ = types.AccountFromBytes([]byte(commons.GoDotEnvVariable("TOWALLETSECRET")))

	fromTokenAccount, err := GetTokenAccount(mintPK, fromWallet.PublicKey)
	if err != nil {
		return nil
	}

	toTokenAccount, err := GetTokenAccount(mintPK, toWallet.PublicKey)
	if err != nil {
		toTokenAccount, _, err = CreateTokenAccount(mintPK, []byte(commons.GoDotEnvVariable("TOWALLETSECRET")))
		if err != nil {
			return nil
		}
	}

	c := client.NewClient(rpc.TestnetRPCEndpoint)

	res, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		return nil
	}
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        toWallet.PublicKey,
			RecentBlockhash: res.Blockhash,
			Instructions: []types.Instruction{
				tokenprog.TransferChecked(tokenprog.TransferCheckedParam{
					From:     *fromTokenAccount,
					To:       *toTokenAccount,
					Mint:     mintPK,
					Auth:     fromWallet.PublicKey,
					Signers:  []common.PublicKey{},
					Amount:   1,
					Decimals: 0,
				}),
			},
		}),
		Signers: []types.Account{toWallet, fromWallet},
	})
	if err != nil {
		return nil
	}

	txhash, err := c.SendTransactionWithConfig(context.TODO(), tx, client.SendTransactionConfig{
		SkipPreflight:       false,
		PreflightCommitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		return nil
	}

	wsClient, err := ws.Connect(context.Background(), rp.TestNet_WS)
	sub, err := wsClient.SignatureSubscribe(
		solana.MustSignatureFromBase58(txhash),
		rp.CommitmentFinalized,
	)

	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv()
		if err != nil {
			return &txhash
		}
		if got.Value.Err != nil {
			return nil
		} else {
			return &txhash
		}
	}

}
