package api

import (
	"errors"
	"net/http"
	"testing"
)

var createPipelineTemplateV2Mock func(interface{}) (*http.Response, error)

type mockedAPIClient struct {
}

func (client mockedAPIClient) createPipelineTemplateV2(values interface{}) (*http.Response, error) {
	return createPipelineTemplateV2Mock(values)
}

func TestV2CreatePipelineTemplateReturnErrorFromApi(t *testing.T) {
	createPipelineTemplateV2Mock = func(interface{}) (*http.Response, error) {
		return nil, errors.New("This is an error")
	}
	var mock mockedAPIClient

	err := V2CreatePipelineTemplate(mock, "useless")

	if nil == err {
		t.Fatal("TestV2CreatePipelineTemplateReturnErrorFromApi failed: expecting error, got nil")
	}

	t.Log("TestV2CreatePipelineTemplateReturnErrorFromApi PASS")
}
