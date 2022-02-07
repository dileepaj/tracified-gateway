package authForIssuer

import (
	"fmt"
	"log"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/stellar/go/xdr"
)

func CheckPayment(royaltyXDR model.RoyaltyXDR) (bool) {

	var txe xdr.Transaction

	fmt.Println(royaltyXDR.XDR)

	var result bool

		//decode the XDR
	err := xdr.SafeUnmarshalBase64(royaltyXDR.XDR, &txe)
	if err != nil {
		log.Println("Error while SafeUnmarshalBase64" + err.Error())
	}else{
		//check if the payment operation is there
		if *txe.Memo.Text == "Royalty Payed"{
			result = true
		}else{
			result = false
		}
	}
	return result
}