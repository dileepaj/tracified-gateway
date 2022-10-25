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
func InsertAndFindArtifactID(artifactID string) (uint64, error) {
	var artifactMapID uint64
	object := dao.Connection{}
	artifactMap, errInArtifactMap := object.GetArtifactMapID(artifactID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errInArtifactMap != nil {
		logrus.Info("Error when retrieving artifact id from DB " + errInArtifactMap.Error())
	}
	if artifactMap == nil {
		logrus.Error("Artifact ID " + artifactID + " is not recorded in the DB")
		data, errWhenGettingTheSequence := object.GetNextSequenceValue("ARTIFACTID")
		if errWhenGettingTheSequence != nil {
			logrus.Error("Error when taking the sequence no Error : " + errWhenGettingTheSequence.Error())
			return 0, errors.New("Error when taking the sequence no Error : " + errWhenGettingTheSequence.Error())
		}
		insertArtifactMap := model.ArtifactIDMap{
			ArtifactID:   artifactID,
			MapID:      data.SequenceValue,
		}
		errWhenInsertingToResourceMap := object.InsertToArtifactIDMAP(insertArtifactMap)
		if errWhenInsertingToResourceMap != nil {
			logrus.Error("Inserting to artifact map ID was failed" + errWhenInsertingToResourceMap.Error())
			return 0, errors.New("Inserting to artifact map ID was failed")
		}
		artifactMapID = data.SequenceValue
	} else {
		logrus.Info("Artifact ID " + artifactID + " is recorded in the DB")
		artifactMapData := artifactMap.(model.ArtifactIDMap)
		artifactMapID = artifactMapData.MapID
	}
	return artifactMapID, nil
}