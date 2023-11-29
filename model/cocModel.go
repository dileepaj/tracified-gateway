package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type COCState struct {
	Id                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	IssuerPublicKey   string             `json:"issuerpublickey" bson:"issuerpublickey"`
	SenderPublicKey   string             `json:"senderpublickkey" bson:"senderpublickkey"`
	COCAssetName      string             `json:"cocassetname" bson:"cocassetname"`
	ReceiverPublicKey string             `json:"receiverpublickey" bson:"receiverpublickey"`
	CurrentCOCOwner   string             `json:"currentcocowner" bson:"currentcocowner"`
	COCStatus         uint8              `json:"cocstatus" bson:"cocstatus" `
	Timestamp         string             `json:"timestamp" bson:"timestamp"`
	TenantID          string             `json:"tenantid" bson:"tenantid"`
	TransferAmount    uint64             `json:"transferamount" bson:"transferamount"`
	AssetType         string             `json:"assettype" bson:"assettype"`
	BatchName         string             `json:"batchname" bson:"batchname"`
	ProductName       string             `json:"productname" bson:"productname"`
	COCTxn            string 			 `json:"coctxn" bson:"coctxn"`
}

type UpdateCOCState struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	COCStatus uint8              `json:"cocstatus" bson:"cocstatuws" `
}

type COCStateDBResponse struct {
	Id        string `json:"id" bson:"_id,omitempty"`
	COCStatus uint8  `json:"cocstatus" bson:"cocstatus" `
}

type UpdateCOCOwner struct {
	Id              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CurrentCOCOwner string             `json:"currentcocowner" bson:"currentcocowner"`
	NewOwner        string             `json:"newowner" bson:"newowner"`
}

type COCPreviousOwner struct {
	Id              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CurrentCOCOwner string             `json:"currentcocowner" bson:"currentcocowner"`
	COCStatus       uint8              `json:"cocstatus" bson:"cocstatus" `
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
