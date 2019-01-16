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
	Identifier   string `json:"Identifier"`
	TdpId string `json:"TdpId"`
	// ProfileID      string
	TxnHash string `json:"TxnHash"`
	// ProfileHash    string
	PublicKey string `json:"PublicKey"`
	TxnType string `json:"TxnType"`
	XDR string `json:"XDR"`
	Status string `json:"Status"`
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