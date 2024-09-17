package gcore

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/floatingip/v1/floatingips"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	InstanceDeleting        int = 1200
	InstanceCreatingTimeout int = 1200
	InstancePoint               = "instances"

	InstanceVMStateActive  = "active"
	InstanceVMStateStopped = "stopped"
)

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext:      resourceInstanceCreate,
		ReadContext:        resourceInstanceRead,
		UpdateContext:      resourceInstanceUpdate,
		DeleteContext:      resourceInstanceDelete,
		Description:        "Represent instance",
		DeprecationMessage: "!> **WARNING:** This resource is deprecated, please use 'gcore_instancev2' instead",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, InstanceID, err := ImportStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(InstanceID)

				return []*schema.ResourceData{d}, nil
			},
		},

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
			"flavor_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name_templates": &schema.Schema{
				Type:          schema.TypeList,
				Optional:      true,
				Deprecated:    "Use name_template instead",
				ConflictsWith: []string{"name_template"},
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"name_template": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name_templates"},
			},
			"volume": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Set:      volumeUniqueID,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"source": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Currently available only 'existing-volume' value",
							ValidateDiagFunc: func(val interface{}, key cty.Path) diag.Diagnostics {
								v := val.(string)
								if types.VolumeSource(v) == types.ExistingVolume {
									return diag.Diagnostics{}
								}
								return diag.Errorf("wrong source type %s, now available values is '%s'", v, types.ExistingVolume)
							},
						},
						"boot_index": {
							Type:        schema.TypeInt,
							Description: "If boot_index==0 volumes can not detached",
							Optional:    true,
						},
						"type_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"volume_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"attachment_tag": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"delete_on_termination": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"interface": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: fmt.Sprintf("Available value is '%s', '%s', '%s', '%s'", types.SubnetInterfaceType, types.AnySubnetInterfaceType, types.ExternalInterfaceType, types.ReservedFixedIpType),
						},
						"order": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Order of attaching interface",
						},
						"network_id": {
							Type:        schema.TypeString,
							Description: "required if type is 'subnet' or 'any_subnet'",
							Optional:    true,
							Computed:    true,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "required if type is 'subnet'",
							Optional:    true,
							Computed:    true,
						},
						// nested map is not supported, in this case, you do not need to use the list for the map
						"fip_source": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"existing_fip_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "required if type is  'reserved_fixed_ip'",
							Optional:    true,
						},
						"security_groups": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "list of security group IDs",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
			"keypair_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Firewalls list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Firewall unique id",
							Required:    true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"metadata": &schema.Schema{
				Type:          schema.TypeList,
				Optional:      true,
				Deprecated:    "Use metadata_map instead",
				ConflictsWith: []string{"metadata_map"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"metadata_map": &schema.Schema{
				Type:          schema.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"metadata"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"configuration": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"userdata": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "**Deprecated**",
				Deprecated:    "Use user_data instead",
				ConflictsWith: []string{"user_data"},
			},
			"user_data": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"userdata"},
			},
			"allow_app_ports": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"flavor": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vm_state": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: fmt.Sprintf("Current vm state, use %s to stop vm and %s to start", InstanceVMStateStopped, InstanceVMStateActive),
				ValidateFunc: validation.StringInSlice([]string{
					InstanceVMStateActive, InstanceVMStateStopped,
				}, true),
			},
			"addresses": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"net": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addr": {
										Type:     schema.TypeString,
										Required: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Instance creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	clientv1, err := CreateClient(provider, d, InstancePoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	clientv2, err := CreateClient(provider, d, InstancePoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := instances.CreateOpts{SecurityGroups: []gcorecloud.ItemID{}}

	createOpts.Flavor = d.Get("flavor_id").(string)
	createOpts.Password = d.Get("password").(string)
	createOpts.Username = d.Get("username").(string)
	createOpts.Keypair = d.Get("keypair_name").(string)
	createOpts.ServerGroupID = d.Get("server_group").(string)

	if userData, ok := d.GetOk("userdata"); ok {
		createOpts.UserData = userData.(string)
	} else if userData, ok := d.GetOk("user_data"); ok {
		createOpts.UserData = userData.(string)
	}

	name := d.Get("name").(string)
	if len(name) > 0 {
		createOpts.Names = []string{name}
	}

	if nameTemplatesRaw, ok := d.GetOk("name_templates"); ok {
		nameTemplates := nameTemplatesRaw.([]interface{})
		if len(nameTemplates) > 0 {
			NameTemp := make([]string, len(nameTemplates))
			for i, nametemp := range nameTemplates {
				NameTemp[i] = nametemp.(string)
			}
			createOpts.NameTemplates = NameTemp
		}
	} else if nameTemplate, ok := d.GetOk("name_template"); ok {
		createOpts.NameTemplates = []string{nameTemplate.(string)}
	}

	createOpts.AllowAppPorts = d.Get("allow_app_ports").(bool)

	currentVols := d.Get("volume").(*schema.Set).List()
	if len(currentVols) > 0 {
		vs, err := extractVolumesMap(currentVols)
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.Volumes = vs
	}

	ifs := d.Get("interface").([]interface{})
	// sort interfaces by 'order' key to attach it in right order
	sort.Sort(instanceInterfaces(ifs))
	if len(ifs) > 0 {
		ifaces, err := extractInstanceInterfacesMap(ifs)
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.Interfaces = ifaces
	}

	if metadata, ok := d.GetOk("metadata"); ok {
		if len(metadata.([]interface{})) > 0 {
			md, err := extractKeyValue(metadata.([]interface{}))
			if err != nil {
				return diag.FromErr(err)
			}
			createOpts.Metadata = &md
		}
	} else if metadataRaw, ok := d.GetOk("metadata_map"); ok {
		md := extractMetadataMap(metadataRaw.(map[string]interface{}))
		createOpts.Metadata = &md
	}

	configuration := d.Get("configuration")
	if len(configuration.([]interface{})) > 0 {
		conf, err := extractKeyValue(configuration.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.Configuration = &conf
	}

	log.Printf("[DEBUG] Interface create options: %+v", createOpts)
	results, err := instances.Create(clientv2, createOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	InstanceID, err := tasks.WaitTaskAndReturnResult(clientv1, taskID, true, InstanceCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(clientv1, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		Instance, err := instances.ExtractInstanceIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve Instance ID from task info: %w", err)
		}
		return Instance, nil
	},
	)
	log.Printf("[DEBUG] Instance id (%s)", InstanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(InstanceID.(string))
	resourceInstanceRead(ctx, d, m)

	log.Printf("[DEBUG] Finish Instance creating (%s)", InstanceID)
	return diags
}

func resourceInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Instance reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	instanceID := d.Id()
	log.Printf("[DEBUG] Instance id = %s", instanceID)

	client, err := CreateClient(provider, d, InstancePoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	instance, err := instances.Get(client, instanceID).Extract()
	if err != nil {
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			log.Printf("[WARN] Removing instance %s because resource doesn't exist anymore", d.Id())
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

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

	currentVolumes := extractVolumesIntoMap(d.Get("volume").(*schema.Set).List())

	extVolumes := make([]interface{}, 0, len(instance.Volumes))
	for _, vol := range instance.Volumes {
		v, ok := currentVolumes[vol.ID]
		// todo fix it
		if !ok {
			v = make(map[string]interface{})
			v["volume_id"] = vol.ID
			v["source"] = types.ExistingVolume.String()
		}

		v["id"] = vol.ID
		v["delete_on_termination"] = vol.DeleteOnTermination
		extVolumes = append(extVolumes, v)
	}

	if err := d.Set("volume", schema.NewSet(volumeUniqueID, extVolumes)); err != nil {
		return diag.FromErr(err)
	}

	instancePorts, err := instances.ListPortsAll(client, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}
	secGroups := prepareSecurityGroups(instancePorts)

	if err := d.Set("security_group", secGroups); err != nil {
		return diag.FromErr(err)
	}

	ifs, err := instances.ListInterfacesAll(client, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	interfaces, err := extractInstanceInterfaceIntoMap(d.Get("interface").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	var cleanInterfaces []interface{}
	for ifOrder, iface := range ifs {
		if len(iface.IPAssignments) == 0 {
			continue
		}

		for _, assignment := range iface.IPAssignments {
			subnetID := assignment.SubnetID

			// bad idea, but what to do
			var iOpts instances.InterfaceOpts
			var orderedIOpts OrderedInterfaceOpts
			var ok bool
			// we need to match our interfaces with api's interfaces
			// but we don't have any unique values, that's why we use exactly that list of keys
			for _, k := range []string{subnetID, iface.PortID, iface.NetworkID, types.ExternalInterfaceType.String()} {
				if orderedIOpts, ok = interfaces[k]; ok {
					iOpts = orderedIOpts.InterfaceOpts
					break
				}
			}

			i := make(map[string]interface{})
			if !ok {
				orderedIOpts = OrderedInterfaceOpts{Order: ifOrder}
			} else {
				i["type"] = iOpts.Type.String()
			}

			var name string
			if iface.Name != nil {
				name = *iface.Name
			}
			i["name"] = name
			i["network_id"] = iface.NetworkID
			i["subnet_id"] = subnetID
			i["port_id"] = iface.PortID
			i["order"] = orderedIOpts.Order
			if len(iface.FloatingIPDetails) > 0 {
				i["fip_source"] = types.ExistingFloatingIP
				i["existing_fip_id"] = iface.FloatingIPDetails[0].ID
			}
			i["ip_address"] = assignment.IPAddress.String()

			if port, err := findInstancePort(iface.PortID, instancePorts); err == nil {
				sgs := make([]string, len(port.SecurityGroups))
				for i, sg := range port.SecurityGroups {
					sgs[i] = sg.ID
				}
				i["security_groups"] = sgs
			}

			cleanInterfaces = append(cleanInterfaces, i)
		}
	}
	if err := d.Set("interface", cleanInterfaces); err != nil {
		return diag.FromErr(err)
	}

	if metadataRaw, ok := d.GetOk("metadata"); ok {
		metadata := metadataRaw.([]interface{})
		sliced := make([]map[string]string, len(metadata))
		for i, data := range metadata {
			d := data.(map[string]interface{})
			mdata := make(map[string]string, 2)
			md, err := instances.MetadataGet(client, instanceID, d["key"].(string)).Extract()
			if err != nil {
				return diag.Errorf("cannot get metadata with key: %s. Error: %s", instanceID, err)
			}
			mdata["key"] = md.Key
			mdata["value"] = md.Value
			sliced[i] = mdata
		}
		d.Set("metadata", sliced)
	} else {
		metadata := d.Get("metadata_map").(map[string]interface{})
		newMetadata := make(map[string]interface{}, len(metadata))
		for k := range metadata {
			md, err := instances.MetadataGet(client, instanceID, k).Extract()
			if err != nil {
				return diag.Errorf("cannot get metadata with key: %s. Error: %s", instanceID, err)
			}
			newMetadata[k] = md.Value
		}
		if err := d.Set("metadata_map", newMetadata); err != nil {
			return diag.FromErr(err)
		}
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

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Instance updating")
	instanceID := d.Id()
	log.Printf("[DEBUG] Instance id = %s", instanceID)
	config := m.(*Config)
	provider := config.Provider
	client, err := CreateClient(provider, d, InstancePoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	fipClient, err := CreateClient(provider, d, floatingIPsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		nameTemplates := d.Get("name_templates").([]interface{})
		nameTemplate := d.Get("name_template").(string)
		if len(nameTemplate) == 0 && len(nameTemplates) == 0 {
			opts := instances.RenameInstanceOpts{
				Name: d.Get("name").(string),
			}
			if _, err := instances.RenameInstance(client, instanceID, opts).Extract(); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("flavor_id") {
		flavor_id := d.Get("flavor_id").(string)
		results, err := instances.Resize(client, instanceID, instances.ChangeFlavorOpts{FlavorID: flavor_id}).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		taskID := results.Tasks[0]
		log.Printf("[DEBUG] Task id (%s)", taskID)
		taskState, err := tasks.WaitTaskAndReturnResult(client, taskID, true, InstanceCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			return taskInfo.State, nil
		},
		)
		log.Printf("[DEBUG] Task state (%s)", taskState)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("metadata") {
		omd, nmd := d.GetChange("metadata")
		if len(omd.([]interface{})) > 0 {
			for _, data := range omd.([]interface{}) {
				d := data.(map[string]interface{})
				k := d["key"].(string)
				err := instances.MetadataDelete(client, instanceID, k).Err
				if err != nil {
					return diag.Errorf("cannot delete metadata key: %s. Error: %s", k, err)
				}
			}
		}
		if len(nmd.([]interface{})) > 0 {
			var MetaData []instances.MetadataOpts
			for _, data := range nmd.([]interface{}) {
				d := data.(map[string]interface{})
				var md instances.MetadataOpts
				md.Key = d["key"].(string)
				md.Value = d["value"].(string)
				MetaData = append(MetaData, md)
			}
			createOpts := instances.MetadataSetOpts{
				Metadata: MetaData,
			}
			err := instances.MetadataCreate(client, instanceID, createOpts).Err
			if err != nil {
				return diag.Errorf("cannot create metadata. Error: %s", err)
			}
		}
	} else if d.HasChange("metadata_map") {
		omd, nmd := d.GetChange("metadata_map")
		if len(omd.(map[string]interface{})) > 0 {
			for k := range omd.(map[string]interface{}) {
				err := instances.MetadataDelete(client, instanceID, k).Err
				if err != nil {
					return diag.Errorf("cannot delete metadata key: %s. Error: %s", k, err)
				}
			}
		}
		if len(nmd.(map[string]interface{})) > 0 {
			var MetaData []instances.MetadataOpts
			for k, v := range nmd.(map[string]interface{}) {
				md := instances.MetadataOpts{
					Key:   k,
					Value: v.(string),
				}
				MetaData = append(MetaData, md)
			}
			createOpts := instances.MetadataSetOpts{
				Metadata: MetaData,
			}
			err := instances.MetadataCreate(client, instanceID, createOpts).Err
			if err != nil {
				return diag.Errorf("cannot create metadata. Error: %s", err)
			}
		}
	}

	if d.HasChange("interface") {
		iList, err := instances.ListInterfacesAll(client, instanceID)
		if err != nil {
			return diag.FromErr(err)
		}

		var fips []instances.FloatingIP
		for _, i := range iList {
			if len(i.FloatingIPDetails) == 0 {
				continue
			}

			for _, fip := range i.FloatingIPDetails {
				fips = append(fips, fip)
			}
		}

		ifsOldRaw, ifsNewRaw := d.GetChange("interface")

		ifsOld := ifsOldRaw.([]interface{})
		ifsNew := ifsNewRaw.([]interface{})

		for _, i := range ifsOld {
			iface := i.(map[string]interface{})
			var opts instances.InterfaceOpts
			opts.PortID = iface["port_id"].(string)
			opts.IpAddress = iface["ip_address"].(string)

			log.Printf("[DEBUG] detach interface: %+v", opts)
			results, err := instances.DetachInterface(client, instanceID, opts).Extract()
			if err != nil {
				return diag.FromErr(err)
			}
			taskID := results.Tasks[0]
			_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, InstanceCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
				taskInfo, err := tasks.Get(client, string(task)).Extract()
				if err != nil {
					return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w, task: %+v", task, err, taskInfo)
				}
				return nil, nil
			},
			)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		sort.Sort(instanceInterfaces(ifsNew))
		for _, i := range ifsNew {
			iface := i.(map[string]interface{})

			iType := types.InterfaceType(iface["type"].(string))
			opts := instances.InterfaceInstanceCreateOpts{
				InterfaceOpts: instances.InterfaceOpts{Type: iType},
			}

			switch iType {
			case types.SubnetInterfaceType:
				opts.SubnetID = iface["subnet_id"].(string)
			case types.AnySubnetInterfaceType:
				opts.NetworkID = iface["network_id"].(string)
			case types.ReservedFixedIpType:
				opts.PortID = iface["port_id"].(string)
			}

			rawSgsID := iface["security_groups"].([]interface{})
			sgs := make([]gcorecloud.ItemID, len(rawSgsID))
			for i, sgID := range rawSgsID {
				sgs[i] = gcorecloud.ItemID{ID: sgID.(string)}
			}
			opts.SecurityGroups = sgs

			log.Printf("[DEBUG] attach interface: %+v", opts)
			results, err := instances.AttachInterface(client, instanceID, opts).Extract()
			if err != nil {
				return diag.Errorf("cannot attach interface: %s. Error: %s", iType, err)
			}

			taskID := results.Tasks[0]
			log.Printf("[DEBUG] attach interface taskID: %s", taskID)
			if err = tasks.WaitForStatus(client, string(taskID), tasks.TaskStateFinished, InstanceCreatingTimeout, true); err != nil {
				return diag.FromErr(err)
			}
		}

		for _, fip := range fips {
			log.Printf("[DEBUG] Reassign floatin IP %s to fixed IP %s port id %s", fip.FloatingIPAddress, fip.FixedIPAddress, fip.PortID)
			mm := make(map[string]string)
			for _, i := range fip.Metadata {
				mm[i.Key] = i.Value
			}

			_, err := floatingips.Assign(fipClient, fip.ID, floatingips.CreateOpts{
				PortID:         fip.PortID,
				FixedIPAddress: fip.FixedIPAddress,
				Metadata:       mm,
			}).Extract()

			if err != nil {
				return diag.Errorf("cannot reassign floating IP %s to fixed IP %s port id %s. Error: %v", fip.FloatingIPAddress, fip.FixedIPAddress, fip.PortID, err)
			}
		}
	}

	if d.HasChange("volume") {
		vClient, err := CreateClient(provider, d, volumesPoint, versionPointV1)
		if err != nil {
			return diag.FromErr(err)
		}

		oldVolumesRaw, newVolumesRaw := d.GetChange("volume")
		oldVolumes := extractInstanceVolumesMap(oldVolumesRaw.(*schema.Set).List())
		newVolumes := extractInstanceVolumesMap(newVolumesRaw.(*schema.Set).List())

		vOpts := volumes.InstanceOperationOpts{InstanceID: d.Id()}
		for vid := range oldVolumes {
			if isAttached := newVolumes[vid]; isAttached {
				// mark as already attached
				newVolumes[vid] = false
				continue
			}
			if _, err := volumes.Detach(vClient, vid, vOpts).Extract(); err != nil {
				return diag.FromErr(err)
			}
		}

		// range over not attached volumes
		for vid, ok := range newVolumes {
			if ok {
				if _, err := volumes.Attach(vClient, vid, vOpts).Extract(); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChange("vm_state") {
		state := d.Get("vm_state").(string)
		switch state {
		case InstanceVMStateActive:
			if _, err := instances.Start(client, instanceID).Extract(); err != nil {
				return diag.FromErr(err)
			}
			startStateConf := &resource.StateChangeConf{
				Target:     []string{InstanceVMStateActive},
				Refresh:    ServerV2StateRefreshFunc(client, instanceID),
				Timeout:    d.Timeout(schema.TimeoutCreate),
				Delay:      10 * time.Second,
				MinTimeout: 3 * time.Second,
			}
			_, err = startStateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for instance (%s) to become active: %s", d.Id(), err)
			}
		case InstanceVMStateStopped:
			if _, err := instances.Stop(client, instanceID).Extract(); err != nil {
				return diag.FromErr(err)
			}
			stopStateConf := &resource.StateChangeConf{
				Target:     []string{InstanceVMStateStopped},
				Refresh:    ServerV2StateRefreshFunc(client, instanceID),
				Timeout:    d.Timeout(schema.TimeoutCreate),
				Delay:      10 * time.Second,
				MinTimeout: 3 * time.Second,
			}
			_, err = stopStateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for instance (%s) to become inactive(stopped): %s", d.Id(), err)
			}
		}
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))
	log.Println("[DEBUG] Finish Instance updating")
	return resourceInstanceRead(ctx, d, m)
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Instance deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	instanceID := d.Id()
	log.Printf("[DEBUG] Instance id = %s", instanceID)

	client, err := CreateClient(provider, d, InstancePoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	var delOpts instances.DeleteOpts
	results, err := instances.Delete(client, instanceID, delOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, InstanceDeleting, func(task tasks.TaskID) (interface{}, error) {
		_, err := instances.Get(client, instanceID).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete instance with ID: %s", instanceID)
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
	log.Printf("[DEBUG] Finish of Instance deleting")
	return diags
}

// ServerV2StateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a gcorecloud instance.
func ServerV2StateRefreshFunc(client *gcorecloud.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(gcorecloud.ErrDefault404); ok {
				return s, "DELETED", nil
			}
			return nil, "", err
		}

		return s, s.VMState, nil
	}
}

func findInstancePort(portID string, ports []instances.InstancePorts) (instances.InstancePorts, error) {
	for _, port := range ports {
		if port.ID == portID {
			return port, nil
		}
	}

	return instances.InstancePorts{}, fmt.Errorf("port not found")
}

func prepareSecurityGroups(ports []instances.InstancePorts) []interface{} {
	sgs := make(map[string]string)
	for _, port := range ports {
		for _, sg := range port.SecurityGroups {
			sgs[sg.ID] = sg.Name
		}
	}

	secGroups := make([]interface{}, 0, len(sgs))
	for sgID, sgName := range sgs {
		s := make(map[string]interface{})
		s["id"] = sgID
		s["name"] = sgName
		secGroups = append(secGroups, s)
	}
	return secGroups
}
