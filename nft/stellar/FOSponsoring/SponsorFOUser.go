package fosponsoring

import (
	"fmt"
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

var tx *txnbuild.Transaction
var decrpytIssuerSecretKey string

func SponsorExisitingFOUsers(txnPayload model.TransactionData) (string, error) {
	fmt.Println("----------------test 6")
	client := commons.GetHorizonClient()
	object := dao.Connection{}

	data, err := object.GetIssuerSK(txnPayload.AccountIssuer).Then(func(data interface{}) interface{} {
		fmt.Println("data from getting account details ", data)
		return data
	}).Await()
	if err != nil {
		log.Println(err)
		return "", err
	} else {
		Keys := data.([]model.TransactionDataKeys)
		//decrypt the secret key
		decrpytIssuerSecretKey = commons.Decrypt(Keys[0].AccountIssuerSK)
		fmt.Println("decryt issuer sk ", decrpytIssuerSecretKey)
		if data == nil {
			logrus.Error("PublicKey is not found in gateway datastore")
		}
	}

	beginSponsorship := txnbuild.BeginSponsoringFutureReserves{
		SponsoredID:   txnPayload.FOUser,
		SourceAccount: txnPayload.AccountIssuer,
	}

	manageDataType := txnbuild.ManageData{
		Name:          "Type",
		Value:         []byte(txnPayload.Type),
		SourceAccount: txnPayload.AccountIssuer,
	}

	endSponsorship := txnbuild.EndSponsoringFutureReserves{
		SourceAccount: txnPayload.FOUser,
	}

	manageDataIdentifier := txnbuild.ManageData{
		Name:          "Identifier",
		Value:         []byte(txnPayload.Identifier),
		SourceAccount: txnPayload.AccountIssuer,
	}

	manageDataProductName := txnbuild.ManageData{
		Name:          "productName",
		Value:         []byte(txnPayload.ProductName),
		SourceAccount: txnPayload.AccountIssuer,
	}

	manageDataProductID := txnbuild.ManageData{
		Name:          "productId",
		Value:         []byte(txnPayload.ProductID),
		SourceAccount: txnPayload.AccountIssuer,
	}

	manageDataAppAccount := txnbuild.ManageData{
		Name:          "appAccount",
		Value:         []byte(txnPayload.AppAccount),
		SourceAccount: txnPayload.AccountIssuer,
	}

	manageDataCurrentStage := txnbuild.ManageData{
		Name:          "CurrentStage",
		Value:         []byte(txnPayload.CurrentStage),
		SourceAccount: txnPayload.AccountIssuer,
	}

	manageDataTypeName := txnbuild.ManageData{
		Name:          "TypeName",
		Value:         []byte(txnPayload.TypeName),
		SourceAccount: txnPayload.AccountIssuer,
	}

	manageDataTimeStamp := txnbuild.ManageData{
		Name:          "Timestamp",
		Value:         []byte(txnPayload.TimeStamp),
		SourceAccount: txnPayload.AccountIssuer,
	}

	manageDataDataHash := txnbuild.ManageData{
		Name:          "dataHash",
		Value:         []byte(txnPayload.DataHash),
		SourceAccount: txnPayload.AccountIssuer,
	}

	manageDataFromIdentifier := txnbuild.ManageData{
		Name:          "FromIdentifier",
		Value:         []byte(txnPayload.FromIdentifier),
		SourceAccount: txnPayload.AccountIssuer,
	}

	manageDataToIdentifier := txnbuild.ManageData{
		Name:          "ToIdentifiers",
		Value:         []byte(txnPayload.ToIdentifiers),
		SourceAccount: txnPayload.AccountIssuer,
	}

	manageDataToFromIdentifier1 := txnbuild.ManageData{
		Name:          "FromIdentifier1",
		Value:         []byte(txnPayload.FromIdentifier1),
		SourceAccount: txnPayload.AccountIssuer,
	}

	manageDataToFromIdentifier2 := txnbuild.ManageData{
		Name:          "FromIdentifier2",
		Value:         []byte(txnPayload.FromIdentifier2),
		SourceAccount: txnPayload.AccountIssuer,
	}

	manageDataToPreviousStage := txnbuild.ManageData{
		Name:          "PreviousStage",
		Value:         []byte(txnPayload.PreviousStage),
		SourceAccount: txnPayload.AccountIssuer,
	}

	accountRequest := horizonclient.AccountRequest{AccountID: txnPayload.AccountIssuer}
	sourceAccount, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Fatal(err)
	}

	if txnPayload.Type == "0" {
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        &sourceAccount,
				IncrementSequenceNum: true,
				Operations: []txnbuild.Operation{&beginSponsorship, &manageDataType, &endSponsorship, &beginSponsorship, &manageDataIdentifier, &endSponsorship,
					&beginSponsorship, &manageDataProductName, &endSponsorship, &beginSponsorship, &manageDataProductID, &endSponsorship,
					&beginSponsorship, &manageDataAppAccount, &endSponsorship, &beginSponsorship, &manageDataCurrentStage, &endSponsorship,
					&beginSponsorship, &manageDataTypeName, &endSponsorship, &beginSponsorship, &manageDataTimeStamp, &endSponsorship},
				BaseFee:       constants.MinBaseFee,
				Memo:          nil,
				Preconditions: txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
			},
		)
		if err != nil {
			log.Fatal("Error while trying to build tranaction: ", err)
		}
	} else if txnPayload.Type == "2" {
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        &sourceAccount,
				IncrementSequenceNum: true,
				Operations: []txnbuild.Operation{&beginSponsorship, &manageDataType, &endSponsorship, &beginSponsorship, &manageDataIdentifier, &endSponsorship,
					&beginSponsorship, &manageDataProductName, &endSponsorship, &beginSponsorship, &manageDataProductID, &endSponsorship,
					&beginSponsorship, &manageDataAppAccount, &endSponsorship, &beginSponsorship, &manageDataCurrentStage, &endSponsorship,
					&beginSponsorship, &manageDataTypeName, &endSponsorship, &beginSponsorship, &manageDataTimeStamp, &endSponsorship, &beginSponsorship, &manageDataDataHash, &endSponsorship},
				BaseFee:       constants.MinBaseFee,
				Memo:          nil,
				Preconditions: txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
			},
		)
		if err != nil {
			log.Fatal("Error while trying to build tranaction: ", err)
		}
	} else if txnPayload.Type == "6" {
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        &sourceAccount,
				IncrementSequenceNum: true,
				Operations: []txnbuild.Operation{&beginSponsorship, &manageDataType, &endSponsorship, &beginSponsorship, &manageDataIdentifier, &endSponsorship,
					&beginSponsorship, &manageDataProductName, &endSponsorship, &beginSponsorship, &manageDataProductID, &endSponsorship,
					&beginSponsorship, &manageDataAppAccount, &endSponsorship,
					&beginSponsorship, &manageDataTypeName, &endSponsorship, &beginSponsorship, &manageDataTimeStamp, &endSponsorship, &beginSponsorship, &manageDataFromIdentifier, &endSponsorship},
				BaseFee:       constants.MinBaseFee,
				Memo:          nil,
				Preconditions: txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
			},
		)
		if err != nil {
			log.Fatal("Error while trying to build tranaction: ", err)
		}
	} else if txnPayload.Type == "5" {
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        &sourceAccount,
				IncrementSequenceNum: true,
				Operations: []txnbuild.Operation{&beginSponsorship, &manageDataType, &endSponsorship, &beginSponsorship, &manageDataToIdentifier, &endSponsorship,
					&beginSponsorship, &manageDataProductName, &endSponsorship, &beginSponsorship, &manageDataProductID, &endSponsorship,
					&beginSponsorship, &manageDataAppAccount, &endSponsorship,
					&beginSponsorship, &manageDataTypeName, &endSponsorship, &beginSponsorship, &manageDataTimeStamp, &endSponsorship, &beginSponsorship, &manageDataFromIdentifier, &endSponsorship},
				BaseFee:       constants.MinBaseFee,
				Memo:          nil,
				Preconditions: txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
			},
		)
		if err != nil {
			log.Fatal("Error while trying to build tranaction: ", err)
		}
	} else if txnPayload.Type == "7" {
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        &sourceAccount,
				IncrementSequenceNum: true,
				Operations: []txnbuild.Operation{&beginSponsorship, &manageDataType, &endSponsorship, &beginSponsorship, &manageDataIdentifier, &endSponsorship,
					&beginSponsorship, &manageDataProductName, &endSponsorship, &beginSponsorship, &manageDataProductID, &endSponsorship,
					&beginSponsorship, &manageDataAppAccount, &endSponsorship, &beginSponsorship, &manageDataCurrentStage, &endSponsorship,
					&beginSponsorship, &manageDataTypeName, &endSponsorship, &beginSponsorship, &manageDataTimeStamp, &endSponsorship, &beginSponsorship, &manageDataToFromIdentifier1, &endSponsorship,
					&beginSponsorship, &manageDataToFromIdentifier2, &endSponsorship},
				BaseFee:       constants.MinBaseFee,
				Memo:          nil,
				Preconditions: txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
			},
		)
		if err != nil {
			log.Fatal("Error while trying to build tranaction: ", err)
		}

	} else if txnPayload.Type == "9" {
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        &sourceAccount,
				IncrementSequenceNum: true,
				Operations: []txnbuild.Operation{&beginSponsorship, &manageDataType, &endSponsorship, &beginSponsorship, &manageDataIdentifier, &endSponsorship,
					&beginSponsorship, &manageDataToPreviousStage, &endSponsorship, &beginSponsorship, &manageDataCurrentStage, &endSponsorship,
					&beginSponsorship, &manageDataAppAccount, &endSponsorship},
				BaseFee:       constants.MinBaseFee,
				Memo:          nil,
				Preconditions: txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
			},
		)
		if err != nil {
			log.Fatal("Error while trying to build tranaction: ", err)
		}

	} else {

	}

	sposorerSK := decrpytIssuerSecretKey
	sponsorerKeypair, _ := keypair.ParseFull(sposorerSK)

	txe64, err := tx.Sign(commons.GetStellarNetwork(), sponsorerKeypair)
	if err != nil {
		hError := err.(*horizonclient.Error)
		log.Fatal("Error when submitting the transaction : ", hError)
	}

	txe, err := txe64.Base64()
	if err != nil {
		logger := utilities.NewCustomLogger()
		logger.LogWriter("Error converting to B64 : "+err.Error(), constants.ERROR)
		return txe, err
	}

	return txe, nil
}
