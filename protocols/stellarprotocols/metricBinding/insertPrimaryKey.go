package metricBinding

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

/*
 For retrieving mapped artifact id from DB or creating a new one with the next sequence value
*/
func InsertAndFindPrimaryKeyID(primaryKeyId string) (uint64, error) {
	var primaryKeyIdMapValue uint64
	object := dao.Connection{}
	artifactMap, errInArtifactMap := object.GetPrimaryKeyMapID(primaryKeyId).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errInArtifactMap != nil {
		logrus.Info("Error when retrieving artifact id from DB(insertArtifactDetails.go) " + errInArtifactMap.Error())
	}
	if artifactMap == nil {
		logrus.Error("Artifact ID " + primaryKeyId + " is not recorded in the DB(insertArtifactDetails.go)")
		data, errWhenGettingTheSequence := object.GetNextSequenceValue("PRIMARYKEYID")
		if errWhenGettingTheSequence != nil {
			logrus.Error("Error when taking the sequence number(insertArtifactDetails.go) Error : " + errWhenGettingTheSequence.Error())
			return 0, errors.New("error when taking the sequence no Error : " + errWhenGettingTheSequence.Error())
		}
		insertArtifactMap := model.PrimaryKeyMap{
			PrimaryKeyID: primaryKeyId,
			MapID:        data.SequenceValue,
		}
		errWhenInsertingToResourceMap := object.InsertToPrimaryKeyIdMap(insertArtifactMap)
		if errWhenInsertingToResourceMap != nil {
			logrus.Error("Inserting to primary key map ID was failed(insertArtifactDetails.go)" + errWhenInsertingToResourceMap.Error())
			return 0, errors.New("inserting to artifact map ID was failed")
		}
		primaryKeyIdMapValue = data.SequenceValue
	} else {
		logrus.Info("Primary key ID " + primaryKeyId + " is recorded in the DB")
		artifactMapData := artifactMap.(model.ArtifactTemplateId)
		primaryKeyIdMapValue = artifactMapData.MapID
	}
	return primaryKeyIdMapValue, nil
}
