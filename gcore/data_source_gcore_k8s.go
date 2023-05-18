package gcore

import (
	"context"
	"log"
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v1/clusters"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v1/pools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceK8s() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceK8sRead,
		Description: "Represent k8s cluster with one default pool.",
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
			"cluster_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"fixed_network": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"fixed_subnet": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_healing_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"external_dns_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"master_lb_floating_ip_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"keypair": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"pool": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"min_node_count": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_node_count": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"node_count": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"docker_volume_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"docker_volume_size": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"uuid": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"stack_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"node_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_reason": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_addresses": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"node_addresses": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"container_version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"discovery_url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_status_reason": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"faults": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"master_flavor_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_template_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": &schema.Schema{
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

func dataSourceK8sRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start K8s reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, K8sPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterID := d.Get("cluster_id").(string)
	cluster, err := clusters.Get(client, clusterID).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cluster.UUID)

	d.Set("name", cluster.Name)
	d.Set("fixed_network", cluster.FixedNetwork)
	d.Set("fixed_subnet", cluster.FixedSubnet)
	d.Set("master_lb_floating_ip_enabled", cluster.FloatingIPEnabled)
	d.Set("keypair", cluster.KeyPair)
	d.Set("node_count", cluster.NodeCount)
	d.Set("status", cluster.Status)
	d.Set("status_reason", cluster.StatusReason)
	d.Set("external_dns_enabled", cluster.ExternalDNSEnabled)

	masterAddresses := make([]string, len(cluster.MasterAddresses))
	for i, addr := range cluster.MasterAddresses {
		masterAddresses[i] = addr.String()
	}
	if err := d.Set("master_addresses", masterAddresses); err != nil {
		return diag.FromErr(err)
	}

	nodeAddresses := make([]string, len(cluster.NodeAddresses))
	for i, addr := range cluster.NodeAddresses {
		nodeAddresses[i] = addr.String()
	}
	if err := d.Set("node_addresses", nodeAddresses); err != nil {
		return diag.FromErr(err)
	}

	d.Set("container_version", cluster.ContainerVersion)
	d.Set("api_address", cluster.APIAddress.String())
	d.Set("user_id", cluster.UserID)
	d.Set("discovery_url", cluster.DiscoveryURL.String())

	d.Set("health_status", cluster.HealthStatus)
	if err := d.Set("health_status_reason", cluster.HealthStatusReason); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("faults", cluster.Faults); err != nil {
		return diag.FromErr(err)
	}

	d.Set("master_flavor_id", cluster.MasterFlavorID)
	d.Set("cluster_template_id", cluster.ClusterTemplateID)
	d.Set("version", cluster.Version)
	d.Set("updated_at", cluster.UpdatedAt.Format(time.RFC850))
	d.Set("created_at", cluster.CreatedAt.Format(time.RFC850))

	var pool pools.ClusterPool
	for _, p := range cluster.Pools {
		if p.IsDefault {
			pool = p
		}
	}

	p := make(map[string]interface{})
	p["uuid"] = pool.UUID
	p["name"] = pool.Name
	p["flavor_id"] = pool.FlavorID
	p["min_node_count"] = pool.MinNodeCount
	p["max_node_count"] = pool.MaxNodeCount
	p["node_count"] = pool.NodeCount
	p["docker_volume_type"] = pool.DockerVolumeType.String()
	p["docker_volume_size"] = pool.DockerVolumeSize
	p["stack_id"] = pool.StackID
	p["created_at"] = pool.CreatedAt.Format(time.RFC850)

	if err := d.Set("pool", []interface{}{p}); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish K8s reading")
	return diags
}
