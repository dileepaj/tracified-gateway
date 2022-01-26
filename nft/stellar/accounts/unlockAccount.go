package accounts

// import (
// 	"log"

// 	"github.com/dileepaj/tracified-gateway/commons"
// 	"github.com/stellar/go/build"
// )

// /*UnlockNFTIssuerAccount
// @desc - Unlock the issuer account
// @params - Current Issuer PK
// */
// func UnlockIssuingAccount(CurrentIssuerPK string) error {

// 	txn, err := build.Transaction(
// 		commons.GetHorizonNetwork(),
// 		build.SourceAccount{AddressOrSeed: CurrentIssuerPK},
// 		build.AutoSequence{SequenceProvider: commons.GetHorizonClient()},
// 		build.SetOptions(build.MasterWeight(255)),
// 	)

// 	if err != nil {
// 		log.Fatal(err)

// 	}
// 	signTxn, err := txn.Sign("SCHSOQDUY2BFFAKO3SPK6WEX4QRTRFUQ7KCI4T6VU4UUGBSLEQYFZCK3")
// 	if err != nil {
// 		log.Fatal("Error when submitting the transaction : ", " hError")

// 	}
// 	encodedTxn, err := signTxn.Base64()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	//submit transaction
// 	respn, err := commons.GetHorizonClient().SubmitTransaction(encodedTxn)
// 	if err != nil {
// 		log.Fatal("Error submitting transaction:", err)

// 	}
// 	log.Println("Hash for unlocking is", respn.Hash)
// 	return nil
// }
