package services

import (
	"strconv"
	// "time"
	// "fmt"
	// "github.com/stellar/go/xdr"
	"log"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"

	// "fmt"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
)

//CheckTempOrphan ...
func CheckTempOrphan() {
	var clientList []string

	//if multiple clients do exist then query and append them
	// clientList = append(clientList, "GDWVYNO3S4TW25LZLCXH7DU5AWDI3OEDJWLOPG6I65RO5DFDSBZNOBB2")
	// clientList = append(clientList, "GB4CUDE7YBVI3VRC76RR4JPXI5KX62ANJWARFR47D5QTX46CSDKWRAKB")
	clientList = append(clientList, "GDNMRYEJ6NI2CVOPU3SDL4HYZWY535CZV7SC5MKEU2Z5PGTULMLAF2L5")
	object := dao.Connection{}

	//loop through clients
	for _, address := range clientList {
		//load horizon account
		account, err := horizon.DefaultPublicNetClient.LoadAccount(address)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Current Sequence for address:", address)
		log.Println(account.Sequence)

		seq, er := strconv.Atoi(account.Sequence)
		if er != nil {

		}

		stop := false //for infinite loop
		//loop through sequence incrementally and see match
		for i := seq + 1; ; i++ {

			p := object.GetSpecialForPkAndSeq(address, int64(i))
			p.Then(func(data interface{}) interface{} {
				result := data.(model.TransactionCollectionBody)
				var UserTxnHash string

				///HARDCODED CREDENTIALS
				publicKey := constants.PublicKey
				secretKey := constants.SecretKey
				switch result.TxnType {
				case "0":
					display := stellarExecuter.ConcreteSubmitXDR{XDR: result.XDR}
					response := display.SubmitXDR(false, result.TxnType)
					UserTxnHash = response.TXNID
					if response.Error.Code == 400 {
						break
					}

					var PreviousTXNBuilder build.ManageDataBuilder
					PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(""))

					//BUILD THE GATEWAY XDR
					tx, err := build.Transaction(
						build.PublicNetwork,
						build.SourceAccount{publicKey},
						build.AutoSequence{horizon.DefaultPublicNetClient},
						build.SetData("Type", []byte("G"+result.TxnType)),
						PreviousTXNBuilder,
						build.SetData("CurrentTXN", []byte(UserTxnHash)),
					)

					//SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
					GatewayTXE, err := tx.Sign(secretKey)
					if err != nil {
						break
					}

					//CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
					txeB64, err := GatewayTXE.Base64()
					if err != nil {
						break
					}

					//SUBMIT THE GATEWAY'S SIGNED XDR
					display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
					response1 := display1.SubmitXDR(false, "G"+result.TxnType)

					if response1.Error.Code == 400 {
						break
					}

					result.TxnHash = response1.TXNID
					result.Status = "done"
					///INSERT INTO TRANSACTION COLLECTION
					err2 := object.InsertTransaction(result)
					if err2 != nil {
						break
					}

				case "2":

					var PreviousTXNBuilder build.ManageDataBuilder

					// var PreviousTxn string
					p := object.GetLastTransactionbyIdentifier(result.Identifier)
					p.Then(func(data interface{}) interface{} {
						///ASSIGN PREVIOUS MANAGE DATA BUILDER
						res := data.(model.TransactionCollectionBody)
						PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(res.TxnHash))
						result.PreviousTxnHash = res.TxnHash
						return nil
					}).Catch(func(error error) error {

						return error
					})
					p.Await()

					display := stellarExecuter.ConcreteSubmitXDR{XDR: result.XDR}
					response := display.SubmitXDR(false, result.TxnType)
					UserTxnHash = response.TXNID
					if response.Error.Code == 400 {
						break
					}
					//BUILD THE GATEWAY XDR
					tx, err := build.Transaction(
						build.PublicNetwork,
						build.SourceAccount{publicKey},
						build.AutoSequence{horizon.DefaultPublicNetClient},
						build.SetData("Type", []byte("G"+result.TxnType)),
						PreviousTXNBuilder,
						build.SetData("CurrentTXN", []byte(UserTxnHash)),
					)

					//SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
					GatewayTXE, err := tx.Sign(secretKey)
					if err != nil {
						break
					}

					//CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
					txeB64, err := GatewayTXE.Base64()
					if err != nil {
						break
					}

					//SUBMIT THE GATEWAY'S SIGNED XDR
					display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
					response1 := display1.SubmitXDR(false, "G"+result.TxnType)

					if response1.Error.Code == 400 {
						break
					}

					result.TxnHash = response1.TXNID
					result.Status = "done"
					///INSERT INTO TRANSACTION COLLECTION
					err2 := object.InsertTransaction(result)
					if err2 != nil {
						break
					}
				}
				return nil
			}).Catch(func(error error) error {
				// return error
				log.Println("No transactions in the scheduler")
				stop = true //to break loop
				return nil
			})
			p.Await()

			if stop {
				break
			}
		}

	}

}
