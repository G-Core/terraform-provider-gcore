package gcore

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	ai "github.com/G-Core/gcorelabscloud-go/gcore/ai/v1/ais"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/securitygroups"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata/v1/metadata"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	volumesV2 "github.com/G-Core/gcorelabscloud-go/gcore/volume/v2/volumes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	AIClusterDeletingTimeout int = 1200
	AIClusterCreatingTimeout int = 1200
	AIClusterSuspendTimeout  int = 300

	AIClusterPoint = "ai/clusters"
	TaskPoint      = "tasks"
)

const (
	SuspendedStatus = "SUSPENDED"
	ActiveStatus    = "ACTIVE"
)

const (
	SgStateAssigned   = "assigned"
	SgStateUnassigned = "unassigned"
)

const (
	StatusErrorMessage = "cluster status is not '%s'"
)

func resourceAICluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAIClusterCreate,
		ReadContext:   resourceAIClusterRead,
		UpdateContext: resourceAIClusterUpdate,
		DeleteContext: resourceAIClusterDelete,
		Description:   "Represent instance",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, clusterID, err := ImportStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(clusterID)

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"region_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				DiffSuppressFunc: suppressDiffRegionID,
			},
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Description: "AI Cluster Name",
				Required:    true,
			},
			"cluster_status": {
				Type:        schema.TypeString,
				Description: "AI Cluster status",
				Optional:    true,
				Computed:    true,
			},
			"task_id": {
				Type:        schema.TypeString,
				Description: "Task ID associated with the cluster",
				Computed:    true,
			},
			"task_status": {
				Type:        schema.TypeString,
				Description: "Task status",
				Computed:    true,
			},
			"creator_task_id": {
				Type:        schema.TypeString,
				Description: "Task that created this entity",
				Computed:    true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Description: "Datetime when the cluster was created",
				Computed:    true,
			},
			"flavor": {
				Type:        schema.TypeString,
				Description: "Flavor ID (name)",
				Required:    true,
			},
			"image_id": {
				Type:        schema.TypeString,
				Description: "Image ID",
				Required:    true,
			},
			"image_name": {
				Type:        schema.TypeString,
				Description: "Image name",
				Computed:    true,
			},
			"volume": {
				Type:        schema.TypeSet,
				Description: "List of volumes attached to the cluster",
				Optional:    true,
				Set:         aiVolumeHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_id": {
							Type:        schema.TypeString,
							Description: "Volume ID",
							Optional:    true,
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Volume name",
							Optional:    true,
							Computed:    true,
						},
						"source": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							ValidateDiagFunc: func(val interface{}, key cty.Path) diag.Diagnostics {
								v := val.(string)
								if types.VolumeSource(v) == types.ExistingVolume || types.VolumeSource(v) == types.Image {
									return diag.Diagnostics{}
								}
								return diag.Errorf("wrong source type %s, now available values are '%s', '%s'", v, types.ExistingVolume, types.Image)
							},
							Description: "Currently available only value",
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Volume status",
							Computed:    true,
						},
						"size": {
							Type:        schema.TypeInt,
							Description: "Volume size, GiB",
							Computed:    true,
							Optional:    true,
						},
						"created_at": {
							Type:        schema.TypeString,
							Description: "Datetime when the volume was created",
							Computed:    true,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Description: "Datetime when the volume was last updated",
							Computed:    true,
						},
						"volume_type": {
							Type:        schema.TypeString,
							Description: "Volume type",
							Computed:    true,
							Optional:    true,
						},
						"image_id": {
							Type:        schema.TypeString,
							Description: "Volume ID. Mandatory if volume is pre-existing volume",
							Optional:    true,
						},
						"attachments": {
							Type:        schema.TypeSet,
							Description: "Attachment list",
							Computed:    true,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_id": {
										Type:        schema.TypeString,
										Description: "Instance ID",
										Computed:    true,
									},
									"instance_name": {
										Type:        schema.TypeString,
										Description: "Instance name (if attached and server name is known)",
										Computed:    true,
									},
									"attachment_id": {
										Type:        schema.TypeString,
										Description: "ID of attachment object",
										Computed:    true,
									},
									"volume_id": {
										Type:        schema.TypeString,
										Description: "Volume ID",
										Computed:    true,
									},
									"device": {
										Type:        schema.TypeString,
										Description: "Block device name in guest",
										Computed:    true,
									},
									"attached_at": {
										Type:        schema.TypeString,
										Description: "Attachment creation datetime",
										Computed:    true,
									},
								},
							},
						},
						"volume_image_metadata": {
							Type:        schema.TypeMap,
							Description: "Image information for volumes that were created from image",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"creator_task_id": {
							Type:        schema.TypeString,
							Description: "Task that created this entity",
							Computed:    true,
						},
					},
				},
			},
			"security_group": {
				Type:        schema.TypeSet,
				Description: "Security groups attached to the cluster",
				Optional:    true,
				Set:         aiSgHashID,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Security group ID",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"interface": {
				Type:        schema.TypeList,
				Description: "Networks managed by user and associated with the cluster",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Description: "Network type",
							Optional:    true,
							ValidateDiagFunc: func(val interface{}, key cty.Path) diag.Diagnostics {
								v := val.(string)
								if types.InterfaceType(v) == types.ExternalInterfaceType || types.InterfaceType(v) == types.SubnetInterfaceType {
									return diag.Diagnostics{}
								}
								return diag.Errorf("wrong source type %s, now available values are '%s', '%s'", v, types.ExternalInterfaceType, types.SubnetInterfaceType)
							},
						},
						"network_id": {
							Type:        schema.TypeString,
							Description: "Network ID",
							Optional:    true,
							Computed:    true,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "Port is assigned to IP address from the subnet",
							Optional:    true,
						},
						"port_id": {
							Type:        schema.TypeString,
							Description: "Network ID the subnet belongs to. Port will be plugged in this network",
							Optional:    true,
						},
					},
				},
			},
			"keypair_name": {
				Type:        schema.TypeString,
				Description: "Ssh keypair name",
				Optional:    true,
			},
			"user_data": {
				Type:        schema.TypeString,
				Description: "String in base64 format. Must not be passed together with 'username' or 'password'. Examples of the user_data: https://cloudinit.readthedocs.io/en/latest/topics/examples.html",
				Optional:    true,
				Computed:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "A name of a new user in the Linux instance. It may be passed with a 'password' parameter",
				Optional:    true,
				Computed:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "A password for baremetal instance. This parameter is used to set a password for the Admin user on a Windows instance, a default user or a new user on a Linux instance",
				Optional:    true,
				Computed:    true,
			},
			"cluster_metadata": {
				Type:        schema.TypeMap,
				Description: "Cluster metadata (simple key-value pairs)",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"poplar_servers": {
				Type:        schema.TypeList,
				Description: "Poplar servers",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Description: "Instance ID",
							Computed:    true,
						},
						"flavor_id": {
							Type:        schema.TypeString,
							Description: "",
							Computed:    true,
						},
						"flavor": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"metadata": {
							Type:        schema.TypeMap,
							Description: "VM metadata",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"instance_name": {
							Type:        schema.TypeString,
							Description: "Instance name",
							Computed:    true,
						},
						"instance_description": {
							Type:        schema.TypeString,
							Description: "Instance description",
							Computed:    true,
						},
						"instance_created": {
							Type:        schema.TypeString,
							Description: "Datetime when instance was created",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "VM status",
							Computed:    true,
						},
						"vm_state": {
							Type:        schema.TypeString,
							Description: "Virtual machine state (active)",
							Computed:    true,
						},
						"task_state": {
							Type:        schema.TypeString,
							Description: "Task state",
							Computed:    true,
						},
						"volumes": {
							Type:        schema.TypeSet,
							Description: "List of volumes",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Description: "Volume ID",
										Computed:    true,
									},
									"delete_on_termination": {
										Type:        schema.TypeBool,
										Description: "Whether the volume is deleted together with the VM",
										Computed:    true,
									},
								},
							},
						},
						"addresses": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_name": {
										Type:        schema.TypeString,
										Description: "Interface network",
										Computed:    true,
									},
									"address": {
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
						"security_groups": {
							Type:        schema.TypeSet,
							Description: "Server security group",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "Security group name",
										Computed:    true,
									},
								},
							},
						},
						"creator_task_id": {
							Type:        schema.TypeString,
							Description: "Task that created this entity",
							Computed:    true,
						},
						"task_id": {
							Type:        schema.TypeString,
							Description: "Active task. If None, action has been performed immediately in the request itself.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func extractAIClusterInterfacesMap(interfaces []interface{}) ([]instances.InterfaceInstanceCreateOpts, error) {
	Interfaces := make([]instances.InterfaceInstanceCreateOpts, len(interfaces))
	for index, iface := range interfaces {
		ifaceMap := iface.(map[string]interface{})
		var IfaceOpts instances.InterfaceOpts
		switch {
		case ifaceMap["type"].(string) == string(types.ExternalInterfaceType):
			{
				IfaceOpts.Type = types.ExternalInterfaceType
			}
		case ifaceMap["type"].(string) == string(types.AnySubnetInterfaceType):
			{
				IfaceOpts.Type = types.AnySubnetInterfaceType
				IfaceOpts.NetworkID = ifaceMap["network_id"].(string)
			}
		case ifaceMap["type"].(string) == string(types.SubnetInterfaceType):
			{
				IfaceOpts.Type = types.SubnetInterfaceType
				IfaceOpts.NetworkID = ifaceMap["network_id"].(string)
				IfaceOpts.SubnetID = ifaceMap["subnet_id"].(string)
			}
		case ifaceMap["type"].(string) == string(types.ReservedFixedIpType):
			{
				IfaceOpts.Type = types.ReservedFixedIpType
				IfaceOpts.SubnetID = ifaceMap["port_id"].(string)
			}
		}
		Interfaces[index] = instances.InterfaceInstanceCreateOpts{
			InterfaceOpts: IfaceOpts,
		}
	}
	return Interfaces, nil
}

