package gcore

import (
	"context"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSubnet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubnetRead,
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
			"network_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"enable_dhcp": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"connect_to_network_router": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"dns_nameservers": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"host_routes": &schema.Schema{
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
			"gateway_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Subnet reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, subnetPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	networkID := d.Get("network_id").(string)
	subnetsOpts := &subnets.ListOpts{NetworkID: networkID}

	if metadataK, ok := d.GetOk("metadata_k"); ok {
		subnetsOpts.MetadataK = metadataK.(string)
	}
	if metadataRaw, ok := d.GetOk("metadata_kv"); ok {
		typedMetadataKV := make(map[string]string, len(metadataRaw.(map[string]interface{})))
		for k, v := range metadataRaw.(map[string]interface{}) {
			typedMetadataKV[k] = v.(string)
		}
		subnetsOpts.MetadataKV = typedMetadataKV
	}

	snets, err := subnets.ListAll(client, *subnetsOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	var found bool
	var subnet subnets.Subnet
	for _, sn := range snets {
		if sn.Name == name {
			subnet = sn
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("subnet with name %s not found", name)
	}

	d.SetId(subnet.ID)
	d.Set("name", subnet.Name)
	d.Set("enable_dhcp", subnet.EnableDHCP)
	d.Set("cidr", subnet.CIDR.String())
	d.Set("network_id", subnet.NetworkID)

	metadataReadOnly := make([]map[string]interface{}, 0, len(subnet.Metadata))
	if len(subnet.Metadata) > 0 {
		for _, metadataItem := range subnet.Metadata {
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

	dns := make([]string, len(subnet.DNSNameservers))
	for i, ns := range subnet.DNSNameservers {
		dns[i] = ns.String()
	}
	d.Set("dns_nameservers", dns)

	hrs := make([]map[string]string, len(subnet.HostRoutes))
	for i, hr := range subnet.HostRoutes {
		hR := map[string]string{"destination": "", "nexthop": ""}
		hR["destination"] = hr.Destination.String()
		hR["nexthop"] = hr.NextHop.String()
		hrs[i] = hR
	}
	d.Set("host_routes", hrs)
	d.Set("region_id", subnet.RegionID)
	d.Set("project_id", subnet.ProjectID)
	d.Set("gateway_ip", subnet.GatewayIP.String())

	d.Set("connect_to_network_router", true)
	if subnet.GatewayIP == nil {
		d.Set("connect_to_network_router", false)
		d.Set("gateway_ip", "disable")
	}

	log.Println("[DEBUG] Finish Subnet reading")
	return diags
}
