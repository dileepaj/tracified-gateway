package businessFacades

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

type bitString string

// BuildSocialImpactFormula --> this method take the farmulaJSOn from backend and map the equationID to 4 byte value and store it DB,
// generate  executionTemplete object via passing FCL query to FCL library,
// base on the executionTemplete objcet store the formula in stellar blckchin
func BuildSocialImpactFormula(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var formulaJSON model.BuildFormula

	err := json.NewDecoder(r.Body).Decode(&formulaJSON)
	if err != nil {
		logrus.Error(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while decoding the body")
		return
	}

	object := dao.Connection{}
	formulaMap,err5:=object.GetFormulaMapID(formulaJSON.FormulaID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err5 != nil {
		logrus.Info("Unable to connect gateway datastore ",err5)
		w.WriteHeader(http.StatusNotFound)
	}
	if formulaMap != nil {
		logrus.Error("Formula Id is in gateway datastore")
		w.WriteHeader(http.StatusNoContent)
		response := model.Error{Code: http.StatusNoContent, Message: "Formula Id is in gateway datastore"}
		json.NewEncoder(w).Encode(response)
		return
	}
	


	data,err:=object.GetNextSequenceValue("FORMULAID")
	if err != nil {
		fmt.Println("GetNextSequenceValu was failed" + err.Error())
	}

	fmt.Println("Sequen numb ++++++++++++++++ ",data)


	n := strconv.FormatInt(data.SequenceValue, 2)
	// fmt.Println(n)
	str := fmt.Sprintf("%032s", n)
	fmt.Println("srt    ",str)
	bitValue := bitString(str)
	byteValue := bitValue.AsByteSlice()
	fmt.Println(byteValue)
	fmt.Println(binary.Size(bitValue.AsByteSlice()))
	for _, nq := range byteValue {
		// fmt.Printf("% 09b", nq) // prints 00000000 11111101
		fmt.Println("vvv ----------------      ", nq)
	}
	fmt.Println("vvvvvvvvvv----------: ", string(byteValue[:]))


	formulaIDMap:=model.FormulaIDMap{
		FormulaID: formulaJSON.FormulaID,
		MapID:     data.SequenceValue,
		ByteID:    byteValue,
		BitString: str,
	}

	err6:=object.InsertFormulaIDMap(formulaIDMap)
	if err6 != nil {
		
	}

	publicKey := constants.PublicKey
	secretKey := constants.SecretKey
	tracifiedAccount, err := keypair.ParseFull(secretKey)
	client := commons.GetHorizonClient()
	pubaccountRequest := horizonclient.AccountRequest{AccountID: publicKey}
	pubaccount, err := client.AccountDetail(pubaccountRequest)


	PreviousTXNBuilder := txnbuild.ManageData{
		Name:  "new",
		Value: byteValue,
	}

	if err != nil {
		logrus.Println(err)
	}
	// BUILD THE GATEWAY XDR
	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &pubaccount,
		IncrementSequenceNum: true,
		Operations:           []txnbuild.Operation{&PreviousTXNBuilder},
		BaseFee:              txnbuild.MinBaseFee,
		Memo:                 nil,
		Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
	})
	if err != nil {
		logrus.Println("Error while buliding XDR " + err.Error())
	}
	// SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
	GatewayTXE, err := tx.Sign(commons.GetStellarNetwork(), tracifiedAccount)
	if err != nil {
		logrus.Println("Error while getting GatewayTXE by secretKey " + err.Error())
	}
	// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
	resp, err := client.SubmitTransaction(GatewayTXE)

	fmt.Println("has----  ", resp.Hash)

	url := "https://horizon-testnet.stellar.org/transactions/" + resp.Hash + "/operations"
	result, err := http.Get(url)
	if err != nil {
		logrus.Error("Unable to reach Stellar network", url)
	}
	if result.StatusCode != 200 {
		logrus.Error(result)
	}
	defer result.Body.Close()
	assertInfo, err := ioutil.ReadAll(result.Body)
	if err != nil {
		logrus.Error(err)
	}

	var raw map[string]interface{}
	var raw1 map[string]interface{}

	json.Unmarshal(assertInfo, &raw)

	out, err := json.Marshal(raw["_embedded"])
	if err != nil {
		logrus.Error("Unable to marshal embedded")
	}
	err = json.Unmarshal(out, &raw1)
	if err != nil {
		logrus.Error("Unable to unmarshal  json.Unmarshal(out, &raw1)")
	}
	out1, _ := json.Marshal(raw1["records"])
	json.Unmarshal(out1, &raw1)

	keysBody := out1
	keys := make([]PublicKeyPOCOC, 0)
	err = json.Unmarshal(keysBody, &keys)
	if err != nil {
		logrus.Error("Unable to unmarshal keys data")
	}
	fmt.Println("--raw1----- ", keys)
	for i := range keys {
		fmt.Println(keys[i])
		acceptTxn_byteData, err := base64.StdEncoding.DecodeString(keys[i].Value)
		if err != nil {
			logrus.Error("Unable to base64 decoding")
		}
		acceptTxn := string(acceptTxn_byteData)
		fmt.Println("acceptTxn----------: " + acceptTxn)
		for _, nq := range acceptTxn_byteData {
			// fmt.Printf("% 09b", nq) // prints 00000000 11111101
			fmt.Println("vvv retrive----------------      ", nq)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(formulaJSON)
	return
}

func mapEquationIDInToNumber() {
}

func stringToBin(s string) (binString string) {
	for _, c := range s {
		binString = fmt.Sprintf("%s%b", binString, c)
	}
	return
}

func (b bitString) AsByteSlice() []byte {
	var out []byte
	var str string

	for i := len(b); i > 0; i -= 8 {
		if i-8 < 0 {
			str = string(b[0:i])
		} else {
			str = string(b[i-8 : i])
		}
		v, err := strconv.ParseUint(str, 2, 8)
		if err != nil {
			panic(err)
		}
		out = append([]byte{byte(v)}, out...)
	}
	return out
}

func (b bitString) AsHexSlice() []string {
	var out []string
	byteSlice := b.AsByteSlice()
	for _, b := range byteSlice {
		out = append(out, "0x"+hex.EncodeToString([]byte{b}))
	}
	return out
}
