package api

import (
	"errors"
	"net/http"
	"testing"
)

var createPipelineTemplateV2Mock func(interface{}) (*http.Response, error)
var getPipelineTemplateV2Mock func(string) (map[string]interface{}, *http.Response, error)

type mockedAPIClient struct {
}

func (client mockedAPIClient) createPipelineTemplateV2(template interface{}) (*http.Response, error) {
	return createPipelineTemplateV2Mock(template)
}

func (client mockedAPIClient) getPipelineTemplateV2(templateID string) (map[string]interface{}, *http.Response, error) {
	return getPipelineTemplateV2Mock(templateID)
}

// TestV2CreatePipelineTemplateReturnErrorFromApi test that
// V2CreatePipelineTemplate return error from api
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

// TestV2CreatePipelineTemplateReturnErrorByStatusCode test that
// V2CreatePipelineTemplate return error from api if status code is different
// from accepted once
func TestV2CreatePipelineTemplateReturnErrorByStatusCode(t *testing.T) {
	createPipelineTemplateV2Mock = func(interface{}) (*http.Response, error) {
		response := http.Response{
			StatusCode: 400,
		}

		return &response, nil
	}
	var mock mockedAPIClient

	err := V2CreatePipelineTemplate(mock, "useless")

	if nil == err {
		t.Fatal("TestV2CreatePipelineTemplateReturnErrorByStatusCode failed: expecting error, got nil")
	}

	if err.Error() != "Encountered an error saving template, status code: 400" {
		t.Fatalf(`TestV2CreatePipelineTemplateReturnErrorByStatusCode failed: expecting "Encountered an error saving template, status code: 400", got %v`, err.Error())
	}

	t.Log("TestV2CreatePipelineTemplateReturnErrorByStatusCode PASS")
}

func TestV2GetPipelineTemplateReturnNoSuchEntityError(t *testing.T) {
	getPipelineTemplateV2Mock = func(string) (map[string]interface{}, *http.Response, error) {
		return nil, nil, errors.New("Undefined error")
	}

	var mock mockedAPIClient

	err := V2GetPipelineTemplate(mock, "john-doe", nil)

}
