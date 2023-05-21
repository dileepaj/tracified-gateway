package workflowmodel

type WorkflowData struct {
	TenantID string `json:"tenantId"`
	Revision int    `json:"revision"`
	Workflow []Workflow
}

type Workflow struct {
	Stages         []Stage `json:"stages"`
	Name           string  `json:"name"`
	EcommerceStage string  `json:"ecommerceStage"`
}

type Stage struct {
	Name    string `json:"name"`
	StageID string `json:"stageId"`
}
