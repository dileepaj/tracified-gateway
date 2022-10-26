package metricBinding

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

/*
 For retrieving mapped metric id from DB or creating a new one with the next sequence value
*/
func InsertAndFindMetricID(metricID string, metricName string) (uint64,int, error) {
	var metricMapID uint64
	object := dao.Connection{}
	metricMap, errInMetricMap := object.GetMetricMapID(metricID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errInMetricMap != nil {
		logrus.Info("Error when retrieving metric id from DB " + errInMetricMap.Error())
	}
	if metricMap == nil {
		logrus.Error("Metric ID " + metricID + " is not recorded in the DB")
		data, errWhenGettingTheSequence := object.GetNextSequenceValue("METRICID")
		if errWhenGettingTheSequence != nil {
			logrus.Error("Error when taking the sequence no Error : " + errWhenGettingTheSequence.Error())
			return 0,500, errors.New("Error when taking the sequence no Error: " + errWhenGettingTheSequence.Error())
		}
		insertMerticMap := model.MetricMapDetails{
			MetricID:   metricID,
			MetricName: metricName,
			MapID:      data.SequenceValue,
		}
		errWhenInsertingToResourceMap := object.InsertMetricMapID(insertMerticMap)
		if errWhenInsertingToResourceMap != nil {
			logrus.Error("Inserting to metric map ID was failed" + errWhenInsertingToResourceMap.Error())
			return 0,500, errors.New("Inserting to metric map ID was failed Error: " + errWhenInsertingToResourceMap.Error())
		}
		metricMapID = data.SequenceValue
	} else {
		logrus.Info("Metric ID " + metricID + " is recorded in the DB")
		metricMapData := metricMap.(model.MetricMapDetails)
		metricMapID = metricMapData.MapID
	}
	return metricMapID,200, nil
}