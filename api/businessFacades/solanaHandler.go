package businessFacades

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/nft/solana"
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
	log.Println("------------------------------inside gateway minting process-------------------")
	var TrustLineResponseNFT model.NFTSolana
	var NFTcollectionObj model.NFTWithTransactionSolana
	var MarketplaceNFTNFTcollectionObj model.MarketPlaceNFT
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&TrustLineResponseNFT)
	if err != nil {
		panic(err)
	}
	log.Println(TrustLineResponseNFT)
	if TrustLineResponseNFT.OwnerPK != "" && TrustLineResponseNFT.Asset_code != "" && TrustLineResponseNFT.NFTURL != "" {
		log.Println("\n..................... BEGIN MINTING NFT ...................")
		var fromWalletSecret = []byte{10, 75, 10, 90, 145, 78, 142, 248, 104, 3, 36, 7, 69, 207, 109, 98, 82, 58, 146, 202, 44, 188, 70, 70, 64, 173, 35, 130, 18, 133, 107, 236, 231, 43, 70, 165, 182, 191, 162, 242, 126, 119, 49, 3, 231, 43, 249, 47, 228, 225, 70, 91, 254, 22, 160, 42, 20, 186, 184, 196, 240, 151, 157, 207}
		mintPK, ownerPK, mintedTxHash, ATA, err := solana.MintSolana(fromWalletSecret, TrustLineResponseNFT.Asset_code, TrustLineResponseNFT.NFTURL)
		log.Println("\nMINTED PK", mintPK)
		log.Println("OWNER PK", ownerPK)
		log.Println("TX HASH", *mintedTxHash)
		log.Println("ATA Tracified", ATA)
		log.Println("\n..................... END MINTING NFT ...................")
		if err == nil {
			// var toWalletSecret = []byte{47, 163, 68, 180, 12, 82, 124, 0, 101, 163, 250, 17, 181, 250, 63, 165, 179, 85, 112, 117, 245, 102, 63, 181, 48, 68, 190, 193, 178, 112, 227, 57, 17, 239, 150, 83, 192, 134, 121, 241, 161, 240, 133, 128, 9, 112, 247, 2, 71, 181, 138, 177, 227, 201, 12, 225, 164, 158, 122, 91, 176, 169, 10, 147}

			// transferTXHash := solana.Transfer(*mintPK)
			// log.Println("\nTX HASH", *transferTXHash)
			NFTcollectionObj = model.NFTWithTransactionSolana{
				Identifier:                       common.PublicKey(*ATA).String(),
				Categories:                       TrustLineResponseNFT.Categories,
				Collection:                       TrustLineResponseNFT.Collection,
				ImageBase64:                      TrustLineResponseNFT.NFTURL,
				NFTTXNhash:                       *mintedTxHash,
				NftTransactionExistingBlockchain: "Solana",
				NftIssuingBlockchain:             TrustLineResponseNFT.NFTBlockChain,
				Timestamp:                        "00-00-00",
				NftURL:                           TrustLineResponseNFT.NFTLinks,
				NftContentName:                   TrustLineResponseNFT.Asset_code,
				NftContent:                       "TRACIFIED SOLANA",
				MinterPK:                         common.PublicKey(*mintPK).String(),
				OwnerPK:                          TrustLineResponseNFT.OwnerPK,
				Description:                      TrustLineResponseNFT.Description,
				Copies:                           TrustLineResponseNFT.Copies,
				NFTArtistName:                    TrustLineResponseNFT.ArtistName,
				NFTArtistURL:                     TrustLineResponseNFT.ArtistLink,
			}

			MarketplaceNFTNFTcollectionObj = model.MarketPlaceNFT{
				Identifier:                       common.PublicKey(*ownerPK).String(),
				Categories:                       TrustLineResponseNFT.Categories,
				Collection:                       TrustLineResponseNFT.Collection,
				ImageBase64:                      TrustLineResponseNFT.NFTURL,
				NFTTXNhash:                       *mintedTxHash,
				NftTransactionExistingBlockchain: "Solana",
				NftIssuingBlockchain:             TrustLineResponseNFT.NFTBlockChain,
				Timestamp:                        "00-00-00",
				NftURL:                           TrustLineResponseNFT.NFTLinks,
				NftContentName:                   TrustLineResponseNFT.Asset_code,
				NftContent:                       "Solana",
				NFTArtistName:                    TrustLineResponseNFT.ArtistName,
				NFTArtistURL:                     TrustLineResponseNFT.ArtistLink,
				Description:                      TrustLineResponseNFT.Description,
				Copies:                           TrustLineResponseNFT.Copies,
				OriginPK:                         common.PublicKey(*ownerPK).String(),
				InitialDistributorPK:             common.PublicKey(*ATA).String(),
				InitialIssuerPK:                  common.PublicKey(*mintPK).String(),
				MainAccountPK:                    commons.GoDotEnvVariable("NFTSTELLARISSUERPUBLICKEYK"),
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
			err1, err2 := object.InsertSolanaNFT(NFTcollectionObj, MarketplaceNFTNFTcollectionObj)
			if err1 != nil && err2 != nil {
				log.Error("NFT not inserted : ", err1, err2)
			} else {
				log.Error("NFT inserted to the collection")
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(NFTCeactedResponse)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Error{Message: "Can not issue NFT1"}
			json.NewEncoder(w).Encode(response)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Can not issue NFT2"}
		json.NewEncoder(w).Encode(response)
	}
}

func RetrieveSolanaMinter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	object := dao.Connection{}
	p := object.GetNFTMinterPKSolana(vars["ImageBase64"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.NFTWithTransactionSolana)
		res := model.Minter{NFTIssuerPK: result.MinterPK, NFTTxnHash: result.NFTTXNhash, NFTIdentifier: result.Identifier}
		log.Println("-----------------minter pk and hash ------------", res)
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

// func TransferNFT(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

// 	vars := mux.Vars(r)
// 	mint:=vars["InitialIssuerPK"]
//      mintPubKey := *mint.PublicKey() ;
// 	p := solana.Transfer(mintPubKey)

// 	if p!=nil{
// 		w.WriteHeader(http.StatusBadRequest)
// 			response := model.Error{Message: "NFT Transferred"}
// 			json.NewEncoder(w).Encode(response)
// 	}else{
// 		w.WriteHeader(http.StatusBadRequest)
// 			response := model.Error{Message: "Error occured when transferring"}
// 			json.NewEncoder(w).Encode(response)
// 	}

// }
