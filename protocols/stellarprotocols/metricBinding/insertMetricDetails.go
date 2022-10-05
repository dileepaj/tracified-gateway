package metricbinding

import (
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/sirupsen/logrus"
)

func InsertAndFindMetricID(metricID string, metricName string) (uint64, error) {
	var metricMapID uint64
	object := dao.Connection{}

	metricMap, errInMetricMap := object.GetMetricMapID(metricID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errInMetricMap != nil {
		logrus.Error("Error when retrieving metric id from DB " + errInMetricMap.Error())
	}
	if metricMap == nil {
		// logrus.Error("Metric ID is not recorded in the DB")
		// data, errWhenGettingTheSequence := object.GetNextSequenceValue("METRICID")
		// if errWhenGettingTheSequence != nil {
		// 	logrus.Error("Error when taking the sequence no Error : " + errWhenGettingTheSequence.Error())
		// 	return 0, errors.New("Error when taking the sequence no Error : " + errWhenGettingTheSequence.Error())
		// }

		//insert to metric map
		// insertMerticMap := model.MetricMapDetails{
		// 	MetricID:   metricID,
		// 	MetricName: metricName,
		// 	MapID:      data.SequenceValue,
		// }

		//errWhenInsertingToResourceMap := object.InsertMetricMapID(insertMerticMap)

	}

	return metricMapID, nil
}
