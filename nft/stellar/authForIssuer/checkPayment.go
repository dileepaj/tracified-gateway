package authForIssuer

import (
	"fmt"
	"log"

	"github.com/stellar/go/xdr"
)

func CheckPayment(paymentXDR string) (bool) {

	var txe xdr.TransactionEnvelope

	fmt.Println(paymentXDR)

	//decode the XDR
	err := xdr.SafeUnmarshalBase64(paymentXDR, &txe)
	if err != nil {
		log.Println("Error while SafeUnmarshalBase64" + err.Error())
	}

	//check if the payment operation is there
	if *txe.Tx.Memo.Text == "Royalty Payed"{
		return true
	}else{
		return false
	}
	
	
}