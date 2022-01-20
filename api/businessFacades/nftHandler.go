package businessFacades

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/nft/stellar"
	"github.com/sirupsen/logrus"
)

func MintNFTStellar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var TrustLineResponseNFT model.TrustLineResponseNFT
	var NFTcollectionObj model.NFTWithTransaction
	var MarketplaceNFTNFTcollectionObj model.MarketPlaceNFT
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

			MarketplaceNFTNFTcollectionObj=model.MarketPlaceNFT{
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
				PreviousOwnerNFTPK: 			  "",
				CurrentOwnerNFTPK: 				  TrustLineResponseNFT.DistributorPublickKey,
				OriginPK: 						  commons.GoDotEnvVariable("NFTSTELLARISSUERPUBLICKEYK"),
				SellingStatus:					  false,
			}

			NFTCeactedResponse := model.NFTCreactedResponse{
				NFTTxnHash: NFTtxnhash,
				TDPTxnHash: TrustLineResponseNFT.TDPtxnhash,
				NFTName:    TrustLineResponseNFT.Asset_code,
			}
				//after that call the insertTODB NFT Details to NFT  (gateway(NFTstellar collection)
			object := dao.Connection{}
			err1 ,err2:= object.InsertStellarNFT(NFTcollectionObj,MarketplaceNFTNFTcollectionObj)
			if err1 != nil && err2!=nil {
			fmt.Println("NFT not inserted : ", err1,err2)
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

func RetriveAllNFTForSell(w http.ResponseWriter, r *http.Request)  {
	var response model.Error;
	var result []model.MarketPlaceNFT
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	key1, error := r.URL.Query()["perPage"]

	if !error || len(key1[0]) < 1 {
		logrus.Error("Url Parameter 'perPage' is missing")
		return
	}

	key2, error := r.URL.Query()["page"]

	if !error || len(key2[0]) < 1 {
		logrus.Error("Url Parameter 'page' is missing")
		return
	}

	key3, error := r.URL.Query()["NoPage"]

	if !error || len(key3[0]) < 1 {
		logrus.Error("Url Parameter 'NoPage' is missing")
		return
	}

	perPage, err := strconv.Atoi(key1[0])
	if err != nil {
		logrus.Error("Query parameter error" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		response = model.Error{Code:http.StatusBadRequest,Message: "The parameter should be an integer " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	page, err := strconv.Atoi(key2[0])
	if err != nil {
		logrus.Error("Query parameter error" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		response = model.Error{Code:http.StatusBadRequest,Message: "The parameter should be an integer " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	object := dao.Connection{}
	var SellingStatus=true
	qdata, err := object.GetAllSellingNFTStellar_Paginated(SellingStatus,perPage, page).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		logrus.Error("Unable to connect gateway datastore")
		w.WriteHeader(http.StatusNotFound)
		response = model.Error{Code:http.StatusNotFound,Message: "Unable to connect gateway datastore"}
		json.NewEncoder(w).Encode(response)
		return
	}
	if qdata == nil {
		logrus.Error("Selling NFTs are not found in gateway datastore")
		w.WriteHeader(http.StatusNoContent)
		response = model.Error{Code:http.StatusNoContent,Message: "Identifier is not found in gateway datastore"}
		json.NewEncoder(w).Encode(response)
		return
	}
	res := qdata.(model.MarketPlaceNFTTrasactionWithCount)

		for _, TxnBody := range res.MarketPlaceNFTItems {
		
			temp := model.MarketPlaceNFT{
				Identifier:                       TxnBody.Identifier,
				TxnType:                          "TDP",
				TDPID:							  TxnBody.TDPID,
				TDPTxnHash:                       TxnBody.TDPTxnHash,
				DataHash:                         TxnBody.DataHash,
				NFTTXNhash:                       TxnBody.NFTTXNhash,
				NftTransactionExistingBlockchain: "Stellar",
				NftIssuingBlockchain:             TxnBody.NftIssuingBlockchain,
				Timestamp:                        "00-00-00",
				NftAssetName:                     TxnBody.NftAssetName,
				NftContentName:                   TxnBody.NftAssetName,
				NftContent:                       TxnBody.NftContent,
				TrustLineCreatedAt:               TxnBody.TrustLineCreatedAt,
				ProductName:					  TxnBody.ProductName,
				PreviousOwnerNFTPK: 			  TxnBody.PreviousOwnerNFTPK,
				CurrentOwnerNFTPK: 				  TxnBody.CurrentOwnerNFTPK,
				OriginPK: 						  TxnBody.OriginPK,
				SellingStatus:					  TxnBody.SellingStatus,
				Amount:                           TxnBody.Amount,
				Price:                            TxnBody.Price,
				}

			result = append(result, temp)
			fmt.Println("-----------------",result)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(result)
			return 
		}
		if err != nil {
			if(response.Message!=""&&response.Code!=0){
			logrus.Error(response.Message)
			w.WriteHeader(response.Code)
			 json.NewEncoder(w).Encode(response)
			}else{
				logrus.Error("No Transactions Found in Gateway DataStore")
				w.WriteHeader(http.StatusNoContent)
				json.NewEncoder(w).Encode(model.Error{Code:http.StatusNoContent,Message:"No Transactions Found in Gateway DataStore"})
			}
		}
}
func UpdateSellingStatus(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	key1, error := r.URL.Query()["currentPK"]

	if !error || len(key1[0]) < 1 {
		logrus.Error("Url Parameter 'perPage' is missing")
		return
	}

	key2, error := r.URL.Query()["previousPK"]

	if !error || len(key2[0]) < 1 {
		logrus.Error("Url Parameter 'page' is missing")
		return
	}

	key3, error := r.URL.Query()["txnHash"]

	if !error || len(key3[0]) < 1 {
		logrus.Error("Url Parameter 'perPage' is missing")
		return
	}

	key4, error := r.URL.Query()["sellingStatus"]

	if !error || len(key4[0]) < 1 {
		logrus.Error("Url Parameter 'page' is missing")
		return
	}
	currentPK := key1[0]
	previousPK := key2[0]
	txnHash := key3[0]
	sellingStatus := key4[0]

	fmt.Println(currentPK,previousPK,txnHash,sellingStatus)

	object := dao.Connection{}
	//get the current document
	_, err1 := object.GetNFTByNFTTxn(txnHash).Then(func(data interface{}) interface{}{
		selection := data.(model.MarketPlaceNFT)
		fmt.Println("----------------------------------", selection)
		err2 := object.UpdateSellingStatus(selection, currentPK, previousPK, sellingStatus)
		if err2 != nil{
			w.Header().Set("Content-Type", "application/json;")
			w.WriteHeader(http.StatusBadRequest)
			result := apiModel.SubmitXDRSuccess{
				Status: "Error when updating the selling status",
			}
			json.NewEncoder(w).Encode(result)
		}else{
			w.Header().Set("Content-Type", "application/json;")
			w.WriteHeader(http.StatusOK)
			result := apiModel.SubmitXDRSuccess{
				Status: "Selling status updated successfully",
			}
		json.NewEncoder(w).Encode(result)
		}
		return data
	}).Await()
	
	if err1 != nil{
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error when fetching the protocol from Datastore or protocol does not exists in the Datastore",
		}
		fmt.Println(err1)
		json.NewEncoder(w).Encode(result)
	}
}

//ecfaa63f84ca9345a2220e6c21e3ed06e0323b060a335523304f95455377f4b6