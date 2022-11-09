package services

import (
	//"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stellar/go/xdr"

	// "github.com/stellar/go/clients/horizon"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
)

// CheckCOCStatus
func CheckCOCStatus() {
	log.Debug("----------------------------------- CheckCOCStatus -------------------------------------")
	object := dao.Connection{}
	p := object.GetCOCbyStatus("pending")
	p.Then(func(data interface{}) interface{} {
		result := data.([]model.COCCollectionBody)
		// temp:=result
		var temp model.COCCollectionBody
		temp.Status="expired"
		for i := 0; i < len(result); i++ {
			if result[i].Status == "pending" {
				var txe xdr.Transaction
				err := xdr.SafeUnmarshalBase64(result[i].AcceptXdr, &txe)
				if err != nil {
					log.Error("Error @SafeUnmarshalBase64 @CheckCOCStatus" + err.Error())
				}
				if int64(txe.TimeBounds().MaxTime) < time.Now().Unix() {
					// result[i].Status="expired"
					err1:=object.UpdateCOC(result[i],temp)
					if err1!= nil{
						log.Error("Error @UpdateCOC" + err1.Error())
					}
					log.Info("Expired")
				}else{
					// fmt.Println("Not Expired")
				}

			}
		}
		return nil
	}).Catch(func(error error) error {
		log.Error("Error @GetCOCbyStatus " + error.Error())
		return error
	})
	p.Await()
}

func CheckCOCExpired() {
	log.Debug("---------------------------- CheckCOCExpired ----------------------------")
	object := dao.Connection{}
	p := object.GetCOCbyStatus("pending")
	p.Then(func(data interface{}) interface{} {
		result := data.([]model.COCCollectionBody)
		// temp:=result
		var temp model.COCCollectionBody
		temp.Status="expired"
		for i := 0; i < len(result); i++ {
			if result[i].Status == "pending" {
				var txe xdr.Transaction
				err := xdr.SafeUnmarshalBase64(result[i].AcceptXdr, &txe)
				if err != nil {
					log.Error("Error @SafeUnmarshalBase64 @CheckCOCExpired " + err.Error())
				}
				if int64(txe.TimeBounds().MaxTime) < time.Now().Unix() {
					// result[i].Status="expired"
					err1:=object.UpdateCOC(result[i],temp)
					if err1!= nil{
						log.Error("Error @UpdateCOC @CheckCOCExpired "+err1.Error())
					}
					// fmt.Println("Expired")
				}else{
					// fmt.Println("Not Expired")
				}

			}
		}
		return nil
	}).Catch(func(error error) error {
		log.Error("Error @GetCOCbyStatus @CheckCOCExpired "+error.Error())
		return error
	})
	p.Await()
}

