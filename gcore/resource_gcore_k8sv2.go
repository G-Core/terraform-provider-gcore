package gcore

import (
	"context"
	"fmt"
	"log"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/clusters"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/pools"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	K8sPoint         = "k8s/clusters"
	tasksPoint       = "tasks"
	K8sCreateTimeout = 3600
)

var k8sCreateTimeout = time.Second * time.Duration(K8sCreateTimeout)

func resourceK8sV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceK8sV2Create,
		ReadContext:   resourceK8sV2Read,
		UpdateContext: resourceK8sV2Update,
		DeleteContext: resourceK8sV2Delete,
		Description:   "Represent k8s cluster with one default pool.",
		Timeouts: &schema.ResourceTimeout{
			Create: &k8sCreateTimeout,
			Update: &k8sCreateTimeout,
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, clusterName, err := ImportStringParser(d.Id())
				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.Set("name", clusterName)
				d.SetId(clusterName)
				return []*schema.ResourceData{d}, nil
			},
		},

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
				ForceNew: true,
			},
			"fixed_network": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"fixed_subnet": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet should have a router",
			},
			"pods_ip_pool": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"services_ip_pool": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"keypair": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pool": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"flavor_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"min_node_count": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"max_node_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"boot_volume_type": {
							Type:        schema.TypeString,
							Description: "Available values are 'standard', 'ssd_hiiops', 'cold', 'ultra'.",
							Optional:    true,
							Computed:    true,
						},
						"boot_volume_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"auto_healing_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"is_public_ipv4": {
							Type:     schema.TypeBool,
							Optional: true,
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

func resourceK8sV2Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start k8s cluster creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, K8sPoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := clusters.CreateOpts{
		Name:         d.Get("name").(string),
		FixedNetwork: d.Get("fixed_network").(string),
		FixedSubnet:  d.Get("fixed_subnet").(string),
		KeyPair:      d.Get("keypair").(string),
		Version:      d.Get("version").(string),
	}

	if podsIP, ok := d.GetOk("pods_ip_pool"); ok {
		gccidr, err := parseCIDRFromString(podsIP.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		opts.PodsIPPool = &gccidr
	}

	if svcIP, ok := d.GetOk("services_ip_pool"); ok {
		gccidr, err := parseCIDRFromString(svcIP.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		opts.ServicesIPPool = &gccidr
	}

	for _, poolRaw := range d.Get("pool").([]interface{}) {
		pool := poolRaw.(map[string]interface{})
		opts.Pools = append(opts.Pools, pools.CreateOpts{
			Name:               pool["name"].(string),
			FlavorID:           pool["flavor_id"].(string),
			MinNodeCount:       pool["min_node_count"].(int),
			MaxNodeCount:       pool["max_node_count"].(int),
			BootVolumeSize:     pool["boot_volume_size"].(int),
			BootVolumeType:     volumes.VolumeType(pool["boot_volume_type"].(string)),
			AutoHealingEnabled: pool["auto_healing_enabled"].(bool),
			IsPublicIPv4:       pool["is_public_ipv4"].(bool),
		})
	}

	results, err := clusters.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)

	tasksClient, err := CreateClient(provider, d, tasksPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	clusterName, err := tasks.WaitTaskAndReturnResult(tasksClient, taskID, true, K8sCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		cluster, err := clusters.Get(client, opts.Name).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot create k8s cluster with name: %s. Error: %w", opts.Name, err)
		}
		return cluster.Name, nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(clusterName.(string))
	resourceK8sV2Read(ctx, d, m)

	log.Printf("[DEBUG] Finish k8s cluster creating (%s)", clusterName)
	return diags
}

func resourceK8sV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start k8s cluster reading")
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
	if err := d.Set("pool", ps); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish k8s cluster reading")
	return diags
}

func resourceK8sV2Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start k8s cluster updating")
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, K8sPoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	tasksClient, err := CreateClient(provider, d, tasksPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterName := d.Get("name").(string)

	if d.HasChange("version") {
		upgradeOpts := clusters.UpgradeOpts{
			Version: d.Get("version").(string),
		}
		results, err := clusters.Upgrade(client, clusterName, upgradeOpts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		taskID := results.Tasks[0]
		log.Printf("[DEBUG] Task id (%s)", taskID)
		_, err = tasks.WaitTaskAndReturnResult(tasksClient, taskID, true, K8sCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
			return nil, nil
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("pool") {
		add, upd, del := diffK8sV2ClusterPoolChange(d.GetChange("pool"))
		for _, pool := range del {
			poolName := pool["name"].(string)
			log.Printf("[DEBUG] Removing pool (%s)", poolName)

			results, err := pools.Delete(client, clusterName, poolName).Extract()
			if err != nil {
				return diag.FromErr(err)
			}

			taskID := results.Tasks[0]
			log.Printf("[DEBUG] Task id (%s)", taskID)
			_, err = tasks.WaitTaskAndReturnResult(tasksClient, taskID, true, K8sCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
				return nil, nil
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, pool := range upd {
			poolName := pool["name"].(string)
			log.Printf("[DEBUG] Updating pool (%s)", poolName)

			opts := pools.UpdateOpts{
				AutoHealingEnabled: pool["auto_healing_enabled"].(bool),
				MinNodeCount:       pool["min_node_count"].(int),
				MaxNodeCount:       pool["max_node_count"].(int),
			}
			_, err := pools.Update(client, clusterName, poolName, opts).Extract()
			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, pool := range add {
			poolName := pool["name"].(string)
			log.Printf("[DEBUG] Creating pool (%s)", poolName)

			opts := pools.CreateOpts{
				Name:               pool["name"].(string),
				FlavorID:           pool["flavor_id"].(string),
				MinNodeCount:       pool["min_node_count"].(int),
				MaxNodeCount:       pool["max_node_count"].(int),
				BootVolumeSize:     pool["boot_volume_size"].(int),
				BootVolumeType:     volumes.VolumeType(pool["boot_volume_type"].(string)),
				AutoHealingEnabled: pool["auto_healing_enabled"].(bool),
				IsPublicIPv4:       pool["is_public_ipv4"].(bool),
			}
			results, err := pools.Create(client, clusterName, opts).Extract()
			if err != nil {
				return diag.FromErr(err)
			}

			taskID := results.Tasks[0]
			log.Printf("[DEBUG] Task id (%s)", taskID)
			_, err = tasks.WaitTaskAndReturnResult(tasksClient, taskID, true, K8sCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
				return nil, nil
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	diags := resourceK8sV2Read(ctx, d, m)
	log.Printf("[DEBUG] Finish k8s cluster updating (%s)", clusterName)
	return diags
}

func resourceK8sV2Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start k8s cluster deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, K8sPoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterName := d.Get("name").(string)
	results, err := clusters.Delete(client, clusterName).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)

	tasksClient, err := CreateClient(provider, d, tasksPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = tasks.WaitTaskAndReturnResult(tasksClient, taskID, true, K8sCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := clusters.Get(client, clusterName).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete k8s cluster with name: %s", clusterName)
		}
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			return nil, nil
		default:
			return nil, err
		}
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[DEBUG] Finish k8s cluster deleting")
	return diags
}

func diffK8sV2ClusterPoolChange(old, new interface{}) ([]map[string]interface{}, []map[string]interface{}, []map[string]interface{}) {
	oldmap := map[string]map[string]interface{}{}
	for _, o := range old.([]interface{}) {
		v := o.(map[string]interface{})
		oldmap[v["name"].(string)] = v
	}
	newmap := map[string]map[string]interface{}{}
	for _, n := range new.([]interface{}) {
		v := n.(map[string]interface{})
		newmap[v["name"].(string)] = v
	}

	var add, upd, del []map[string]interface{}
	for k, v := range newmap {
		if _, ok := oldmap[k]; ok {
			upd = append(upd, v)
		} else {
			add = append(add, v)
		}
	}
	for k, v := range oldmap {
		if _, ok := newmap[k]; !ok {
			del = append(del, v)
		}
	}
	return add, upd, del
}
