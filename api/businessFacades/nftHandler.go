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
	"github.com/dileepaj/tracified-gateway/nft/stellar/authForIssuer"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	//"go.mongodb.org/mongo-driver/x/mongo/driver/auth"
)

/*MintNFTStellar
@desc - Call the IssueNft method and store new NFT details in the DB
@params - ResponseWriter,Request
*/
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
		var NFTtxnhash, NftContent, err = stellar.IssueNft(TrustLineResponseNFT.IssuerPublicKey, TrustLineResponseNFT.DistributorPublickKey, TrustLineResponseNFT.Asset_code, TrustLineResponseNFT.TDPtxnhash)
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
				MainAccountPK:                    commons.GoDotEnvVariable("NFTSTELLARISSUERPUBLICKEYK"),
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
				PreviousOwnerNFTPK:               "TRACIFIED",
				CurrentOwnerNFTPK:                TrustLineResponseNFT.DistributorPublickKey,
				SellingStatus:                    "NOTFORSALE",
			}

			NFTCeactedResponse := model.NFTCreactedResponse{
				NFTTxnHash:         NFTtxnhash,
				TDPTxnHash:         TrustLineResponseNFT.TDPtxnhash,
				NFTName:            TrustLineResponseNFT.Asset_code,
				NFTIssuerPublicKey: TrustLineResponseNFT.IssuerPublicKey,
			}
			object := dao.Connection{}
			err1, err2 := object.InsertStellarNFT(NFTcollectionObj, MarketplaceNFTNFTcollectionObj)
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

/*RetriveNFTByStatusAndPK
@desc - Call the GetAllSellingNFTStellar_Paginated method and get all the NFT relevent to the selling status and Public key
@params - ResponseWriter,Request
*/
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

	qdata, err := object.GetAllSellingNFTStellar_Paginated(sellingstatus[0], distributorPK[0]).Then(func(data interface{}) interface{} {
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
		response = model.Error{Code: http.StatusNoContent, Message: "Selling NFTs are not found in gateway datastore"}
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
			InitialIssuerPK:                  TxnBody.InitialIssuerPK,
			InitialDistributorPK:             TxnBody.InitialDistributorPK,
			MainAccountPK:                    TxnBody.MainAccountPK,
		}

		result = append(result, temp)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

/*UpdateSellingStatus
@desc - Update the selling status in the collection when selling the NFT
@params - ResponseWriter,Request
*/
func UpdateSellingStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	key1, error := r.URL.Query()["sellingStatus"]

	if !error || len(key1[0]) < 1 {
		logrus.Error("Url Parameter 'sellingStatus' is missing")
		return
	}

	key2, error := r.URL.Query()["amount"]

	if !error || len(key2[0]) < 1 {
		logrus.Error("Url Parameter 'amount' is missing")
		return
	}

	key3, error := r.URL.Query()["price"]

	if !error || len(key3[0]) < 1 {
		logrus.Error("Url Parameter 'price' is missing")
		return
	}

	key4, error := r.URL.Query()["nfthash"]

	if !error || len(key4[0]) < 1 {
		logrus.Error("Url Parameter 'nfthash' is missing")
		return
	}
	sellingStatus := key1[0]
	amount := key2[0]
	price := key3[0]
	nfthash := key4[0]
	object := dao.Connection{}
	_, err1 := object.GetNFTByNFTTxn(nfthash).Then(func(data interface{}) interface{} {
		selection := data.(model.MarketPlaceNFT)
		err2 := object.UpdateSellingStatus(selection, sellingStatus, amount, price)
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
			Status: "Error when fetching the NFT details from Datastore or NFT hash does not exists in the Datastore",
		}
		log.Println(err1)
		json.NewEncoder(w).Encode(result)
	}
}

/*UpdateBuyingStatus
@desc - Update the selling status in the collection when buying the NFT
@params - ResponseWriter,Request
*/
func UpdateBuyingStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	key1, error := r.URL.Query()["sellingStatus"]

	if !error || len(key1[0]) < 1 {
		logrus.Error("Url Parameter 'sellingStatus' is missing")
		return
	}

	key2, error := r.URL.Query()["currentPK"]

	if !error || len(key2[0]) < 1 {
		logrus.Error("Url Parameter 'currentPK' is missing")
		return
	}

	key3, error := r.URL.Query()["previousPK"]

	if !error || len(key3[0]) < 1 {
		logrus.Error("Url Parameter 'previousPK' is missing")
		return
	}

	key4, error := r.URL.Query()["nfthash"]

	if !error || len(key4[0]) < 1 {
		logrus.Error("Url Parameter 'nfthash' is missing")
		return
	}
	sellingStatus := key1[0]
	currentPK := key2[0]
	previousPK := key3[0]
	nfthash := key4[0]
	object := dao.Connection{}
	_, err1 := object.GetNFTByNFTTxn(nfthash).Then(func(data interface{}) interface{} {
		selection := data.(model.MarketPlaceNFT)
		err2 := object.UpdateBuyingStatus(selection, sellingStatus, currentPK, previousPK)
		if err2 != nil {
			w.Header().Set("Content-Type", "application/json;")
			w.WriteHeader(http.StatusBadRequest)
			result := apiModel.SubmitXDRSuccess{
				Status: "Error when updating the buying status",
			}
			json.NewEncoder(w).Encode(result)
		} else {
			w.Header().Set("Content-Type", "application/json;")
			w.WriteHeader(http.StatusOK)
			result := apiModel.SubmitXDRSuccess{
				Status: "Buying status updated successfully",
			}
			json.NewEncoder(w).Encode(result)
		}
		return data
	}).Await()
	if err1 != nil {
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error when fetching the NFT details from Datastore or NFT hash does not exists in the Datastore",
		}
		log.Println(err1)
		json.NewEncoder(w).Encode(result)
	}
}

