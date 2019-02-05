package builder

import (
	"fmt"
	"strings"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/xdr"

	// "github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/model"
)

// type InsertData struct{}

// type XDRSubmitInterface interface {
// 	InsertDataHash() model.InsertDataResponse
// }

type AbstractXDRSubmiter struct {
	TxnBody []model.TransactionCollectionBody
}

func (AP *AbstractXDRSubmiter) SubmitGenesis() bool {
	var Done bool
	object := dao.Connection{}
	var copy []model.TransactionCollectionBody

	var UserTxnHashes []string
	// var PreviousTxn []string
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
		// TDP.PreviousTxnHash=
		TxnType := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		Identifier := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		TxnBody.Identifier = Identifier
		TxnBody.TxnType = TxnType
		TxnBody.Status = "pending"
		// TxnBody.TdpId=
		AP.TxnBody[i].Identifier = Identifier
		AP.TxnBody[i].TxnType = TxnType

		fmt.Println(Identifier)
		p := object.GetLastTransactionbyIdentifier(Identifier)
		p.Then(func(data interface{}) interface{} {
			///ASSIGN PREVIOU y S MANAGE DATA BUILDER

			result := data.(model.TransactionCollectionBody)
			TxnBody.PreviousTxnHash = result.TxnHash
			fmt.Println("Previous TXN: " + result.TxnHash)

			return nil
		}).Catch(func(error error) error {
			///ASSIGN PREVIOUS MANAGE DATA BUILDER - LEAVE IT EMPTY
			fmt.Println("Previous TXN: ERR ")
			return error
		})
		p.Await()
		copy = append(copy, TxnBody)
		///INSERT INTO TRANSACTION COLLECTION
		err1 := object.InsertTransaction(TxnBody)
		if err1 != nil {
			TDP.Status = "failed"
		}
		//SUBMIT THE FIRST XDR SIGNED BY THE USER
		display := stellarExecuter.ConcreteSubmitXDR{XDR: TxnBody.XDR}
		result := display.SubmitXDR()
		UserTxnHashes = append(UserTxnHashes, result.TXNID)

		if result.Error.Code != 404 {
			Done = true
			// return Done
		}
	}

	go func() {
		for i, TxnBody := range AP.TxnBody {

			var PreviousTXNBuilder build.ManageDataBuilder

			PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(copy[i].PreviousTxnHash))

			//BUILD THE GATEWAY XDR
			tx, err := build.Transaction(
				build.TestNetwork,
				build.SourceAccount{publicKey},
				build.AutoSequence{horizon.DefaultTestNetClient},
				PreviousTXNBuilder,
				build.SetData("CurrentTXN", []byte(UserTxnHashes[i])),
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

				upd := model.TransactionCollectionBody{
					TxnHash:         response1.TXNID,
					Status:          "done",
				}
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
	// Done=true
	return Done
}


func (AP *AbstractXDRSubmiter) SubmitData() bool {
	var Done bool
	object := dao.Connection{}
	var copy []model.TransactionCollectionBody

	var UserTxnHashes []string
	// var PreviousTxn []string
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
		// TDP.PreviousTxnHash=
		TxnType := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		Identifier := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		TxnBody.Identifier = Identifier
		TxnBody.TxnType = TxnType
		TxnBody.Status = "pending"
		// TxnBody.TdpId=
		AP.TxnBody[i].Identifier = Identifier
		AP.TxnBody[i].TxnType = TxnType

		fmt.Println(Identifier)
		p := object.GetLastTransactionbyIdentifier(Identifier)
		p.Then(func(data interface{}) interface{} {
			///ASSIGN PREVIOU y S MANAGE DATA BUILDER

			result := data.(model.TransactionCollectionBody)
			TxnBody.PreviousTxnHash = result.TxnHash
			fmt.Println("Previous TXN: " + result.TxnHash)

			return nil
		}).Catch(func(error error) error {
			///ASSIGN PREVIOUS MANAGE DATA BUILDER - LEAVE IT EMPTY
			fmt.Println("Previous TXN: ERR ")
			return error
		})
		p.Await()
		copy = append(copy, TxnBody)
		///INSERT INTO TRANSACTION COLLECTION
		err1 := object.InsertTransaction(TxnBody)
		if err1 != nil {
			TDP.Status = "failed"
		}
		//SUBMIT THE FIRST XDR SIGNED BY THE USER
		display := stellarExecuter.ConcreteSubmitXDR{XDR: TxnBody.XDR}
		result := display.SubmitXDR()
		UserTxnHashes = append(UserTxnHashes, result.TXNID)

		if result.Error.Code != 404 {
			Done = true
			// return Done
		}
	}

	go func() {
		for i, TxnBody := range AP.TxnBody {

			var PreviousTXNBuilder build.ManageDataBuilder

			PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(copy[i].PreviousTxnHash))

			//BUILD THE GATEWAY XDR
			tx, err := build.Transaction(
				build.TestNetwork,
				build.SourceAccount{publicKey},
				build.AutoSequence{horizon.DefaultTestNetClient},
				PreviousTXNBuilder,
				build.SetData("CurrentTXN", []byte(UserTxnHashes[i])),
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

				upd := model.TransactionCollectionBody{
					TxnHash:         response1.TXNID,
					Status:          "done",
				}
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
	// Done=true
	return Done
}

func (AP *AbstractXDRSubmiter) SubmitSplit() bool {
	var Done bool
	object := dao.Connection{}
	var copy []model.TransactionCollectionBody

	var UserSplitTxnHashes []string
	var ParentIdentifier string
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
		TDP.PublicKey = txe.SourceAccount.Address()
		TxnType := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		Identifier := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		TDP.Identifier = Identifier
		AP.TxnBody[i].Identifier = Identifier
		AP.TxnBody[i].TxnType = TxnType

		if i == 0 {
			ParentIdentifier = Identifier

			p := object.GetLastTransactionbyIdentifier(Identifier)
			p.Then(func(data interface{}) interface{} {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER
				result := data.(model.TransactionCollectionBody)
				PreviousTxn = result.TxnHash
				return nil
			}).Catch(func(error error) error {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
				//DUE TO THE CHILD HAVING A NEW IDENTIFIER
				PreviousTxn = ""
				return error
			})
			p.Await()
		}

		TDP.TxnType = TxnType
		TDP.Status = "pending"

		copy = append(copy, TDP)

		///INSERT INTO TRANSACTION COLLECTION
		err1 := object.InsertTransaction(TDP)
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

		for i, TxnBody := range AP.TxnBody {
			var PreviousTXNBuilder build.ManageDataBuilder

			////GET THE PREVIOUS TRANSACTION FOR THE IDENTIFIER

			if i == 0 {
				PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(PreviousTxn))
			} else {
				p := object.GetLastTransactionbyIdentifier(ParentIdentifier)
				p.Then(func(data interface{}) interface{} {
					///ASSIGN PREVIOUS MANAGE DATA BUILDER
					result := data.(model.TransactionCollectionBody)
					PreviousTxn = result.TxnHash
					return nil
				}).Catch(func(error error) error {
					///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
					//DUE TO THE CHILD HAVING A NEW IDENTIFIER
					PreviousTxn = ""
					return error
				})
				p.Await()
				PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(PreviousTxn))
			}

			//BUILD THE GATEWAY XDR
			tx, err := build.Transaction(
				build.TestNetwork,
				build.SourceAccount{publicKey},
				build.AutoSequence{horizon.DefaultTestNetClient},
				PreviousTXNBuilder,
				build.SetData("CurrentTXN", []byte(UserSplitTxnHashes[i])),
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

				upd := model.TransactionCollectionBody{TxnHash: response1.TXNID, Status: "done"}
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

func (AP *AbstractXDRSubmiter) SubmitMerge() bool {
	var Done bool
	object := dao.Connection{}
	var copy model.TransactionCollectionBody
	///HARDCODED CREDENTIALS
	publicKey := constants.PublicKey
	secretKey := constants.SecretKey
	// var result model.SubmitXDRResponse

	for _, TxnBody := range AP.TxnBody {

		var TDP model.TransactionCollectionBody
		var txe xdr.Transaction
		var PreviousTXN build.ManageDataBuilder

		//decode the XDR
		err := xdr.SafeUnmarshalBase64(TxnBody.XDR, &txe)
		if err != nil {
			fmt.Println(err)
		}

		//GET THE TYPE AND IDENTIFIER FROM THE XDR
		TDP.PublicKey = txe.SourceAccount.Address()
		TxnType := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		Identifier := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		TDP.Identifier = Identifier
		TDP.TxnType = TxnType
		TDP.Status = "pending"

		copy = TDP

		///INSERT INTO TRANSACTION COLLECTION
		err1 := object.InsertTransaction(TDP)
		if err1 != nil {
			TDP.Status = "failed"
		}

		//SUBMIT THE FIRST XDR SIGNED BY THE USER
		display := stellarExecuter.ConcreteSubmitXDR{XDR: TxnBody.XDR}
		result := display.SubmitXDR()

		if result.Error.Code != 404 {
			Done = true
			// return Done
		}
		go func() {
			////GET THE PREVIOUS TRANSACTION FOR THE IDENTIFIER
			p := object.GetLastTransactionbyIdentifier(Identifier)
			p.Then(func(data interface{}) interface{} {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER
				result := data.(model.TransactionCollectionBody)
				PreviousTXN = build.SetData("PreviousTXN", []byte(result.TxnHash))
				return nil
			}).Catch(func(error error) error {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER - LEAVE IT EMPTY
				PreviousTXN = build.SetData("PreviousTXN", []byte(""))
				return error
			})
			p.Await()

			//BUILD THE GATEWAY XDR
			tx, err := build.Transaction(
				build.TestNetwork,
				build.SourceAccount{publicKey},
				build.AutoSequence{horizon.DefaultTestNetClient},
				PreviousTXN,
				build.SetData("CurrentTXN", []byte(result.TXNID)),
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

				upd := model.TransactionCollectionBody{TxnHash: response1.TXNID, Status: "done"}
				err2 := object.UpdateTransaction(copy, upd)
				if err2 != nil {
					TxnBody.Status = "pending"
				} else {
					TxnBody.Status = "done"
				}
				Done = true
			}
		}()
	}
	return Done
}

func (AP *AbstractXDRSubmiter) SubmitTransfer() bool {
	var Done bool
	object := dao.Connection{}
	var copy model.TransactionCollectionBody
	///HARDCODED CREDENTIALS
	publicKey := "GDYVWLGWSCYCG7QDO5IFDQXJAUDCPBYUV2EPPDFUYEZNIOA5FQ4JAXB6"
	secretKey := "SC45GBPQJ53QCYODXFGBCYR5FWTKMZBNF4ZTKDTUFY76FZTIKIU7WGTS"
	// var result model.SubmitXDRResponse

	for _, TxnBody := range AP.TxnBody {

		var TDP model.TransactionCollectionBody
		var txe xdr.Transaction
		var PreviousTXN build.ManageDataBuilder

		//decode the XDR
		err := xdr.SafeUnmarshalBase64(TxnBody.XDR, &txe)
		if err != nil {
			fmt.Println(err)
		}

		//GET THE TYPE AND IDENTIFIER FROM THE XDR
		TDP.PublicKey = txe.SourceAccount.Address()
		TxnType := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		Identifier := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		TDP.Identifier = Identifier
		TDP.TxnType = TxnType
		TDP.Status = "pending"

		copy = TDP

		///INSERT INTO TRANSACTION COLLECTION
		err1 := object.InsertTransaction(TDP)
		if err1 != nil {
			TDP.Status = "failed"
		}

		//SUBMIT THE FIRST XDR SIGNED BY THE USER
		display := stellarExecuter.ConcreteSubmitXDR{XDR: TxnBody.XDR}
		result := display.SubmitXDR()

		if result.Error.Code != 404 {
			Done = true
			// return Done
		}
		go func() {
			////GET THE PREVIOUS TRANSACTION FOR THE IDENTIFIER
			p := object.GetLastTransactionbyIdentifier(Identifier)
			p.Then(func(data interface{}) interface{} {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER
				result := data.(model.TransactionCollectionBody)
				PreviousTXN = build.SetData("PreviousTXN", []byte(result.TxnHash))
				return nil
			}).Catch(func(error error) error {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER - LEAVE IT EMPTY
				PreviousTXN = build.SetData("PreviousTXN", []byte(""))
				return error
			})
			p.Await()

			//BUILD THE GATEWAY XDR
			tx, err := build.Transaction(
				build.TestNetwork,
				build.SourceAccount{publicKey},
				build.AutoSequence{horizon.DefaultTestNetClient},
				PreviousTXN,
				build.SetData("CurrentTXN", []byte(result.TXNID)),
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

				upd := model.TransactionCollectionBody{TxnHash: response1.TXNID, Status: "done"}
				err2 := object.UpdateTransaction(copy, upd)
				if err2 != nil {
					TxnBody.Status = "pending"
				} else {
					TxnBody.Status = "done"
				}
				Done = true
			}
		}()
	}
	return Done
}

func XDRSubmitter(TDP []model.TransactionCollectionBody) (bool, model.SubmitXDRResponse) {
	object := dao.Connection{}
	var copy model.TransactionCollectionBody
	var ret model.SubmitXDRResponse

	for i := 0; i < len(TDP); i++ {
		TDP[i].Status = "Pending"
		var txe xdr.Transaction
		err := xdr.SafeUnmarshalBase64(TDP[i].XDR, &txe)
		if err != nil {
			fmt.Println(err)
		}

		TDP[i].PublicKey = txe.SourceAccount.Address()
		TxnType := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		Identifier := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[2].Body.ManageDataOp.DataValue), "&")
		TDP[i].Identifier = Identifier
		TDP[i].TxnType = TxnType
		TDP[i].Status = "pending"

		copy = TDP[i]
		err1 := object.InsertTransaction(TDP[i])
		if err1 != nil {
			TDP[i].Status = "failed"
		}

	}
	for i := 0; i < len(TDP); i++ {
		display := stellarExecuter.ConcreteSubmitXDR{XDR: TDP[i].XDR}

		response := display.SubmitXDR()
		ret = response
		if response.Error.Code == 404 {
			TDP[i].Status = "pending"
		} else {
			TDP[i].TxnHash = response.TXNID

			upd := model.TransactionCollectionBody{TxnHash: response.TXNID, Status: "done"}
			err2 := object.UpdateTransaction(copy, upd)
			if err2 != nil {
				TDP[i].Status = "pending"
			} else {
				TDP[i].Status = "done"
			}
		}
	}

	return true, ret
}
