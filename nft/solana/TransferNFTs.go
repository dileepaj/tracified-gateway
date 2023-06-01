package solana

import (
	"context"
	"errors"
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/gagliardetto/solana-go"
	rp "github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/associated_token_account"
	"github.com/portto/solana-go-sdk/program/token"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

func TransferNFTs(fromWalletSecret string, from string, to string, mintpk string) (*string, error) {

	var fromWallet, _ = types.AccountFromBase58(fromWalletSecret)

	c := client.NewClient(commons.GetSolanaNetwork())

	ata, _, err := common.FindAssociatedTokenAddress(fromWallet.PublicKey, common.PublicKeyFromString(mintpk))
	if err != nil {
		return nil, err
	}

	toata, _, err := common.FindAssociatedTokenAddress(common.PublicKeyFromString(to), common.PublicKeyFromString(mintpk))
	if err != nil {
		return nil, err
	}

	recentBlockhashResponse, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		return nil, err
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{fromWallet},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        fromWallet.PublicKey,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
			Instructions: []types.Instruction{
				associated_token_account.CreateAssociatedTokenAccount(associated_token_account.CreateAssociatedTokenAccountParam{
					Funder:                 fromWallet.PublicKey,
					Owner:                  common.PublicKeyFromString(to),
					Mint:                   common.PublicKeyFromString(mintpk),
					AssociatedTokenAccount: toata,
				}),
				token.TransferChecked(token.TransferCheckedParam{
					From:     ata,
					To:       toata,
					Mint:     common.PublicKeyFromString(mintpk),
					Auth:     fromWallet.PublicKey,
					Amount:   1,
					Decimals: 0,
				}),
			},
		}),
	})
	if err != nil {
		return nil, err
	}

	sign, err := c.SendTransactionWithConfig(context.TODO(), tx, client.SendTransactionConfig{
		SkipPreflight:       false,
		PreflightCommitment: rpc.CommitmentFinalized,
	})

	if err != nil {
		return nil, err
	}

	wsClient, err := ws.Connect(context.Background(), commons.GetSolanaRPC())
	sub, err := wsClient.SignatureSubscribe(
		solana.MustSignatureFromBase58(sign),
		rp.CommitmentFinalized,
	)

	if err != nil {
		return nil, err
	}
	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv()
		if err != nil {
			return &sign, nil
		}
		if got.Value.Err != nil {
			log.Println(errors.New("transaction confirmation failed"))
		} else {
			return &sign, nil
		}
	}

}
