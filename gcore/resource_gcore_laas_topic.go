package gcore

import (
	"context"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/laas/v1/laas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const laasPoint = "laas"

func resourceLaaSTopic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLaaSTopicCreate,
		ReadContext:   resourceLaaSTopicRead,
		DeleteContext: resourceLaaSTopicDelete,
		Description:   "Represent LaaS topic",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, topicName, err := ImportStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(topicName)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"region_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				DiffSuppressFunc: suppressDiffRegionID,
			},
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLaaSTopicCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LaaS topic creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, laasPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := laas.CreateTopicOpts{Name: d.Get("name").(string)}
	topic, err := laas.CreateTopic(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(topic.Name)

	resourceLaaSTopicRead(ctx, d, m)

	log.Printf("[DEBUG] Finish LaaS topic creating (%s)", topic.Name)
	return diags
}

func resourceLaaSTopicRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LaaS topic reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	topicName := d.Id()
	log.Printf("[DEBUG] Topic id = %s", topicName)

	client, err := CreateClient(provider, d, laasPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	topics, err := laas.ListTopicAll(client)
	if err != nil {
		return diag.Errorf("cannot get topic's list. Error: %s", err.Error())
	}
	var topic laas.Topic
	for _, t := range topics {
		if t.Name == topicName {
			topic = t
			break
		}
	}
	if topic.Name == "" {
		return diag.Errorf("cant find topic with name %s", topicName)
	}
	d.Set("name", topic.Name)

	log.Println("[DEBUG] Finish LaaS topic reading")
	return diags
}

func resourceLaaSTopicDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LaaS topic deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	topicName := d.Id()
	log.Printf("[DEBUG] Topic id = %s", topicName)

	client, err := CreateClient(provider, d, laasPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := laas.DeleteTopic(client, topicName).ExtractErr(); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[DEBUG] Finish of LaaS topic deleting")
	return diags
}
