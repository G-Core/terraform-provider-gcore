package gcore

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"slices"
	"sort"
	"strconv"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
	instancesV2 "github.com/G-Core/gcorelabscloud-go/gcore/instance/v2/instances"
	typesV2 "github.com/G-Core/gcorelabscloud-go/gcore/instance/v2/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/securitygroups"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	volumesV2 "github.com/G-Core/gcorelabscloud-go/gcore/volume/v2/volumes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	instanceOperationTimeout = 1200
)

func resourceInstanceV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceV2Create,
		ReadContext:   resourceInstanceV2Read,
		UpdateContext: resourceInstanceV2Update,
		DeleteContext: resourceInstanceDelete,
		Description: `
Gcore Instance offer a flexible, powerful, and scalable solution for hosting applications and services.
Designed to meet a wide range of computing needs, our instances ensure optimal performance, reliability, and security for
your applications.`,
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
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Project ID, only one of project_id or project_name should be set",
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"region_id": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Region ID, only one of region_id or region_name should be set",
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				DiffSuppressFunc: suppressDiffRegionID,
			},
			"project_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Project name, only one of project_id or project_name should be set",
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region name, only one of region_id or region_name should be set",
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
			},
			"flavor_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Flavor ID",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the instance.",
				Computed:    true,
			},
			"name_template": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance name template. You can use forms 'ip_octets', 'two_ip_octets', 'one_ip_octet'",
			},
			"volume": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Description: `
List of volumes for the instance. You can detach the volume from the instance by removing the
volume from the instance resource. You cannot detach the boot volume. You can attach a data volume
by adding the volume resource inside an instance resource.`,
				Set: volumeUniqueID,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the volume",
							Computed:    true,
						},
						"boot_index": {
							Type:        schema.TypeInt,
							Description: "If boot_index==0 volumes can not detached",
							Optional:    true,
						},
						"type_name": {
							Type:        schema.TypeString,
							Description: "Volume type name",
							Computed:    true,
						},
						"image_id": {
							Type:        schema.TypeString,
							Description: "Image ID for the volume",
							Computed:    true,
						},
						"size": {
							Type:        schema.TypeInt,
							Description: "Size of the volume in GiB",
							Computed:    true,
						},
						"volume_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"attachment_tag": {
							Type:        schema.TypeString,
							Description: "Tag for the volume attachment",
							Computed:    true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delete_on_termination": {
							Type:        schema.TypeBool,
							Description: "Delete volume on termination",
							Computed:    true,
						},
					},
				},
			},
			"interface": &schema.Schema{
				Type:     schema.TypeSet,
				Set:      instanceInterfaceUniqueID,
				Required: true,
				Description: `
List of interfaces for the instance. You can detach the interface from the instance by removing the
interface from the instance resource and attach the interface by adding the interface resource
inside an instance resource.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: fmt.Sprintf("Available value is '%s', '%s', '%s', '%s'", types.SubnetInterfaceType, types.AnySubnetInterfaceType, types.ExternalInterfaceType, types.ReservedFixedIpType),
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of interface, should be unique for the instance",
						},
						"order": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Order of attaching interface",
						},
						"ip_family": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP family for the interface, available values are 'dual', 'ipv4' and 'ipv6'",
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
						"existing_fip_id": {
							Type:        schema.TypeString,
							Description: "The id of the existing floating IP that will be attached to the interface",
							Optional:    true,
						},
						"port_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "required if type is  'reserved_fixed_ip'",
							Optional:    true,
						},
						"security_groups": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "list of security group IDs, they will be attached to exact interface",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "IP address for the interface.",
						},
					},
				},
			},
			"keypair_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the keypair to use for the instance",
			},
			"server_group": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the server group to use for the instance",
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Description: `
For Linux instances, 'username' and 'password' are used to create a new user.
When only 'password' is provided, it is set as the password for the default user of the image. 'user_data' is ignored
when 'password' is specified. For Windows instances, 'username' cannot be specified. Use the 'password' field to set
the password for the 'Admin' user on Windows. Use the 'user_data' field to provide a script to create new users
on Windows. The password of the Admin user cannot be updated via 'user_data'`,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: `
For Linux instances, 'username' and 'password' are used to create a new user. For Windows
instances, 'username' cannot be specified. Use 'password' field to set the password for the 'Admin' user on Windows.`,
			},
			"metadata_map": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Create one or more metadata items for the instance",
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
			"configuration": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Parameters for the application template from the marketplace",
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
			"user_data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: `
String in base64 format. For Linux instances, 'user_data' is ignored when 'password' field is provided.
For Windows instances, Admin user password is set by 'password' field and cannot be updated via 'user_data'
`,
			},
			"allow_app_ports": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Description: `If true, application ports will be allowed in the security group for instances created
				from the marketplace application template`,
			},
			"flavor": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "Flavor details, RAM, vCPU, etc.",
				Computed:    true,
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Status of the instance",
				Computed:    true,
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
				Type:        schema.TypeList,
				Description: "List of instance addresses",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"net": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addr": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP address",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the address",
									},
								},
							},
						},
					},
				},
			},
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceInstanceV2Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	if userData, ok := d.GetOk("user_data"); ok {
		createOpts.UserData = userData.(string)
	}

	name := d.Get("name").(string)
	if len(name) > 0 {
		createOpts.Names = []string{name}
	}

	if nameTemplate, ok := d.GetOk("name_template"); ok {
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

	ifs := d.Get("interface").(*schema.Set).List()
	// sort interfaces by 'order' key to attach it in right order
	sort.Sort(instanceInterfaces(ifs))
	if len(ifs) > 0 {
		ifaces, err := extractInstanceInterfacesMapV2(ifs)
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.Interfaces = ifaces
	}

	if metadataRaw, ok := d.GetOk("metadata_map"); ok {
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
	resourceInstanceV2Read(ctx, d, m)

	log.Printf("[DEBUG] Finish Instance creating (%s)", InstanceID)
	return diags
}

func resourceInstanceV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	clientVol, err := CreateClient(provider, d, volumesPoint, versionPointV1)
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
		}
		v["id"] = vol.ID
		v["delete_on_termination"] = vol.DeleteOnTermination

		volume, err := volumes.Get(clientVol, vol.ID).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		v["size"] = volume.Size
		v["type_name"] = volume.VolumeType.String()

		extVolumes = append(extVolumes, v)
	}

	if err := d.Set("volume", schema.NewSet(volumeUniqueID, extVolumes)); err != nil {
		return diag.FromErr(err)
	}

	instancePorts, err := instances.ListPortsAll(client, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}
	ifs, err := instances.ListInterfacesAll(client, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	statesInterface := d.Get("interface").(*schema.Set)
	interfaces, err := extractInstanceInterfaceIntoMapV2(statesInterface.List())
	if err != nil {
		return diag.FromErr(err)
	}

	var cleanInterfaces []interface{}
	for ifOrder, iface := range ifs {
		if len(iface.IPAssignments) == 0 {
			continue
		}

		// we need to match our interfaces with api's interfaces
		// but we don't have any unique values, that's why we use exactly that list of keys
		ifaceName := iface.Name
		if ifaceName == nil {
			log.Printf("[WARN] Interface for instance %s missing name. Using PortID as identifier.", instanceID)
			generatedName := fmt.Sprintf("interface_%s", iface.PortID)
			ifaceName = &generatedName
		}

		for _, assignment := range iface.IPAssignments {
			subnetID := assignment.SubnetID

			// bad idea, but what to do
			var iOpts instances.InterfaceOpts
			var orderedIOpts OrderedInterfaceOpts
			var ok bool

			if orderedIOpts, ok = interfaces[*ifaceName]; ok {
				iOpts = orderedIOpts.InterfaceOpts
			}

			i := make(map[string]interface{})
			if !ok {
				orderedIOpts = OrderedInterfaceOpts{Order: ifOrder}
			} else {
				i["type"] = iOpts.Type.String()
			}

			i["network_id"] = iface.NetworkID
			i["subnet_id"] = subnetID
			i["port_id"] = iface.PortID
			i["name"] = *ifaceName
			i["order"] = orderedIOpts.Order
			if len(iface.FloatingIPDetails) > 0 {
				i["existing_fip_id"] = iface.FloatingIPDetails[0].ID
			}
			i["ip_address"] = assignment.IPAddress.String()

			if port, err := findInstancePort(iface.PortID, instancePorts); err == nil {
				sgs := make([]interface{}, len(port.SecurityGroups))
				for i, sg := range port.SecurityGroups {
					sgs[i] = sg.ID
				}
				i["security_groups"] = schema.NewSet(sgUniqueIDs, sgs)
			}

			cleanInterfaces = append(cleanInterfaces, i)
		}
	}
	if err := d.Set("interface", schema.NewSet(instanceInterfaceUniqueID, cleanInterfaces)); err != nil {
		return diag.FromErr(err)
	}

	metadataMap, metadataReadOnly := PrepareMetadata(instance.MetadataDetailed)

	if err = d.Set("metadata_map", metadataMap); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("metadata_read_only", metadataReadOnly); err != nil {
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

func resourceInstanceV2Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Instance updating")
	instanceID := d.Id()
	log.Printf("[DEBUG] Instance id = %s", instanceID)
	config := m.(*Config)
	provider := config.Provider
	client, err := CreateClient(provider, d, InstancePoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	clientV2, err := CreateClient(provider, d, InstancePoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	clientSg, err := CreateClient(provider, d, securityGroupPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		nameTemplate := d.Get("name_template").(string)
		if len(nameTemplate) == 0 {
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

	if d.HasChange("metadata_map") {
		omd, nmd := d.GetChange("metadata_map")
		if len(omd.(map[string]interface{})) > 0 {
			for k := range omd.(map[string]interface{}) {
				err := instancesV2.MetadataItemDelete(clientV2, instanceID, instancesV2.MetadataItemOpts{Key: k}).Err
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
		instancePorts, err := instances.ListPortsAll(client, instanceID)
		if err != nil {
			return diag.FromErr(err)
		}

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

		ifsOld := ifsOldRaw.(*schema.Set)
		ifsNew := ifsNewRaw.(*schema.Set)

		// we have to create separate sets for old and new interfaces by name, to be able to match
		// interfaces which wasn't changed. We do it because new set doesn't contain portID.
		// port id is needed to reassign security groups
		ifsSetByNameOld := schema.NewSet(instanceInterfaceUniqueIDByName, ifsOld.List())
		ifsSetByNameNew := schema.NewSet(instanceInterfaceUniqueIDByName, ifsNew.List())

		ifsForUpdate := schema.NewSet(instanceInterfaceUniqueIDByName, []interface{}{})

		for _, i := range ifsOld.Difference(ifsNew).List() {
			// if name left the same in new set, we can skip detaching
			if ifsSetByNameNew.Contains(i) {
				ifsForUpdate.Add(i)
				continue
			}

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

		ifsNewSorted := ifsNew.Difference(ifsOld).List()
		sort.Sort(instanceInterfaces(ifsNewSorted))
		for _, i := range ifsNewSorted {
			// if it is completely new interface we need to attach it
			if !ifsSetByNameOld.Contains(i) {
				if err := attachNewInterface(i, client, instanceID); err != nil {
					return diag.FromErr(err)
				}
				continue
			}

			iface := i.(map[string]interface{})

			var portID string
			// try to find port id from old interfaces
			for _, iOld := range ifsForUpdate.List() {
				interfaceOld := iOld.(map[string]interface{})
				if interfaceOld["name"] == iface["name"] {
					portID = interfaceOld["port_id"].(string)
					break
				}
			}

			log.Println("[DEBUG] Reassign security groups")
			port, err := findInstancePort(portID, instancePorts)
			if err != nil {
				log.Println("[DEBUG] Port not found")
				continue
			}

			// detach what should be detached
			sgToDetach := make([]string, 0)
			for _, sg := range port.SecurityGroups {
				if !iface["security_groups"].(*schema.Set).Contains(sg.ID) {
					sgToDetach = append(sgToDetach, sg.Name)
				}
			}
			detachOpts := instances.SecurityGroupOpts{
				PortsSecurityGroupNames: []instances.PortSecurityGroupNames{{
					PortID:             &portID,
					SecurityGroupNames: sgToDetach,
				}},
			}
			if err := instances.UnAssignSecurityGroup(client, instanceID, detachOpts).ExtractErr(); err != nil {
				log.Printf("[WARNING] Cannot detach security groups: %v", err)
			}

			// attach what should be attached
			sgToAttach := make([]string, 0)
			for _, sg := range iface["security_groups"].(*schema.Set).List() {
				if !slices.ContainsFunc(port.SecurityGroups, func(s gcorecloud.ItemIDName) bool {
					return s.ID == sg.(string)
				}) {
					// get the name of the security group
					secGroup, err := securitygroups.Get(clientSg, sg.(string)).Extract()
					if err != nil {
						log.Printf("[WARNING] Cannot get security group %s: %v", sg, err)
						continue
					}
					sgToAttach = append(sgToAttach, secGroup.Name)
				}
			}
			attachOpts := instances.SecurityGroupOpts{
				PortsSecurityGroupNames: []instances.PortSecurityGroupNames{{
					PortID:             &portID,
					SecurityGroupNames: sgToAttach,
				}},
			}
			if err := instances.AssignSecurityGroup(client, instanceID, attachOpts).ExtractErr(); err != nil {
				log.Printf("[WARNING] Cannot attach security groups: %v", err)
			}
		}
	}

	if d.HasChange("volume") {
		vClient, err := CreateClient(provider, d, volumesPoint, versionPointV2)
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
			results, err := volumesV2.Detach(vClient, vid, vOpts).Extract()
			if err != nil {
				return diag.FromErr(err)
			}

			taskID := results.Tasks[0]
			if err := waitInstanceOperation(client, taskID); err != nil {
				return diag.FromErr(err)
			}
		}

		// range over not attached volumes
		for vid, ok := range newVolumes {
			if ok {
				results, err := volumesV2.Attach(vClient, vid, vOpts).Extract()
				if err != nil {
					return diag.FromErr(err)
				}

				taskID := results.Tasks[0]
				if err := waitInstanceOperation(client, taskID); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChange("vm_state") {
		state := d.Get("vm_state").(string)
		opts := instancesV2.ActionOpts{}
		switch state {
		case InstanceVMStateActive:
			opts.Action = typesV2.InstanceActionTypeStart
		case InstanceVMStateStopped:
			opts.Action = typesV2.InstanceActionTypeStop
		}

		results, err := instancesV2.Action(clientV2, instanceID, opts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		taskID := results.Tasks[0]
		if err := waitInstanceOperation(client, taskID); err != nil {
			return diag.FromErr(err)
		}
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))
	log.Println("[DEBUG] Finish Instance updating")
	return resourceInstanceV2Read(ctx, d, m)
}

func instanceInterfaceUniqueID(i interface{}) int {
	e := i.(map[string]interface{})
	h := md5.New()
	securitygroupsRaw := e["security_groups"].(*schema.Set).List()
	var securitygroups string
	for _, sg := range securitygroupsRaw {
		securitygroups += sg.(string)
	}
	io.WriteString(h, e["name"].(string))
	io.WriteString(h, securitygroups)
	return int(binary.BigEndian.Uint64(h.Sum(nil)))
}

func instanceInterfaceUniqueIDByName(i interface{}) int {
	e := i.(map[string]interface{})
	h := md5.New()
	io.WriteString(h, e["name"].(string))
	return int(binary.BigEndian.Uint64(h.Sum(nil)))
}

func sgUniqueIDs(i interface{}) int {
	e := i.(string)
	h := md5.New()
	io.WriteString(h, e)
	return int(binary.BigEndian.Uint64(h.Sum(nil)))
}

func attachNewInterface(i interface{}, client *gcorecloud.ServiceClient, instanceID string) error {
	iface := i.(map[string]interface{})
	iType := types.InterfaceType(iface["type"].(string))
	ifaceName := iface["name"].(string)
	opts := instances.InterfaceInstanceCreateOpts{
		InterfaceOpts: instances.InterfaceOpts{
			Name:     &ifaceName,
			Type:     iType,
			IPFamily: types.IPFamilyType(iface["ip_family"].(string)),
		},
	}

	switch iType {
	case types.SubnetInterfaceType:
		opts.SubnetID = iface["subnet_id"].(string)
	case types.AnySubnetInterfaceType:
		opts.NetworkID = iface["network_id"].(string)
	case types.ReservedFixedIpType:
		opts.PortID = iface["port_id"].(string)
	}

	rawSgsID := iface["security_groups"].(*schema.Set).List()
	sgs := make([]gcorecloud.ItemID, len(rawSgsID))
	for i, sgID := range rawSgsID {
		sgs[i] = gcorecloud.ItemID{ID: sgID.(string)}
	}
	opts.SecurityGroups = sgs

	log.Printf("[DEBUG] attach interface: %+v", opts)
	results, err := instances.AttachInterface(client, instanceID, opts).Extract()
	if err != nil {
		return fmt.Errorf("cannot attach interface: %s. Error: %s", iType, err)
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] attach interface taskID: %s", taskID)
	if err = tasks.WaitForStatus(client, string(taskID), tasks.TaskStateFinished, InstanceCreatingTimeout, true); err != nil {
		return err
	}
	return nil
}

func waitInstanceOperation(client *gcorecloud.ServiceClient, taskID tasks.TaskID) error {
	_, err := tasks.WaitTaskAndReturnResult(client, taskID, true, instanceOperationTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		return nil, nil
	})
	return err
}
