package businessFacades

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/authentication"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/configs"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/nft/stellar"
	"github.com/dileepaj/tracified-gateway/nft/stellar/accounts"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

/*
MintNFTStellar
@desc - Call the IssueNft method and store new NFT details in the DB
@params - ResponseWriter,Request
*/
var dt = time.Now()

func MintNFTStellar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	dt := time.Now()
	var TrustLineResponseNFT model.TrustLineResponseNFT
	var NFTcollectionObj model.NFTWithTransaction
	var MarketplaceNFTNFTcollectionObj model.MarketPlaceNFT
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&TrustLineResponseNFT)
	if err != nil {
		log.Println(err)
	}
	if TrustLineResponseNFT.IssuerPublicKey != "" && TrustLineResponseNFT.TrustLineCreatedAt != "" && TrustLineResponseNFT.DistributorPublickKey != "" && TrustLineResponseNFT.Asset_code != "" && TrustLineResponseNFT.NFTURL != "" && TrustLineResponseNFT.Successfull {
		var NFTtxnhash, NftContent, err = stellar.IssueNft(TrustLineResponseNFT.IssuerPublicKey, TrustLineResponseNFT.DistributorPublickKey, TrustLineResponseNFT.Asset_code, TrustLineResponseNFT.NFTURL)
		if err == nil {
			NFTcollectionObj = model.NFTWithTransaction{
				Identifier:                       TrustLineResponseNFT.IssuerPublicKey,
				Categories:                       TrustLineResponseNFT.Categories,
				Collection:                       TrustLineResponseNFT.Collection,
				ImageBase64:                      TrustLineResponseNFT.NFTURL,
				NFTTXNhash:                       NFTtxnhash,
				NftTransactionExistingBlockchain: "Stellar",
				NftIssuingBlockchain:             TrustLineResponseNFT.NFTBlockChain,
				Timestamp:                        dt.Format("01-02-2006 15:04:05"),
				NftURL:                           TrustLineResponseNFT.NFTLinks,
				NftContentName:                   TrustLineResponseNFT.Asset_code,
				NftContent:                       NftContent,
				CuurentIssuerPK:                  TrustLineResponseNFT.IssuerPublicKey,
				InitialDistributorPublickKey:     TrustLineResponseNFT.DistributorPublickKey,
				InitialIssuerPK:                  TrustLineResponseNFT.IssuerPublicKey,
				MainAccountPK:                    commons.GoDotEnvVariable("NFTSTELLARISSUERPUBLICKEYK"),
				TrustLineCreatedAt:               TrustLineResponseNFT.TrustLineCreatedAt,
				Description:                      TrustLineResponseNFT.Description,
				Copies:                           TrustLineResponseNFT.Copies,
				NFTArtistName:                    TrustLineResponseNFT.ArtistName,
				NFTArtistURL:                     TrustLineResponseNFT.ArtistLink,
			}

			MarketplaceNFTNFTcollectionObj = model.MarketPlaceNFT{
				Identifier:                       TrustLineResponseNFT.IssuerPublicKey,
				Categories:                       TrustLineResponseNFT.Categories,
				Collection:                       TrustLineResponseNFT.Collection,
				ImageBase64:                      TrustLineResponseNFT.NFTURL,
				NFTTXNhash:                       NFTtxnhash,
				NftTransactionExistingBlockchain: "Stellar",
				NftIssuingBlockchain:             TrustLineResponseNFT.NFTBlockChain,
				Timestamp:                        dt.Format("01-02-2006 15:04:05"),
				NftURL:                           TrustLineResponseNFT.NFTLinks,
				NftContentName:                   TrustLineResponseNFT.Asset_code,
				NftContent:                       NftContent,
				NFTArtistName:                    TrustLineResponseNFT.ArtistName,
				NFTArtistURL:                     TrustLineResponseNFT.ArtistLink,
				Description:                      TrustLineResponseNFT.Description,
				Copies:                           TrustLineResponseNFT.Copies,
				OriginPK:                         TrustLineResponseNFT.DistributorPublickKey,
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
				TDPTxnHash:         TrustLineResponseNFT.NFTURL,
				NFTName:            TrustLineResponseNFT.Asset_code,
				NFTIssuerPublicKey: TrustLineResponseNFT.IssuerPublicKey,
			}
			object := dao.Connection{}
			err1, err2 := object.InsertStellarNFT(NFTcollectionObj, MarketplaceNFTNFTcollectionObj)
			if err1 != nil && err2 != nil {
				log.Error("NFT not inserted : ", err1, err2)
			}
			if err1 == nil && err2 != nil {
				log.Error("NFT not inserted into StellarNFT Collection : ", err2)
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

func MintWalletNFTStellar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if configs.JWTAuthenticationEnabledForMintingWalletNFT {
		permissionStatus := authentication.WalletUserHasPermissionToMint(r.Header.Get("Authorization"))
		if !permissionStatus.Status {
			commons.JSONErrorReturn(w, r, "", http.StatusUnauthorized, "Status Unauthorized")
			return
		}
	}
	dt := time.Now()
	var TrustLineResponseNFT model.TrustLineResponseNFT
	var NFTcollectionObj model.NFTWithTransaction
	var MarketplaceNFTNFTcollectionObj model.MarketPlaceNFT
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&TrustLineResponseNFT)
	if err != nil {
		log.Println(err)
	}
	if TrustLineResponseNFT.IssuerPublicKey != "" && TrustLineResponseNFT.TrustLineCreatedAt != "" && TrustLineResponseNFT.DistributorPublickKey != "" && TrustLineResponseNFT.Asset_code != "" && TrustLineResponseNFT.NFTURL != "" && TrustLineResponseNFT.Successfull {
		var NFTtxnhash, NftContent, err = stellar.IssueNft(TrustLineResponseNFT.IssuerPublicKey, TrustLineResponseNFT.DistributorPublickKey, TrustLineResponseNFT.Asset_code, TrustLineResponseNFT.NFTURL)
		if err == nil {
			NFTcollectionObj = model.NFTWithTransaction{
				Identifier:                       TrustLineResponseNFT.IssuerPublicKey,
				Categories:                       TrustLineResponseNFT.Categories,
				Collection:                       TrustLineResponseNFT.Collection,
				ImageBase64:                      TrustLineResponseNFT.NFTURL,
				NFTTXNhash:                       NFTtxnhash,
				NftTransactionExistingBlockchain: "Stellar",
				NftIssuingBlockchain:             TrustLineResponseNFT.NFTBlockChain,
				Timestamp:                        dt.Format("01-02-2006 15:04:05"),
				NftURL:                           TrustLineResponseNFT.NFTLinks,
				NftContentName:                   TrustLineResponseNFT.Asset_code,
				NftContent:                       NftContent,
				CuurentIssuerPK:                  TrustLineResponseNFT.IssuerPublicKey,
				InitialDistributorPublickKey:     TrustLineResponseNFT.DistributorPublickKey,
				InitialIssuerPK:                  TrustLineResponseNFT.IssuerPublicKey,
				MainAccountPK:                    commons.GoDotEnvVariable("NFTSTELLARISSUERPUBLICKEYK"),
				TrustLineCreatedAt:               TrustLineResponseNFT.TrustLineCreatedAt,
				Description:                      TrustLineResponseNFT.Description,
				Copies:                           TrustLineResponseNFT.Copies,
				NFTArtistName:                    TrustLineResponseNFT.ArtistName,
				NFTArtistURL:                     TrustLineResponseNFT.ArtistLink,
			}

			MarketplaceNFTNFTcollectionObj = model.MarketPlaceNFT{
				Identifier:                       TrustLineResponseNFT.IssuerPublicKey,
				Categories:                       TrustLineResponseNFT.Categories,
				Collection:                       TrustLineResponseNFT.Collection,
				ImageBase64:                      TrustLineResponseNFT.NFTURL,
				NFTTXNhash:                       NFTtxnhash,
				NftTransactionExistingBlockchain: "Stellar",
				NftIssuingBlockchain:             TrustLineResponseNFT.NFTBlockChain,
				Timestamp:                        dt.Format("01-02-2006 15:04:05"),
				NftURL:                           TrustLineResponseNFT.NFTLinks,
				NftContentName:                   TrustLineResponseNFT.Asset_code,
				NftContent:                       NftContent,
				NFTArtistName:                    TrustLineResponseNFT.ArtistName,
				NFTArtistURL:                     TrustLineResponseNFT.ArtistLink,
				Description:                      TrustLineResponseNFT.Description,
				Copies:                           TrustLineResponseNFT.Copies,
				OriginPK:                         TrustLineResponseNFT.DistributorPublickKey,
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
				TDPTxnHash:         TrustLineResponseNFT.NFTURL,
				NFTName:            TrustLineResponseNFT.Asset_code,
				NFTIssuerPublicKey: TrustLineResponseNFT.IssuerPublicKey,
			}
			object := dao.Connection{}
			err1, err2 := object.InsertStellarNFT(NFTcollectionObj, MarketplaceNFTNFTcollectionObj)
			if err1 != nil && err2 != nil {
				log.Error("NFT not inserted : ", err1, err2)
			}
			if err1 == nil && err2 != nil {
				log.Error("NFT not inserted into StellarNFT Collection : ", err2)
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

/*
RetriveNFTByStatusAndPK
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
			Collection:                       TxnBody.Collection,
			Categories:                       TxnBody.Categories,
			ImageBase64:                      TxnBody.ImageBase64,
			NFTTXNhash:                       TxnBody.NFTTXNhash,
			NftTransactionExistingBlockchain: "Stellar",
			NftIssuingBlockchain:             TxnBody.NftIssuingBlockchain,
			Timestamp:                        dt.Format("01-02-2006 15:04:05"),
			NftURL:                           TxnBody.NftURL,
			NftContentName:                   TxnBody.NftContentName,
			NftContent:                       TxnBody.NftContent,
			NFTArtistName:                    TxnBody.NFTArtistName,
			NFTArtistURL:                     TxnBody.NFTArtistURL,
			TrustLineCreatedAt:               TxnBody.TrustLineCreatedAt,
			Description:                      TxnBody.Description,
			Copies:                           TxnBody.Copies,
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

/*
UpdateSellingStatus
@desc - Update the selling status in the collection when selling the NFT
@params - ResponseWriter,Request
*/
func UpdateSellingStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	key1, error := r.URL.Query()["Price"]

	if !error || len(key1[0]) < 1 {
		logrus.Error("Url Parameter 'Price' is missing")
		return
	}

	key2, error := r.URL.Query()["Status"]

	if !error || len(key2[0]) < 1 {
		logrus.Error("Url Parameter 'Status' is missing")
		return
	}

	key3, error := r.URL.Query()["Amount"]

	if !error || len(key3[0]) < 1 {
		logrus.Error("Url Parameter 'Amount' is missing")
		return
	}

	key4, error := r.URL.Query()["NFTTxnHash"]

	if !error || len(key4[0]) < 1 {
		logrus.Error("Url Parameter 'NFTTxnHash' is missing")
		return
	}
	price := key1[0]
	sellingStatus := key2[0]
	amount := key3[0]
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

/*
UpdateBuyingStatus
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

/*
CreateNFTIssuerAccount
@desc - Create issuer accounts, add credentials to the DB and send the PK as the response
@params - ResponseWriter,Request
*/
func GetNFTIssuerAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var NFTIssuerPK, EncodedNFTIssuerSK, encSK, err = accounts.CreateIssuerAccount()
	if err != nil && NFTIssuerPK == "" && EncodedNFTIssuerSK == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	} else {
		var NFTKeys = model.NFTKeys{
			PublicKey: NFTIssuerPK,
			SecretKey: encSK,
		}
		//adding the credentials to the DB
		object := dao.Connection{}
		err := object.InsertStellarNFTKeys(NFTKeys)
		if err != nil {
			log.Println(err)
		}
		//send the response
		result := model.NFTIssuerAccount{
			NFTIssuerPK: NFTIssuerPK,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

func GetLastNFTbyIdentifier(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	object := dao.Connection{}
	p := object.GetLastNFTbyInitialDistributorPK(vars["InitialDistributorPK"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.MarketPlaceNFT)
		log.Println(result)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "InitialDistributorPK Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()

}

func RetrieveStellarTxn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	object := dao.Connection{}
	p := object.GetNFTTxnForStellar(vars["ImageBase64"], vars["blockchain"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.NFTWithTransaction)
		res := model.StellarMintTXN{NFTTxnHash: result.NFTTXNhash}
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

func FundAndGetAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	key1, error := r.URL.Query()["publickey"]

	if !error || len(key1[0]) < 1 {
		logrus.Error("Url Parameter 'publickey' is missing")
		return
	}
	publickey := key1[0]
	var txn, err = accounts.FundAccount(publickey)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error when updating the buying status",
		}
		json.NewEncoder(w).Encode(result)
	} else {
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusOK)
		result := model.PublicKey{
			PublicKey: txn,
		}
		json.NewEncoder(w).Encode(result)
	}

}

func GetSponsorAccountXDR(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	key1, error := r.URL.Query()["publickey"]

	if !error || len(key1[0]) < 1 {
		logrus.Error("Url Parameter 'publickey' is missing")
		return
	}

	key2, error := r.URL.Query()["nftName"]

	if !error || len(key2[0]) < 1 {
		logrus.Error("Url Parameter 'nftName' is missing")
		return
	}

	key3, error := r.URL.Query()["issuer"]

	if !error || len(key3[0]) < 1 {
		logrus.Error("Url Parameter 'issuer' is missing")
		return
	}
	publickey := key1[0]
	nftname := key2[0]
	issuer := key3[0]

	var txn, err = stellar.SponsorCreateAccount(publickey, nftname, issuer)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error when updating the buying status",
		}
		json.NewEncoder(w).Encode(result)
	} else {
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusOK)
		result := model.XDRRuri{
			XDR: txn,
		}
		logrus.Println("XDR been passed to frontend : ", result)
		json.NewEncoder(w).Encode(result)
	}

}

func SponsorAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	key1, error := r.URL.Query()["publickey"]

	if !error || len(key1[0]) < 1 {
		logrus.Error("Url Parameter 'publickey' is missing")
		return
	}

	publickey := key1[0]

	var txn, err = stellar.SponsorAccount(publickey)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error when updating the buying status",
		}
		json.NewEncoder(w).Encode(result)
	} else {
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusOK)
		result := model.XDRRuri{
			XDR: txn,
		}
		logrus.Println("XDR been passed to frontend : ", result)
		json.NewEncoder(w).Encode(result)
	}

}

func GetSponsorTrustXDR(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	key1, error := r.URL.Query()["publickey"]

	if !error || len(key1[0]) < 1 {
		logrus.Error("Url Parameter 'publickey' is missing")
		return
	}

	key2, error := r.URL.Query()["nftName"]

	if !error || len(key2[0]) < 1 {
		logrus.Error("Url Parameter 'nftName' is missing")
		return
	}

	key3, error := r.URL.Query()["issuer"]

	if !error || len(key3[0]) < 1 {
		logrus.Error("Url Parameter 'issuer' is missing")
		return
	}
	publickey := key1[0]
	nftname := key2[0]
	issuer := key3[0]

	var txn, err = stellar.SponsorTrust(publickey, nftname, issuer)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error when updating the buying status",
		}
		json.NewEncoder(w).Encode(result)
	} else {
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusOK)
		result := model.XDRRuri{
			XDR: txn,
		}
		logrus.Println("XDR been passed to frontend : ", result)
		json.NewEncoder(w).Encode(result)
	}

}
