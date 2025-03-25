package gcore

import (
	"context"
	"fmt"
	"log"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/loadbalancers"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata/v1/metadata"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLoadBalancerV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerV2Create,
		ReadContext:   resourceLoadBalancerV2Read,
		UpdateContext: resourceLoadBalancerV2Update,
		DeleteContext: resourceLoadBalancerDelete,
		Description:   "Represent load balancer without nested listener",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(LoadBalancerResourceTimeoutMinutes * time.Minute),
			Delete: schema.DefaultTimeout(LoadBalancerResourceTimeoutMinutes * time.Minute),
			Update: schema.DefaultTimeout(LoadBalancerResourceTimeoutMinutes * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, lbID, err := ImportStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(lbID)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "ID of the desired project to create load balancer in. Alternative for `project_name`. One of them should be specified.",
				Optional:    true,
				ForceNew:    true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"region_id": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "ID of the desired region to create load balancer in. Alternative for `region_name`. One of them should be specified.",
				Optional:    true,
				ForceNew:    true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				DiffSuppressFunc: suppressDiffRegionID,
			},
			"project_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the desired project to create load balancer in. Alternative for `project_id`. One of them should be specified.",
				Optional:    true,
				ForceNew:    true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the desired region to create load balancer in. Alternative for `region_id`. One of them should be specified.",
				Optional:    true,
				ForceNew:    true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the load balancer.",
				Required:    true,
			},
			"flavor": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Desired flavor to be used for load balancer. By default, `lb1-1-2` will be used. ",
				Optional:    true,
			},
			"vip_network_id": &schema.Schema{
				Type: schema.TypeString,
				Description: "ID of the desired network. " +
					"Can be used with vip_subnet_id, in this case Load Balancer will be created in specified subnet, otherwise in most free subnet. " +
					"Note: add all created `gcore_subnet` resources within the network with this id to the `depends_on` to be sure that `gcore_loadbalancerv2` will be destroyed first",
				Optional: true,
				ForceNew: true,
			},
			"vip_subnet_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "ID of the desired subnet. Should be used together with vip_network_id.",
				Optional:    true,
				ForceNew:    true,
			},
			"vip_address": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Load balancer IP address. IP address will be changed when load balancer will be recreated if `vip_port_id` is not specified.",
				Computed:    true,
			},
			"vip_port_id": &schema.Schema{
				Type: schema.TypeString,
				Description: "Load balancer Port ID. It might be ID of the already created Reserved Fixed IP, otherwise we will create port automatically in specified `vip_network_id`/`vip_subnet_id`. " +
					"It is an alternative for specifying `vip_network_id`/`vip_subnet_id`.",
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vrrp_ips": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:        schema.TypeString,
							Description: "IP address of the LB instance.",
							Computed:    true,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "Subnet ID of the LB instance.",
							Computed:    true,
						},
					},
				},
			},
			"vip_ip_family": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: fmt.Sprintf("Available values are '%s', '%s', '%s'", types.IPv4IPFamilyType, types.IPv6IPFamilyType, types.DualStackIPFamilyType),
				ValidateDiagFunc: func(val interface{}, key cty.Path) diag.Diagnostics {
					v := val.(string)
					switch types.IPFamilyType(v) {
					case types.IPv4IPFamilyType, types.IPv6IPFamilyType, types.DualStackIPFamilyType:
						return diag.Diagnostics{}
					}
					return diag.Errorf("wrong type %s, available values are '%s', '%s', '%s'", v, types.IPv4IPFamilyType, types.IPv6IPFamilyType, types.DualStackIPFamilyType)
				},
			},
			"preferred_connectivity": &schema.Schema{
				Type:        schema.TypeString,
				Description: fmt.Sprintf("Available values are '%s', '%s'", types.PreferredConnectivityL2, types.PreferredConnectivityL3),
				Optional:    true,
				Computed:    true,
			},
			"additional_vips": &schema.Schema{
				Type:        schema.TypeList,
				Description: "Load Balancer additional VIPs",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:        schema.TypeString,
							Description: "Load Balancer additional VIP",
							Computed:    true,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "Load Balancer additional VIP subnet ID",
							Computed:    true,
						},
					},
				},
			},
			"last_updated": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Datetime when load balancer was updated at the last time.",
				Computed:    true,
			},
			"metadata_map": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Metadata map to apply to the load balancer.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"metadata_read_only": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of metadata items.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Description: "Key of the metadata (tag) item.",
							Computed:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "Value of the metadata (tag) item.",
							Computed:    true,
						},
						"read_only": {
							Type:        schema.TypeBool,
							Description: "Is the current key read-only or not.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceLoadBalancerV2Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LoadBalancer creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, LoadBalancersPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := loadbalancers.CreateOpts{
		Name:                  d.Get("name").(string),
		VipNetworkID:          d.Get("vip_network_id").(string),
		VipSubnetID:           d.Get("vip_subnet_id").(string),
		VipPortID:             d.Get("vip_port_id").(string),
		VIPIPFamily:           types.IPFamilyType(d.Get("vip_ip_family").(string)),
		PreferredConnectivity: types.PreferredConnectivityType(d.Get("preferred_connectivity").(string)),
	}

	if metadataRaw, ok := d.GetOk("metadata_map"); ok {
		meta, err := utils.MapInterfaceToMapString(metadataRaw)
		if err != nil {
			return diag.FromErr(err)
		}
		opts.Metadata = meta
	}

	lbFlavor := d.Get("flavor").(string)
	if len(lbFlavor) != 0 {
		opts.Flavor = &lbFlavor
	}
	timeout := int(d.Timeout(schema.TimeoutCreate).Seconds())
	rc := GetConflictRetryConfig(timeout)
	results, err := loadbalancers.Create(client, opts, &gcorecloud.RequestOpts{
		ConflictRetryAmount:   rc.Amount,
		ConflictRetryInterval: rc.Interval,
	}).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	lbID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, timeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		lbID, err := loadbalancers.ExtractLoadBalancerIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve LoadBalancer ID from task info: %w", err)
		}
		return lbID, nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(lbID.(string))
	resourceLoadBalancerV2Read(ctx, d, m)

	log.Printf("[DEBUG] Finish LoadBalancer creating (%s)", lbID)
	return diags
}

