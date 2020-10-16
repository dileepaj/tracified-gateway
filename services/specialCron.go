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

	// if multiple clients do exist then query and append them
	clientList = append(clientList, "GDWVYNO3S4TW25LZLCXH7DU5AWDI3OEDJWLOPG6I65RO5DFDSBZNOBB2")
	clientList = append(clientList, "GB4CUDE7YBVI3VRC76RR4JPXI5KX62ANJWARFR47D5QTX46CSDKWRAKB")
	clientList = append(clientList, "GBCBEUOCIYWD7TF3BPPYT66EJSIDO2ITRDWP4SQSBPAMEORKDRRHPKCF")
	
	//new keys
	clientList = append(clientList, "GDYMWT6NSLBTAXV6CJJH2MFVQZZ7DNEM5J3OVMIEA6LQ3TBRR2GIQZXN")
	clientList = append(clientList, "GCY4TYJOER457PGOS6DQNBCULUHA6EALXG6DRKBUAI44ORRIFDNMKYXB")
	clientList = append(clientList, "GA4RORXUCHI6HPUHKWMKN6SGYUHJVN2D5K3NBDPN7YE33HLIDGAGMBUY")
	clientList = append(clientList, "GDDTPYXF7VWV54GVOAVM2IPOA2HZYCYYEZYM5BI6NMGDYRUNXWPAS2XY")
	clientList = append(clientList, "GA4LVTLFF5377F3KXYPCONMN35HWEDTWQCEIHZBWZCZFFE354UISIE5F")
	clientList = append(clientList, "GCKTRM6BP2ZJSN7DAYNTPGH6BDWIQMH7OMHEGCSMSMHCTMQ4LUHCPDFX")
	clientList = append(clientList, "GDULLUKSH37SJKBMQXXYNJ4FAJMYKSBFGVCA2TQ5UAX7CHSAWLLQY3GL")
	clientList = append(clientList, "GAHVESB27GLES6YONUBTZDBKJYAQIK57KQJNCDKMKK5T5A6G76HPZEOY")

	//06-11-2019
	clientList = append(clientList, "GCYQRGGEB2W5V6VXN5OJFOIQOH7MDINX5RQMATVJUWKPC2A7PNQNTVU6")

	//01-05-2020
	clientList = append(clientList, "GC23K7ED7A5UZL42ZD6XJYHHDCJ4HWJTKEHTGFEDDA4RBOCIHFGO735Z")
	
	//16-10-2020 - for PVH demo
	clientList = append(clientList, "GANDBWSTTYGTG6VQJBUMYYG53WFOPPJRVIPYR6YP2UIKNS27GCOFD6WT")
	clientList = append(clientList, "GDULLUKSH37SJKBMQXXYNJ4FAJMYKSBFGVCA2TQ5UAX7CHSAWLLQY3GL")
	clientList = append(clientList, "GC6ELFQGPY42OSU43XPBTYS4KJR7GE74ZVM26JPYKC6P5FCO2K2GNEUI")
	clientList = append(clientList, "GA4XZCQCLIKJ2BAVPWZ2ES3V7V5TDBZFF6DLVOTXK5CT755XGNU2ARSL")

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
						PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(""))

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
				case "9":

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
						PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(""))

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
