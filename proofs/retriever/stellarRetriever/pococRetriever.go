package stellarRetriever

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	// "github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/stellar/go/support/log"
)

/*
type ConcretePOCOC struct {
	Txn string
	// DBCOC xdr.Transaction
	BCCOC xdr.Transaction
}
*/

/*RetrievePOCOC - WORKING MODEL
@author - Azeem Ashraf
@desc - Retrieves the COC Txn from stellar using the TXN ID
@params - XDR
*/

/*

func (db *ConcretePOCOC) RetrievePOCOC() (xdr.Transaction, bool, string, string, string) {

	var CurrentTxn string
	CurrentTxn = db.Txn

	timestamp := ""
	ledger := ""
	feePaid := ""
	var txe xdr.Transaction
	//RETRIEVE GATEWAY SIGNED TXN
	// result, err := http.Get("https://horizon.stellar.orgtransactions/" + db.Txn + "/operations")
	// if err != nil {

	// } else {
	// 	data, _ := ioutil.ReadAll(result.Body)

	// 	if result.StatusCode == 200 {
	// 		var raw map[string]interface{}
	// 		json.Unmarshal(data, &raw)
	// 		// raw["count"] = 2
	// 		out, _ := json.Marshal(raw["_embedded"])
	// 		var raw1 map[string]interface{}
	// 		json.Unmarshal(out, &raw1)
	// 		out1, _ := json.Marshal(raw1["records"])

	// 		keysBody := out1
	// 		keys := make([]PublicKey, 0)
	// 		json.Unmarshal(keysBody, &keys)

	// 		// Gtype:=Base64DecEnc("Decode", keys[0].Value)
	// 		// PreviousTxn = Base64DecEnc("Decode", keys[1].Value)
	// 		CurrentTxn = Base64DecEnc("Decode", keys[2].Value)

	//RETRIEVE THE USER SIGNED TXN USING THE CURRENT TXN IN GATEWAY SIGNED TRANSACTION
	result, err := http.Get("https://horizon.stellar.orgtransactions/" + CurrentTxn)
	if err != nil {
		return txe, false, timestamp, ledger, feePaid
	} else {
		data, _ := ioutil.ReadAll(result.Body)

		if result.StatusCode == 200 {
			var raw map[string]interface{}
			json.Unmarshal(data, &raw)

			fmt.Println(raw["envelope_xdr"])
			fmt.Println("HAHAHAHAAHAHAH")
			timestamp = fmt.Sprintf("%s", raw["created_at"])
			ledger = fmt.Sprintf("%.0f", raw["ledger"])
			feePaid = fmt.Sprintf("%.0f", raw["fee_paid"])

			err := xdr.SafeUnmarshalBase64(fmt.Sprintf("%s", raw["envelope_xdr"]), &txe)
			if err != nil {
			}

		} else {
			return txe, false, timestamp, ledger, feePaid

		}

	}
	// 	}

	// }

	return txe, true, timestamp, ledger, feePaid

}
*/
type ConcretePOCOCNew struct {
	Txn string
	// DBCOC xdr.Transaction
	BCCOC XDR
}

type XDR struct {
	SourceAccount string
	Identifier    string
	AssetCode     string
	AssetAmount   float64
	Destination   string
}

func (db *ConcretePOCOCNew) RetrievePOCOCNew() (XDR, bool, string, string, string) {

	CurrentTxn := db.Txn
	timestamp := ""
	ledger := ""
	feePaid := ""

	var txe XDR
	result, err := http.Get(commons.GetHorizonClient().HorizonURL + "transactions/" + CurrentTxn)
	if err != nil {
		return txe, false, timestamp, ledger, feePaid
	} else {
		data, _ := ioutil.ReadAll(result.Body)

		if result.StatusCode == 200 {
			var raw map[string]interface{}
			json.Unmarshal(data, &raw)

			timestamp = fmt.Sprintf("%s", raw["created_at"])
			ledger = fmt.Sprintf("%.0f", raw["ledger"])
			feePaid = fmt.Sprintf("%s", raw["fee_charged"])

			result, err := http.Get(commons.GetHorizonClient().HorizonURL + "transactions/" + CurrentTxn + "/operations")

			if err != nil {
				log.Error("Error while url failed " + err.Error())
			}

			data1, err := ioutil.ReadAll(result.Body)
			if err != nil {
				log.Error("Error while read response " + err.Error())
			}
			var raw1 map[string]interface{}

			err = json.Unmarshal(data1, &raw1)
			if err != nil {
				log.Error("Error while json.Unmarshal(data, &raw) " + err.Error())
			}

			out, err := json.Marshal(raw1["_embedded"])
			if err != nil {
				log.Error("Error while json marshal _embedded " + err.Error())
			}
			var raw2 map[string]interface{}
			err = json.Unmarshal(out, &raw2)
			if err != nil {
				log.Error("Error while json.Unmarshal(out, &raw1) " + err.Error())
			}
			out1, err := json.Marshal(raw2["records"])
			if err != nil {
				log.Error("Error while json marshal records " + err.Error())
			}
			keysBody := out1
			keys := make([]PublicKeyProofhash, 0)
			err = json.Unmarshal(keysBody, &keys)

			if err != nil {
				log.Error("Error while json.Unmarshal(keysBody, &keys) " + err.Error())
			}

			txe.SourceAccount = string(keys[1].Source_account)
			log.Info("Source Account: " + txe.SourceAccount)

			AssetCode_byteData, err := base64.StdEncoding.DecodeString(keys[3].Value)
			if err != nil {
				log.Error("Error while base64.StdEncoding.DecodeString " + err.Error())
			}

			txe.AssetCode = string(AssetCode_byteData)
			log.Info("Asset Code: " + txe.AssetCode)

			AssetAmount_byteData, err := base64.StdEncoding.DecodeString(keys[4].Value)
			if err != nil {
				log.Error("Error while base64.StdEncoding.DecodeString " + err.Error())
			}

			txe.AssetAmount, err = strconv.ParseFloat(string(AssetAmount_byteData), 64)
			if err != nil {
				log.Error("string to float64 conversion error")
			}

			log.Info("Asset Amount: " + fmt.Sprintf("%f", txe.AssetAmount))

			Identifier_byteData, err := base64.StdEncoding.DecodeString(keys[1].Value)
			if err != nil {
				log.Error("Error while base64.StdEncoding.DecodeString " + err.Error())
			}

			txe.Identifier = string(Identifier_byteData)
			log.Info("Identifier: " + txe.Identifier)

			Receiver_byteData, err := base64.StdEncoding.DecodeString(keys[2].Value)
			if err != nil {
				log.Error("Error while base64.StdEncoding.DecodeString " + err.Error())
			}

			txe.Destination = string(Receiver_byteData)
			log.Info("Identifier: " + txe.Identifier)

		} else {
			return txe, false, timestamp, ledger, feePaid

		}

	}

	return txe, true, timestamp, ledger, feePaid
}

type PublicKeyProofhash struct {
	Name           string
	Value          string
	Source_account string
}
