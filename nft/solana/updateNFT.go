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
	"github.com/portto/solana-go-sdk/pkg/pointer"
	"github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

func UpdateNFT(fromWalletSecret string, mintkey string, code_name string, code_uri string, symbol string) (string, error) {
	var fromWallet, _ = types.AccountFromBase58(fromWalletSecret)

	c := client.NewClient(commons.GetSolanaNetwork())

	// mint address
	nft := common.PublicKeyFromString(mintkey)

	tokenMetadataPubkey, err := token_metadata.GetTokenMetaPubkey(nft)
	if err != nil {
		log.Fatalf("failed to find a valid token metadata, err: %v", err)

	}

	recentBlockhashResponse, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{fromWallet},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        fromWallet.PublicKey,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
			Instructions: []types.Instruction{
				token_metadata.UpdateMetadataAccount(token_metadata.UpdateMetadataAccountParam{
					MetadataAccount: tokenMetadataPubkey,
					UpdateAuthority: fromWallet.PublicKey,
					Data: &token_metadata.Data{
						Name:                 code_name,
						Symbol:               symbol,
						Uri:                  code_uri,
						SellerFeeBasisPoints: 10000,
						// Creators: &[]token_metadata.Creator{
						// 	{
						// 		Address:  fromWallet.PublicKey,
						// 		Verified: false,
						// 		Share:    100,
						// 	},
						// },
					},
					NewUpdateAuthority:  &fromWallet.PublicKey,
					PrimarySaleHappened: pointer.Get[bool](true),
				}),
			},
		}),
	})
	if err != nil {
		log.Fatalf("failed to new a tx, err: %v", err)
	}

	sign, err := c.SendTransactionWithConfig(context.TODO(), tx, client.SendTransactionConfig{
		SkipPreflight:       false,
		PreflightCommitment: rpc.CommitmentFinalized,
		MaxRetries:          0,
	})
	if err != nil {
		return "", err
	}
	wsClient, err := ws.Connect(context.Background(), commons.GetSolanaRPC())
	sub, err := wsClient.SignatureSubscribe(
		solana.MustSignatureFromBase58(sign),
		rp.CommitmentFinalized,
	)

	if err != nil {
		return "", err
	}
	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv()
		if err != nil {
			return "", err
		}
		if got.Value.Err != nil {
			log.Println(errors.New("transaction confirmation failed"))
		} else {
			log.Println("update signatured", sign)
			return sign, nil
		}
	}
}