func resourceLoadBalancerV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LoadBalancer reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, LoadBalancersPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	lb, err := loadbalancers.Get(client, d.Id(), nil).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("project_id", lb.ProjectID)
	d.Set("region_id", lb.RegionID)
	d.Set("name", lb.Name)
	d.Set("flavor", lb.Flavor.FlavorName)
	d.Set("vip_port_id", lb.VipPortID)
	d.Set("vrrp_ips", lb.VrrpIPs)
	d.Set("vip_ip_family", lb.VipIPFamilyType)
	d.Set("preferred_connectivity", lb.PreferredConnectivity)

	if lb.VipAddress != nil {
		d.Set("vip_address", lb.VipAddress.String())
	}

	fields := []string{"vip_network_id", "vip_subnet_id"}
	revertState(d, &fields)

	vrrpIps := make([]map[string]string, len(lb.VrrpIPs))
	for i, vrrpIp := range lb.VrrpIPs {
		v := map[string]string{"subnet_id": "", "ip_address": ""}
		v["subnet_id"] = vrrpIp.SubnetID
		v["ip_address"] = vrrpIp.IpAddress.String()
		vrrpIps[i] = v
	}
	d.Set("vrrp_ips", vrrpIps)

	additionalVIPs := make([]map[string]string, len(lb.AdditionalVips))
	for i, vip := range lb.AdditionalVips {
		v := map[string]string{"subnet_id": "", "ip_address": ""}
		v["subnet_id"] = vip.SubnetID
		v["ip_address"] = vip.IpAddress.String()
		additionalVIPs[i] = v
	}
	d.Set("additional_vips", additionalVIPs)

	metadataMap := make(map[string]string)
	metadataReadOnly := make([]map[string]interface{}, 0, len(lb.Metadata))

	if len(lb.Metadata) > 0 {
		for _, metadataItem := range lb.Metadata {
			if !metadataItem.ReadOnly {
				metadataMap[metadataItem.Key] = metadataItem.Value
			}
			metadataReadOnly = append(metadataReadOnly, map[string]interface{}{
				"key":       metadataItem.Key,
				"value":     metadataItem.Value,
				"read_only": metadataItem.ReadOnly,
			})
		}
	}

	if err := d.Set("metadata_map", metadataMap); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata_read_only", metadataReadOnly); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish LoadBalancer reading")
	return diags
}

func resourceLoadBalancerV2Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LoadBalancer updating")
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, LoadBalancersPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	var change bool
	updateOpts := loadbalancers.UpdateOpts{}
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
		change = true
	}
	if d.HasChange("preferred_connectivity") {
		updateOpts.PreferredConnectivity = types.PreferredConnectivityType(d.Get("preferred_connectivity").(string))
		change = true
	}
	if change {
		_, err = loadbalancers.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	if d.HasChange("metadata_map") {
		_, nmd := d.GetChange("metadata_map")

		meta, err := utils.MapInterfaceToMapString(nmd.(map[string]interface{}))
		if err != nil {
			return diag.Errorf("cannot get metadata. Error: %s", err)
		}

		err = metadata.MetadataReplace(client, d.Id(), meta).Err
		if err != nil {
			return diag.Errorf("cannot update metadata. Error: %s", err)
		}
	}

	if d.HasChange("flavor") {
		flavor := d.Get("flavor").(string)
		timeout := int(d.Timeout(schema.TimeoutUpdate).Seconds())
		rc := GetConflictRetryConfig(timeout)
		results, err := loadbalancers.Resize(client, d.Id(), loadbalancers.ResizeOpts{
			Flavor: flavor,
		}, &gcorecloud.RequestOpts{
			ConflictRetryAmount:   rc.Amount,
			ConflictRetryInterval: rc.Interval,
		}).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		taskID := results.Tasks[0]
		log.Printf("[DEBUG] Task id (%s)", taskID)
		taskState, err := tasks.WaitTaskAndReturnResult(client, taskID, true, timeout, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			return taskInfo.State, nil
		})
		log.Printf("[DEBUG] Task state (%s)", taskState)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Println("[DEBUG] Finish LoadBalancer updating")
	return resourceLoadBalancerV2Read(ctx, d, m)
}