/*CreateNFTIssuerAccount
@desc - Create issuer accounts, add credentials to the DB and send the PK as the response
@params - ResponseWriter,Request
*/
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
		//send currentPK to authorize
		err1 := authForIssuer.SetAuth(NFTIssuerPK)
		if err1 != nil {
			panic(err1)
		}
		//send the response
		result := model.NFTIssuerAccount{
			NFTIssuerPK: NFTIssuerPK,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

/*UnlockNFTIssuerAccount
@desc - Unlock the issuer account
@params - ResponseWriter,Request
*/
// func UnlockNFTIssuerAccount(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	key1, error := r.URL.Query()["currentPK"]

// 	if !error || len(key1[0]) < 1 {
// 		logrus.Error("Url Parameter 'perPage' is missing")
// 		return
// 	}

// 	currentPK := key1[0]
// 	log.Println(currentPK)

// 	err := accounts.UnlockIssuingAccount(currentPK)
// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json;")
// 		w.WriteHeader(http.StatusBadRequest)
// 		result := apiModel.SubmitXDRSuccess{
// 			Status: "Error when unlocking account",
// 		}
// 		log.Println(err)
// 		json.NewEncoder(w).Encode(result)
// 	} else {
// 		w.Header().Set("Content-Type", "application/json;")
// 		w.WriteHeader(http.StatusOK)
// 		result := apiModel.SubmitXDRSuccess{
// 			Status: "Current Issuer Account unlocked successfully",
// 		}
// 		json.NewEncoder(w).Encode(result)
// 	}

// }

/*LockNFTIssuerAccount
@desc - Lock the issuer account
@params - ResponseWriter,Request
*/
// func LockNFTIssuerAccount(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	key1, error := r.URL.Query()["currentPK"]

// 	if !error || len(key1[0]) < 1 {
// 		logrus.Error("Url Parameter 'perPage' is missing")
// 		return
// 	}

// 	currentPK := key1[0]
// 	log.Println(currentPK)
// 	err := accounts.LockIssuingAccount(currentPK)
// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json;")
// 		w.WriteHeader(http.StatusBadRequest)
// 		result := apiModel.SubmitXDRSuccess{
// 			Status: "Error when locking account",
// 		}
// 		log.Println(err)
// 		json.NewEncoder(w).Encode(result)
// 	} else {
// 		w.Header().Set("Content-Type", "application/json;")
// 		w.WriteHeader(http.StatusOK)
// 		result := apiModel.SubmitXDRSuccess{
// 			Status: "Current Issuer Account locked successfully",
// 		}
// 		json.NewEncoder(w).Encode(result)
// 	}

// }

func SetAuthTrust(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	key1, error := r.URL.Query()["currentPK"]
	if !error || len(key1[0]) < 1 {
		logrus.Error("Url Parameter 'perPage' is missing")
		return
	}
	fmt.Println("------------------", key1)

	key2, error := r.URL.Query()["trustor"]

	if !error || len(key2[0]) < 1 {
		logrus.Error("Url Parameter 'trustor' is missing")
		return
	}
	fmt.Println("------------------", key2)

	key3, error := r.URL.Query()["code"]

	if !error || len(key3[0]) < 1 {
		logrus.Error("Url Parameter 'code' is missing")
		return
	}
	fmt.Println("------------------", key3)

	key4, error := r.URL.Query()["XDR"]

	if !error || len(key4[0]) < 1 {
		logrus.Error("Url Parameter 'XDR' is missing")
		return
	}
	fmt.Println("------------------", key4)

	currentPK := key1[0]
	trustor := key2[0]
	code := key3[0]
	paymentXDR := key4[0]

	log.Println(currentPK)
	log.Println(trustor)
	log.Println(code)
	log.Println(paymentXDR)

	royaltyXDR := model.RoyaltyXDR{
		XDR: paymentXDR,
	}

	var result = authForIssuer.CheckPayment(royaltyXDR)
	if !result {
		log.Println("Error while Checking the royalty payment")
		fmt.Println(result)
	} else {
		fmt.Println(result)

		//calling the autorization of the trustline

		// err := authForIssuer.AuthTrust(NFTBuyTrustlineAuth.CurrentPK, NFTBuyTrustlineAuth.Trustor, NFTBuyTrustlineAuth.Code)
		// if err != nil {
		// 	w.Header().Set("Content-Type", "application/json;")
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	result := apiModel.SubmitXDRSuccess{
		// 		Status: "Error when unlocking account",
		// 	}
		// 	log.Println(err)
		// 	json.NewEncoder(w).Encode(result)
		// } else {
		// 	w.Header().Set("Content-Type", "application/json;")
		// 	w.WriteHeader(http.StatusOK)
		// 	result := apiModel.SubmitXDRSuccess{
		// 		Status: "Current Issuer Account unlocked successfully",
		// 	}
		// 	json.NewEncoder(w).Encode(result)
		// }
	}
}
