package gcore

import (
	"context"
	"log"
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/faas/v1/faas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFaaSNamespace() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFaaSNamespaceRead,
		Description: "Represent FaaS namespace",
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
			"envs": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceFaaSNamespaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FaaS namespace reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	nsName := d.Get("name").(string)
	log.Printf("[DEBUG] namespace = %s", nsName)

	client, err := CreateClient(provider, d, faasPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	ns, err := faas.GetNamespace(client, nsName).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ns.Name)
	d.Set("name", ns.Name)
	d.Set("description", ns.Description)
	d.Set("status", ns.Status)
	d.Set("created_at", ns.CreatedAt.Format(time.RFC3339))

	if err := d.Set("envs", ns.Envs); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish FaaS namespace reading")
	return diags
}
