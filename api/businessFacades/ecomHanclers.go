package businessFacades

import (
	"encoding/base64"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
    "regexp"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stellar/go/strkey"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/gorilla/mux"
)

type transaction struct {
	transaction string
}

type tdpToTransaction struct {
	TdpId string `json:"tdpId"`
}

func GetTransactionId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	object := dao.Connection{}
	p := object.GetTransactionForTdpId(vars["id"])
	fmt.Println(vars["id"])
	p.Then(func(data interface{}) interface{} {
		TxnHash := (data.(model.TransactionCollectionBody)).TxnHash

		mapD := map[string]string{"transaction": TxnHash}
		mapB, _ := json.Marshal(mapD)
		fmt.Println(string(mapB))
		// trans := transaction{transaction:TxnHash}
		// s := fmt.Sprintf("%v", trans)

		encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
		text := (string(encoded))
		result := model.TransactionId{Txnhash: TxnHash,
			Url: "https://www.stellar.org/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
				text + "%3D%3D&network=public"}

		// res := TDP{TdpId: result.TdpId}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "TDP ID Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()

}

func GetTransactionsForTDP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	var result []model.TransactionIds
	object := dao.Connection{}
	p := object.GetAllTransactionForTdpId(vars["id"])
	fmt.Println(vars["id"])
	p.Then(func(data interface{}) interface{} {
		res := data.([]model.TransactionCollectionBody)
		for _, TxnBody := range res {
			TxnHash := TxnBody.TxnHash

			mapD := map[string]string{"transaction": TxnHash}
			mapB, _ := json.Marshal(mapD)
			fmt.Println(string(mapB))
			// trans := transaction{transaction:TxnHash}
			// s := fmt.Sprintf("%v", trans)

			encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
			text := (string(encoded))
			temp := model.TransactionIds{Txnhash: TxnHash,
				Url: "https://www.stellar.org/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
					text + "%3D%3D&network=public",
				Identifier: TxnBody.Identifier}

			result = append(result, temp)
		}

		// res := TDP{TdpId: result.TdpId}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "TDP ID Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()

}

func GetTransactionsForTdps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	fmt.Println("lol")
	var TDPs apiModel.GetTransactionId

	if r.Header == nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "No Header present!",
		}
		json.NewEncoder(w).Encode(result)

		return
	}

	if r.Header.Get("Content-Type") == "" {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "No Content-Type present!",
		}
		json.NewEncoder(w).Encode(result)

		return
	}

	// fmt.Println(TDP)
	err := json.NewDecoder(r.Body).Decode(&TDPs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		fmt.Println(err)
		return
	}
	// fmt.Println(TDPs)
	object := dao.Connection{}

	var resultArray []model.TransactionIds
	var identifer string
	var arrSize = len(TDPs.TdpID)
	for i := 0; i < arrSize; i++ {

		p := object.GetTransactionForTdpId(TDPs.TdpID[i])
		fmt.Println(TDPs.TdpID[i])
		p.Then(func(data interface{}) interface{} {

			if TDPs.TdpID[i] == "" {
				temp := model.TransactionIds{Txnhash: "Not Found", Url: "Not Found", Identifier: "Not Found", TdpId: TDPs.TdpID[i]}
				resultArray = append(resultArray, temp)
			} else {
				Txn := data.(model.TransactionCollectionBody)
				mapD := map[string]string{"transaction": Txn.TxnHash}
				mapB, _ := json.Marshal(mapD)
				fmt.Println(Txn.ProfileID)
				// trans := transaction{transaction:TxnHash}
				// s := fmt.Sprintf("%v", trans)
				identifer = Txn.Identifier
				encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
				text := (string(encoded))
				temp := model.TransactionIds{Txnhash: Txn.TxnHash,
					Url: "https://www.stellar.org/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
						text + "%3D%3D&network=public", Identifier: Txn.Identifier, TdpId: TDPs.TdpID[i]}

				resultArray = append(resultArray, temp)
			}

			// res := TDP{TdpId: result.TdpId}
			return nil
		}).Catch(func(error error) error {
			temp := model.TransactionIds{Txnhash: "Not Found", Url: "Not Found", Identifier: "Not Found", TdpId: TDPs.TdpID[i]}
			resultArray = append(resultArray, temp)

			return error
		})
		p.Await()

		q := object.GetPogTransaction(identifer)
		fmt.Println(identifer)
		q.Then(func(data interface{}) interface{} {

			// if TDPs.TdpID[i] == "" {
			// 	temp := model.TransactionIds{Txnhash: "Not Found", Url: "Not Found", Identifier: "Not Found", TdpId: TDPs.TdpID[i]}
			// 	resultArray = append(resultArray, temp)
			// } else {
			Txn := data.(model.TransactionCollectionBody)
			mapD := map[string]string{"transaction": Txn.TxnHash}
			mapB, _ := json.Marshal(mapD)
			fmt.Println(Txn.TxnHash)
			// trans := transaction{transaction:TxnHash}
			// s := fmt.Sprintf("%v", trans)

			encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
			text := (string(encoded))
			temp := model.TransactionIds{Txnhash: Txn.TxnHash,
				Url: "https://www.stellar.org/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
					text + "%3D%3D&network=public", Identifier: Txn.Identifier, TdpId: TDPs.TdpID[i]}

			resultArray = append(resultArray, temp)
			// }

			// res := TDP{TdpId: result.TdpId}
			return nil
		}).Catch(func(error error) error {
			temp := model.TransactionIds{Txnhash: "Not Found", Url: "Not Found", Identifier: "Not Found", TdpId: TDPs.TdpID[i]}
			resultArray = append(resultArray, temp)

			return error
		})
		q.Await()
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultArray)

}

