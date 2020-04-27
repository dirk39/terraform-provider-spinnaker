package api

import (
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	gate "github.com/spinnaker/spin/cmd/gateclient"
)

const (
	//V2ErrCodeNoSuchEntityException error code for v2 templates
	V2ErrCodeNoSuchEntityException = "NoSuchEntityException"
)

type gatewayClient interface {
	createPipelineTemplateV2(interface{}) (*http.Response, error)
}

// V2CreatePipelineTemplate create a pipeline template
func V2CreatePipelineTemplate(client gatewayClient, template interface{}) error {
	resp, err := client.createPipelineTemplateV2(template)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Encountered an error saving template, status code: %d", resp.StatusCode)
	}

	return nil
}

func V2GetPipelineTemplate(client *gate.GatewayClient, templateID string, dest interface{}) error {
	successPayload, resp, err := client.V2PipelineTemplatesControllerApi.GetUsingGET2(client.Context, templateID, nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("%s", ErrCodeNoSuchEntityException)
		}
		return fmt.Errorf("Encountered an error getting pipeline template %s, %s",
			templateID,
			err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error getting pipeline template %s, status code: %d",
			templateID,
			resp.StatusCode,
		)
	}

	if successPayload == nil {
		return fmt.Errorf(ErrCodeNoSuchEntityException)
	}

	if err := mapstructure.Decode(successPayload, dest); err != nil {
		return err
	}

	return nil
}

func V2DeletePipelineTemplate(client *gate.GatewayClient, templateID string) error {
	_, resp, err := client.V2PipelineTemplatesControllerApi.DeleteUsingDELETE1(client.Context, templateID, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Encountered an error deleting pipeline template %s, status code: %d",
			templateID,
			resp.StatusCode)
	}

	return nil
}

func V2UpdatePipelineTemplate(client *gate.GatewayClient, templateID string, template interface{}) error {
	resp, err := client.V2PipelineTemplatesControllerApi.UpdateUsingPOST1(client.Context, templateID, template, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Encountered an error updating pipeline template %s, status code: %d",
			templateID,
			resp.StatusCode)
	}

	return nil
}
