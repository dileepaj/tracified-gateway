package apiModel

import(
	"main/model"
)

type GetCOCCollection struct {
	Collection model.COCCollectionBody
}

type GetMultiCOCCollection struct {
	Collection []model.COCCollectionBody
}

type InsertCOCCollectionResponse struct {
	Message string
	Body model.COCCollectionBody
}


type InsertTransactionCollectionResponse struct {
	Message string
	Body model.TransactionCollectionBody
}