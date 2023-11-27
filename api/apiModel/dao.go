package apiModel

import (
	"github.com/dileepaj/tracified-gateway/model"
)

type GetCOCCollection struct {
	Collection model.COCCollectionBody
}

type GetMultiCOCCollection struct {
	Collection []model.COCCollectionBody
}

type InsertCOCCollectionResponse struct {
	Message string
	Body    model.COCCollectionBody
}

type InsertTransactionCollectionResponse struct {
	Message string
	Body    model.TransactionCollectionBody
}

type GetSubAccountStatus struct {
	User        string   `json:"user"`
	SubAccounts []string `json:"subAccounts"`
	Receivers   []string `json:"receivers"`
}

type GetTransactionId struct {
	TdpID []string `json:"TdpID"`
}

type GetSubAccountStatusResponse struct {
	// Message string `json:"message"`
	SubAccount string `json:"subAccount"`
	Receiver   string `json:"receiver"`
	SequenceNo string `json:"sequenceNo"`
	Available  bool   `json:"available"`
	Expiration bool   `json:"expiration"`
	Operation  string `json:"operation"`
}

type InsertorganizationCollectionResponse struct {
	Message string
}

type InsertTestimonialCollectionResponse struct {
	Message string
}

// mapvalues - numbered identifiers (tracified mapped real identifier to unique id)
// identifiers - real identifier
type IdentifierModel struct {
	MapValue   string
	Identifier string
	Type       string
	ProductId  string
}

type TDPOperationRequest struct {
	Index      int
	Identifier string
	TDPID      string
	XDR        string
	Status     string
}

type CommonBadResponse struct {
	Message string
}
type CommonSuccessMessage struct {
	Message string
}

type COCTransferResponse struct {
	Response []model.COCState
}
