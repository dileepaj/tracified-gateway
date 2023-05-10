package controller

import (
	"net/http"
	"time"

	responseDtos "github.com/dileepaj/tracified-gateway/apiDemo/model/dtos/responseDtos"
	"github.com/dileepaj/tracified-gateway/utilities"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	resp := responseDtos.HealthCheckResponse{
		Note:    "Tracified nft backend up and running",
		Time:    time.Now().Format("Mon Jan _2 15:04:05 2006"),
		Version: "0",
	}
	utilities.SuccessResponse[responseDtos.HealthCheckResponse](w, resp)
	return
}
