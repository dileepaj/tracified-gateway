package services

import (
	"fmt"
	"github.com/stellar/go/xdr"
	// "github.com/stellar/go/clients/horizon"

	"github.com/tracified-gateway/dao"
	"github.com/tracified-gateway/model"
)

func CheckCOCStatus() {
	fmt.Println("NEW STUFF")
	object := dao.Connection{}
	p := object.GetCOCbyStatus("pending")
	p.Then(func(data interface{}) interface{} {
		result := data.([]model.COCCollectionBody)
		// temp:=result
		for i := 0; i < len(result); i++ {
			if result[i].Status == "pending" {
				var txe xdr.Transaction
				err := xdr.SafeUnmarshalBase64(result[i].AcceptXdr, &txe)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(i)
				fmt.Println(txe.TimeBounds.MaxTime)

			}
		}
		return nil
	}).Catch(func(error error) error {
		return error
	})
	p.Await()
}
