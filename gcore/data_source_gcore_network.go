package gcore

import (
	"context"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/availablenetworks"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkRead,
		Description: "Represent network. A network is a software-defined network in a cloud computing infrastructure",
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
			"mtu": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "'vlan' or 'vxlan' network type is allowed. Default value is 'vxlan'",
			},
			"external": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"shared": &schema.Schema{
				Type:     schema.TypeBool,
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

func dataSourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Network reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, networksPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	clientShared, err := CreateClient(provider, d, sharedNetworksPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)

	metaOpts := &networks.ListOpts{}

	if metadataK, ok := d.GetOk("metadata_k"); ok {
		metaOpts.MetadataK = metadataK.(string)
	}

	if metadataRaw, ok := d.GetOk("metadata_kv"); ok {
		typedMetadataKV := make(map[string]string, len(metadataRaw.(map[string]interface{})))
		for k, v := range metadataRaw.(map[string]interface{}) {
			typedMetadataKV[k] = v.(string)
		}
		metaOpts.MetadataKV = typedMetadataKV
	}

	nets, err := networks.ListAll(client, *metaOpts)

	if err != nil {
		return diag.FromErr(err)
	}

	// todo refactor, also refactor inner func
	var rawNetwork map[string]interface{}
	network, found := findNetworkByName(name, nets)
	if !found {
		// trying to find among shared networks
		nets, err := availablenetworks.ListAll(clientShared, nil)
		if err != nil {
			return diag.FromErr(err)
		}
		network, found := findSharedNetworkByName(name, nets)
		if !found {
			return diag.Errorf("network with name %s not found", name)
		}

		rawNetwork, err = StructToMap(network)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		rawNetwork, err = StructToMap(network)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(rawNetwork["id"].(string))
	d.Set("name", rawNetwork["name"])
	d.Set("mtu", rawNetwork["mtu"])
	d.Set("type", rawNetwork["type"])
	d.Set("region_id", rawNetwork["region_id"])
	d.Set("project_id", rawNetwork["project_id"])
	d.Set("external", rawNetwork["external"])
	d.Set("shared", rawNetwork["shared"])

	metadataReadOnly := make([]map[string]interface{}, 0, len(network.Metadata))
	if len(network.Metadata) > 0 {
		for _, metadataItem := range network.Metadata {
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

	log.Println("[DEBUG] Finish Network reading")
	return diags
}
