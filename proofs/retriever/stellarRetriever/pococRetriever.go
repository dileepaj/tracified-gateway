package stellarRetriever

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	// "github.com/dileepaj/tracified-gateway/model"

	"github.com/stellar/go/xdr"
)

type ConcretePOCOC struct {
	Txn string
	// DBCOC xdr.Transaction
	BCCOC xdr.Transaction
}

/*RetrievePOCOC - WORKING MODEL
@author - Azeem Ashraf
@desc - Retrieves the COC Txn from stellar using the TXN ID
@params - XDR
*/
func (db *ConcretePOCOC) RetrievePOCOC() xdr.Transaction {

	var CurrentTxn string
	CurrentTxn=db.Txn

	var txe xdr.Transaction
	//RETRIEVE GATEWAY SIGNED TXN
	// result, err := http.Get("https://horizon.stellar.org/transactions/" + db.Txn + "/operations")
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
			result, err := http.Get("https://horizon.stellar.org/transactions/" + CurrentTxn)
			if err != nil {

			} else {
				data, _ := ioutil.ReadAll(result.Body)

				if result.StatusCode == 200 {
					var raw map[string]interface{}
					json.Unmarshal(data, &raw)

					fmt.Println(raw["envelope_xdr"])
					fmt.Println("HAHAHAHAAHAHAH")
					err := xdr.SafeUnmarshalBase64(fmt.Sprintf("%s",raw["envelope_xdr"]), &txe)
					if err != nil {
					} 

				} else {

				}

			}
	// 	}


	// }

	return txe
}
