package metadatawriters

import (
	"encoding/base64"
	"os"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/expertformulas/ethereum/deploy"
	"github.com/dileepaj/tracified-gateway/model"
	ethereumsocialimpact "github.com/dileepaj/tracified-gateway/services/ethereumServices/ethereumSocialImpact"
	"github.com/sirupsen/logrus"
)

// ! Relevant map ID will be checked on the route handler to see whether the metric contract is already deployed or not
// ! All the other failed database calls will be handled in the handler function
func MetricMetadataContractDeployer(element model.MetricMetadataReq, metricMapID string, ethMetricMetadataObj model.EthereumMetricBind) error {
	object := dao.Connection{}
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

	ethMetricMetadataObj.TemplateString = templateB64
	ethMetricMetadataObj.BINstring = binString
	ethMetricMetadataObj.ABIstring = abiString

	errWhenUpdatingMetricDetails := object.UpdateEthereumMetricStatus(ethMetricMetadataObj.MetricID, ethMetricMetadataObj.TransactionUUID, ethMetricMetadataObj)
	if errWhenUpdatingMetricDetails != nil {
		logrus.Error("Error when updating the metric metadata contract details : ", errWhenUpdatingMetricDetails)
		return errWhenUpdatingMetricDetails
	}

	errWhenDeploying := ethereumsocialimpact.DeployMetricContract(ethMetricMetadataObj)
	if errWhenDeploying != nil {
		logrus.Error("Error when sending to the metric metadata contract to deployer : ", errWhenDeploying)
		return errWhenDeploying
	}

	return nil

}
