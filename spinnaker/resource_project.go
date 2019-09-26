package spinnaker

import (
	"strings"

	"github.com/marcoreni/terraform-provider-spinnaker/spinnaker/api"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mitchellh/mapstructure"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"applications": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"clusters": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account": {
							Type:     schema.TypeString,
							Required: true,
						},
						"applications": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"detail": {
							Type:     schema.TypeString,
							Default:  "*",
							Optional: true,
						},
						"stack": {
							Type:     schema.TypeString,
							Default:  "*",
							Optional: true,
						},
					},
				},
				Optional: true,
			},
			"pipelines": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"application": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
			},
		},
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Exists: resourceProjectExists,
	}
}

// Project represents the Gate API schema
type Project struct {
	ID     string         `json:"id,omitempty"`
	Name   string         `json:"name"`
	Email  string         `json:"email"`
	Config *ProjectConfig `json:"config"`
}

type ProjectConfig struct {
	Applications    []string                 `json:"applications"`
	Clusters        []*ProjectCluster        `json:"clusters"`
	PipelineConfigs []*ProjectPipelineConfig `json:"pipelineConfigs"`
}

type ProjectCluster struct {
	Account      string   `json:"account"`
	Applications []string `json:"applications"`
	Detail       string   `json:"detail"`
	Stack        string   `json:"stack"`
}

type ProjectPipelineConfig struct {
	ID          string `json:"pipelineConfigId"`
	Application string `json:"application"`
}

func resourceProjectCreate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client

	prj := projectFromResource(data)
	if err := api.CreateProject(client, prj.Name, prj); err != nil {
		return err
	}

	return resourceProjectRead(data, meta)
}

func resourceProjectRead(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client

	projectName := data.Get("name").(string)
	var prj Project
	if err := api.GetProject(client, projectName, &prj); err != nil {
		return err
	}

	return readProject(data, &prj)
}

func resourceProjectUpdate(data *schema.ResourceData, meta interface{}) error {
	// the Project update in spinnaker is an simple upsert
	clientConfig := meta.(gateConfig)
	client := clientConfig.client

	prj := projectFromResource(data)
	if err := api.CreateProject(client, prj.Name, prj); err != nil {
		return err
	}

	return resourceProjectRead(data, meta)
}

func resourceProjectDelete(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	projectID := data.Id()
	projectName := data.Get("name").(string)

	return api.DeleteProject(client, projectID, projectName)
}

func resourceProjectExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	projectName := data.Get("name").(string)

	var prj Project
	if err := api.GetProject(client, projectName, &prj); err != nil {
		errmsg := err.Error()
		if strings.Contains(errmsg, "not found") {
			return false, nil
		}
		return false, err
	}

	if prj.ID == "" {
		return false, nil
	}

	return true, nil
}

func projectFromResource(data *schema.ResourceData) *Project {
	prj := &Project{
		ID:    data.Id(),
		Name:  data.Get("name").(string),
		Email: data.Get("email").(string),
		Config: &ProjectConfig{
			Applications:    []string{},
			Clusters:        []*ProjectCluster{},
			PipelineConfigs: []*ProjectPipelineConfig{},
		},
	}

	for _, app := range data.Get("applications").([]interface{}) {
		prj.Config.Applications = append(prj.Config.Applications, app.(string))
	}

	for _, o := range data.Get("clusters").([]interface{}) {
		cl := &ProjectCluster{}

		if err := mapstructure.Decode(o, cl); err != nil {
			panic(err)
		}
		prj.Config.Clusters = append(prj.Config.Clusters, cl)
	}

	for _, o := range data.Get("pipelines").([]interface{}) {
		pc := &ProjectPipelineConfig{}

		if err := mapstructure.Decode(o, pc); err != nil {
			panic(err)
		}
		prj.Config.PipelineConfigs = append(prj.Config.PipelineConfigs, pc)
	}

	return prj
}

func readProject(data *schema.ResourceData, project *Project) error {
	data.SetId(project.ID)
	data.Set("name", project.Name)
	data.Set("email", project.Email)
	data.Set("applications", project.Config.Applications)

	var clusters []map[string]interface{}
	for _, cluster := range project.Config.Clusters {
		clusters = append(clusters, map[string]interface{}{
			"account":      cluster.Account,
			"applications": cluster.Applications,
			"detail":       cluster.Detail,
			"stack":        cluster.Stack,
		})
	}
	data.Set("clusters", clusters)

	var pipelines []map[string]interface{}
	for _, pipeline := range project.Config.PipelineConfigs {
		pipelines = append(clusters, map[string]interface{}{
			"id":          pipeline.ID,
			"application": pipeline.Application,
		})
	}
	data.Set("pipelines", pipelines)

	return nil
}
