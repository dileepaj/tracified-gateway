package metricBinding

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

func InsertAndFindActivityID(activityId string, activityName string, metricId string, stageId string) (uint64, error) {
	var acitivityMapID uint64
	object := dao.Connection{}

	activityMap, errInActivityMap := object.GetActivityMapID(activityId).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errInActivityMap != nil {
		logrus.Info("Error when retrieving activity ID " + errInActivityMap.Error())
	}
	if activityMap != nil {
		logrus.Info("Activity ID " + activityId + " is recorded in the DB")
		activityMapData := activityMap.(model.ActivityMapDetails)
		acitivityMapID = activityMapData.MapID
	} else {
		logrus.Error("Activity ID " + activityId + " is not recorded in the DB")
		data, errWhenGettingTheSequence := object.GetNextSequenceValue("ACTIVITYID")
		if errWhenGettingTheSequence != nil {
			logrus.Error("Error when taking the sequence no Error : " + errWhenGettingTheSequence.Error())
			return 0, errors.New("Error when taking the sequence no Error : " + errWhenGettingTheSequence.Error())
		}

		insertActivityDetails := model.ActivityMapDetails{
			ActivityID: activityId,
			Name:       activityName,
			MetricID:   metricId,
			StageID:    stageId,
			MapID:      data.SequenceValue,
		}

		errWhenInsertActivityDetails := object.InsertActivityID(insertActivityDetails)
		if errWhenInsertActivityDetails != nil {
			logrus.Error("Inserting to activity map ID was failed" + errWhenInsertActivityDetails.Error())
			return 0, errors.New("Inserting to activity map ID was failed")
		}

		acitivityMapID = data.SequenceValue
	}

	return acitivityMapID, nil
}
