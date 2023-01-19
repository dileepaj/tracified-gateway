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

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
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
	object := dao.Connection{} 
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
			p:=object.GetIdentifMap(dataValue);
			p.Then(func(data interface{}) interface{} {
				mapIdentifier:=data.(apiModel.IdentifierModel)
				txe.Identifier =mapIdentifier.Identifier
				return nil
			}).Catch(func(error error) error {
				txe.Identifier = dataValue
				logrus.Error("Identifier is not in the identierMap DB",error)
				return nil
			}).Await()
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
		case "realidentifer":
			txe.RealIdentifier = dataValue
			break
		case "mapidentifier":
			txe.MapIdentifier = dataValue
			break
		case "mapidentifier1":
			txe.MapIdentifier1 = dataValue
			break
		case "mapidentifier2":
			txe.MapIdentifier2 = dataValue
			break
	}
	return txe
}
