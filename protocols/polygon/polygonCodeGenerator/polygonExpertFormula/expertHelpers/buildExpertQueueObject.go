package experthelpers

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/services/rabbitmq"
	"github.com/dileepaj/tracified-gateway/utilities"
)

func BuildExpertQueueObjectAndSendToQueue(formulaObj model.EthereumExpertFormula, queueType string, status string) error {
	logger := utilities.NewCustomLogger()
	queueObject := model.SendToQueue{
		EthereumExpertFormula: formulaObj,
		Type:                  queueType,
		Status:                status,
	}

	//add to queue
	errWhenSendingToQueue := rabbitmq.SendToQueue(queueObject)
	if errWhenSendingToQueue != nil {
		logger.LogWriter("Error when sending the polygon request to queue "+errWhenSendingToQueue.Error(), constants.ERROR)
		return errors.New(errWhenSendingToQueue.Error())
	}

	return nil
}
