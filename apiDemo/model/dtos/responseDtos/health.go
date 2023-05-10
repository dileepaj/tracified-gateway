package responsedtos

type HealthCheckResponse struct {
	Note    string `json:"note"`
	Time    string `json:"time"`
	Version string `json:"version"`
}