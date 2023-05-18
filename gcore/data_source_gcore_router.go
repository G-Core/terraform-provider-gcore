package gcore

import (
	"context"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/router/v1/routers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRouter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRouterRead,
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
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_gateway_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_snat": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_fixed_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"subnet_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"interfaces": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"routes": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"nexthop": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv4 address to forward traffic to if it's destination IP matches 'destination' CIDR",
						},
					},
				},
			},
		},
	}
}

func dataSourceRouterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Router reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, RouterPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	rs, err := routers.ListAll(client, routers.ListOpts{})
	if err != nil {
		return diag.FromErr(err)
	}

	var found bool
	var router routers.Router
	for _, r := range rs {
		if r.Name == name {
			router = r
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("router with name %s not found", name)
	}

	d.SetId(router.ID)
	d.Set("name", router.Name)
	d.Set("status", router.Status)

	if len(router.ExternalGatewayInfo.ExternalFixedIPs) > 0 {
		egi := make(map[string]interface{}, 4)
		egilst := make([]map[string]interface{}, 1)
		egi["enable_snat"] = router.ExternalGatewayInfo.EnableSNat
		egi["network_id"] = router.ExternalGatewayInfo.NetworkID

		efip := make([]map[string]string, len(router.ExternalGatewayInfo.ExternalFixedIPs))
		for i, fip := range router.ExternalGatewayInfo.ExternalFixedIPs {
			tmpfip := make(map[string]string, 1)
			tmpfip["ip_address"] = fip.IPAddress
			tmpfip["subnet_id"] = fip.SubnetID
			efip[i] = tmpfip
		}
		egi["external_fixed_ips"] = efip

		egilst[0] = egi
		d.Set("external_gateway_info", egilst)
	}

	ifs := make([]map[string]interface{}, 0, len(router.Interfaces))
	for _, iface := range router.Interfaces {
		for _, subnet := range iface.IPAssignments {
			smap := make(map[string]interface{}, 6)
			smap["port_id"] = iface.PortID
			smap["network_id"] = iface.NetworkID
			smap["mac_address"] = iface.MacAddress.String()
			smap["type"] = "subnet"
			smap["subnet_id"] = subnet.SubnetID
			smap["ip_address"] = subnet.IPAddress.String()
			ifs = append(ifs, smap)
		}
	}
	d.Set("interfaces", ifs)

	rss := make([]map[string]string, len(router.Routes))
	for i, r := range router.Routes {
		rmap := make(map[string]string, 2)
		rmap["destination"] = r.Destination.String()
		rmap["nexthop"] = r.NextHop.String()
		rss[i] = rmap
	}
	d.Set("routes", rss)

	log.Println("[DEBUG] Finish router reading")
	return diags
}
