package spinnaker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/marcoreni/terraform-provider-spinnaker/spinnaker/api"
)

func resourcePipelineTemplateV2() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"template": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressEquivalentPipelineTemplateDiffs,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Create: resourcePipelineTemplateCreateV2,
		Read:   resourcePipelineTemplateReadV2,
		Update: resourcePipelineTemplateUpdateV2,
		Delete: resourcePipelineTemplateDeleteV2,
		Exists: resourcePipelineTemplateExistsV2,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

type templateRead struct {
	ID string `json:"id"`
}

func resourcePipelineTemplateCreateV2(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	var templateName string
	template := data.Get("template").(string)

	d, err := yaml.YAMLToJSON([]byte(template))
	if err != nil {
		return err
	}

	var jsonContent map[string]interface{}
	if err = json.NewDecoder(bytes.NewReader(d)).Decode(&jsonContent); err != nil {
		return fmt.Errorf("Error decoding json: %s", err.Error())
	}

	if _, ok := jsonContent["schema"]; !ok {
		return fmt.Errorf("Pipeline save command currently only supports pipeline template configurations")
	}

	templateName = jsonContent["id"].(string)

	log.Println("[DEBUG] Making request to spinnaker")
	if err := api.V2CreatePipelineTemplate(client, jsonContent); err != nil {
		log.Printf("[DEBUG] Error response from spinnaker: %s", err.Error())
		return err
	}

	log.Printf("[DEBUG] Created template successfully")
	data.SetId(templateName)
	return resourcePipelineTemplateReadV2(data, meta)
}

func resourcePipelineTemplateReadV2(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	templateName := data.Id()

	t := make(map[string]interface{})
	if err := api.V2GetPipelineTemplate(client, templateName, &t); err != nil {
		if err.Error() == api.V2ErrCodeNoSuchEntityException {
			data.SetId("")
			return nil
		}
		return err
	}

	// Remove timestamp from response
	delete(t, "updateTs")
	delete(t, "lastModifiedBy")

	jsonContent, err := json.Marshal(t)
	if err != nil {
		return err
	}

	raw, err := yaml.JSONToYAML(jsonContent)
	if err != nil {
		return err
	}
	data.Set("name", t["id"].(string))
	data.Set("template", string(raw))
	data.Set("url", fmt.Sprintf("spinnaker://%s", t["id"].(string)))
	data.SetId(t["id"].(string))

	return nil
}

func resourcePipelineTemplateUpdateV2(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	var templateName string
	template := data.Get("template").(string)

	d, err := yaml.YAMLToJSON([]byte(template))
	if err != nil {
		return err
	}

	var jsonContent map[string]interface{}
	if err = json.NewDecoder(bytes.NewReader(d)).Decode(&jsonContent); err != nil {
		return fmt.Errorf("Error decoding json: %s", err.Error())
	}

	if _, ok := jsonContent["schema"]; !ok {
		return fmt.Errorf("Pipeline save command currently only supports pipeline template configurations")
	}

	templateName = jsonContent["id"].(string)

	if err := api.V2UpdatePipelineTemplate(client, templateName, jsonContent); err != nil {
		return err
	}

	data.SetId(templateName)
	return resourcePipelineTemplateReadV2(data, meta)
}

func resourcePipelineTemplateDeleteV2(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	templateName := data.Id()

	if err := api.V2DeletePipelineTemplate(client, templateName); err != nil {
		return err
	}

	data.SetId("")
	return nil
}

func resourcePipelineTemplateExistsV2(data *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	templateName := data.Id()

	var t templateRead
	if err := api.V2GetPipelineTemplate(client, templateName, &t); err != nil {
		if err.Error() == api.V2ErrCodeNoSuchEntityException {
			return false, nil
		}
		return false, err
	}

	if t.ID == templateName {
		return true, nil
	}

	return false, nil
}

func suppressEquivalentPipelineTemplateDiffs(k, old, new string, d *schema.ResourceData) bool {
	equivalent, err := areEqualJSON(old, new)
	if err != nil {
		return false
	}

	return equivalent
}

func areEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	log.Printf("[DEBUG] s1: %s", s1)
	err = yaml.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	log.Printf("[DEBUG] s2: %s", s2)
	err = yaml.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}
