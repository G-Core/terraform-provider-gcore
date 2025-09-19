package gcore

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/G-Core/gcorelabscloud-go/gcore/utils"

	"github.com/G-Core/gcorelabscloud-go/gcore/floatingip/v1/floatingips"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFloatingIP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFloatingIPRead,
		Description: "A floating IP is a static IP address that points to one of your Instances. It allows you to redirect network traffic to any of your Instances in the same datacenter.",
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
			"floating_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ExactlyOneOf: []string{
					"fixed_ip_address",
					"floating_ip_address",
				},
				ValidateDiagFunc: func(val interface{}, key cty.Path) diag.Diagnostics {
					v := val.(string)
					ip := net.ParseIP(v)
					if ip != nil {
						return diag.Diagnostics{}
					}

					return diag.FromErr(fmt.Errorf("%q must be a valid ip, got: %s", key, v))
				},
			},
			"fixed_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ExactlyOneOf: []string{
					"fixed_ip_address",
					"floating_ip_address",
				},
				ValidateDiagFunc: func(val interface{}, key cty.Path) diag.Diagnostics {
					v := val.(string)
					ip := net.ParseIP(v)
					if ip != nil {
						return diag.Diagnostics{}
					}

					return diag.FromErr(fmt.Errorf("%q must be a valid ip, got: %s", key, v))
				},
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"router_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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
			"tags": &schema.Schema{
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

func dataSourceFloatingIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FloatingIP reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, floatingIPsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	ipAddr := d.Get("floating_ip_address").(string)
	if ipAddr == "" {
		ipAddr = d.Get("fixed_ip_address").(string)
	}

	metaOpts := &floatingips.ListOpts{}

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

	ips, err := floatingips.ListAll(client, *metaOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	var found bool
	var floatingIP floatingips.FloatingIPDetail
	for _, ip := range ips {
		if ip.FloatingIPAddress.String() == ipAddr {
			floatingIP = ip
			found = true
			break
		}

		if ip.FixedIPAddress.String() == ipAddr {
			floatingIP = ip
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("floatingIP %s not found", ipAddr)
	}

	d.SetId(floatingIP.ID)
	if floatingIP.FixedIPAddress != nil {
		d.Set("fixed_ip_address", floatingIP.FixedIPAddress.String())
	} else {
		d.Set("fixed_ip_address", "")
	}

	d.Set("project_id", floatingIP.ProjectID)
	d.Set("region_id", floatingIP.RegionID)
	d.Set("status", floatingIP.Status)
	d.Set("port_id", floatingIP.PortID)
	d.Set("router_id", floatingIP.RouterID)
	d.Set("floating_ip_address", floatingIP.FloatingIPAddress.String())

	metadataReadOnly := make([]map[string]interface{}, 0, len(floatingIP.Metadata))
	if len(floatingIP.Metadata) > 0 {
		for _, metadataItem := range floatingIP.Metadata {
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

	if err := d.Set("tags", metadataReadOnly); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish FloatingIP reading")
	return diags
}
