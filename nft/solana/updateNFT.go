package solana

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/gagliardetto/solana-go"
	rp "github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/pkg/pointer"
	"github.com/portto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
	"github.com/sirupsen/logrus"
)

func UpdateNFTs(marketplaceNFT model.UpdateableNFT) error {
	object := dao.Connection{}

	log.Println(marketplaceNFT)
	if marketplaceNFT.BatchId != "" && marketplaceNFT.TenantId != "" {
		var result model.NFTWithTransactionSolana
		p := object.GetExistingSolanaNFT(marketplaceNFT.BatchId, marketplaceNFT.ProductId, marketplaceNFT.TenantId)
		p.Then(func(data interface{}) interface{} {
			result = data.(model.NFTWithTransactionSolana)
			return nil
		}).Catch(func(error error) error {
			logrus.Error("Error while GetTransactionForTdpIdSequence " + error.Error())
			return error
		}).Await()

		if result.MinterPK != "" {
			encodedString := base64.StdEncoding.EncodeToString([]byte(marketplaceNFT.BatchId))
			url := constants.NFTBackend + `/nft/timeline/html/hash/` + marketplaceNFT.ProductId + `/` + encodedString
			logrus.Info("url", url)

			responseBody, err := commons.MakeGetRequest(url)
			if err != nil {
				logrus.Error("Error:", err)
				return err
			}
			var timelineHtml model.HTMLTimelineHashGenerationResponse
			err = json.Unmarshal(responseBody, &timelineHtml)
			if err != nil {
				logrus.Error("Cannot unmarshal response" + err.Error())
				return err
			}

			var hash string = timelineHtml.TimelineHtmlHash + "-" + result.ImageBase64
			fmt.Println("new hash ", hash)
			var stringver string
			WALLETSECRET := (commons.GoDotEnvVariable("WALLETSECRET"))
			updateTXNX, err := UpdateNFT(WALLETSECRET, result.MinterPK, result.NftContentName, hash, "UNFT")
			if err == nil {
				ver, errAtoi := strconv.Atoi(result.Version)
				if errAtoi != nil {
					return errAtoi
				}
				ver++
				stringver = strconv.Itoa(ver)
				updatednft := model.UpdateableNFT{
					BatchId:   marketplaceNFT.BatchId,
					ProductId: marketplaceNFT.ProductId,
					TenantId:  marketplaceNFT.TenantId,
					SvgHash:   timelineHtml.TimelineHtmlHash,
					Version:   stringver,
					TxnHash:   updateTXNX,
					MinterPK:  result.MinterPK,
				}
				errversion := object.InsertSolanaNFTVersions(updatednft) // this might turn into an update function
				if errversion == nil {
					errupdate := object.UpdateNFTSolana(updatednft)
					if errupdate != nil {
						return errupdate
					}
				} else {
					err := errors.New("Couldnt create a version record in the gateway")
					return err
				}

				return nil
			} else {
				err := errors.New("Something went wrong")
				return err
			}
		} else {
			err := errors.New("No NFT for Batch minted previously!")
			return err
		}

	} else {
		err := errors.New("No NFT for Batch minted previously!")
		return err
	}
}

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
