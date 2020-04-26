package api

import (
	"net/http"

	gate "github.com/spinnaker/spin/cmd/gateclient"
)

//Client is the type able to comunicate with Spinnaker
type Client struct {
	*gate.GatewayClient
}

func (c Client) createPipelineTemplateV2(template interface{}) (*http.Response, error) {
	return c.GatewayClient.V2PipelineTemplatesControllerApi.CreateUsingPOST1(c.GatewayClient.Context, template, nil)
}

// InitAPIClient return initialized spinnaker api client
func InitAPIClient(gatewayClient *gate.GatewayClient) Client {
	return Client{
		gatewayClient,
	}

}
