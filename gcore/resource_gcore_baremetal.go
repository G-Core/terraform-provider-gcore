package gcore

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/baremetal/v1/bminstances"
	"github.com/G-Core/gcorelabscloud-go/gcore/floatingip/v1/floatingips"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	BmInstanceDeleting        int = 1200
	BmInstanceCreatingTimeout int = 3600
	BmInstancePoint               = "bminstances"
)

var bmCreateTimeout = time.Second * time.Duration(BmInstanceCreatingTimeout)

func resourceBmInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBmInstanceCreate,
		ReadContext:   resourceBmInstanceRead,
		UpdateContext: resourceBmInstanceUpdate,
		DeleteContext: resourceBmInstanceDelete,
		Description:   "Represent baremetal instance",
		Timeouts: &schema.ResourceTimeout{
			Create: &bmCreateTimeout,
		},
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
				Description:      "Project ID, only one of project_id or project_name should be set",
			},
			"region_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				DiffSuppressFunc: suppressDiffRegionID,
				Description:      "Region ID, only one of region_id or region_name should be set",
			},
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				Description: "Project name, only one of project_id or project_name should be set",
			},
			"region_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				Description: "Region name, only one of region_id or region_name should be set",
			},
			"flavor_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the flavor (type of server configuration). This field is required. Example: 'bm1-hf-medium-4x1nic'",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the baremetal server. If not provided, it will be generated automatically. Example: 'bm-server-01'",
			},
			"name_templates": &schema.Schema{
				Type:          schema.TypeList,
				Optional:      true,
				Deprecated:    "Use name_template instead",
				ConflictsWith: []string{"name_template"},
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "Deprecated. List of baremetal names which will be changed by template",
			},
			"name_template": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name_templates"},
				Description:   "The template used to generate server names. You can use forms 'ip_octets', 'two_ip_octets', 'one_ip_octet'. Example: 'server-${ip_octets}'",
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"image_id",
					"apptemplate_id",
				},
				Description: "The ID of the image to use. The image will be used to provision the bare metal server. Provide either 'image_id' or 'apptemplate_id', but not both",
			},
			"apptemplate_id": {
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"image_id",
					"apptemplate_id",
				},
				Description: "The ID of the application template to use. Provide either 'apptemplate_id' or 'image_id', but not both",
			},
			"interface": &schema.Schema{
				Type: schema.TypeList,
				//Set:      interfaceUniqueID,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: fmt.Sprintf("The type of the network interface. Available value is '%s', '%s', '%s', '%s'", types.SubnetInterfaceType, types.AnySubnetInterfaceType, types.ExternalInterfaceType, types.ReservedFixedIpType),
						},
						"is_parent": {
							Type:        schema.TypeBool,
							Computed:    true,
							Optional:    true,
							Description: "Indicates whether this interface is the parent. If not set will be calculated after creation. Trunk interface always attached first. Can't detach interface if is_parent true. Fields affect only on creation",
						},
						"order": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Order of attaching interface. Trunk (parent) interface always attached first, fields affect only on creation",
						},
						"network_id": {
							Type:        schema.TypeString,
							Description: "The network ID to attach the interface to. Required if type is 'subnet' or 'any_subnet'",
							Optional:    true,
							Computed:    true,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "The subnet ID to attach the interface to. Required if type is 'subnet'",
							Optional:    true,
							Computed:    true,
						},
						"port_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The port ID for reserved fixed IP. Required if type is  'reserved_fixed_ip'",
							Optional:    true,
						},
						// nested map is not supported, in this case, you do not need to use the list for the map
						"fip_source": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The source of floating IP. Can be 'new' or 'existing'",
						},
						"existing_fip_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the existing floating IP that will be attached to the interface",
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "The IP address for the interface",
						},
					},
				},
			},
			"keypair_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the SSH keypair to use for the baremetal",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The password for accessing the baremetal server. This parameter is used to set a password for the 'Admin' user on a Windows instance, a default user or a new user on a Linux instance",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A name of a new user in the Linux instance. It may be passed with a 'password' parameter",
			},
			"metadata": &schema.Schema{
				Type:          schema.TypeList,
				Optional:      true,
				Deprecated:    "Use metadata_map instead",
				ConflictsWith: []string{"metadata_map"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Metadata key",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Metadata value",
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
				Description: "A map of metadata items. Key-value pairs for instance metadata. Example: {'environment': 'production', 'owner': 'user'}",
			},
			"app_config": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Parameters for the application template from the marketplace. This could include parameters required for app setup. Example: {'shadowsocks_method': 'chacha20-ietf-poly1305', 'shadowsocks_password': '123'}",
			},
			"user_data": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User data string in base64 format. This is passed to the instance at launch. For Linux instances, 'user_data' is ignored when 'password' field is provided. For Windows instances, Admin user password is set by 'password' field and cannot be updated via 'user_data'",
			},
			"flavor": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Details about the flavor (server configuration) including RAM, vCPU, etc.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the baremetal server.",
			},
			"vm_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the virtual machine",
			},
			"addresses": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"net": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addr": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The IP address of the interface",
									},
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The type of IP address",
									},
								},
							},
						},
					},
				},
			},
			"last_updated": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The date and time when the baremetal server was last updated",
			},
		},
	}
}

func resourceBmInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start BaremetalInstance creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, BmInstancePoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	ifs := d.Get("interface").([]interface{})
	// sort interfaces by 'is_parent' at first and by 'order' key to attach it in right order
	sort.Sort(instanceInterfaces(ifs))
	newInterface := make([]bminstances.InterfaceOpts, len(ifs))
	for i, iface := range ifs {
		raw := iface.(map[string]interface{})
		newIface := bminstances.InterfaceOpts{
			Type:      types.InterfaceType(raw["type"].(string)),
			NetworkID: raw["network_id"].(string),
			SubnetID:  raw["subnet_id"].(string),
			PortID:    raw["port_id"].(string),
		}

		fipSource := raw["fip_source"].(string)
		fipID := raw["existing_fip_id"].(string)
		if fipSource != "" {
			newIface.FloatingIP = &bminstances.CreateNewInterfaceFloatingIPOpts{
				Source:             types.FloatingIPSource(fipSource),
				ExistingFloatingID: fipID,
			}
		}
		newInterface[i] = newIface
	}

	log.Printf("[DEBUG] Baremetal interfaces: %+v", newInterface)
	opts := bminstances.CreateOpts{
		Flavor:        d.Get("flavor_id").(string),
		ImageID:       d.Get("image_id").(string),
		AppTemplateID: d.Get("apptemplate_id").(string),
		Keypair:       d.Get("keypair_name").(string),
		Password:      d.Get("password").(string),
		Username:      d.Get("username").(string),
		UserData:      d.Get("user_data").(string),
		AppConfig:     d.Get("app_config").(map[string]interface{}),
		Interfaces:    newInterface,
	}

	name := d.Get("name").(string)
	if len(name) > 0 {
		opts.Names = []string{name}
	}

	if nameTemplatesRaw, ok := d.GetOk("name_templates"); ok {
		nameTemplates := nameTemplatesRaw.([]interface{})
		if len(nameTemplates) > 0 {
			NameTemp := make([]string, len(nameTemplates))
			for i, nametemp := range nameTemplates {
				NameTemp[i] = nametemp.(string)
			}
			opts.NameTemplates = NameTemp
		}
	} else if nameTemplate, ok := d.GetOk("name_template"); ok {
		opts.NameTemplates = []string{nameTemplate.(string)}
	}

	if metadata, ok := d.GetOk("metadata"); ok {
		if len(metadata.([]interface{})) > 0 {
			md, err := extractKeyValue(metadata.([]interface{}))
			if err != nil {
				return diag.FromErr(err)
			}
			opts.Metadata = &md
		}
	} else if metadataRaw, ok := d.GetOk("metadata_map"); ok {
		md := extractMetadataMap(metadataRaw.(map[string]interface{}))
		opts.Metadata = &md
	}

	results, err := bminstances.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]

	InstanceID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, BmInstanceCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
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
	log.Printf("[DEBUG] Baremetal Instance id (%s)", InstanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(InstanceID.(string))
	resourceBmInstanceRead(ctx, d, m)

	log.Printf("[DEBUG] Finish Baremetal Instance creating (%s)", InstanceID)
	return diags
}

func resourceBmInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Baremetal Instance reading")
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
		return diag.Errorf("cannot get instance with ID: %s. Error: %s", instanceID, err)
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

	ifs, err := instances.ListInterfacesAll(client, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	interfaces, err := extractInstanceInterfaceIntoMap(d.Get("interface").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	if len(ifs) == 0 {
		return diag.Errorf("interface not found")
	}

	var cleanInterfaces []interface{}
	for _, iface := range ifs {
		for _, assignment := range iface.IPAssignments {
			subnetID := assignment.SubnetID

			// bad idea, but what to do
			var iOpts OrderedInterfaceOpts
			var orderedIOpts OrderedInterfaceOpts
			var ok bool
			// we need to match our interfaces with api's interfaces
			// but with don't have any unique value, that's why we use exactly that list of keys
			for _, k := range []string{subnetID, iface.PortID, iface.NetworkID, types.ExternalInterfaceType.String()} {
				if orderedIOpts, ok = interfaces[k]; ok {
					iOpts = orderedIOpts
					break
				}
			}

			if !ok {
				continue
			}

			i := make(map[string]interface{})

			i["type"] = iOpts.Type.String()
			i["network_id"] = iface.NetworkID
			i["subnet_id"] = subnetID
			i["port_id"] = iface.PortID
			i["is_parent"] = true
			i["order"] = iOpts.Order
			if iOpts.FloatingIP != nil {
				i["fip_source"] = iOpts.FloatingIP.Source.String()
				i["existing_fip_id"] = iOpts.FloatingIP.ExistingFloatingID
			}
			i["ip_address"] = assignment.IPAddress.String()

			cleanInterfaces = append(cleanInterfaces, i)
		}

		for _, iface1 := range iface.SubPorts {
			for _, assignment := range iface1.IPAssignments {
				subnetID := assignment.SubnetID

				// bad idea, but what to do
				var iOpts OrderedInterfaceOpts
				var orderedIOpts OrderedInterfaceOpts
				var ok bool
				// we need to match our interfaces with api's interfaces
				// but with don't have any unique value, that's why we use exactly that list of keys
				for _, k := range []string{subnetID, iface1.PortID, iface1.NetworkID, types.ExternalInterfaceType.String()} {
					if orderedIOpts, ok = interfaces[k]; ok {
						iOpts = orderedIOpts
						break
					}
				}

				if !ok {
					continue
				}

				i := make(map[string]interface{})

				i["type"] = iOpts.Type.String()
				i["network_id"] = iface1.NetworkID
				i["subnet_id"] = subnetID
				i["port_id"] = iface1.PortID
				i["is_parent"] = false
				i["order"] = iOpts.Order
				if iOpts.FloatingIP != nil {
					i["fip_source"] = iOpts.FloatingIP.Source.String()
					i["existing_fip_id"] = iOpts.FloatingIP.ExistingFloatingID
				}
				i["ip_address"] = assignment.IPAddress.String()

				cleanInterfaces = append(cleanInterfaces, i)
			}
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

	fields := []string{"user_data", "app_config"}
	revertState(d, &fields)

	log.Println("[DEBUG] Finish Instance reading")
	return diags
}

func resourceBmInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Baremetal Instance updating")
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
			if isInterfaceContains(iface, ifsNew) {
				log.Println("[DEBUG] Skipped, dont need detach")
				continue
			}

			if iface["is_parent"].(bool) {
				return diag.Errorf("could not detach trunk interface")
			}

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

		currentIfs, err := instances.ListInterfacesAll(client, d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		sort.Sort(instanceInterfaces(ifsNew))
		for _, i := range ifsNew {
			iface := i.(map[string]interface{})
			if isInterfaceContains(iface, ifsOld) {
				log.Println("[DEBUG] Skipped, dont need attach")
				continue
			}
			if isInterfaceAttached(currentIfs, iface) {
				continue
			}

			iType := types.InterfaceType(iface["type"].(string))
			opts := instances.InterfaceOpts{Type: iType}
			switch iType {
			case types.SubnetInterfaceType:
				opts.SubnetID = iface["subnet_id"].(string)
			case types.AnySubnetInterfaceType:
				opts.NetworkID = iface["network_id"].(string)
			case types.ReservedFixedIpType:
				opts.PortID = iface["port_id"].(string)
			}

			log.Printf("[DEBUG] attach interface: %+v", opts)
			results, err := instances.AttachInterface(client, instanceID, opts).Extract()
			if err != nil {
				return diag.Errorf("cannot attach interface: %s. Error: %s", iType, err)
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

		currentIfs, err = instances.ListInterfacesAll(client, d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		desired := map[string]string{} // fipID -> portID
		ifsCfg := d.Get("interface").([]interface{})
		for _, raw := range ifsCfg {
			iface := raw.(map[string]interface{})
			fipID, _ := iface["existing_fip_id"].(string)
			if fipID == "" {
				continue
			}
			var portID string
			switch types.InterfaceType(iface["type"].(string)) {
			case types.ReservedFixedIpType:
				portID = iface["port_id"].(string)
			case types.SubnetInterfaceType:
				for _, ci := range currentIfs {
					for _, a := range ci.IPAssignments {
						if a.SubnetID == iface["subnet_id"].(string) {
							portID = ci.PortID
							break
						}
					}
				}
			case types.AnySubnetInterfaceType:
				for _, ci := range currentIfs {
					if ci.NetworkID == iface["network_id"].(string) {
						portID = ci.PortID
						break
					}
				}
			}
			if portID != "" {
				desired[fipID] = portID
			}
		}

		current := map[string]string{}
		for _, ci := range currentIfs {
			for _, f := range ci.FloatingIPDetails {
				current[f.ID] = ci.PortID
			}
		}

		for fipID, havePort := range current {
			if _, ok := desired[fipID]; !ok && havePort != "" {
				log.Printf("[DEBUG] Unassign floating IP %s (not in desired)", fipID)
				if _, err := floatingips.UnAssign(fipClient, fipID).Extract(); err != nil {
					return diag.Errorf("failed to unassign floating IP %s: %v", fipID, err)
				}
			}
		}

		for fipID, wantPort := range desired {
			if havePort, ok := current[fipID]; !ok || havePort != wantPort {
				log.Printf("[DEBUG] Assign floating IP %s -> port %s", fipID, wantPort)
				if _, err := floatingips.Assign(fipClient, fipID, floatingips.CreateOpts{PortID: wantPort}).Extract(); err != nil {
					return diag.Errorf("failed to assign floating IP %s to port %s: %v", fipID, wantPort, err)
				}
			}
		}
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))
	log.Println("[DEBUG] Finish Instance updating")
	return resourceBmInstanceRead(ctx, d, m)
}

func resourceBmInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Baremetal Instance deleting")
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
	delOpts.DeleteFloatings = true

	results, err := instances.Delete(client, instanceID, delOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, BmInstanceDeleting, func(task tasks.TaskID) (interface{}, error) {
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
