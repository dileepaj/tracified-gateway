package ethereuemmetricbind

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/vendor/github.com/sirupsen/logrus"
)

// to get the metric map id from the DB

func GetMetricMapId(metricID string) (uint64, error) {
	var metricMapID uint64
	object := dao.Connection{}

	metricIDMap, errInMetricIDMap := object.GetEthMetricMapID(metricID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errInMetricIDMap != nil {
		logrus.Info("Error when retrieving metric id from DB. Error: " + errInMetricIDMap.Error())
	}
	if metricIDMap == nil {
		logrus.Error("Artifact ID " + metricID + " is not recorded in the DB")
		data, errWhenGettingTheSequence := object.GetNextSequenceValue("ETHMETRICID")
		if errWhenGettingTheSequence != nil {
			logrus.Error("Error when taking the sequence number. Error : " + errWhenGettingTheSequence.Error())
			return 0, errors.New("error when taking the sequence no Error : " + errWhenGettingTheSequence.Error())
		}
		metricIDMapDetails := model.EthMetricIDMap{
			MetricID: metricID,
			MapID:    data.SequenceValue,
		}
		errWhenInsertingToMetricIdMap := object.InsertEthMetricIDMap(metricIDMapDetails)
		if errWhenInsertingToMetricIdMap != nil {
			logrus.Error("Inserting to metric ID map was failed. Error: " + errWhenInsertingToMetricIdMap.Error())
			return 0, errors.New("inserting metric map ID was failed")
		}
		metricMapID = data.SequenceValue
	} else {
		logrus.Info("Metric ID " + metricID + " is recorded in the DB")
		metricIDMapData := metricIDMap.(model.EthMetricIDMap)
		metricMapID = metricIDMapData.MapID
	}
	return metricMapID, nil
}
