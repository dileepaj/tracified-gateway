package businesslogic

import (
	"encoding/base64"
	"encoding/json"
	"strconv"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/apiDemo/dao"
	"github.com/dileepaj/tracified-gateway/utilities"
)

func insertIdentifierMap(identifier, realIdentifier string) {
	object := dao.Connection{}
	logger := utilities.NewCustomLogger()
	var identifierModel apiModel.IdentifierModel
	rawDecodedText, err := base64.StdEncoding.DecodeString(identifier)
	if err != nil {
		logger.LogWriter("Decode String failed" + err.Error(), 3)
	}

	var jsonID apiModel.Identifier
	json.Unmarshal([]byte(rawDecodedText), &jsonID)
	identifierModel.MapValue = realIdentifier
	identifierModel.Identifier = jsonID.Id
	identifierModel.Type = jsonID.Type
	err3, code := object.InsertIdentifier(identifierModel)
	if err3 != nil {
		logger.LogWriter("Identifier map failed"+err3.Error() + " ErrorCode:"+ strconv.Itoa(code), 3)
	}
}
