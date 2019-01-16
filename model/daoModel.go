package model
import(
	"github.com/stellar/go/build"

)

type COCCollectionBody struct {
	TxnHash string
	Sender   string
	SubAccount string
	SequenceNo int
	Receiver string
	AcceptXdr string
	RejectXdr string
	AcceptTxn string
	RejectTxn string
	Identifier string
	Status    string
	
}
type TransactionCollectionBody struct {
	Identifier   string `json:"identifier"`
	TdpId string `json:"tdpId"`
	// ProfileID      string
	TxnHash string `json:"txnHash"`
	// ProfileHash    string
	PublicKey string `json:"publicKey"`
	TxnType string `json:"txnType"`
	XDR string `json:"xdr"`
	Status string `json:"status"`
}
type XDR struct{
	XDR build.TransactionMutator
}


type LastTxnResponse struct {

	LastTxn string
	
}

type COCCollectionList struct {
	List []COCCollectionBody
}
type TransactionCollectionList struct {
	List []TransactionCollectionBody
}
type TransactionUpdate struct{
	Selector TransactionCollectionBody
	Update TransactionCollectionBody
}