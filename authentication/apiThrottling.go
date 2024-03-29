package authentication

import (
	"errors"
	"strconv"
	"time"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

// get the request entity, request entity type, limit
func APIThrottler(element model.API_ThrottlerRequest, insertToDB bool) (error, int, int) {
	object := dao.Connection{}
	var totalReqCount model.RequestCount
	// get existing records for the requested entity type and entity
	allRecordsWithinTimeDuration, err := object.GetRequestAmount(element.RequestEntityType, element.RequestEntity, element.FromTime, element.ToTime).Then(func(data interface{}) interface{} {
		totalReqCount = model.RequestCount{TotalRequestCount: data.(int64)}
		return data
	}).Await()
	if err != nil {
		logrus.Info("Unable to connect gateway datastore to get API counters " + err.Error())
		return errors.New("Unable to connect gateway datastore to get API counters "), 500, -1
	} else {
		// all records within the array is empty or less than the allowed count
		if allRecordsWithinTimeDuration == nil || totalReqCount.TotalRequestCount < int64(element.AllowedAmount) {
			// current location assume to UTC
			time.Local = time.UTC
			insertNewReqObject := model.ThrottlerRecord{
				RequestEntityType: element.RequestEntityType,
				RequestEntity:     element.RequestEntity,
				Timestamp:         time.Now(),
				FormulaID:         element.FormulaID,
			}
			if insertToDB {
				// allow to execute this and pass the transaction and insert it to the collection with the current timestamp
				errWhenInsertingNewReq := object.InsertToAPIThrottler(insertNewReqObject)
				if errWhenInsertingNewReq != nil {
					logrus.Error("Error when inserting the new request to the API throttler " + errWhenInsertingNewReq.Error())
					return errors.New("Error when inserting the new request to the API throttler " + errWhenInsertingNewReq.Error()), 500, -1
				}
			}
		} else if allRecordsWithinTimeDuration != nil || totalReqCount.TotalRequestCount >= int64(element.AllowedAmount) {
			// the limit exceeded
			logrus.Error("API request exceeded the given request limit of ", element.AllowedAmount)
			return errors.New("API request exceeded the given request limit of " + strconv.Itoa(element.AllowedAmount)), 429, int(totalReqCount.TotalRequestCount)
		}
		return nil, 200, int(totalReqCount.TotalRequestCount)
	}
}
