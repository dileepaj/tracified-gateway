package httpresponse

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/apiDemo/model/dtos/responsedtos"
)

// Move to model interface location or add under Response dtos
type ResultResponse struct {
	StatusCode int `json:"Status"`
	Response   any `json:"Response"`
}

type ResultType interface {
	responsedtos.HealthCheckResponse | string
}

func SuccessResponse[T ResultType](w http.ResponseWriter, result T) {
	w.WriteHeader(http.StatusOK)
	response := ResultResponse{
		StatusCode: http.StatusOK,
		Response:   result,
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		// logs.ErrorLogger.Println(err)
	}
}
