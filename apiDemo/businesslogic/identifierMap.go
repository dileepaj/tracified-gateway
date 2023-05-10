package businesslogic

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/apiDemo/dao"
)

func insertIdentifierMap(identifier, realIdentifier string) {
	object := dao.Connection{}
	var identifierModel apiModel.IdentifierModel
	rawDecodedText, err := base64.StdEncoding.DecodeString(identifier)
	if err != nil {
		fmt.Println("Decode String failed" + err.Error())
	}

	var jsonID apiModel.Identifier
	json.Unmarshal([]byte(rawDecodedText), &jsonID)
	identifierModel.MapValue = realIdentifier
	identifierModel.Identifier = jsonID.Id
	identifierModel.Type = jsonID.Type
	err3, code := object.InsertIdentifier(identifierModel)
	if err3 != nil {
		fmt.Println("identifier map failed"+err3.Error(), " ErrorCode:", code)
	}
}
