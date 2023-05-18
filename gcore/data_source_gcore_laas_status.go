package gcore

import (
	"context"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/laas/v1/laas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLaaSStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLaaSStatusRead,
		Description: "Represent LaaS hosts",
		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"region_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				DiffSuppressFunc: suppressDiffRegionID,
			},
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
			},
			"namespace": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_initialized": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceLaaSStatusRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LaaS status reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, laasPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	status, err := laas.GetStatus(client).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(getUniqueID(d))
	d.Set("namespace", status.Namespace)
	d.Set("is_initialized", status.IsInitialized)

	log.Println("[DEBUG] Finish LaaS status reading")
	return diags
}
