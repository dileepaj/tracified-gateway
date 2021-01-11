package businessFacades

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthCheckResponse struct {
	Note    string  `json:"note"`
	Time    string  `json:"time"`
	Version string `json:"version"`
}


func HealthCheck(w http.ResponseWriter, r *http.Request){
	resp := HealthCheckResponse{
		Note: "Gateway up and running",
		Time:    time.Now().Format("Mon Jan _2 15:04:05 2006"),
		Version: "Not Found",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
	}
}
