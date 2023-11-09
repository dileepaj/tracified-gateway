package managedatatostellarprotocol

import (
	"encoding/base64"
	"strings"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/sirupsen/logrus"
)

func MemoConversion(base64DataMemo string) {
	byteArrayMemo, errMemo := base64.StdEncoding.DecodeString(base64DataMemo)

	if errMemo != nil {
		logrus.Error("error:", errMemo)
	}
	if len(byteArrayMemo) != 28 {
		logrus.Error("Error: The length of the byte array is not 28")
	}
	commons.ByteArryToHexString(byteArrayMemo[0:10], "1. Manifest of the memo:")
	commons.ByteArryToInt64(byteArrayMemo[10:18], "2. Formula ID:")
	commons.ByteArryToInt32(byteArrayMemo[18:22], "3. No of variables:")
	commons.ByteArryToHexString(byteArrayMemo[22:28], "4. For future use:")
}

func FormulaIdentityConversion(keyFormulaIdentity string, base64DataFormulaIdentity string) {
	byteArrayFormulaIdentity, errFormulaIdentity := base64.StdEncoding.DecodeString(base64DataFormulaIdentity)
	if errFormulaIdentity != nil {
		logrus.Error("error:", errFormulaIdentity)
	}
	if len(byteArrayFormulaIdentity) != 64 {
		logrus.Error("Error: The length of the value field is not 64")
	}

	logrus.Println("Fields in the value field of the Formula Identity Manage Data:")
	commons.ByteArryToString(byteArrayFormulaIdentity[0:8], "1. AuthorIdentity:")
	commons.ByteArryToInt64(byteArrayFormulaIdentity[8:28], "2. Formula Name:")
	commons.ByteArryToHexString(byteArrayFormulaIdentity[28:64], "3. FutureUse:")
}

func AuthorIdentityConversion(keyAuthorIdentity string, base64DataAuthorIdentity string) {
	byteArrayFormulaIdentity, errFormulaIdentity := base64.StdEncoding.DecodeString(base64DataAuthorIdentity)
	if errFormulaIdentity != nil {
		logrus.Error("error:", errFormulaIdentity)
	}
	if len(byteArrayFormulaIdentity) != 64 {
		logrus.Error("Error: The length of the value field is not 64")
	}
	commons.ByteArryToHexString(byteArrayFormulaIdentity[0:64], "Author Primary key 2nd 64 bytes:")
}

func VariableConversion(keyVariable string, base64DataVariable string) {

	logrus.Info(" - Variable Name Field")
	logrus.Info("name(text): ", keyVariable)
	logrus.Info("Variable name actual length: ", len(keyVariable))
	if len(keyVariable) != 64 {
		logrus.Error("Error: The length of the key field is not 64")
	}
	logrus.Info("\t1. Description: ", keyVariable[0:40])
	logrus.Info("\t\t Actual description: ", strings.Split(keyVariable[0:40], "/")[0])
	logrus.Info("\t2. For future use: ", keyVariable[40:64])

	logrus.Info()
	logrus.Info(" - Variable Value Field")
	
	byteArrayVariable, errVariable := base64.StdEncoding.DecodeString(base64DataVariable)
	if errVariable != nil {
		logrus.Error("error:", errVariable)
	}
	logrus.Println("Fields in the value field of the Variable Manage Data:")
	commons.ByteArryToHexString(byteArrayVariable[0:1], "1. Value Type:")
	commons.ByteArryToInt64(byteArrayVariable[1:9], "2. Value ID:")	
	commons.ByteArryToString(byteArrayVariable[9:29], "3. Variable Name:")
	logrus.Println("\t\tActual variable name: ", strings.Split(string(byteArrayVariable[9:29]), "/" )[0])
	commons.ByteArryToHexString(byteArrayVariable[29:30], "4. Data Type:")	
	commons.ByteArryToInt16(byteArrayVariable[30:32], "5. Unit:")	
	commons.ByteArryToHexString(byteArrayVariable[32:33], "6. Precision:")	
	commons.ByteArryToHexString(byteArrayVariable[33:64], "7. For future use:")
}