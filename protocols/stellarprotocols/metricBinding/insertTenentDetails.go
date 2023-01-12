package metricBinding

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

/*
 For retrieving mapped tenant id from DB or creating a new one with the next sequence value
*/

func InsertAndFindTenantID(tenentUUID string) (uint64, error) {
	var tenentMapID uint64
	object := dao.Connection{}

	tenentMap, errInTenentMap := object.GetTenentMapID(tenentUUID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errInTenentMap != nil {
		logrus.Info("Error when retrieving tenent ID " + errInTenentMap.Error())
	}
	if tenentMap != nil {
		logrus.Info("Tenent ID " + tenentUUID + " is recorded in the DB")
		tenentMapData := tenentMap.(model.TenentMapDetails)
		tenentMapID = tenentMapData.MapID
	} else {
		logrus.Error("Tenent ID " + tenentUUID + " is not recorded in the DB")
		data, errWhenGettingTheSequence := object.GetNextSequenceValue("TENENTID")
		if errWhenGettingTheSequence != nil {
			logrus.Error("Error when taking the sequence no Error : " + errWhenGettingTheSequence.Error())
			return 0, errors.New("error when taking the sequence no Error : " + errWhenGettingTheSequence.Error())
		}

		insertTenentDetails := model.TenentMapDetails{
			TenentID: tenentUUID,
			MapID:    data.SequenceValue,
		}

		errWhenInsertingTenentDetails := object.InsertTenentMapID(insertTenentDetails)
		if errWhenInsertingTenentDetails != nil {
			logrus.Error("Inserting to tenent map ID was failed" + errWhenInsertingTenentDetails.Error())
			return 0, errors.New("inserting to tenent map ID was failed" + errWhenInsertingTenentDetails.Error())
		}

		tenentMapID = data.SequenceValue
	}

	return tenentMapID, nil
}
