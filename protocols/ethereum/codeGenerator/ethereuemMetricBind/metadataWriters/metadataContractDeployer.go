package metadatawriters

import (
	"encoding/base64"
	"os"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum/deploy"
	"github.com/dileepaj/tracified-gateway/services"
	"github.com/sirupsen/logrus"
)

//! Relevant map ID will be checked on the route handler to see whether the metric contract is already deployed or not
//! All the other failed database calls will be handled in the handler function
func MetricMetadataContractDeployer(element model.MetricMetadataReq, metricMapID string) error {

	reqType := "METRIC"

	//generate the contract
	metadataContractTemplate, errWhenGeneratingContract := MetricMetadataWriter(element, metricMapID)
	if errWhenGeneratingContract != nil {
		logrus.Error("Error when generating contract: ", errWhenGeneratingContract)
		return errWhenGeneratingContract
	}

	contractName := "Metric_" + metricMapID + "_Metadata"
	fo, errInOutput := os.Create(commons.GoDotEnvVariable("METRICCONTRACTLOCATION") + "/" + contractName + `.sol`)
	if errInOutput != nil {
		logrus.Error("Error when generating metadata contract file: ", errInOutput)
		return errInOutput
	}
	defer fo.Close()
	_, errInWritingOutput := fo.Write([]byte(metadataContractTemplate))
	if errInWritingOutput != nil {
		logrus.Error("Error when writing into the metadata contract file: ", errInWritingOutput)
		return errInWritingOutput
	}

	//generate ABI
	abiString, errWhenGeneratingABI := deploy.GenerateABI(contractName, reqType)
	if errWhenGeneratingABI != nil {
		logrus.Error("Error when generating ABI for metric metadata contract : ", errWhenGeneratingABI)
		return errWhenGeneratingABI
	}

	//generate BIN
	binString, errWhenGeneratingBIN := deploy.GenerateBIN(contractName, reqType)
	if errWhenGeneratingBIN != nil {
		logrus.Error("Error when generating BIN for metric metadata contract : ", errWhenGeneratingBIN)
		return errWhenGeneratingBIN
	}

	templateB64 := base64.StdEncoding.EncodeToString([]byte(metadataContractTemplate))

	ethMetricMetadataObj := model.EthereumMetricBind{
		Type:              "METADATA",
		MetricID:          element.MetricId,
		MetricName:        element.MetricName,
		Metric:            element.Metric,
		ContractName:      contractName,
		TemplateString:    templateB64,
		BINstring:         binString,
		ABIstring:         abiString,
		Timestamp:         time.Now().String(),
		ContractAddress:   "",
		TransactionHash:   "",
		TransactionTime:   "",
		TransactionUUID:   "",
		TransactionSender: commons.GoDotEnvVariable("ETHEREUMPUBKEY"),
		User:              element.User,
		ErrorMessage:      "",
		Status:            "",
	}

	buildQueueObject := model.SendToQueue{
		EthereumMetricBind: ethMetricMetadataObj,
		Type:               "ETHMETRICBIND",
		User:               element.User,
		Status:             "QUEUE",
	}

	errWhenSendingToQueue := services.SendToQueue(buildQueueObject)
	if errWhenSendingToQueue != nil {
		logrus.Error("Error when sending to the metric metadata contract to queue : ", errWhenSendingToQueue)
		return errWhenSendingToQueue
	}

	return nil

}
