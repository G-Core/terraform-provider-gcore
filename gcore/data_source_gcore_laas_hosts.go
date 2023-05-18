package gcore

import (
	"context"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/laas/v1/laas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLaaSHosts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLaaSHostsRead,
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
			"opensearch": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"kafka": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceLaaSHostsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LaaS hosts reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, laasPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	kafkaHosts, err := laas.ListKafkaHosts(client).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	openSearchHosts, err := laas.ListOpenSearchHosts(client).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(getUniqueID(d))
	if err := d.Set("kafka", kafkaHosts); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("opensearch", openSearchHosts); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish LaaS hosts reading")
	return diags
}
