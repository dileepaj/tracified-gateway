package solana

import (
	"context"
	"errors"
	"log"

	//"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/gagliardetto/solana-go"
	rp "github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/pkg/pointer"
	"github.com/portto/solana-go-sdk/program/associated_token_account"
	"github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/program/system"
	"github.com/portto/solana-go-sdk/program/token"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

func MintSolana(fromWalletSecret string, code_name string, code_url string, royalty string, owner string) (*common.PublicKey, *common.PublicKey, *string, *common.PublicKey, error) {
	var fromWallet, _ = types.AccountFromBase58(fromWalletSecret)

	// royaltyvalue, err := strconv.Atoi(royalty)
	// if err != nil {
	// 	return nil, nil, nil, nil, err
	// }

	c := client.NewClient(commons.GetSolanaNetwork())

	mint := types.NewAccount()

	ata, _, err := common.FindAssociatedTokenAddress(fromWallet.PublicKey, mint.PublicKey)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tokenMetadataPubkey, err := token_metadata.GetTokenMetaPubkey(mint.PublicKey)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tokenMasterEditionPubkey, err := token_metadata.GetMasterEdition(mint.PublicKey)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	mintAccountRent, err := c.GetMinimumBalanceForRentExemption(context.Background(), token.MintAccountSize)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	recentBlockhashResponse, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{mint, fromWallet},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        fromWallet.PublicKey,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
			Instructions: []types.Instruction{
				system.CreateAccount(system.CreateAccountParam{
					From:     fromWallet.PublicKey,
					New:      mint.PublicKey,
					Owner:    common.TokenProgramID,
					Lamports: mintAccountRent,
					Space:    token.MintAccountSize,
				}),

				token.InitializeMint(token.InitializeMintParam{
					Decimals:   0,
					Mint:       mint.PublicKey,
					MintAuth:   fromWallet.PublicKey,
					FreezeAuth: &fromWallet.PublicKey,
				}),

				token_metadata.CreateMetadataAccount(token_metadata.CreateMetadataAccountParam{
					Metadata:                tokenMetadataPubkey,
					Mint:                    mint.PublicKey,
					MintAuthority:           fromWallet.PublicKey,
					Payer:                   fromWallet.PublicKey,
					UpdateAuthority:         fromWallet.PublicKey,
					UpdateAuthorityIsSigner: true,
					IsMutable:               true,
					MintData: token_metadata.Data{
						Name:                 code_name,
						Symbol:               "TRAC",
						Uri:                  code_url,
						SellerFeeBasisPoints: 0,
						// Creators: &[]token_metadata.Creator{
						// 	{
						// 		Address:  common.PublicKeyFromString(owner),
						// 		Verified: true,
						// 		Share:    uint8(royaltyvalue),
						// 	},
						// },
					},
				}),

				associated_token_account.CreateAssociatedTokenAccount(associated_token_account.CreateAssociatedTokenAccountParam{
					Funder:                 fromWallet.PublicKey,
					Owner:                  fromWallet.PublicKey,
					Mint:                   mint.PublicKey,
					AssociatedTokenAccount: ata,
				}),
				token.MintTo(token.MintToParam{
					Mint:   mint.PublicKey,
					To:     ata,
					Auth:   fromWallet.PublicKey,
					Amount: 1,
				}),

				token_metadata.CreateMasterEditionV3(token_metadata.CreateMasterEditionParam{
					Edition:         tokenMasterEditionPubkey,
					Mint:            mint.PublicKey,
					UpdateAuthority: fromWallet.PublicKey,
					MintAuthority:   fromWallet.PublicKey,
					Metadata:        tokenMetadataPubkey,
					Payer:           fromWallet.PublicKey,
					MaxSupply:       pointer.Get[uint64](0),
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
		MaxRetries:          0,
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
