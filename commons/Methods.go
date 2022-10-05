package commons

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
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
