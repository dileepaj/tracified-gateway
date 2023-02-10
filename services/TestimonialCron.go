package services

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stellar/go/xdr"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
)

func CheckTestimonialStatus() {
	if (commons.GoDotEnvVariable("LOGSTYPE")=="DEBUG"){
	log.Debug("----------------------------------- CheckTestimonialStatus -------------------------------------")
	}
	object := dao.Connection{}
	p := object.GetTestimonialbyStatus(model.Pending.String())
	p.Then(func(data interface{}) interface{} {
		result := data.([]model.Testimonial)
		// temp:=result
		var temp model.Testimonial
		temp.Status = model.Expired.String()
		for i := 0; i < len(result); i++ {
			// fmt.Println(temp[i])

			if result[i].Status == model.Pending.String() {
				var txe xdr.Transaction
				err := xdr.SafeUnmarshalBase64(result[i].AcceptXDR, &txe)
				if err != nil {
					if (commons.GoDotEnvVariable("LOGSTYPE")=="CRON"){
						log.Error("Error @SafeUnmarshalBase64 @CheckTestimonialStatus" + err.Error())
					}
				}
				// fmt.Println(i)
				// fmt.Println(txe.TimeBounds.MaxTime)
				if int64(txe.TimeBounds().MaxTime) < time.Now().Unix() {
					// result[i].Status="expired"
					err1 := object.UpdateTestimonial(result[i], temp)
					if err1 != nil {
						if (commons.GoDotEnvVariable("LOGSTYPE")=="CRON"){
							log.Error("Error @UpdateTestimonial" + err1.Error())
						}
					}
				} else {
					// fmt.Println("Not Expired")
				}

			}
		}
		return nil
	}).Catch(func(error error) error {
		if (commons.GoDotEnvVariable("LOGSTYPE")=="CRON"){
			log.Error("Error @GetTestimonialbyStatus " + error.Error())
		}
		return error
	})
	p.Await()
}
