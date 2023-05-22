package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model/workflowmodel"
)

type WorkflowService struct {
	TenantID string
}

func (wf *WorkflowService) GetWorkflowByTenantId() (workflowmodel.WorkflowData, error) {
	var workflowData workflowmodel.WorkflowData
	url := commons.GoDotEnvVariable("TRACIFIED_ADMIN") + "/public/latestWorkflow/" + wf.TenantID

	// Send GET request
	response, err := http.Get(url)
	if err != nil {
		return workflowData, fmt.Errorf("failed to send GET request: %s", err)
	}
	defer response.Body.Close()

	// Check API response status code
	if response.StatusCode != http.StatusOK {
		return workflowData, fmt.Errorf("API request failed with status code: %d", response.StatusCode)
	}
	// Read response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return workflowData, fmt.Errorf("failed to read response body: %s", err)
	}

	// Unmarshal the API response into the struct
	err = json.Unmarshal(body, &workflowData)
	if err != nil {
		return workflowData, fmt.Errorf("failed to unmarshal API response: %s", err)
	}

	return workflowData, nil
}