func checkAIClusterStatus(client *gcorecloud.ServiceClient, clusterID string, desiredStatus string) error {
	cluster, err := ai.Get(client, clusterID).Extract()
	if err != nil {
		return err
	}
	if cluster.ClusterStatus != desiredStatus {
		return fmt.Errorf(StatusErrorMessage, desiredStatus)
	}
	return nil
}

func validateCreateOpts(createOpts *ai.CreateOpts) error {
	if isBmFlavor(createOpts.Flavor) && len(createOpts.Volumes) > 0 {
		return errors.New("volumes are not supported for baremetal poplar servers")
	}
	if !isBmFlavor(createOpts.Flavor) && len(createOpts.Volumes) == 0 {
		return errors.New("at least one image volume is required for vm poplar cluster")
	}
	if !isBmFlavor(createOpts.Flavor) && len(createOpts.Volumes) > 0 && createOpts.Volumes[0].Source != types.Image {
		return errors.New("the first volume must be image volume for vm poplar cluster")
	}
	var imageSourceCount int
	for _, volume := range createOpts.Volumes {
		if volume.Source == types.Image {
			imageSourceCount++
			if volume.ImageID != createOpts.ImageID {
				return fmt.Errorf("cluster image '%s' is not equal boot bolume image '%s'", createOpts.ImageID, volume.ImageID)
			}
		}
	}
	if !isBmFlavor(createOpts.Flavor) && imageSourceCount > 1 {
		return errors.New("only one image volume is allowed")
	}
	if len(createOpts.Interfaces) == 0 {
		return errors.New("At least one interface is required")
	}
	var extInterfaceCounter int
	for _, iface := range createOpts.Interfaces {
		if iface.InterfaceOpts.Type == types.ExternalInterfaceType {
			extInterfaceCounter++
		}
	}
	if extInterfaceCounter > 1 {
		return errors.New("only one external interface is allowed")
	}

	return nil
}

