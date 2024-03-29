package stellarRetriever

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/stellar/go/xdr"
)

type ConcreteStellarTransaction struct {
	Txnhash string
}

func (stxn *ConcreteStellarTransaction) RetrieveTransaction() (*model.StellarTransaction, error) {
	result, err := http.Get(commons.GetHorizonClient().HorizonURL + "transactions/" + stxn.Txnhash)
	if err != nil {
		return nil, err
	}
	data, _ := ioutil.ReadAll(result.Body)
	if result.StatusCode == 200 {
		var txn model.StellarTransaction
		error := json.Unmarshal(data, &txn)
		if error != nil {
			return nil, error
		}
		return &txn, nil
	}
	return nil, errors.New("Transaction is not valid")
}

func (stxn *ConcreteStellarTransaction) RetrieveOperations() (*model.StellarOperations, error) {
	result, err := http.Get(commons.GetHorizonClient().HorizonURL + "transactions/" + stxn.Txnhash + "/operations")
	if err != nil {
		return nil, err
	}
	data, _ := ioutil.ReadAll(result.Body)
	if result.StatusCode == 200 {
		var oprn model.StellarOperations
		error := json.Unmarshal(data, &oprn)
		if error != nil {
			return nil, error
		}
		return &oprn, nil
	}
	return nil, errors.New("Transaction is not valid")
}

func (stxn *ConcreteStellarTransaction) GetTransactionCollection() (*model.TransactionCollectionBody, error) {
	var txe model.TransactionCollectionBody
	txn, err := stxn.RetrieveTransaction()
	if err != nil {
		return nil, err
	}
	oprn, err := stxn.RetrieveOperations()
	if err != nil {
		return nil, err
	}
	
	sqnc, err := strconv.Atoi(txn.SourceAccountSequence)
	if err == nil {
		txe.SequenceNo = int64(sqnc)
	}
	txe.PublicKey = txn.SourceAccount
	txe.TxnHash = txn.Hash
	txe.XDR = txn.EnvelopeXdr
	txe.Status = "stellar-" + commons.GetHorizonClientNetworkName()
	mapAPIOperations(&txe, *oprn)
	// checked whether the transaction type is a parent. Parent transaction data does not have mange data called Identifier
	if txe.TxnType == "5" && txe.ProfileID != "parent" {
		txe.Identifier = txe.FromIdentifier1
	}

	return &txe, nil
}

func MapXDROperations(txe *model.TransactionCollectionBody, md []xdr.Operation) (*model.TransactionCollectionBody) {
	for _, operation := range md {
		if operation.Body.Type != xdr.OperationTypeManageData {
			continue
		} else {
			mapDataOperations(txe, string(*&operation.Body.ManageDataOp.DataName), string(*operation.Body.ManageDataOp.DataValue), false)
		}
	}
	return txe
}

func mapAPIOperations(txe *model.TransactionCollectionBody, sp model.StellarOperations) (*model.TransactionCollectionBody) {
	for _, record := range sp.Embedded.Records {
		if strings.ToLower(record.Type) != "manage_data" {
			continue
		} else {
			mapDataOperations(txe, record.Name, record.Value, true)
		}
	}
	return txe
}
		
func mapDataOperations(txe *model.TransactionCollectionBody, dataType string, dataValue string, isBase64Value bool) (*model.TransactionCollectionBody) {
	if isBase64Value {
		decoded, err := base64.StdEncoding.DecodeString(dataValue)
		if err == nil {
			dataValue = string(decoded)
		}
	}
	switch strings.ToLower(strings.ReplaceAll(dataType, " ", "")) {
		case "type": 
			re := regexp.MustCompile("[0-9]+")
			txe.TxnType = re.FindAllString(dataValue, -1)[0]
			break
		case "identifier": 
				txe.Identifier = dataValue
			break
		case "productname": 
			txe.ProductName = dataValue
			break
		case "productid": 
			txe.ProductID = dataValue
			break
		case "appaccount": 
			txe.AppAccount = dataValue
			break
		case "datahash": 
			txe.DataHash = dataValue
			break
		case "fromidentifier": 
			txe.FromIdentifier1 = dataValue
			break
		case "toidentifiers": 
			txe.ToIdentifier = dataValue
			break
		case "fromidentifier1": 
			txe.FromIdentifier1 = dataValue
			break
		case "fromidentifier2": 
			txe.FromIdentifier2 = dataValue
			break
		case "assetcode": 
			txe.ItemCode = dataValue
			break
		case "assetamount": 
			txe.ItemAmount = dataValue
			break
		case "previousstage": 
			txe.PreviousStage = dataValue
			break
		case "currentstage": 
			txe.CurrentStage = dataValue
			break
		case "currenttxn": 
			txe.CurrentTxnHash = dataValue
			break
		case "previoustxn": 
			txe.PreviousTxnHash = dataValue
			break
		case "profileid": 
			txe.ProfileID = dataValue
			break
		case "previousprofile": 
			txe.PreviousSplitProfile = dataValue
			break
		case "mergeid": 
			txe.MergeID = dataValue
			break
		case "typename": 
			txe.TypeName = dataValue
			break
		case "geolocation": 
			txe.GeoLocation = dataValue
			break
		case "timestamp": 
			txe.Timestamp = dataValue
			break
		case "tenantname": 
			txe.TenantNameBase64 = dataValue
			break
		case "tenantnamebase64": 
			txe.TenantNameBase64 = dataValue
			break
		case "productnamebase64": 
			txe.ProductName = dataValue
			break
		case "batchname": 
			txe.RealIdentifier = dataValue
			break
		case "tenantid": 
			txe.TenantID = dataValue
			break
	}
	return txe
}
