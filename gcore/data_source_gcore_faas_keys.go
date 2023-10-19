package gcore

import (
	"context"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/faas/v1/faas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	faasKeysPoint = "faas/keys"
)

func dataSourceFaaSKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFaaSKeyRead,
		Description: "Represent FaaS API keys",
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
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"functions": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceFaaSKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FaaS API key reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	prodiver := config.Provider
	keyName := d.Get("name").(string)
	log.Printf("[DEBUG] API key = %s", keyName)

	client, err := CreateClient(prodiver, d, faasKeysPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	key, err := faas.GetKey(client, keyName).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(key.Name)
	d.Set("name", key.Name)
	d.Set("description", key.Description)
	d.Set("status", key.Status)
	d.Set("created_at", key.CreatedAt)
	d.Set("expire", key.Expire)
	fs := make([]map[string]any, len(key.Functions))
	for idx, f := range key.Functions {
		fs[idx] = map[string]any{
			"name":      f.Name,
			"namespace": f.Namespace,
		}
	}

	if err := d.Set("functions", fs); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish FaaS API key reading")
	return diags
}
