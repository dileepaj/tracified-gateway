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
		logrus.Info("Error when retrieving artifact id from DB(insertArtifactDetails.go) " + errInArtifactMap.Error())
	}
	if artifactMap == nil {
		logrus.Error("Artifact ID " + artifactID + " is not recorded in the DB(insertArtifactDetails.go)")
		data, errWhenGettingTheSequence := object.GetNextSequenceValue("ARTIFACTID")
		if errWhenGettingTheSequence != nil {
			logrus.Error("Error when taking the sequence number(insertArtifactDetails.go) Error : " + errWhenGettingTheSequence.Error())
			return 0, errors.New("error when taking the sequence no Error : " + errWhenGettingTheSequence.Error())
		}
		insertArtifactMap := model.ArtifactIDMap{
			ArtifactID:   artifactID,
			MapID:      data.SequenceValue,
		}
		errWhenInsertingToResourceMap := object.InsertToArtifactIDMAP(insertArtifactMap)
		if errWhenInsertingToResourceMap != nil {
			logrus.Error("Inserting to artifact map ID was failed(insertArtifactDetails.go)" + errWhenInsertingToResourceMap.Error())
			return 0, errors.New("inserting to artifact map ID was failed")
		}
		artifactMapID = data.SequenceValue
	} else {
		logrus.Info("Artifact ID " + artifactID + " is recorded in the DB")
		artifactMapData := artifactMap.(model.ArtifactIDMap)
		artifactMapID = artifactMapData.MapID
	}
	return artifactMapID, nil
}