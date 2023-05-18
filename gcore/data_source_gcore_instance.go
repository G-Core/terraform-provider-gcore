package gcore

import (
	"context"
	"log"
	"strconv"

	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceRead,
		Description: "Represent instance. Could be used with baremetal also",
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
			"flavor_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Set:      volumeUniqueID,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delete_on_termination": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"interface": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"security_group": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Firewalls list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"metadata": &schema.Schema{
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
					},
				},
			},
			"flavor": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"addresses": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"net": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addr": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Instance reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, InstancePoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	insts, err := instances.ListAll(client, instances.ListOpts{Name: name})
	if err != nil {
		return diag.FromErr(err)
	}

	var found bool
	var instance instances.Instance
	for _, l := range insts {
		if l.Name == name {
			instance = l
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("instance with name %s not found", name)
	}

	d.SetId(instance.ID)
	d.Set("name", instance.Name)

	d.Set("flavor_id", instance.Flavor.FlavorID)
	d.Set("status", instance.Status)
	d.Set("vm_state", instance.VMState)

	flavor := make(map[string]interface{}, 4)
	flavor["flavor_id"] = instance.Flavor.FlavorID
	flavor["flavor_name"] = instance.Flavor.FlavorName
	flavor["ram"] = strconv.Itoa(instance.Flavor.RAM)
	flavor["vcpus"] = strconv.Itoa(instance.Flavor.VCPUS)
	d.Set("flavor", flavor)

	extVolumes := make([]interface{}, 0, len(instance.Volumes))
	for _, vol := range instance.Volumes {
		v := make(map[string]interface{})
		v["volume_id"] = vol.ID
		v["delete_on_termination"] = vol.DeleteOnTermination
		extVolumes = append(extVolumes, v)
	}

	if err := d.Set("volume", schema.NewSet(volumeUniqueID, extVolumes)); err != nil {
		return diag.FromErr(err)
	}

	ifs, err := instances.ListInterfacesAll(client, instance.ID)
	log.Printf("instance data source interfaces: %+v", ifs)
	if err != nil {
		return diag.FromErr(err)
	}
	var cleanInterfaces []interface{}
	for _, iface := range ifs {
		if len(iface.IPAssignments) == 0 {
			continue
		}

		for _, assignment := range iface.IPAssignments {
			subnetID := assignment.SubnetID

			i := make(map[string]interface{})

			i["network_id"] = iface.NetworkID
			i["subnet_id"] = subnetID
			i["port_id"] = iface.PortID
			i["ip_address"] = iface.IPAssignments[0].IPAddress.String()

			cleanInterfaces = append(cleanInterfaces, i)
		}
	}
	if err := d.Set("interface", cleanInterfaces); err != nil {
		return diag.FromErr(err)
	}

	sliced := make([]map[string]interface{}, 0, len(instance.Metadata))
	for k, data := range instance.Metadata {
		mdata := make(map[string]interface{}, 2)
		mdata["key"] = k
		mdata["value"] = data
		sliced = append(sliced, mdata)
	}
	if err := d.Set("metadata", sliced); err != nil {
		return diag.FromErr(err)
	}

	secGrps := make([]map[string]interface{}, 0, len(instance.SecurityGroups))
	for _, sg := range instance.SecurityGroups {
		i := make(map[string]interface{})
		i["name"] = sg.Name
		secGrps = append(secGrps, i)
	}
	if err := d.Set("security_group", secGrps); err != nil {
		return diag.FromErr(err)
	}

	addresses := []map[string][]map[string]string{}
	for _, data := range instance.Addresses {
		d := map[string][]map[string]string{}
		netd := make([]map[string]string, len(data))
		for i, iaddr := range data {
			ndata := make(map[string]string, 2)
			ndata["type"] = iaddr.Type.String()
			ndata["addr"] = iaddr.Address.String()
			netd[i] = ndata
		}
		d["net"] = netd
		addresses = append(addresses, d)
	}
	if err := d.Set("addresses", addresses); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish Instance reading")
	return diags
}
