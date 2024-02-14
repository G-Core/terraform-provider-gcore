package gcore

import (
	"context"
	"log"
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/clusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceK8sV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceK8sV2Read,
		Description: "Represent k8s cluster with one default pool.",
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"region_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				DiffSuppressFunc: suppressDiffRegionID,
			},
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": {
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cni": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"fixed_network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fixed_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"keypair": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"min_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"boot_volume_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available values are 'standard', 'ssd_hiiops', 'cold', 'ultra'.",
						},
						"boot_volume_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"auto_healing_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_public_ipv4": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator_task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceK8sV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start K8s reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, K8sPoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterName := d.Get("name").(string)
	cluster, err := clusters.Get(client, clusterName).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cluster.Name)

	d.Set("name", cluster.Name)
	d.Set("fixed_network", cluster.FixedNetwork)
	d.Set("fixed_subnet", cluster.FixedSubnet)
	d.Set("keypair", cluster.KeyPair)
	d.Set("version", cluster.Version)
	d.Set("node_count", cluster.NodeCount)
	d.Set("flavor_id", cluster.FlavorID)
	d.Set("status", cluster.Status)
	d.Set("is_public", cluster.IsPublic)
	d.Set("created_at", cluster.CreatedAt.Format(time.RFC850))
	d.Set("creator_task_id", cluster.CreatorTaskID)
	d.Set("task_id", cluster.TaskID)

	if cluster.CNI != nil {
		v := map[string]interface{}{
			"provider": cluster.CNI.Provider.String(),
		}
		d.Set("cni", []interface{}{v})
	}

	var ps []map[string]interface{}
	for _, pool := range cluster.Pools {
		ps = append(ps, map[string]interface{}{
			"name":                 pool.Name,
			"flavor_id":            pool.FlavorID,
			"min_node_count":       pool.MinNodeCount,
			"max_node_count":       pool.MaxNodeCount,
			"node_count":           pool.NodeCount,
			"boot_volume_type":     pool.BootVolumeType.String(),
			"boot_volume_size":     pool.BootVolumeSize,
			"auto_healing_enabled": pool.AutoHealingEnabled,
			"is_public_ipv4":       pool.IsPublicIPv4,
			"status":               pool.Status,
			"created_at":           pool.CreatedAt.Format(time.RFC850),
		})
	}
	if err := d.Set("pools", ps); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish K8s reading")
	return diags
}
