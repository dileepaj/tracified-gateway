package metricBinding

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

/*
des - Insert and find the workflow ID
*/

func InsertAndFindWorkflowId(workflowID string) (uint64, error) {
	var workflowId uint64
	object := dao.Connection{}

	workflowMap, errInWorkflowMap := object.GetWorkflowMapID(workflowID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errInWorkflowMap != nil {
		logrus.Error("Error when retrieving workflow ID " + errInWorkflowMap.Error())
	}
	if workflowMap != nil {
		logrus.Info("Workflow " + workflowID + " is already recorded in the DB")
		workflowMapData := workflowMap.(model.WorkflowMap)
		workflowId = workflowMapData.MapID
	} else {
		logrus.Error("Workflow " + workflowID + " is not recorded in the DB")
		data, errWhenGettingSequence := object.GetNextSequenceValue("WORKFLOWID")
		if errWhenGettingSequence != nil {
			logrus.Info("Error when taking sequence number " + errWhenGettingSequence.Error())
			return 0, errors.New("Error when taking sequence number " + errWhenGettingSequence.Error())
		}

		insertWorkflowDetails := model.WorkflowMap{
			WorkflowID: workflowID,
			MapID:      data.SequenceValue,
		}

		errWhenInsertingWorkflowDetails := object.InsertToWorkflowIDMap(insertWorkflowDetails)
		if errWhenInsertingWorkflowDetails != nil {
			logrus.Error("Inserting to workflow map ID was failed " + errWhenInsertingWorkflowDetails.Error())
			return 0, errors.New("Inserting to workflow map ID was failed " + errWhenInsertingWorkflowDetails.Error())
		}

		workflowId = data.SequenceValue
	}
	return workflowId, nil
}
