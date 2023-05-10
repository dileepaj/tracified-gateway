package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dileepaj/tracified-gateway/apiDemo/model/dtos/responsedtos"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	resp := responsedtos.HealthCheckResponse{
		Note:    "Tracified nft backend up and running",
		Time:    time.Now().Format("Mon Jan _2 15:04:05 2006"),
		Version: "0",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	return
}
