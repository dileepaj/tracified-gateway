package ethereuemmetricbind

import (
	"net/http"

	"github.com/dileepaj/tracified-gateway/model"
	activityWriters "github.com/dileepaj/tracified-gateway/protocols/ethereum/codeGenerator/ethereuemMetricBind/activityContractWriters"
	metadataWriters "github.com/dileepaj/tracified-gateway/protocols/ethereum/codeGenerator/ethereuemMetricBind/metadataWriters"
	"github.com/dileepaj/tracified-gateway/vendor/github.com/sirupsen/logrus"
)

func SmartContractHandlerForMetric(w http.ResponseWriter, r *http.Request, metricBindJson model.MetricDataBindingRequest) {

	// create the metadata object
	metaDataObj := model.MetricMetadataReq{
		MetricId:     metricBindJson.Metric.ID,
		MetricName:   metricBindJson.Metric.Name,
		TenetId:      metricBindJson.Metric.TenantId,
		NoOfFormulas: len(metricBindJson.Metric.MetricActivities),
		TrustNetPK:   metricBindJson.User.Publickey,
		Metric:       metricBindJson.Metric,
		User:         metricBindJson.User,
	}

	// get the activity list of the metric
	activities := metricBindJson.Metric.MetricActivities

	// get the status of the metric metadata contract
	status, errWhenGettingMetadataContractStatus := GetMetricSmartContractStatus(metricBindJson.Metric.ID, "METADATA")
	if errWhenGettingMetadataContractStatus != nil {
		// TODO: handle the error
		logrus.Info(errWhenGettingMetadataContractStatus)
	}

	if status == "" || status == "FAILED" {
		// deploy the smart contract for meta data
		errWhenDeployingMetaDataSmartContract := metadataWriters.MetricMetadataContractDeployer(metaDataObj, metricBindJson.Metric.ID)
		if errWhenDeployingMetaDataSmartContract != nil {
			// TODO: handle the error
			logrus.Info(errWhenDeployingMetaDataSmartContract)
		}
	}

	// get the status of the metric metadata contract
	status, errWhenGettingMetadataContractStatus = GetMetricSmartContractStatus(metricBindJson.Metric.ID, "METADATA")
	if errWhenGettingMetadataContractStatus != nil {
		// TODO: handle the error
		logrus.Info(errWhenGettingMetadataContractStatus)
	}

	// check if the status is SUCCESS or not, if SUCCESS then proceed to create the smart contract for the metric activities
	// TODO: handle status QUEUE
	if status == "SUCCESS" {
		if len(activities) > 0 {
			for _, activity := range activities {
				// TODO: validation to check whether the contract is already deployed or not(metric collection) using metric id and formula id
				// if it is deployed then skip the deployment 
				// check the index of the loop, if it is not 0 then check the index-1 activity contract status from the metric collection
				// if the status is SUCCESS then deploy the activity contract

				errWhenDeployingActivityContract := activityWriters.ActivityContractDeployer(metricBindJson.Metric.ID, activity, metricBindJson.Metric.Name, metricBindJson.Metric, metricBindJson.User)
				if errWhenDeployingActivityContract != nil {
					// TODO: handle the error
					logrus.Info(errWhenDeployingActivityContract)
				}
			}
		}
	}

}
