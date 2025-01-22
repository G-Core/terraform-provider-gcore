package gcore

import (
	"context"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/inference/v3/flavors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceInferenceFlavor() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInferenceFlavorRead,
		Description: "Represent Inference flavor.",
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
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"cpu": &schema.Schema{
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"memory": &schema.Schema{
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"gpu": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gpu_memory": &schema.Schema{
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"gpu_model": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_gpu_shared": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"gpu_compute_capability": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceInferenceFlavorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start flavor reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, inferenceDeploymentPoint, versionPointV3)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)

	flavor, err := flavors.GetFlavor(client, name).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(flavor.Name)
	d.Set("cpu", flavor.Cpu)
	d.Set("memory", flavor.Memory)
	d.Set("gpu", flavor.Gpu)
	d.Set("gpu_memory", flavor.GpuMemory)
	d.Set("gpu_model", flavor.GpuModel)
	d.Set("is_gpu_shared", flavor.IsGpuShared)
	d.Set("gpu_compute_capability", flavor.GpuComputeCapability)

	log.Println("[DEBUG] Finish flavor reading")
	return diags
}
