package gcore

import (
	"context"
	"fmt"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/utils"

	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/loadbalancers"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceLoadBalancerV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLoadBalancerV2Read,
		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "ID of the project in which load balancer was created.",
				Optional:    true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"region_id": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "ID of the region in which load balancer was created.",
				Optional:    true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				DiffSuppressFunc: suppressDiffRegionID,
			},
			"project_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the project in which load balancer was created.",
				Optional:    true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the region in which load balancer was created.",
				Optional:    true,
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
			"vip_address": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Load balancer IP address.",
				Computed:    true,
			},
			"vip_port_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Load balancer Port ID.",
				Computed:    true,
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
				Type:         schema.TypeString,
				Description:  fmt.Sprintf("Available values are '%s', '%s', '%s'", types.IPv4IPFamilyType, types.IPv6IPFamilyType, types.DualStackIPFamilyType),
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(types.IPFamilyType("").StringList(), false),
			},
			"preferred_connectivity": &schema.Schema{
				Type:        schema.TypeString,
				Description: fmt.Sprintf("Available values are '%s', '%s'", types.PreferredConnectivityL2, types.PreferredConnectivityL3),
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
			"metadata_k": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Metadata string of the load balancer.",
				Optional:    true,
			},
			"metadata_kv": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "Metadata map of the load balancer.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"metadata_read_only": &schema.Schema{
				Type:        schema.TypeList,
				Description: "List of metadata items.",
				Computed:    true,
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

func dataSourceLoadBalancerV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LoadBalancer reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, LoadBalancersPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)

	metaOpts := &loadbalancers.ListOpts{}

	if metadataK, ok := d.GetOk("metadata_k"); ok {
		metaOpts.MetadataK = metadataK.(string)
	}

	if metadataRaw, ok := d.GetOk("metadata_kv"); ok {
		meta, err := utils.MapInterfaceToMapString(metadataRaw)
		if err != nil {
			return diag.FromErr(err)
		}
		metaOpts.MetadataKV = meta
	}

	lbs, err := loadbalancers.ListAll(client, *metaOpts)

	if err != nil {
		return diag.FromErr(err)
	}

	var found bool
	var lb loadbalancers.LoadBalancer
	for _, l := range lbs {
		if l.Name == name {
			lb = l
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("load balancer with name %s not found", name)
	}

	vrrpIps := make([]map[string]string, len(lb.VrrpIPs))
	for i, vrrpIp := range lb.VrrpIPs {
		v := map[string]string{"subnet_id": "", "ip_address": ""}
		v["subnet_id"] = vrrpIp.SubnetID
		v["ip_address"] = vrrpIp.IpAddress.String()
		vrrpIps[i] = v
	}

	additionalVIPs := make([]map[string]string, len(lb.AdditionalVips))
	for i, vip := range lb.AdditionalVips {
		v := map[string]string{"subnet_id": "", "ip_address": ""}
		v["subnet_id"] = vip.SubnetID
		v["ip_address"] = vip.IpAddress.String()
		additionalVIPs[i] = v
	}

	d.SetId(lb.ID)
	d.Set("project_id", lb.ProjectID)
	d.Set("region_id", lb.RegionID)
	d.Set("name", lb.Name)
	d.Set("vip_address", lb.VipAddress.String())
	d.Set("vip_port_id", lb.VipPortID)
	d.Set("vrrp_ips", vrrpIps)
	d.Set("vip_ip_family", lb.VipIPFamilyType)
	d.Set("preferred_connectivity", lb.PreferredConnectivity)
	d.Set("additional_vips", additionalVIPs)

	log.Println("[DEBUG] Finish LoadBalancer reading")
	return diags
}
