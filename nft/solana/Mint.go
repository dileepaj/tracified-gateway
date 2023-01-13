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
	"github.com/portto/solana-go-sdk/program/assotokenprog"
	"github.com/portto/solana-go-sdk/program/metaplex/tokenmeta"
	"github.com/portto/solana-go-sdk/program/sysprog"
	"github.com/portto/solana-go-sdk/program/tokenprog"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

func MintSolana(fromWalletSecret string, code_name string, code_url string) (*common.PublicKey, *common.PublicKey, *string, *common.PublicKey, error) {

	var fromWallet, _ = types.AccountFromBase58(fromWalletSecret)

	c := client.NewClient(commons.GetSolanaNetwork())

	mint := types.NewAccount()
	ata, _, err := common.FindAssociatedTokenAddress(fromWallet.PublicKey, mint.PublicKey)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tokenMetadataPubkey, err := tokenmeta.GetTokenMetaPubkey(mint.PublicKey)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tokenMasterEditionPubkey, err := tokenmeta.GetMasterEdition(mint.PublicKey)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	mintAccountRent, err := c.GetMinimumBalanceForRentExemption(context.Background(), tokenprog.MintAccountSize)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	recentBlockhashResponse, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{mint, fromWallet},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        fromWallet.PublicKey,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
			Instructions: []types.Instruction{
				sysprog.CreateAccount(sysprog.CreateAccountParam{
					From:     fromWallet.PublicKey,
					New:      mint.PublicKey,
					Owner:    common.TokenProgramID,
					Lamports: mintAccountRent,
					Space:    tokenprog.MintAccountSize,
				}),
				tokenprog.InitializeMint(tokenprog.InitializeMintParam{
					Decimals:   0,
					Mint:       mint.PublicKey,
					MintAuth:   fromWallet.PublicKey,
					FreezeAuth: &fromWallet.PublicKey,
				}),

				tokenmeta.CreateMetadataAccount(tokenmeta.CreateMetadataAccountParam{
					Metadata:                tokenMetadataPubkey,
					Mint:                    mint.PublicKey,
					MintAuthority:           fromWallet.PublicKey,
					Payer:                   fromWallet.PublicKey,
					UpdateAuthority:         fromWallet.PublicKey,
					UpdateAuthorityIsSigner: true,
					IsMutable:               true,
					MintData: tokenmeta.Data{
						Name:                 code_name,
						Symbol:               "FSMB",
						Uri:                  code_url,
						SellerFeeBasisPoints: 100,
						Creators: &[]tokenmeta.Creator{
							{
								Address:  fromWallet.PublicKey,
								Verified: true,
								Share:    100,
							},
						},
					},
				}),
				assotokenprog.CreateAssociatedTokenAccount(assotokenprog.CreateAssociatedTokenAccountParam{
					Funder:                 fromWallet.PublicKey,
					Owner:                  fromWallet.PublicKey,
					Mint:                   mint.PublicKey,
					AssociatedTokenAccount: ata,
				}),
				tokenprog.MintTo(tokenprog.MintToParam{
					Mint:   mint.PublicKey,
					To:     ata,
					Auth:   fromWallet.PublicKey,
					Amount: 1,
				}),

				tokenmeta.CreateMasterEditionV3(tokenmeta.CreateMasterEditionParam{
					Edition:         tokenMasterEditionPubkey,
					Mint:            mint.PublicKey,
					UpdateAuthority: fromWallet.PublicKey,
					MintAuthority:   fromWallet.PublicKey,
					Metadata:        tokenMetadataPubkey,
					Payer:           fromWallet.PublicKey,
					MaxSupply:       pointer.Uint64(0),
				}),
			},
		}),
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	sign, err := c.SendTransactionWithConfig(context.TODO(), tx, client.SendTransactionConfig{
		SkipPreflight:       false,
		PreflightCommitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	wsClient, err := ws.Connect(context.Background(), commons.GetSolanaRPC())
	sub, err := wsClient.SignatureSubscribe(
		solana.MustSignatureFromBase58(sign),
		rp.CommitmentFinalized,
	)

	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv()
		if err != nil {
			return &mint.PublicKey, &fromWallet.PublicKey, &sign, &ata, nil
		}
		if got.Value.Err != nil {
			log.Println(errors.New("transaction confirmation failed"))
		} else {
			return &mint.PublicKey, &fromWallet.PublicKey, &sign, &ata, nil
		}
	}

}
