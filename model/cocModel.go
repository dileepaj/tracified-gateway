package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type COCState struct {
	Id                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	IssuerPublicKey   string             `json:"issuerpublickey" bson:"issuerpublickey"`
	SenderPublicKey   string             `json:"senderpublickkey" bson:"senderpublickkey"`
	COCAssetName      string             `json:"cocassetname" bson:"cocassetname"`
	ReciverPublickKey string             `json:"reciverpublickkey" bson:"reciverpublickkey"`
	CurrentCOCOwner   string             `json:"currentcocowner" bson:"currentcocowner"`
	COCStatus         uint8              `json:"cocstatus" bson:"cocstatus" `
	Timestamp         string             `json:"timestamp" bson:"timestamp"`
	TenantID          string             `json:"tenantid" bson:"tenantid"`
	TransferAmount    uint64             `json:"transferamount" bson:"transferamount"`
}

type UpdateCOCState struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	COCStatus uint8              `json:"cocstatus" bson:"cocstatuws" `
}

type COCStateDBResponse struct {
	Id        string `json:"id" bson:"_id,omitempty"`
	COCStatus uint8  `json:"cocstatus" bson:"cocstatus" `
}

const (
	COC_TRANSFER_ENABLED  = 1
	COC_TRANSFER_ACCPETED = 2
	COC_TRANSFER_REJECTED = 3
	COC_TRANSFER_MADE     = 4
)

const (
	COC_USERTYPE_SENDER  = 1
	COC_USERTYPE_RECIVER = 2
)
