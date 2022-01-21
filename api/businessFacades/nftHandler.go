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
	"github.com/dileepaj/tracified-gateway/nft/stellar/accounts"
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
		panic(err)
	}
	if TrustLineResponseNFT.IssuerPublicKey != "" && TrustLineResponseNFT.TrustLineCreatedAt != "" && TrustLineResponseNFT.DistributorPublickKey != "" && TrustLineResponseNFT.Asset_code != "" && TrustLineResponseNFT.TDPtxnhash != "" && TrustLineResponseNFT.Successfull {
		var NFTtxnhash, issuerPK, NftContent, err = stellar.IssueNft(TrustLineResponseNFT.IssuerPublicKey, TrustLineResponseNFT.DistributorPublickKey, TrustLineResponseNFT.Asset_code, TrustLineResponseNFT.TDPtxnhash)
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
				SellingStatus:					  "FORSELL",
			}

			NFTCeactedResponse := model.NFTCreactedResponse{
				NFTTxnHash: NFTtxnhash,
				TDPTxnHash: TrustLineResponseNFT.TDPtxnhash,
				NFTName:    TrustLineResponseNFT.Asset_code,
				NFTIssuerPublicKey:TrustLineResponseNFT.IssuerPublicKey,
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
	var SellingStatus="FORSELL"
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
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
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

//Create issuer accounts, add credentials to the DB and send the PK as the response
func CreateNFTIssuerAccount(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var NFTIssuerPK,EncodedNFTIssuerSK,err = accounts.CreateIssuerAccount()
	if (err!=nil && NFTIssuerPK=="" && EncodedNFTIssuerSK==""){
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}else{
		var NFTKeys= model.NFTKeys{
			PublicKey:NFTIssuerPK,
			SecretKey: EncodedNFTIssuerSK,
		}
		//adding the credentials to the DB
		object := dao.Connection{}
		err:= object.InsertStellarNFTKeys(NFTKeys)
		if err!=nil{
			panic(err)
		}
		//send the response
		result := model.NFTIssuerAccount{
			NFTIssuerPK: NFTIssuerPK,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}