func GetTransactionsForPK(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	var result []model.TransactionIds
	object := dao.Connection{}
	p := object.GetAllTransactionForPK(vars["id"])
	fmt.Println(vars["id"])
	p.Then(func(data interface{}) interface{} {
		res := data.([]model.TransactionCollectionBody)
		for _, TxnBody := range res {
			TxnHash := TxnBody.TxnHash

			mapD := map[string]string{"transaction": TxnHash}
			mapB, _ := json.Marshal(mapD)
			fmt.Println(string(mapB))
			// trans := transaction{transaction:TxnHash}
			// s := fmt.Sprintf("%v", trans)

			encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
			text := (string(encoded))
			temp := model.TransactionIds{Txnhash: TxnHash,
				Url: "https://www.stellar.org/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
					text + "%3D%3D&network=public",
				Identifier: TxnBody.Identifier, TdpId: TxnBody.TdpId}

			result = append(result, temp)
		}

		// res := TDP{TdpId: result.TdpId}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "TDP ID Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()

}

func QueryTransactionsByKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var result []model.TransactionIds

	vars := mux.Vars(r)
	object := dao.Connection{}

	switch checkValidVersionByte(vars["key"]) {
	case "pk":
		p := object.GetAllTransactionForPK(vars["key"])
		p.Then(func(data interface{}) interface{} {
			res := data.([]model.TransactionCollectionBody)
			for _, TxnBody := range res {
				TxnHash := TxnBody.TxnHash
				mapD := map[string]string{"transaction": TxnHash}
				mapB, _ := json.Marshal(mapD)
				fmt.Println(string(mapB))

				encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
				text := (string(encoded))
				temp := model.TransactionIds{Txnhash: TxnHash,
					Url: "https://www.stellar.org/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
						text + "%3D%3D&network=public",
					Identifier: TxnBody.Identifier, TdpId: TxnBody.TdpId}
				result = append(result, temp)
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(result)
			return nil
		}).Catch(func(error error) error {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Error{Message: "Public Key Not Found in Gateway DataStore"}
			json.NewEncoder(w).Encode(response)
			return error
		})
		p.Await()
	case "txn":
		q := object.GetAllTransactionForTxId(vars["key"])
		q.Then(func(data interface{}) interface{} {
			res := data.([]model.TransactionCollectionBody)
			for _, TxnBody := range res {
				TxnHash := TxnBody.TxnHash

				mapD := map[string]string{"transaction": TxnHash}
				mapB, _ := json.Marshal(mapD)
				fmt.Println(string(mapB))

				encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
				text := (string(encoded))
				temp := model.TransactionIds{Txnhash: TxnHash,
					Url: "https://www.stellar.org/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
						text + "%3D%3D&network=public",
					Identifier: TxnBody.Identifier, TdpId: TxnBody.TdpId}
				result = append(result, temp)
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(result)
			return nil
		}).Catch(func(error error) error {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Error{Message: "TxnHash Not Found in Gateway DataStore"}
			json.NewEncoder(w).Encode(response)
			return error
		})
		q.Await()

	case "tdpid":
		p := object.GetAllTransactionForTdpId(vars["key"])
		fmt.Println(vars["id"])
		p.Then(func(data interface{}) interface{} {
			res := data.([]model.TransactionCollectionBody)
			for _, TxnBody := range res {
				TxnHash := TxnBody.TxnHash

				mapD := map[string]string{"transaction": TxnHash}
				mapB, _ := json.Marshal(mapD)
				fmt.Println(string(mapB))

				encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
				text := (string(encoded))
				temp := model.TransactionIds{Txnhash: TxnHash,
					Url: "https://www.stellar.org/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
						text + "%3D%3D&network=public",
					Identifier: TxnBody.Identifier, TdpId: TxnBody.TdpId}

				result = append(result, temp)
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(result)
			return nil
		}).Catch(func(error error) error {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Error{Message: "TDP ID Not Found in Gateway DataStore"}
			json.NewEncoder(w).Encode(response)
			return error
		})
		p.Await()
	case "":
		p := object.GetTransactionsbyIdentifier(vars["key"])
		p.Then(func(data interface{}) interface{} {
			res := data.([]model.TransactionCollectionBody)
			for _, TxnBody := range res {
				TxnHash := TxnBody.TxnHash

				mapD := map[string]string{"transaction": TxnHash}
				mapB, _ := json.Marshal(mapD)
				fmt.Println(string(mapB))

				encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
				text := (string(encoded))
				temp := model.TransactionIds{Txnhash: TxnHash,
					Url: "https://www.stellar.org/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
						text + "%3D%3D&network=public",
					Identifier: TxnBody.Identifier, TdpId: TxnBody.TdpId}

				result = append(result, temp)
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(result)
			return nil
		}).Catch(func(error error) error {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Error{Message: "Identifier Not Found in Gateway DataStore"}
			json.NewEncoder(w).Encode(response)
			return error
		})
		p.Await()
	}
}

func RetriveTransactionId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	var result []model.TransactionIds
	object := dao.Connection{}
	p := object.GetAllTransactionForTxId(vars["id"])
	fmt.Println(vars["id"])
	p.Then(func(data interface{}) interface{} {
		res := data.([]model.TransactionCollectionBody)
		for _, TxnBody := range res {
			TxnHash := TxnBody.TxnHash

			mapD := map[string]string{"transaction": TxnHash}
			mapB, _ := json.Marshal(mapD)
			fmt.Println(string(mapB))
			// trans := transaction{transaction:TxnHash}
			// s := fmt.Sprintf("%v", trans)

			encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
			text := (string(encoded))
			temp := model.TransactionIds{Txnhash: TxnHash,
				Url: "https://www.stellar.org/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
					text + "%3D%3D&network=public",
				Identifier: TxnBody.Identifier, TdpId: TxnBody.TdpId}

			result = append(result, temp)
		}

		// res := TDP{TdpId: result.TdpId}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "TDP ID Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()

}


func checkValidVersionByte(key string) string {

	version, er := strkey.Version(key)
	if er != nil {
	}

	if version == strkey.VersionByteAccountID {
		return "pk"
	}

	if version == strkey.VersionByteSeed {
		return "sk"
	}

	// if version == strkey.VersionByteHashTx {
	// 	return "txn"
	// }

	// if version == strkey.VersionByteHashX {
	// return "hash"
	// }

	matched, err := regexp.MatchString(`^[0-9a-f]{64}$`, key)
	if err != nil {
	}
	if matched {
		return "txn"
	}
	matched1, err1 := regexp.MatchString(`^[0-9a-f]{24}$`, key)
	if err1 != nil {
	}
	if matched1 {
		return "tdpid"
	}
	return ""
}
