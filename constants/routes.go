package constants

import "github.com/dileepaj/tracified-gateway/commons"

//TracifiedBackend The URL TO Traceability Staging
// var TracifiedBackend="https://staging.api.tracified.com"
//TracifiedBackend The URL TO Traceability Production
var TracifiedBackend = commons.GoDotEnvVariable("TRACIFIEDBACKEND")

//  var TracifiedBackend="http://localhost:3001"

//RawTDP The Route to Get Raw TDP Data
var RawTDP = "/api/v2/dataPackets/raw?id="

//BackendToken The Token used to Retrieve Raw TDP from Traceability DataStore
//Will Expire in 03/19/2022 @ 12:31pm (UTC).
var BackendToken = commons.GoDotEnvVariable("TRACIFIEDBACKENDTOKEN")
