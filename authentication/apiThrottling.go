package authentication

import (
	"errors"
	"fmt"
	"time"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/relvacode/iso8601"
	"github.com/sirupsen/logrus"
)

//get the request entity, request entity type, limit
func API_Throttler(element model.API_ThrottlerRequest) (error, int) {
	object := dao.Connection{}
	var totalReqCount model.RequestCount
	// currentTime := time.Now()
	// layout := "2006-01-02T15:04:05Z"

	//get exisiting records for the requested entity type and entity
	allRecordsWithinTimeDuration, errWhenGettingRecords := object.GetRequestAmount(element.RequestEntityType, element.RequestEntity, element.FromTime, element.ToTime).Then(func(data interface{}) interface{} {
		totalReqCount = model.RequestCount{TotalRequestCount: data.(int64)}
		return data
	}).Await()
	if errWhenGettingRecords != nil {
		logrus.Info("Unable to connect gateway datastore to get API counters " + errWhenGettingRecords.Error())
	}
	fmt.Println(allRecordsWithinTimeDuration, "----", totalReqCount.TotalRequestCount)
	//all records within the array is empty or less than the allowed count
	if allRecordsWithinTimeDuration == nil || totalReqCount.TotalRequestCount < int64(element.AllowedAmount) {
		//convert time to ISO standard
		fmt.Println(time.Now().String())
		convertedFromTime, errWhenConvertingISO := iso8601.ParseString(time.Now().String())
		if errWhenConvertingISO != nil {
			logrus.Error("Error when converting to ISO standart time " + errWhenConvertingISO.Error())
			return errors.New("Error when converting to ISO standart time:  " + errWhenConvertingISO.Error()), -1
		}
		//build the new insert object
		insertNewReqObject := model.ThrottlerRecord{
			RequestEntityType: element.RequestEntityType,
			RequestEntity:     element.RequestEntity,
			Timestamp:         convertedFromTime,
			FormulaID:         element.FormulaID,
		}

		//allow to execute this and pass the transaction and insert it to the collection with the current timestamp
		errWhenInsertingNewReq := object.InsertToAPIThrottler(insertNewReqObject)
		if errWhenInsertingNewReq != nil {
			logrus.Error("Erron when inserting the new request to the counter " + errWhenInsertingNewReq.Error())
			return errWhenInsertingNewReq, -1
		}
	} else if allRecordsWithinTimeDuration != nil || totalReqCount.TotalRequestCount >= int64(element.AllowedAmount) {
		//the limit exceeded
		logrus.Error("API request exceeded the given request limit of ", element.AllowedAmount)
		return errors.New("API request exceeded the given request limit"), int(totalReqCount.TotalRequestCount)
	}

	return nil, int(totalReqCount.TotalRequestCount)
}
