package commons

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

// work as ternary operator for string
func ValidateStrings(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func JSONErrorReturn(w http.ResponseWriter, r *http.Request, err string, errorCode int, message string) {
	logrus.Error(message, " ", err)
	w.WriteHeader(errorCode)
	response := model.Error{Code: errorCode, Message: message + " " + err}
	json.NewEncoder(w).Encode(response)
	return
}

func ChunkSlice(slice []txnbuild.ManageData, chunkSize int) [][]txnbuild.ManageData {
	var chunks [][]txnbuild.ManageData
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}
