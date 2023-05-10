package httpresponse

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/apiDemo/model/dtos/responseDtos"
)

// Move to model interface location or add under Response dtos
type ResultResponse struct {
	Status   int `json:"Status"`
	Response any `json:"Response"`
}

type ResultType interface {
	apiModel.SubmitXDRSuccess | responseDtos.HealthCheckResponse
}
type Commonresponse struct {
}

func SuccessStatus[T ResultType](w http.ResponseWriter, result T) {
	w.WriteHeader(http.StatusOK)
	response := ResultResponse{
		Status:   http.StatusOK,
		Response: result,
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		// logs.ErrorLogger.Println(err)
	}
}
