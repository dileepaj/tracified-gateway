package authForIssuer

import (
	"encoding/base64"
	"fmt"
	"log"
	//"github.com/dileepaj/tracified-gateway/model"
	//"github.com/stellar/go/xdr"
)

func CheckPayment(hash string) bool {
	log.Println("Inside Check Payment!")
	fmt.Println(hash)

	var result bool

	//decode the hash to get memo
	decodedString, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		log.Println("Error while Decoding" + err.Error())
	} else {
		log.Println(decodedString)
		//check if the payment operation is there
		// 	if decodedString.Memo.text == "Payment has been made!" {
		// 		result = true
		// 	} else {
		// 		result = false
		// 	}
	}
	return result
}
