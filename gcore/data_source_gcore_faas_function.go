package gcore

import (
	"context"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/faas/v1/faas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFaaSFunction() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFaaSFunctionRead,
		Description: "Represent FaaS function",
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
			"namespace": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Namespace of the function",
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
			"runtime": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"code_text": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_instances": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Autoscaling max number of instances",
				Computed:    true,
			},
			"min_instances": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Autoscaling min number of instances",
				Computed:    true,
			},
			"main_method": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Main startup method name",
				Computed:    true,
			},
			"flavor": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"build_message": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"build_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"deploy_status": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceFaaSFunctionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FaaS function reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	fName := d.Get("name").(string)
	nsName := d.Get("namespace").(string)
	log.Printf("[DEBUG] function = %s in %s", fName, nsName)

	client, err := CreateClient(provider, d, faasPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	function, err := faas.GetFunction(client, nsName, fName).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(funcID(function.Name, nsName))
	if err := faaSSetState(d, function); err != nil {
		diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish FaaS function reading")
	return diags
}
