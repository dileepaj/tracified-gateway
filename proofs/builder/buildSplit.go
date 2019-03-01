package builder

import (
	// "encoding/json"
	"fmt"
	// "net/http"
	// "strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/xdr"

	// "github.com/dileepaj/tracified-gateway/api/apiModel"
)



func (AP *AbstractXDRSubmiter) SubmitSplit() bool {
	var Done bool
	object := dao.Connection{}
	var copy []model.TransactionCollectionBody

	var UserSplitTxnHashes []string
	// var ParentIdentifier string
	// var ParentTxn string
	var PreviousTxn string

	///HARDCODED CREDENTIALS
	publicKey := constants.PublicKey
	secretKey := constants.SecretKey
	// var result model.SubmitXDRResponse

	for i, TxnBody := range AP.TxnBody {

		var TDP model.TransactionCollectionBody
		var txe xdr.Transaction

		//decode the XDR
		err := xdr.SafeUnmarshalBase64(TxnBody.XDR, &txe)
		if err != nil {
			fmt.Println(err)
		}

		//GET THE TYPE AND IDENTIFIER FROM THE XDR
		TxnBody.PublicKey = txe.SourceAccount.Address()

		TxnBody.TxnType = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		TxnBody.Identifier = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		TxnBody.ItemCode = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[3].Body.ManageDataOp.DataValue), "&")
		if i == 0 {
			TxnBody.ToIdentifier = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[2].Body.ManageDataOp.DataValue), "&")
		} else {
			TxnBody.FromIdentifier1 = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[2].Body.ManageDataOp.DataValue), "&")
			TxnBody.ItemAmount = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[4].Body.ManageDataOp.DataValue), "&")
		}

		AP.TxnBody[i].Identifier = TxnBody.Identifier
		AP.TxnBody[i].TxnType = TxnBody.TxnType

		//FOR THE SPLIT PARENT RETRIEVE THE PREVIOUS TXN FROM GATEWAY DB
		if i == 0 {
			// ParentIdentifier = Identifier
			p := object.GetLastTransactionbyIdentifier(TxnBody.Identifier)
			p.Then(func(data interface{}) interface{} {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER
				result := data.(model.TransactionCollectionBody)
				PreviousTxn = result.TxnHash
				TxnBody.PreviousTxnHash = result.TxnHash

				fmt.Println(TxnBody.PreviousTxnHash)
				return nil
			}).Catch(func(error error) error {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
				//DUE TO THE CHILD HAVING A NEW IDENTIFIER
				TxnBody.PreviousTxnHash = ""
				return error
			})
			p.Await()
		}
		TxnBody.Status = "pending"

		copy = append(copy, TxnBody)

		///INSERT INTO TRANSACTION COLLECTION
		err1 := object.InsertTransaction(TxnBody)
		if err1 != nil {
			TDP.Status = "failed"
		}

		//SUBMIT THE FIRST XDR SIGNED BY THE USER
		display := stellarExecuter.ConcreteSubmitXDR{XDR: TxnBody.XDR}
		result := display.SubmitXDR()
		UserSplitTxnHashes = append(UserSplitTxnHashes, result.TXNID)

		if result.Error.Code != 404 {
			Done = true
			// return Done
		}
	}
	go func() {

		var SplitParentProfile string
		var PreviousSplitProfile string
		for i, TxnBody := range AP.TxnBody {
			var PreviousTXNBuilder build.ManageDataBuilder

			PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(PreviousTxn))
			TxnBody.PreviousTxnHash = PreviousTxn

			//ASSIGN THE PREVIOUS PROFILE ID USING THE PARENT FOR THE CHILDREN AND A DB CALL FOR PARENT
			if i == 0 {
				PreviousSplitProfile = ""
				SplitParentProfile = TxnBody.ProfileID
			} else {
				PreviousSplitProfile = SplitParentProfile
			}
			//BUILD THE GATEWAY XDR
			tx, err := build.Transaction(
				build.TestNetwork,
				build.SourceAccount{publicKey},
				build.AutoSequence{horizon.DefaultTestNetClient},
				build.SetData("Type", []byte("G"+TxnBody.TxnType)),
				PreviousTXNBuilder,
				build.SetData("CurrentTXN", []byte(UserSplitTxnHashes[i])),
				build.SetData("Identifier", []byte(TxnBody.Identifier)),
				build.SetData("ProfileID", []byte(TxnBody.ProfileID)),
				build.SetData("PreviousProfile", []byte(PreviousSplitProfile)),
			)

			//SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
			GatewayTXE, err := tx.Sign(secretKey)
			if err != nil {
				fmt.Println(err)
			}
			//CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
			txeB64, err := GatewayTXE.Base64()
			if err != nil {
				fmt.Println(err)
			}

			//SUBMIT THE GATEWAY'S SIGNED XDR
			display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
			response1 := display1.SubmitXDR()

			if response1.Error.Code == 404 {
				TxnBody.Status = "pending"
			} else {
				//UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
				TxnBody.TxnHash = response1.TXNID
				if i == 0 {
					PreviousTxn = response1.TXNID
				}
				upd := model.TransactionCollectionBody{
					TxnHash:         response1.TXNID,
					Status:          "done",
					PreviousTxnHash: TxnBody.PreviousTxnHash}
				err2 := object.UpdateTransaction(copy[i], upd)
				if err2 != nil {
					TxnBody.Status = "pending"
				} else {
					TxnBody.Status = "done"
				}
				// Done = true
			}
		}

	}()
	// }
	return Done
}