package businessFacades

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	log "github.com/sirupsen/logrus"
)

/*MintNFTStellar
@desc - Call the IssueNft method and store new NFT details in the DB
@params - ResponseWriter,Request
*/
func MintNFTContract(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Println("------------------------------inside gateway contract process-------------------")
	var ResponseNFT model.NFTContracts
	var NFTcollectionObj model.NFTWithTransactionContracts
	var MarketplaceNFTNFTcollectionObj model.MarketPlaceNFT
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&ResponseNFT)
	if err != nil {
		panic(err)
	}
	log.Println(ResponseNFT)
	if ResponseNFT.OwnerPK != "" && ResponseNFT.Asset_code != "" && ResponseNFT.NFTURL != "" {
		NFTcollectionObj = model.NFTWithTransactionContracts{
			Identifier:                       ResponseNFT.Identifier,
			Categories:                       ResponseNFT.Categories,
			Collection:                       ResponseNFT.Collection,
			ImageBase64:                      ResponseNFT.NFTURL,
			NFTTXNhash:                       ResponseNFT.MintNFTTxn,
			NftTransactionExistingBlockchain: ResponseNFT.NFTBlockChain,
			NftIssuingBlockchain:             ResponseNFT.NFTBlockChain,
			Timestamp:                        "00-00-00",
			NftURL:                           ResponseNFT.NFTLinks,
			NftContentName:                   ResponseNFT.Asset_code,
			NftContent:                       "TRACIFIED Contract Issued",
			OwnerPK:                          ResponseNFT.OwnerPK,
			Description:                      ResponseNFT.Description,
			Copies:                           ResponseNFT.Copies,
			NFTArtistName:                    ResponseNFT.ArtistName,
			NFTArtistURL:                     ResponseNFT.ArtistLink,
			NFTContract:                      ResponseNFT.NFTContract,
			MarketplaceContract:              ResponseNFT.MarketplaceContract,
		}
		log.Println("-------------------------------------first tier----------------------", NFTcollectionObj)
		MarketplaceNFTNFTcollectionObj = model.MarketPlaceNFT{
			Identifier:                       ResponseNFT.Identifier,
			Categories:                       ResponseNFT.Categories,
			Collection:                       ResponseNFT.Collection,
			ImageBase64:                      ResponseNFT.NFTURL,
			NFTTXNhash:                       ResponseNFT.MintNFTTxn,
			NftTransactionExistingBlockchain: ResponseNFT.NFTBlockChain,
			NftIssuingBlockchain:             ResponseNFT.NFTBlockChain,
			Timestamp:                        "00-00-00",
			NftURL:                           ResponseNFT.NFTLinks,
			NftContentName:                   ResponseNFT.Asset_code,
			NftContent:                       "TRACIFIED Contract Issued",
			NFTArtistName:                    ResponseNFT.ArtistName,
			NFTArtistURL:                     ResponseNFT.ArtistLink,
			Description:                      ResponseNFT.Description,
			Copies:                           ResponseNFT.Copies,
			OriginPK:                         ResponseNFT.OwnerPK,
			InitialDistributorPK:             ResponseNFT.MarketplaceContract,
			InitialIssuerPK:                  ResponseNFT.NFTContract,
			MainAccountPK:                    ResponseNFT.OwnerPK,
			TrustLineCreatedAt:               "No trust lines for contracts",
			PreviousOwnerNFTPK:               "TRACIFIED",
			CurrentOwnerNFTPK:                ResponseNFT.OwnerPK,
			SellingStatus:                    "NOTFORSALE",
		}
		log.Println("-------------------------------------second tier----------------------", MarketplaceNFTNFTcollectionObj)
		NFTCeactedResponse := model.NFTCreactedResponse{
			NFTTxnHash:         ResponseNFT.MintNFTTxn,
			TDPTxnHash:         ResponseNFT.NFTURL,
			NFTName:            ResponseNFT.Asset_code,
			NFTIssuerPublicKey: ResponseNFT.NFTContract,
		}
		log.Println("-------------------------------------third tier----------------------", NFTCeactedResponse)

		if NFTcollectionObj.NftIssuingBlockchain == "polygon" && MarketplaceNFTNFTcollectionObj.NftIssuingBlockchain == "polygon" {
			object := dao.Connection{}
			log.Println("-------------------inside polygon handler---------")
			err1, err2 := object.InsertPolygonNFT(NFTcollectionObj, MarketplaceNFTNFTcollectionObj)
			if err1 != nil && err2 != nil {
				log.Error("NFT not inserted : ", err1, err2)
			} else {
				log.Error("NFT inserted to the collection - Polygon")
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(NFTCeactedResponse)
			return
		}

		if NFTcollectionObj.NftIssuingBlockchain == "ethereum" && MarketplaceNFTNFTcollectionObj.NftIssuingBlockchain == "ethereum" {
			object := dao.Connection{}
			err1, err2 := object.InsertEthereumNFT(NFTcollectionObj, MarketplaceNFTNFTcollectionObj)
			if err1 != nil && err2 != nil {
				log.Error("NFT not inserted : ", err1, err2)
			} else {
				log.Error("NFT inserted to the collection - Ethereum")
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(NFTCeactedResponse)
			return
		}

	} else {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Can not issue NFT2"}
		json.NewEncoder(w).Encode(response)
	}
}
