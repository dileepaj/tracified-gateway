package commons

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
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

func ConvertBase64StringToHash256(s string) string {
	decodedKey, _ := base64.StdEncoding.DecodeString(s)
	h := sha256.New()
	h.Write(decodedKey)
	return hex.EncodeToString(h.Sum(nil))
}

func MakeGetRequest(url string) ([]byte, error) {
	// Make GET request
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// ContainsString checks if a given string is present in the array.
func ContainsString(array []string, target string) bool {
	for _, value := range array {
		if value == target {
			return true
		}
	}
	return false
}