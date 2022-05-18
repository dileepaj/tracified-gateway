package solana

import (
	"context"
	"errors"

	"github.com/gagliardetto/solana-go"
	rp "github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/assotokenprog"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

func CreateTokenAccount(mintPK common.PublicKey, toWalletSecret []byte) (*common.PublicKey, *string, error) {

	var toWallet, _ = types.AccountFromBytes(toWalletSecret)

	c := client.NewClient(rpc.TestnetRPCEndpoint)

	ata, _, err := common.FindAssociatedTokenAddress(toWallet.PublicKey, mintPK)

	if err != nil {
		return nil, nil, errors.New("cannot find associated token account")
	}
	// fmt.Println("ata:", ata.ToBase58())

	res, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		return nil, nil, errors.New("get recent block hash error")
	}
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        toWallet.PublicKey,
			RecentBlockhash: res.Blockhash,
			Instructions: []types.Instruction{
				assotokenprog.CreateAssociatedTokenAccount(assotokenprog.CreateAssociatedTokenAccountParam{
					Funder:                 toWallet.PublicKey,
					Owner:                  toWallet.PublicKey,
					Mint:                   mintPK,
					AssociatedTokenAccount: ata,
				}),
			},
		}),
		Signers: []types.Account{toWallet},
	})
	if err != nil {
		return nil, nil, err
	}

	sign, err := c.SendTransactionWithConfig(context.TODO(), tx, client.SendTransactionConfig{
		SkipPreflight:       false,
		PreflightCommitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		return nil, nil, err
	}

	wsClient, err := ws.Connect(context.Background(), rp.TestNet_WS)
	sub, err := wsClient.SignatureSubscribe(
		solana.MustSignatureFromBase58(sign),
		rp.CommitmentFinalized,
	)

	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv()
		if err != nil {
			return &ata, &sign, nil
		}
		if got.Value.Err != nil {
			return &ata, &sign, errors.New("transaction confirmation failed")
		} else {
			return &ata, &sign, nil
		}
	}

}
