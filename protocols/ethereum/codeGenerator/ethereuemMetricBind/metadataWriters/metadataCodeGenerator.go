package metadatawriters

import (
	"encoding/base64"
	"strconv"

	"github.com/dileepaj/tracified-gateway/model"
)

func MetricMetadataWriter(element model.MetricMetadataReq, metricMapID string) (string, error) {

	//contract header
	license := `// SPDX-License-Identifier: MIT`
	version := "\n\n" + `pragma solidity ^0.8.7;`

	//contract declaration
	name := "Metric_" + metricMapID + "_Metadata"
	contractName := "\n\n" + `contract ` + name + `{`

	//metadata structure
	metadataStructComment := "\n\n\t" + `// Metadata structure`
	structHeader := "\n\t" + `struct Metadata {`
	metricIdInStruct := "\n\t\t" + `string metricID;`
	metricNameInStruct := "\n\t\t" + `string metricName; // converted value to base64`
	tenetIdInStruct := "\n\t\t" + `string tenantID;`
	noOfFormulasInStruct := "\n\t\t" + `uint noOfFormulas;`
	trustNetPKInStruct := "\n\t\t" + `string trustNetPK;`
	endOfStruct := "\n\t" + `}`
	metadataStruct := metadataStructComment + structHeader + metricIdInStruct + metricNameInStruct + tenetIdInStruct + noOfFormulasInStruct + trustNetPKInStruct + endOfStruct

	//convert metric name to base64
	metricNameB64 := base64.StdEncoding.EncodeToString([]byte(element.MetricName))

	//metadata deceleration call
	metadataDecelerationComment := "\n\n\t" + `// Metadata declaration`
	metadataDecelerationCall := "\n\t" + `Metadata metadata = Metadata("` + element.MetricId + `","` + metricNameB64 + `","` + element.TenetId + `",` + strconv.Itoa(element.NoOfFormulas) + `,"` + element.TrustNetPK + `");`
	metaDataDeceleration := metadataDecelerationComment + metadataDecelerationCall

	//metadata getter method
	metadataGetterComment := "\n\n\t" + `// Metadata getter`
	metadataGetterFunctionSignature := "\n\t" + `function metadataGetter() public view returns(Metadata memory) {`
	metadataMemoryCopy := "\n\t\t" + `Metadata metadata;`
	endOfGetter := "\n\t" + `}`
	metadataGetter := metadataGetterComment + metadataGetterFunctionSignature + metadataMemoryCopy + metadataGetterReturn + endOfGetter

	metadataCode := license + version + contractName + metadataStruct + metaDataDeceleration + metadataGetter + "\n" + `}`

	return metadataCode, nil
}
