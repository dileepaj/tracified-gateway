package model
import(
	"github.com/stellar/go/build"

)

type COCCollectionBody struct {
	Sender   string
	Receiver string
	Xdr      string
	Identifier string
	Status    string
	
}
type TransactionCollectionBody struct {
	Identifier   string
	TdpID string
	// ProfileID      string
	TxnHash string
	// ProfileHash    string
	PublicKey string
	TxnType string
	XDR string
	Status string
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