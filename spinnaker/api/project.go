package api

import (
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	gate "github.com/spinnaker/spin/cmd/gateclient"
)

func GetProject(client *gate.GatewayClient, projectName string, dest interface{}) error {
	app, resp, err := client.ProjectControllerApi.GetUsingGET1(client.Context, projectName)
	if resp != nil {
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Project '%s' not found\n", projectName)
		} else if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Encountered an error getting project: %s status code: %d\n", err, resp.StatusCode)
		}
	}

	if err != nil {
		return err
	}

	if err := mapstructure.Decode(app, dest); err != nil {
		return err
	}

	return nil
}

func CreateProject(client *gate.GatewayClient, projectName string, project interface{}) error {
	task := map[string]interface{}{
		"job":         []interface{}{map[string]interface{}{"type": "upsertProject", "project": project}},
		"application": "spinnaker",
		"description": fmt.Sprintf("Create Project: %s", projectName),
		"project":     projectName,
	}

	return runTask(client, task)
}

func DeleteProject(client *gate.GatewayClient, projectID, projectName string) error {
	jobSpec := map[string]interface{}{
		"type": "deleteProject",
		"project": map[string]interface{}{
			"id": projectID,
		},
	}

	task := map[string]interface{}{
		"job":         []interface{}{jobSpec},
		"application": "spinnaker",
		"project":     projectName,
		"description": fmt.Sprintf("Delete Project: %s", projectName),
	}

	return runTask(client, task)
}
