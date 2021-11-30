package businessFacades

import (
	"encoding/json"
	"net/http"
	"strings"
	"github.com/gorilla/mux"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/retriever/stellarRetriever"
)

func BlockchainDataRetreiverWithHash(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	var txns []model.TransactionCollectionBody
	for _, hash := range strings.Split(vars["txn"], "-") {
		st := stellarRetriever.ConcreteStellarTransaction{Txnhash: hash}
		txn, _ := st.GetTransactionCollection()
		txns = append(txns, *txn)
	}
	json.NewEncoder(w).Encode(txns);
}

func BlockchainTreeRetreiverWithHash(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	s := stellarRetriever.POCTreeV4{TxnHash: vars["txn"]}
	s.ConstructPOC()
	json.NewEncoder(w).Encode(s);
}
