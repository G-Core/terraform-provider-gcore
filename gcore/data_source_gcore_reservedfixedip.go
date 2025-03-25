package gcore

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/G-Core/gcorelabscloud-go/gcore/reservedfixedip/v1/reservedfixedips"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceReservedFixedIP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReservedFixedIPRead,
		Description: "Represent reserved ips",
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
			"fixed_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: func(val interface{}, key cty.Path) diag.Diagnostics {
					v := val.(string)
					ip := net.ParseIP(v)
					if ip != nil {
						return diag.Diagnostics{}
					}

					return diag.FromErr(fmt.Errorf("%q must be a valid ip, got: %s", key, v))
				},
			},
			"fixed_ipv6_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_v6_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_vip": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"port_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "ID of the port_id underlying the reserved fixed IP",
				Computed:    true,
			},
			"allowed_address_pairs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceReservedFixedIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start ReservedFixedIP reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, reservedFixedIPsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	ipAddr := d.Get("fixed_ip_address").(string)
	ips, err := reservedfixedips.ListAll(client, reservedfixedips.ListOpts{})
	if err != nil {
		return diag.FromErr(err)
	}

	var found bool
	var reservedFixedIP reservedfixedips.ReservedFixedIP
	for _, ip := range ips {
		if ip.FixedIPAddress.String() == ipAddr {
			reservedFixedIP = ip
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("reserved fixed ip %s not found", ipAddr)
	}

	// should we use PortID as id?
	d.SetId(reservedFixedIP.PortID)
	d.Set("project_id", reservedFixedIP.ProjectID)
	d.Set("region_id", reservedFixedIP.RegionID)
	d.Set("status", reservedFixedIP.Status)
	d.Set("fixed_ip_address", reservedFixedIP.FixedIPAddress.String())
	d.Set("fixed_ipv6_address", reservedFixedIP.FixedIPv6Address.String())
	d.Set("subnet_id", reservedFixedIP.SubnetID)
	d.Set("subnet_v6_id", reservedFixedIP.Subnetv6ID)
	d.Set("network_id", reservedFixedIP.NetworkID)
	d.Set("is_vip", reservedFixedIP.IsVip)
	d.Set("port_id", reservedFixedIP.PortID)

	allowedPairs := make([]map[string]interface{}, len(reservedFixedIP.AllowedAddressPairs))
	for i, p := range reservedFixedIP.AllowedAddressPairs {
		pair := make(map[string]interface{})

		pair["ip_address"] = p.IPAddress
		pair["mac_address"] = p.MacAddress

		allowedPairs[i] = pair
	}
	if err := d.Set("allowed_address_pairs", allowedPairs); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish ReservedFixedIP reading")
	return diags
}
