package businessFacades

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"
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
	fmt.Println("---------33333333333333333=-====================----",TrustLineResponseNFT.IssuerPublicKey, TrustLineResponseNFT.DistributorPublickKey, TrustLineResponseNFT.Asset_code, TrustLineResponseNFT.TDPtxnhash)
	if TrustLineResponseNFT.IssuerPublicKey != "" && TrustLineResponseNFT.TrustLineCreatedAt != "" && TrustLineResponseNFT.DistributorPublickKey != "" && TrustLineResponseNFT.Asset_code != "" && TrustLineResponseNFT.TDPtxnhash != "" && TrustLineResponseNFT.Successfull {
		var NFTtxnhash, NftContent, err = stellar.IssueNft(TrustLineResponseNFT.IssuerPublicKey, TrustLineResponseNFT.DistributorPublickKey, TrustLineResponseNFT.Asset_code, TrustLineResponseNFT.TDPtxnhash)
		fmt.Println("---------33333333333333333=-====================----",TrustLineResponseNFT.IssuerPublicKey, TrustLineResponseNFT.DistributorPublickKey, TrustLineResponseNFT.Asset_code, TrustLineResponseNFT.TDPtxnhash)
		if err == nil {
			NFTcollectionObj = model.NFTWithTransaction{
				Identifier:                       TrustLineResponseNFT.Identifier,
				TxnType:                          "TDP",
				TDPID:                            TrustLineResponseNFT.TDPID,
				TDPTxnHash:                       TrustLineResponseNFT.TDPtxnhash,
				DataHash:                         TrustLineResponseNFT.TDPtxnhash,
				NFTTXNhash:                       NFTtxnhash,
				NftTransactionExistingBlockchain: "Stellar",
				NftIssuingBlockchain:             TrustLineResponseNFT.NFTBlockChain,
				Timestamp:                        "00-00-00",
				NftAssetName:                     TrustLineResponseNFT.Asset_code,
				NftContentName:                   TrustLineResponseNFT.Asset_code,
				NftContent:                       NftContent,
				InitialDistributorPublickKey:     TrustLineResponseNFT.DistributorPublickKey,
				InitialIssuerPK:                  TrustLineResponseNFT.IssuerPublicKey,
				MainAccountPK:					  commons.GoDotEnvVariable("NFTSTELLARISSUERPUBLICKEYK"),	
				TrustLineCreatedAt:               TrustLineResponseNFT.TrustLineCreatedAt,
				ProductName:                      TrustLineResponseNFT.ProductName,
			}

			MarketplaceNFTNFTcollectionObj = model.MarketPlaceNFT{
				Identifier:                       TrustLineResponseNFT.Identifier,
				ProductName:                      TrustLineResponseNFT.ProductName,
				TxnType:                          "TDP",
				TDPID:                            TrustLineResponseNFT.TDPID,
				TDPTxnHash:                       TrustLineResponseNFT.TDPtxnhash,
				DataHash:                         TrustLineResponseNFT.TDPtxnhash,
				NFTTXNhash:                       NFTtxnhash,
				NftTransactionExistingBlockchain: "Stellar",
				NftIssuingBlockchain:             TrustLineResponseNFT.NFTBlockChain,
				Timestamp:                        "00-00-00",
				NftAssetName:                     TrustLineResponseNFT.Asset_code,
				NftContentName:                   TrustLineResponseNFT.Asset_code,
				NftContent:                       NftContent,
				InitialDistributorPK:             TrustLineResponseNFT.DistributorPublickKey,
				InitialIssuerPK:                  TrustLineResponseNFT.IssuerPublicKey,
				MainAccountPK:                    commons.GoDotEnvVariable("NFTSTELLARISSUERPUBLICKEYK"),
				TrustLineCreatedAt:               TrustLineResponseNFT.TrustLineCreatedAt,
				PreviousOwnerNFTPK:               "",
				CurrentOwnerNFTPK:                TrustLineResponseNFT.DistributorPublickKey,
				SellingStatus:                    "NOTFORSALE",
			}

			NFTCeactedResponse := model.NFTCreactedResponse{
				NFTTxnHash:         NFTtxnhash,
				TDPTxnHash:         TrustLineResponseNFT.TDPtxnhash,
				NFTName:            TrustLineResponseNFT.Asset_code,
				NFTIssuerPublicKey: TrustLineResponseNFT.IssuerPublicKey,
			}
			//after that call the insertTODB NFT Details to NFT  (gateway(NFTstellar collection)
			object := dao.Connection{}
			err1, err2 := object.InsertStellarNFT(NFTcollectionObj, MarketplaceNFTNFTcollectionObj)
			if err1 != nil && err2 != nil {
				fmt.Println("NFT not inserted : ", err1, err2)
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

func RetriveNFTByStatusAndPK(w http.ResponseWriter, r *http.Request) {
	var response model.Error
	var result []model.MarketPlaceNFT
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	sellingstatus, error := r.URL.Query()["sellingstatus"]

	if !error || len(sellingstatus[0]) < 1 {
		logrus.Error("Url Parameter 'sellingstatus' is missing")
		return
	}

	distributorPK, error := r.URL.Query()["distributorPK"]

	if !error || len(distributorPK[0]) < 1 {
		logrus.Error("Url Parameter 'distributorPK' is missing")
		return
	}

	object := dao.Connection{}

	qdata, err := object.GetAllSellingNFTStellar_Paginated(sellingstatus[0],distributorPK[0]).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		logrus.Error("Unable to connect gateway datastore")
		w.WriteHeader(http.StatusNotFound)
		response = model.Error{Code: http.StatusNotFound, Message: "Unable to connect gateway datastore"}
		json.NewEncoder(w).Encode(response)
		return
	}
	if qdata == nil {
		logrus.Error("Selling NFTs are not found in gateway datastore")
		w.WriteHeader(http.StatusNoContent)
		response = model.Error{Code: http.StatusNoContent, Message: "Identifier is not found in gateway datastore"}
		json.NewEncoder(w).Encode(response)
		return
	}
	res := qdata.(model.MarketPlaceNFTTrasactionWithCount)
	for _, TxnBody := range res.MarketPlaceNFTItems {

		temp := model.MarketPlaceNFT{
			Identifier:                       TxnBody.Identifier,
			TxnType:                          "TDP",
			TDPID:                            TxnBody.TDPID,
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
			ProductName:                      TxnBody.ProductName,
			PreviousOwnerNFTPK:               TxnBody.PreviousOwnerNFTPK,
			CurrentOwnerNFTPK:                TxnBody.CurrentOwnerNFTPK,
			OriginPK:                         TxnBody.OriginPK,
			SellingStatus:                    TxnBody.SellingStatus,
			Amount:                           TxnBody.Amount,
			Price:                            TxnBody.Price,
			InitialIssuerPK: TxnBody.InitialIssuerPK,
			InitialDistributorPK: TxnBody.InitialDistributorPK,
			MainAccountPK: TxnBody.MainAccountPK,
		}

		result = append(result, temp)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}


func UpdateSellingStatus(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	key1, error := r.URL.Query()["currentPK"]

	if !error || len(key1[0]) < 1 {
		logrus.Error("Url Parameter 'currentPK' is missing")
		return
	}

	key2, error := r.URL.Query()["previousPK"]

	if !error || len(key2[0]) < 1 {
		logrus.Error("Url Parameter 'previousPK' is missing")
		return
	}

	key3, error := r.URL.Query()["txnHash"]

	if !error || len(key3[0]) < 1 {
		logrus.Error("Url Parameter 'txnHash' is missing")
		return
	}

	key4, error := r.URL.Query()["sellingStatus"]

	if !error || len(key4[0]) < 1 {
		logrus.Error("Url Parameter 'sellingStatus' is missing")
		return
	}
	currentPK := key1[0]
	previousPK := key2[0]
	txnHash := key3[0]
	sellingStatus := key4[0]

	fmt.Println(currentPK, previousPK, txnHash, sellingStatus)

	object := dao.Connection{}
	//get the current document
	_, err1 := object.GetNFTByNFTTxn(txnHash).Then(func(data interface{}) interface{} {
		selection := data.(model.MarketPlaceNFT)
		fmt.Println("----------------------------------", selection)
		err2 := object.UpdateSellingStatus(selection, currentPK, previousPK, sellingStatus)
		if err2 != nil {
			w.Header().Set("Content-Type", "application/json;")
			w.WriteHeader(http.StatusBadRequest)
			result := apiModel.SubmitXDRSuccess{
				Status: "Error when updating the selling status",
			}
			json.NewEncoder(w).Encode(result)
		} else {
			w.Header().Set("Content-Type", "application/json;")
			w.WriteHeader(http.StatusOK)
			result := apiModel.SubmitXDRSuccess{
				Status: "Selling status updated successfully",
			}
			json.NewEncoder(w).Encode(result)
		}
		return data
	}).Await()

	if err1 != nil {
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
func CreateNFTIssuerAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var NFTIssuerPK, EncodedNFTIssuerSK, err = accounts.CreateIssuerAccount()
	if err != nil && NFTIssuerPK == "" && EncodedNFTIssuerSK == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	} else {
		var NFTKeys = model.NFTKeys{
			PublicKey: NFTIssuerPK,
			SecretKey: EncodedNFTIssuerSK,
		}
		//adding the credentials to the DB
		object := dao.Connection{}
		err := object.InsertStellarNFTKeys(NFTKeys)
		if err != nil {
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

func UnlockNFTIssuerAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	key1, error := r.URL.Query()["currentPK"]

	if !error || len(key1[0]) < 1 {
		logrus.Error("Url Parameter 'perPage' is missing")
		return
	}

	currentPK := key1[0]
	fmt.Println(currentPK)

	err := accounts.UnlockIssuingAccount(currentPK)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error when unlocking account",
		}
		fmt.Println(err)
		json.NewEncoder(w).Encode(result)
	} else {
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusOK)
		result := apiModel.SubmitXDRSuccess{
			Status: "Current Issuer Account unlocked successfully",
		}
		json.NewEncoder(w).Encode(result)
	}

}

//nfttxnhash needs to be passed to the retrieve function and
//then the pk from there brought back to the handler and passed to the unlock function
func LockNFTIssuerAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	key1, error := r.URL.Query()["currentPK"]

	if !error || len(key1[0]) < 1 {
		logrus.Error("Url Parameter 'perPage' is missing")
		return
	}

	currentPK := key1[0]
	fmt.Println(currentPK)
	err := accounts.LockIssuingAccount(currentPK)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error when locking account",
		}
		fmt.Println(err)
		json.NewEncoder(w).Encode(result)
	} else {
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusOK)
		result := apiModel.SubmitXDRSuccess{
			Status: "Current Issuer Account locked successfully",
		}
		json.NewEncoder(w).Encode(result)
	}

}
