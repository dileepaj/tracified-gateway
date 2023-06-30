package businessFacades

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/nft/solana"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/gorilla/mux"
	"github.com/portto/solana-go-sdk/common"
	log "github.com/sirupsen/logrus"
)

/*MintNFTStellar
@desc - Call the IssueNft method and store new NFT details in the DB
@params - ResponseWriter,Request
*/
func MintNFTSolana(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	dt := time.Now()
	var TrustLineResponseNFT model.NFTSolana
	var NFTcollectionObj model.NFTWithTransactionSolana
	var MarketplaceNFTcollectionObj model.MarketPlaceNFT
	logger := utilities.NewCustomLogger()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&TrustLineResponseNFT)
	if err != nil {
		logger.LogWriter("Error when decoding the body : "+err.Error(), constants.ERROR)
	}
	if TrustLineResponseNFT.OwnerPK != "" && TrustLineResponseNFT.Asset_code != "" && TrustLineResponseNFT.NFTURL != "" {
		var WALLETSECRET = (commons.GoDotEnvVariable("WALLETSECRET"))

		mintPK, ownerPK, mintedTxHash, ATA, err := solana.MintSolana(WALLETSECRET, TrustLineResponseNFT.Asset_code, TrustLineResponseNFT.NFTURL)
		if err == nil {

			NFTcollectionObj = model.NFTWithTransactionSolana{
				Identifier:                       common.PublicKey(*ATA).String(),
				Categories:                       TrustLineResponseNFT.Categories,
				Collection:                       TrustLineResponseNFT.Collection,
				ImageBase64:                      TrustLineResponseNFT.NFTURL,
				NFTTXNhash:                       *mintedTxHash,
				NftTransactionExistingBlockchain: "Solana",
				NftIssuingBlockchain:             TrustLineResponseNFT.NFTBlockChain,
				Timestamp:                        dt.Format("01-02-2006 15:04:05"),
				NftURL:                           TrustLineResponseNFT.NFTLinks,
				NftContentName:                   TrustLineResponseNFT.Asset_code,
				NftContent:                       "TRACIFIED SOLANA",
				MinterPK:                         common.PublicKey(*mintPK).String(),
				OwnerPK:                          TrustLineResponseNFT.OwnerPK,
				Description:                      TrustLineResponseNFT.Description,
				Copies:                           TrustLineResponseNFT.Copies,
				NFTArtistName:                    TrustLineResponseNFT.ArtistName,
				NFTArtistURL:                     TrustLineResponseNFT.ArtistLink,
				InitialDistributorPK:             common.PublicKey(*ownerPK).String(),
			}

			MarketplaceNFTcollectionObj = model.MarketPlaceNFT{
				Identifier:                       common.PublicKey(*ATA).String(),
				Categories:                       TrustLineResponseNFT.Categories,
				Collection:                       TrustLineResponseNFT.Collection,
				ImageBase64:                      TrustLineResponseNFT.NFTURL,
				NFTTXNhash:                       *mintedTxHash,
				NftTransactionExistingBlockchain: "Solana",
				NftIssuingBlockchain:             TrustLineResponseNFT.NFTBlockChain,
				Timestamp:                        dt.Format("01-02-2006 15:04:05"),
				NftURL:                           TrustLineResponseNFT.NFTLinks,
				NftContentName:                   TrustLineResponseNFT.Asset_code,
				NftContent:                       "Solana",
				NFTArtistName:                    TrustLineResponseNFT.ArtistName,
				NFTArtistURL:                     TrustLineResponseNFT.ArtistLink,
				Description:                      TrustLineResponseNFT.Description,
				Copies:                           TrustLineResponseNFT.Copies,
				OriginPK:                         common.PublicKey(*ownerPK).String(),
				InitialDistributorPK:             common.PublicKey(*ownerPK).String(),
				InitialIssuerPK:                  common.PublicKey(*mintPK).String(),
				MainAccountPK:                    common.PublicKey(*ownerPK).String(),
				TrustLineCreatedAt:               "No trust lines for solana",
				PreviousOwnerNFTPK:               "TRACIFIED",
				CurrentOwnerNFTPK:                TrustLineResponseNFT.OwnerPK,
				SellingStatus:                    "NOTFORSALE",
			}

			NFTCeactedResponse := model.NFTCreactedResponse{
				NFTTxnHash:         *mintedTxHash,
				TDPTxnHash:         TrustLineResponseNFT.NFTURL,
				NFTName:            TrustLineResponseNFT.Asset_code,
				NFTIssuerPublicKey: common.PublicKey(*mintPK).String(),
			}
			object := dao.Connection{}
			err1, err2 := object.InsertSolanaNFT(NFTcollectionObj, MarketplaceNFTcollectionObj)
			if err1 != nil && err2 != nil {
				log.Error("NFT not inserted : ", err1, err2)
			}
			if err1 == nil && err2 != nil {
				log.Error("NFT not inserted into SolanaNFT Collection : ", err2)
			}
			if err1 != nil && err2 == nil {
				log.Error("NFT not inserted into Marketplace Collection : ", err1)
			} else {
				log.Error("NFT inserted to the collection")
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(NFTCeactedResponse)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Error{Message: "Can not save NFT in DB"}
			json.NewEncoder(w).Encode(response)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Can not issue NFT"}
		json.NewEncoder(w).Encode(response)
	}
}

func RetrieveSolanaMinter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	object := dao.Connection{}
	p := object.GetNFTMinterPKSolana(vars["ImageBase64"], vars["blockchain"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.NFTWithTransactionSolana)
		res := model.Minter{NFTIssuerPK: result.MinterPK, NFTTxnHash: result.NFTTXNhash, NFTIdentifier: result.Identifier, CreatorUserID: result.InitialDistributorPK}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "URL Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()

}

func TransferNFTS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var MarketplaceNFT model.NFTTransfer
	logger := utilities.NewCustomLogger()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&MarketplaceNFT)
	if err != nil {
		logger.LogWriter("Error when decoding body : "+err.Error(), constants.ERROR)
	}
	if MarketplaceNFT.Source != "" {
		var WALLETSECRET = (commons.GoDotEnvVariable("WALLETSECRET"))
		transferTXNX, err := solana.TransferNFTs(WALLETSECRET, MarketplaceNFT.Source, MarketplaceNFT.Destination, MarketplaceNFT.MintPubKey)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(transferTXNX)
			return
		} else {
			transferTXNX, err := solana.TransferNFTsToExistingAccount(WALLETSECRET, MarketplaceNFT.Source, MarketplaceNFT.Destination, MarketplaceNFT.MintPubKey)
			if err == nil {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(transferTXNX)
				return
			} else {
				w.WriteHeader(http.StatusBadRequest)
				response := model.Error{Message: "URL Not Found in Gateway DataStore"}
				json.NewEncoder(w).Encode(response)
				return
			}

		}

	}
}