func validateAttachVolumes(volumes []interface{}) error {
	for _, volume := range volumes {
		if volume.(map[string]interface{})["type"].(string) != string(types.ExistingVolume) {
			return fmt.Errorf("Only '%s' volume type is allowed", types.ExistingVolume)
		}
	}
	return nil
}

func resourceAIClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start AI cluster creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, AIClusterPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := ai.CreateOpts{}
	securityGroupList := d.Get("security_group").(*schema.Set).List()
	securityGroupIDs := make([]gcorecloud.ItemID, len(securityGroupList))
	for sgIndex, sgID := range securityGroupList {
		securityGroupIDs[sgIndex] = gcorecloud.ItemID{ID: sgID.(map[string]interface{})["id"].(string)}
	}
	createOpts.SecurityGroups = securityGroupIDs
	createOpts.Name = d.Get("cluster_name").(string)
	createOpts.Flavor = d.Get("flavor").(string)
	createOpts.Password = d.Get("password").(string)
	createOpts.Username = d.Get("username").(string)
	createOpts.Keypair = d.Get("keypair_name").(string)
	createOpts.ImageID = d.Get("image_id").(string)

	if userData, ok := d.GetOk("userdata"); ok {
		createOpts.UserData = userData.(string)
	}

	currentVols := d.Get("volume").(*schema.Set).List()
	if len(currentVols) > 0 {
		vs, err := extractVolumesMap(currentVols)
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.Volumes = vs
	}

	ifs := d.Get("interface").([]interface{})
	if len(ifs) > 0 {
		ifaces, err := extractAIClusterInterfacesMap(ifs)
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.Interfaces = ifaces
	}

	if metadata, ok := d.GetOk("cluster_metadata"); ok {
		createOpts.Metadata = make(map[string]string)
		for metaKey, metaValue := range metadata.(map[string]interface{}) {
			createOpts.Metadata[metaKey] = metaValue.(string)
		}
	}
	err = validateCreateOpts(&createOpts)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] AI cluster create options: %+v", createOpts)
	results, err := ai.Create(client, createOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	taskClient, err := CreateClient(provider, d, TaskPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)

	clusterID, err := tasks.WaitTaskAndReturnResult(taskClient, taskID, true, AIClusterCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(taskClient, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		clusterID, err := ai.ExtractAIClusterIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve AI cluster ID from task info: %w", err)
		}
		return clusterID, nil
	},
	)
	log.Printf("[DEBUG] AI cluster id (%s)", clusterID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(clusterID.(string))
	resourceAIClusterRead(ctx, d, m)

	log.Printf("[DEBUG] Finish AI cluster creating (%s)", clusterID)
	return diags
}

func resourceAIClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start AI Cluster reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	clusterID := d.Id()
	log.Printf("[DEBUG] AI Cluster id = %s", clusterID)

	client, err := CreateClient(provider, d, AIClusterPoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	cluster, err := ai.Get(client, clusterID).Extract()
	if err != nil {
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			log.Printf("[WARN] Removing AI cluster %s because resource doesn't exist anymore", d.Id())
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}
	err = setAIClusterResourcerData(d, provider, cluster)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] Finish AI cluster reading")
	return diags
}

func attachInterfaceID(attachInterface ai.AttachInterfaceOpts) string {
	switch {
	case attachInterface.Type == types.ExternalInterfaceType:
		{
			return fmt.Sprintf("%v", types.ExternalInterfaceType)
		}
	case attachInterface.Type == types.AnySubnetInterfaceType:
		{
			return fmt.Sprintf("%v/%v", types.AnySubnetInterfaceType, attachInterface.NetworkID)
		}
	case attachInterface.Type == types.SubnetInterfaceType:
		{
			return fmt.Sprintf("%v/%v", types.SubnetInterfaceType, attachInterface.SubnetID)
		}
	case attachInterface.Type == types.ReservedFixedIpType:
		{
			return fmt.Sprintf("%v/%v", types.ReservedFixedIpType, attachInterface.PortID)
		}
	default:
		return "unknown"
	}
}

func areInterfacesUnique(attachInterfaces []ai.AttachInterfaceOpts) bool {
	attachInterfaceMap := make(map[string]ai.AttachInterfaceOpts)
	for _, ifaceOpts := range attachInterfaces {
		attachInterfaceMap[attachInterfaceID(ifaceOpts)] = ifaceOpts
	}
	return len(attachInterfaceMap) == len(attachInterfaces)

}

