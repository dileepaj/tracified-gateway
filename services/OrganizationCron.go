package services

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stellar/go/xdr"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
)

func CheckOrganizationStatus() {
	log.Debug("----------------------------------- CheckOrganizationStatus -------------------------------------")
	// fmt.Println("NEW STUFF")
	object := dao.Connection{}
	p := object.GetTestimonialOrganizationbyStatus(model.Pending.String())
	p.Then(func(data interface{}) interface{} {
		result := data.([]model.TestimonialOrganization)
		// temp:=result
		var temp model.TestimonialOrganization
		temp.Status = model.Expired.String()
		for i := 0; i < len(result); i++ {
			// fmt.Println(temp[i])

			if result[i].Status == model.Pending.String() {
				var txe xdr.Transaction
				err := xdr.SafeUnmarshalBase64(result[i].AcceptXDR, &txe)
				if err != nil {
					log.Error("Error @SafeUnmarshalBase64 @CheckOrganizationStatus" + err.Error())
				}
				// fmt.Println(i)
				// fmt.Println(txe.TimeBounds.MaxTime)
				if int64(txe.TimeBounds.MaxTime) < time.Now().Unix() {
					// result[i].Status="expired"
					err1 := object.Updateorganization(result[i], temp)
					if err1 != nil {
						log.Error("Error @UpdateOrganization" + err1.Error())
					}
					log.Info("Expired")
				} else {
					// fmt.Println("Not Expired")
				}

			}
		}
		return nil
	}).Catch(func(error error) error {
		log.Error("Error @GetOrganizationbyStatus " + error.Error())
		return error
	})
	p.Await()
}
