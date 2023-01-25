package metadatawriters

import (
	"os"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

//! Relevant map ID will be checked on the route handler to see whether the metric contract is already deployed or not
func MetricMetadataContractDeployer(element model.MetricMetadataReq, metricMapID string) error {

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

	//TODO generate ABI

	//TODO generate the BIN

	//TODO add the request to the queue

	return nil

}
