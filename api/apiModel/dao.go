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
}
