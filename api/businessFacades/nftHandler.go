package businessFacades

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/nft/stellar"
)

func GettingChangeTrustXDRForNFT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var TrustLineResponseNFT model.TrustLineResponseNFT
	var NFTcollectionObj model.NFTWithTransaction
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&TrustLineResponseNFT)
	if err != nil {
		fmt.Println("eerrrr")
	}
	fmt.Println("TrustLineResponseNFT", TrustLineResponseNFT)
	//callthe issue NFt method(distributerPK,assetcode,TDPtxnhas) mint and return the (NFTtxnhash,issuerPK,NFTContent)
	if TrustLineResponseNFT.TrustLineCreatedAt != "" && TrustLineResponseNFT.DistributorPublickKey != "" && TrustLineResponseNFT.Asset_code != "" && TrustLineResponseNFT.TDPtxnhash != "" && TrustLineResponseNFT.Successfull {
		var NFTtxnhash, issuerPK, NftContent, err = stellar.IssueNft(TrustLineResponseNFT.DistributorPublickKey, TrustLineResponseNFT.Asset_code, TrustLineResponseNFT.TDPtxnhash)
		if err == nil {
			NFTcollectionObj = model.NFTWithTransaction{
				Identifier:                       TrustLineResponseNFT.Identifier,
				TxnType:                          "TDP",
				TDPID:							  TrustLineResponseNFT.TDPID,
				TDPTxnHash:                       TrustLineResponseNFT.TDPtxnhash,
				DataHash:                         TrustLineResponseNFT.TDPtxnhash,
				NFTTXNhash:                       NFTtxnhash,
				NftTransactionExistingBlockchain: "Stellar",
				NftIssuingBlockchain:             TrustLineResponseNFT.NFTBlockChain,
				Timestamp:                        "00-00-00",
				NftAssetName:                     TrustLineResponseNFT.Asset_code,
				NftContentName:                   TrustLineResponseNFT.Asset_code,
				NftContent:                       NftContent,
				DistributorPublickKey:            TrustLineResponseNFT.DistributorPublickKey,
				IssuerPK:                         issuerPK,
				TrustLineCreatedAt:               TrustLineResponseNFT.TrustLineCreatedAt,
				ProductName:					  TrustLineResponseNFT.ProductName,
			}

			NFTCeactedResponse := model.NFTCreactedResponse{
				NFTTxnHash: NFTtxnhash,
				TDPTxnHash: TrustLineResponseNFT.TDPtxnhash,
				NFTName:    TrustLineResponseNFT.Asset_code,
			}
				//after that call the insertTODB NFT Details to NFT  (gateway(NFTstellar collection)
			object := dao.Connection{}
			err1 := object.InsertStellarNFT(NFTcollectionObj)
			if err1 != nil {
			fmt.Println("NFT not inserted : ", err1)
			} else {
			fmt.Println("NFT inserted to the collection")
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