package services

import (
	"time"
	"fmt"
	"github.com/stellar/go/xdr"
	// "github.com/stellar/go/clients/horizon"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
)

// CheckCOCStatus
func CheckCOCStatus() {
	// fmt.Println("NEW STUFF")
	object := dao.Connection{}
	p := object.GetCOCbyStatus("pending")
	p.Then(func(data interface{}) interface{} {
		result := data.([]model.COCCollectionBody)
		// temp:=result
		var temp model.COCCollectionBody
		temp.Status="expired"
		for i := 0; i < len(result); i++ {
			// fmt.Println(temp[i])

			if result[i].Status == "pending" {
				var txe xdr.Transaction
				err := xdr.SafeUnmarshalBase64(result[i].AcceptXdr, &txe)
				if err != nil {
					// fmt.Println(err)
				}
				// fmt.Println(i)
				// fmt.Println(txe.TimeBounds.MaxTime)
				if int64(txe.TimeBounds.MaxTime) < time.Now().Unix() {
					// result[i].Status="expired"
					err1:=object.UpdateCOC(result[i],temp)
					if err1!= nil{
						fmt.Println(err1)
					}
					fmt.Println("Expired")
				}else{
					// fmt.Println("Not Expired")
				}

			}
		}
		return nil
	}).Catch(func(error error) error {
		return error
	})
	p.Await()
}

func CheckCOCExpired() {
	fmt.Println("NEW STUFF")
	object := dao.Connection{}
	p := object.GetCOCbyStatus("pending")
	p.Then(func(data interface{}) interface{} {
		result := data.([]model.COCCollectionBody)
		// temp:=result
		var temp model.COCCollectionBody
		temp.Status="expired"
		for i := 0; i < len(result); i++ {
			// fmt.Println(temp[i])

			if result[i].Status == "pending" {
				var txe xdr.Transaction
				err := xdr.SafeUnmarshalBase64(result[i].AcceptXdr, &txe)
				if err != nil {
					// fmt.Println(err)
				}
				// fmt.Println(i)
				// fmt.Println(txe.TimeBounds.MaxTime)
				if int64(txe.TimeBounds.MaxTime) < time.Now().Unix() {
					// result[i].Status="expired"
					err1:=object.UpdateCOC(result[i],temp)
					if err1!= nil{
						fmt.Println(err1)
					}
					// fmt.Println("Expired")
				}else{
					// fmt.Println("Not Expired")
				}

			}
		}
		return nil
	}).Catch(func(error error) error {
		return error
	})
	p.Await()
}

