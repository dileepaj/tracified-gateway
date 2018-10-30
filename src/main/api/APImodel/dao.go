package apiModel

import(
	"main/model"
)

type GetCOCCollection struct {
	Collection model.COCCollectionBody
}

type InsertCOCCollectionResponse struct {
	Message string
}