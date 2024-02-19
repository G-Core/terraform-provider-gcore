package gcore

import (
	"context"
	"fmt"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/listeners"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/loadbalancers"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		ReadContext:        dataSourceLoadBalancerRead,
		DeprecationMessage: "!> **WARNING:** This data-source is deprecated and will be removed in the next major version. Use gcore_loadbalancerv2 data-source instead",
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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vip_port_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"listener": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: fmt.Sprintf("Available values is '%s', '%s', '%s', '%s'", types.ProtocolTypeHTTP, types.ProtocolTypeHTTPS, types.ProtocolTypeTCP, types.ProtocolTypeUDP),
						},
						"protocol_port": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"metadata_k": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"metadata_kv": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"metadata_read_only": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"read_only": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceLoadBalancerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	d.SetId(lb.ID)
	d.Set("project_id", lb.ProjectID)
	d.Set("region_id", lb.RegionID)
	d.Set("name", lb.Name)
	d.Set("vip_address", lb.VipAddress.String())
	d.Set("vip_port_id", lb.VipPortID)

	metadataReadOnly := make([]map[string]interface{}, 0, len(lb.Metadata))
	if len(lb.Metadata) > 0 {
		for _, metadataItem := range lb.Metadata {
			metadataReadOnly = append(metadataReadOnly, map[string]interface{}{
				"key":       metadataItem.Key,
				"value":     metadataItem.Value,
				"read_only": metadataItem.ReadOnly,
			})
		}
	}

	if err := d.Set("metadata_read_only", metadataReadOnly); err != nil {
		return diag.FromErr(err)
	}

	listenersClient, err := CreateClient(provider, d, LBListenersPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	newListeners := make([]map[string]interface{}, len(lb.Listeners))
	for i, l := range lb.Listeners {
		listener, err := listeners.Get(listenersClient, l.ID).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		newListeners[i] = map[string]interface{}{
			"id":            listener.ID,
			"name":          listener.Name,
			"protocol":      listener.Protocol.String(),
			"protocol_port": listener.ProtocolPort,
		}
	}
	if err := d.Set("listener", newListeners); err != nil {
		diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish LoadBalancer reading")
	return diags
}
