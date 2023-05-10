package commons

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
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

func ConvertBase64StringToHash256(s string) string {
	decodedKey, _ := base64.StdEncoding.DecodeString(s)
	h := sha256.New()
	h.Write(decodedKey)
	return hex.EncodeToString(h.Sum(nil))
}

func DecodeJSONRequestBody(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while decoding the request body",
		}
		json.NewEncoder(w).Encode(result)
		fmt.Println(err)
		return false
	}
	return true
}

func CheckTaskStatus(tasks []model.Task) []model.Task {
	var failedTasks []model.Task
	for _, task := range tasks {
		if !task.Status {
			failedTasks = append(failedTasks, task)
		}
	}
	return failedTasks
}

func TasksToString(tasks []model.Task) string {
	var tasksString []string
	for _, task := range tasks {
		tasksString = append(tasksString, fmt.Sprintf("%s %v", task.Identifier, task.Status))
	}
	return fmt.Sprintf("[%s]", strings.Join(tasksString, ", "))
}