func contains(attachInterfaces []ai.AttachInterfaceOpts, newAttachInterface ai.AttachInterfaceOpts) bool {
	attachInterfaceMap := make(map[string]ai.AttachInterfaceOpts)
	for _, ifaceOpts := range attachInterfaces {
		attachInterfaceMap[attachInterfaceID(ifaceOpts)] = ifaceOpts
	}
	_, found := attachInterfaceMap[attachInterfaceID(newAttachInterface)]
	return found
}

func map2AttachInterfaceOpts(interfaces []interface{}) []ai.AttachInterfaceOpts {
	result := make([]ai.AttachInterfaceOpts, len(interfaces))
	for index, iface := range interfaces {
		ifaceMap := iface.(map[string]interface{})
		result[index] = ai.AttachInterfaceOpts{
			Type:      types.InterfaceType(ifaceMap["type"].(string)),
			NetworkID: ifaceMap["network_id"].(string),
			SubnetID:  ifaceMap["subnet_id"].(string),
			PortID:    ifaceMap["port_id"].(string),
		}
	}
	return result
}

func AICluserSgRefreshedFunc(client *gcorecloud.ServiceClient, clusterID string, sgID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		aiPorts, err := ai.ListPortsAll(client, clusterID)
		if err != nil {
			return aiPorts, "", err
		}
		portSgSet := make(map[string]bool)
		for _, sgItem := range aiPorts[0].SecurityGroups {
			portSgSet[sgItem.ID] = true
		}
		state := SgStateUnassigned
		if _, found := portSgSet[sgID]; found {
			state = SgStateAssigned
		}
		return aiPorts, state, nil
	}
}

func getDetachOptions(instanceInterfaces []instances.Interface, detachIface ai.AttachInterfaceOpts) (*ai.DetachInterfaceOpts, error) {
	allDetachMap := make(map[string]ai.DetachInterfaceOpts, len(instanceInterfaces))
	for _, instanceIface := range instanceInterfaces {
		if detachIface.Type == types.ExternalInterfaceType && instanceIface.NetworkDetails.External {
			return &ai.DetachInterfaceOpts{
				PortID:    instanceIface.PortID,
				IpAddress: instanceIface.IPAssignments[0].IPAddress.String(),
			}, nil
		}
		for _, ipAssignment := range instanceIface.IPAssignments {
			allDetachMap[ipAssignment.SubnetID] = ai.DetachInterfaceOpts{
				PortID:    instanceIface.PortID,
				IpAddress: ipAssignment.IPAddress.String(),
			}
		}
	}
	if detachOpts, found := allDetachMap[detachIface.SubnetID]; found {
		return &detachOpts, nil
	} else {
		return nil, fmt.Errorf("couldn't found detach options for interface: %v", detachIface)
	}
}

var IsResize bool = false

func resourceAIClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start AI cluster updating")
	config := m.(*Config)
	provider := config.Provider
	clientV1, err := CreateClient(provider, d, AIClusterPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	clientV2, err := CreateClient(provider, d, AIClusterPoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}
	taskClient, err := CreateClient(provider, d, TaskPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	clusterID := d.Id()

	if d.HasChange("cluster_status") {
		newStatus := strings.ToTitle(d.Get("cluster_status").(string))
		switch newStatus {
		case SuspendedStatus:
			results, err := ai.Suspend(clientV1, clusterID).Extract()
			if err != nil {
				return diag.FromErr(err)
			}
			taskID := results.Tasks[0]
			log.Printf("[DEBUG] Waiting ai suspend task to finish. Task id (%s)", taskID)
			_, err = tasks.WaitTaskAndReturnResult(taskClient, taskID, true, AIClusterSuspendTimeout, func(task tasks.TaskID) (interface{}, error) {
				return nil, nil
			})
			if err != nil {
				return diag.FromErr(err)
			}
			err = checkAIClusterStatus(clientV2, clusterID, SuspendedStatus)
			if err != nil {
				return diag.FromErr(err)
			}
		case ActiveStatus:
			results, err := ai.Resume(clientV1, clusterID).Extract()
			if err != nil {
				return diag.FromErr(err)
			}
			taskID := results.Tasks[0]
			log.Printf("[DEBUG] Waiting ai resume task to finish. Task id (%s)", taskID)
			_, err = tasks.WaitTaskAndReturnResult(taskClient, taskID, true, AIClusterSuspendTimeout, func(task tasks.TaskID) (interface{}, error) {
				return nil, nil
			})
			if err != nil {
				return diag.FromErr(err)
			}
			err = checkAIClusterStatus(clientV2, clusterID, ActiveStatus)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// Make resize
	if d.HasChanges("flavor", "image_id", "keypair_name", "user_data", "username", "password") || (d.HasChanges("interface") && isBmFlavor(d.Get("flavor").(string))) {
		IsResize = true
		_, newSGs := d.GetChange("security_group")
		securityGroupList := newSGs.(*schema.Set).List()
		securityGroupIDs := make([]gcorecloud.ItemID, len(securityGroupList))
		for sgIndex, sgID := range securityGroupList {
			securityGroupIDs[sgIndex] = gcorecloud.ItemID{ID: sgID.(map[string]interface{})["id"].(string)}
		}
		_, flavor := d.GetChange("flavor")
		_, image_id := d.GetChange("image_id")
		_, keypairName := d.GetChange("keypair_name")
		_, userData := d.GetChange("user_data")
		_, username := d.GetChange("username")
		_, password := d.GetChange("password")

		resizeOpts := ai.ResizeAIClusterOpts{
			Flavor:         flavor.(string),
			ImageID:        image_id.(string),
			Interfaces:     []instances.InterfaceInstanceCreateOpts{},
			Volumes:        []instances.CreateVolumeOpts{},
			SecurityGroups: securityGroupIDs,
			Keypair:        keypairName.(string),
			Password:       password.(string),
			Username:       username.(string),
			UserData:       userData.(string),
			Metadata:       map[string]string{},
		}
		_, newVolumes := d.GetChange("volume")

		volumeList := newVolumes.([]interface{})
		if len(volumeList) > 0 {
			vs, err := extractVolumesMap(volumeList)
			if err != nil {
				return diag.FromErr(err)
			}
			resizeOpts.Volumes = vs
		}

		_, newIface := d.GetChange("interface")
		interfaceList := newIface.([]interface{})
		if len(interfaceList) > 0 {
			ifaces, err := extractAIClusterInterfacesMap(interfaceList)
			if err != nil {
				return diag.FromErr(err)
			}
			resizeOpts.Interfaces = ifaces
		}
		_, newMetadata := d.GetChange("cluster_metadata")

		for metaKey, metaValue := range newMetadata.(map[string]interface{}) {
			resizeOpts.Metadata[metaKey] = metaValue.(string)
		}

		log.Printf("[DEBUG] AI cluster resize options: %+v", resizeOpts)
		results, err := ai.Resize(clientV1, clusterID, resizeOpts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		taskClient, err := CreateClient(provider, d, TaskPoint, versionPointV1)
		if err != nil {
			return diag.FromErr(err)
		}
		taskID := results.Tasks[0]
		log.Printf("[DEBUG] resize ai cluster task id (%s)", taskID)

		_, err = tasks.WaitTaskAndReturnResult(taskClient, taskID, true, AIClusterCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
			return nil, nil
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("volume") && !IsResize {
		oldVolumes, newVolumes := d.GetChange("volume")
		oldVolumeList := extractInstanceVolumesMap(oldVolumes.(*schema.Set).List())
		newVolumeList := extractInstanceVolumesMap(newVolumes.(*schema.Set).List())
		if isBmFlavor(d.Get("flavor").(string)) && len(newVolumeList) > 0 {
			return diag.FromErr(errors.New("baremetal servers don't support external voluems"))
		}
		poplarInstances := d.Get("poplar_servers").([]interface{})
		if len(poplarInstances) > 1 {
			return diag.FromErr(errors.New("only one vm poplar clusters are supported"))
		}
		instanceID := poplarInstances[0].(map[string]interface{})["instance_id"].(string)
		vClient, err := CreateClient(provider, d, volumesPoint, versionPointV2)
		if err != nil {
			return diag.FromErr(err)
		}
		vOpts := volumes.InstanceOperationOpts{InstanceID: instanceID}
		for vid := range oldVolumeList {
			if isAttached := newVolumeList[vid]; isAttached {
				// mark as already attached
				newVolumeList[vid] = false
				continue
			}
			results, err := volumesV2.Detach(vClient, vid, vOpts).Extract()
			if err != nil {
				return diag.FromErr(err)
			}

			taskID := results.Tasks[0]
			if err := waitInstanceOperation(clientV1, taskID); err != nil {
				return diag.FromErr(err)
			}
		}

		// range over not attached volumes
		for vid, ok := range newVolumeList {
			if ok {
				results, err := volumesV2.Attach(vClient, vid, vOpts).Extract()
				if err != nil {
					return diag.FromErr(err)
				}

				taskID := results.Tasks[0]
				if err := waitInstanceOperation(clientV1, taskID); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChanges("interface") && !IsResize {
		oldIfaces, newIfaces := d.GetChange("interface")
		newAttachInterfaces := map2AttachInterfaceOpts(newIfaces.([]interface{}))
		if len(newAttachInterfaces) == 0 {
			return diag.FromErr(errors.New("no interfaces is configured, at least one is required"))
		}
		oldAttachInterfaces := map2AttachInterfaceOpts(oldIfaces.([]interface{}))
		if !areInterfacesUnique(newAttachInterfaces) {
			return diag.FromErr(errors.New("ai cluster don't support attach equal interfaces"))
		}
		for _, newIface := range newAttachInterfaces {
			if !contains(oldAttachInterfaces, newIface) {
				poplarInstances := d.Get("poplar_servers").([]interface{})
				for _, instance := range poplarInstances {
					instanceID := instance.(map[string]interface{})["instance_id"].(string)
					results, err := ai.AttachAIInstanceInterface(clientV1, instanceID, newIface).Extract()
					if err != nil {
						return diag.FromErr(err)
					}
					taskID := results.Tasks[0]
					log.Printf("[DEBUG] Waiting attach interface to ai instance task to finish. Task id (%s)", taskID)
					_, err = tasks.WaitTaskAndReturnResult(taskClient, taskID, true, AIClusterSuspendTimeout, func(task tasks.TaskID) (interface{}, error) {
						return nil, nil
					})
					if err != nil {
						return diag.FromErr(err)
					}
					log.Printf("[DEBUG] interface has been attached to ai instance. AI instance id(%s)", instanceID)
				}
			}
		}
		instanceClient, err := CreateClient(provider, d, InstancePoint, versionPointV1)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, oldIface := range oldAttachInterfaces {
			if !contains(newAttachInterfaces, oldIface) {
				poplarInstances := d.Get("poplar_servers").([]interface{})
				for _, instance := range poplarInstances {
					instanceID := instance.(map[string]interface{})["instance_id"].(string)
					interfaceList, err := instances.ListInterfacesAll(instanceClient, instanceID)
					if err != nil {
						return diag.FromErr(err)
					}
					detachOpts, err := getDetachOptions(interfaceList, oldIface)
					if err != nil {
						return diag.FromErr(err)
					}

					results, err := ai.DetachAIInstanceInterface(clientV1, instanceID, detachOpts).Extract()
					if err != nil {
						return diag.FromErr(err)
					}
					taskID := results.Tasks[0]
					log.Printf("[DEBUG] Waiting attach interface to ai instance task to finish. Task id (%s)", taskID)
					_, err = tasks.WaitTaskAndReturnResult(taskClient, taskID, true, AIClusterSuspendTimeout, func(task tasks.TaskID) (interface{}, error) {
						return nil, nil
					})
					if err != nil {
						return diag.FromErr(err)
					}
					log.Printf("[DEBUG] interface %+v has been detached from ai instance. AI instance id(%s)", detachOpts, instanceID)
				}
			}
		}
	}

	if d.HasChange("security_group") && !IsResize {
		sgClient, err := CreateClient(provider, d, securityGroupPoint, versionPointV1)
		if err != nil {
			return diag.FromErr(err)
		}
		oldSecGroups, newSecGroups := d.GetChange("security_group")
		newSecGroupSet := newSecGroups.(*schema.Set)
		oldSecGroupSet := oldSecGroups.(*schema.Set)
		for _, addSG := range newSecGroupSet.Difference(oldSecGroupSet).List() {
			sgID := addSG.(map[string]interface{})["id"].(string)
			result, err := securitygroups.Get(sgClient, sgID).Extract()
			if err != nil {
				return diag.FromErr(err)
			}
			opts := instances.SecurityGroupOpts{Name: result.Name}
			log.Printf("[DEBUG] assing to ai cluster sg ID %v, name %v", sgID, result.Name)
			err = ai.AssignSecurityGroup(clientV1, clusterID, opts).ExtractErr()
			if err != nil {
				return diag.FromErr(fmt.Errorf("assign security group ID %v: %w", sgID, err))
			}
			stopWaitConf := retry.StateChangeConf{
				Target:     []string{SgStateAssigned},
				Refresh:    AICluserSgRefreshedFunc(clientV1, clusterID, sgID),
				Timeout:    3 * time.Minute,
				Delay:      10 * time.Second,
				MinTimeout: 5 * time.Second,
			}
			_, err = stopWaitConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for sg (%s) assigned): %s", sgID, err)
			}
		}
		for _, delSG := range oldSecGroupSet.Difference(newSecGroupSet).List() {
			sgID := delSG.(map[string]interface{})["id"].(string)
			result, err := securitygroups.Get(sgClient, sgID).Extract()
			if err != nil {
				return diag.FromErr(err)
			}
			opts := instances.SecurityGroupOpts{Name: result.Name}
			log.Printf("[DEBUG] unassing from ai cluster sg ID %v, name %v", sgID, result.Name)
			err = ai.UnAssignSecurityGroup(clientV1, clusterID, opts).ExtractErr()
			if err != nil {
				return diag.FromErr(fmt.Errorf("unassign security group %v: %w", sgID, err))
			}
			stopWaitConf := retry.StateChangeConf{
				Target:     []string{SgStateUnassigned},
				Refresh:    AICluserSgRefreshedFunc(clientV1, clusterID, sgID),
				Timeout:    3 * time.Minute,
				Delay:      10 * time.Second,
				MinTimeout: 5 * time.Second,
			}
			_, err = stopWaitConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for sg (%s) unassigned): %s", sgID, err)
			}
		}
	}

	if d.HasChange("cluster_metadata") && !IsResize {
		_, newMeta := d.GetChange("cluster_metadata")
		meta := make(map[string]string, len(newMeta.(map[string]interface{})))
		for metaKey, metaVal := range newMeta.(map[string]interface{}) {
			meta[metaKey] = metaVal.(string)
		}
		err = metadata.MetadataReplace(clientV2, clusterID, meta).ExtractErr()
		if err != nil {
			return diag.FromErr(fmt.Errorf("update ai cluser metadata id %v: %w", clusterID, err))
		}
	}

	log.Println("[DEBUG] Finish AI cluster updating")
	return resourceAIClusterRead(ctx, d, m)
}

func resourceAIClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start AI cluster deletion")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	clusterID := d.Id()
	log.Printf("[DEBUG] AI cluster ID = %s", clusterID)

	client, err := CreateClient(provider, d, AIClusterPoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	var delOpts ai.DeleteOpts
	results, err := ai.Delete(client, clusterID, delOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	taskID := results.Tasks[0]
	taskClient, err := CreateClient(provider, d, TaskPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Delete AI cluster task ID (%s)", taskID)

	_, err = tasks.WaitTaskAndReturnResult(taskClient, taskID, true, AIClusterDeletingTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := ai.Get(client, clusterID).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete AI cluster with ID: %s", clusterID)
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
	log.Printf("[DEBUG] Finish of AI cluster deletion")
	return diags
}

// ServerV2StateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a gcorecloud instance.
func PoplarServerV2StateRefreshFunc(client *gcorecloud.ServiceClient, instanceID string) resource.StateRefreshFunc {
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

func prepareAIClusterSecurityGroups(ports []instances.InstancePorts) []interface{} {
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
