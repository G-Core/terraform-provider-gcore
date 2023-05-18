package gcore

import (
	"context"
	"log"
	"strconv"

	"github.com/G-Core/gcorelabscloud-go/gcore/ddos/v1/ddos"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ddosTemplatesPoint = "ddos/profile-templates"
)

func dataSourceDDoSProfileTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDDoSProfileTemplateRead,
		Description: "Represents list of available DDoS protection profile templates",
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
			"region_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
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
			"template_id": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Template id",
				ExactlyOneOf: []string{
					"template_id",
					"name",
				},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Template name",
				ExactlyOneOf: []string{
					"template_id",
					"name",
				},
			},
			"fields": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Additional fields",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Field name",
						},
						"field_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Field type",
						},
						"required": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Field description",
						},
						"default": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"validation_schema": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Json schema to validate field_values",
						},
					},
				},
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Template description",
				Computed:    true,
			},
		},
	}
}

func dataSourceDDoSProfileTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Starts DDoS protection profile template reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, ddosTemplatesPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("template_id").(int)
	name := d.Get("name").(string)
	templates, err := ddos.ListAllProfileTemplates(client)
	if err != nil {
		return diag.FromErr(err)
	}

	var found bool
	var template ddos.ProfileTemplate
	for _, t := range templates {
		if t.ID == id {
			template = t
			found = true
			break
		}

		if t.Name == name {
			template = t
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("DDoS protection profile template not found not by ID %d nor by name %s", id, name)
	}

	d.SetId(strconv.Itoa(template.ID))
	d.Set("name", template.Name)
	d.Set("description", template.Description)
	fields := make([]map[string]interface{}, len(template.Fields))
	for i, f := range template.Fields {
		field := map[string]interface{}{
			"field_type":        f.FieldType,
			"required":          f.Required,
			"id":                f.ID,
			"default":           f.Default,
			"description":       f.Description,
			"name":              f.Name,
			"validation_schema": string(f.ValidationSchema),
		}

		fields[i] = field
	}
	if err := d.Set("fields", fields); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish DDoS protection profile template reading")
	return diags
}
