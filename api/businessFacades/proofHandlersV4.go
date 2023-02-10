package businessFacades

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/retriever/stellarRetriever"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func BlockchainDataRetreiverWithHash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	var txns []model.TransactionCollectionBody
	for _, hash := range strings.Split(vars["txn"], "-") {
		st := stellarRetriever.ConcreteStellarTransaction{Txnhash: hash}
		txn, _ := st.GetTransactionCollection()
		txns = append(txns, *txn)
	}
	json.NewEncoder(w).Encode(txns)
}

func BlockchainTreeRetreiverWithHash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	s := stellarRetriever.POCTreeV4{TxnHash: vars["txn"]}
	s.ConstructPOC()
	pocData := s
	// create a new slice to hold the struct
	var pocNodes []stellarRetriever.POCNode
	// identifiers --> numbered identifiers (made by tracified)
	var identifiers []string

	// convert map[string]*POCNode to struct []stellarRetriever.POCNode
	// iterate through the map
	for _, node := range pocData.Nodes {
		// convert the struct to a JSON string
		jsonData, _ := json.Marshal(node)
		// create a new instance of the POCNode struct
		var pocNode stellarRetriever.POCNode
		// use json.Unmarshal to convert the JSON string to the struct
		json.Unmarshal(jsonData, &pocNode)
		identifiers = append(identifiers, pocNode.Data.Identifier)
		// add the struct to the slice
		pocNodes = append(pocNodes, pocNode)
	}
	// realIdentifiers identifiers that entered by the FO app (Not real)
	var realIdentifiers []apiModel.IdentifierModel
	object := dao.Connection{}
	p := object.GetIdentifierMap(identifiers)
	p.Then(func(data interface{}) interface{} {
		realIdentifiers = data.([]apiModel.IdentifierModel)
		return nil
	}).Catch(func(error error) error {
		logrus.Error("Identifier is not in the identifierMap DB", error)
		return nil
	}).Await()
	// replace rela identifiers with numbered identifiers if exist
	for i, node := range pocNodes {
		if len(realIdentifiers) > 0 {
			for _, realIdentifier := range realIdentifiers {
				if node.Data.Identifier == realIdentifier.MapValue && realIdentifier.MapValue != "" && realIdentifier.Identifier != "" {
					pocNodes[i].Data.Identifier = realIdentifier.Identifier
				}
			}
		}
	}
	// convert struct []stellarRetriever.POCNode to map[string]*POCNode
	for i, pocNode := range pocNodes {
		pocData.Nodes[pocNode.Id] = &pocNodes[i]
	}

	json.NewEncoder(w).Encode(pocData)
}
