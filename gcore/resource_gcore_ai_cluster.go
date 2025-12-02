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
	AIClusterDeletingTimeout int = 7200
	AIClusterCreatingTimeout int = 7200
	AIClusterSuspendTimeout  int = 300

	AIClusterPoint    = "ai/clusters"
	AIClusterGPUPoint = "ai/clusters/gpu"
	TaskPoint         = "tasks"
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
		Description:   "Represent GPU cluster",
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
				Description:      "Project ID, only one of project_id or project_name should be set",
			},
			"region_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				DiffSuppressFunc: suppressDiffRegionID,
				Description:      "Region ID, only one of region_id or region_name should be set",
			},
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				Description: "Project name, only one of project_id or project_name should be set",
			},
			"region_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				Description: "Region name, only one of region_id or region_name should be set",
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Description: "GPU Cluster Name",
				Required:    true,
			},
			"instances_count": {
				Type:        schema.TypeInt,
				Description: "Number of servers in the GPU cluster",
				Optional:    true,
			},
			"cluster_status": {
				Type:        schema.TypeString,
				Description: "GPU Cluster status",
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
				Description: "Datetime when the GPU cluster was created",
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
							Required:    true,
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
				Type:        schema.TypeSet,
				Description: "Networks managed by user and associated with the cluster",
				Required:    true,
				Set:         aiInterfaceHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Description: "Network type",
							Required:    true,
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
							Description: "Network ID the subnet belongs to. Port will be plugged in this network",
							Optional:    true,
							Computed:    true,
						},
						"port_id": {
							Type:        schema.TypeString,
							Description: "Port is assigned to IP address from the subnet",
							Computed:    true,
						},
					},
				},
			},
			"keypair_name": {
				Type:        schema.TypeString,
				Description: "The name of the SSH keypair to use for the GPU servers",
				Optional:    true,
			},
			"user_data": {
				Type:        schema.TypeString,
				Description: "User data string in base64 format. This is passed to the instance at launch. For Linux instances, 'user_data' is ignored when 'password' field is provided. For Windows instances, Admin user password is set by 'password' field and cannot be updated via 'user_data'",
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
				Description: "A password for servers in GPU cluster. This parameter is used to set a password for the Admin user on a Windows instance, a default user or a new user on a Linux instance",
				Optional:    true,
				Computed:    true,
			},
			"cluster_metadata": {
				Type:        schema.TypeMap,
				Description: "A map of metadata items. Key-value pairs for GPU cluster metadata. Example: {'environment': 'production', 'owner': 'user'}",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"poplar_servers": {
				Type:        schema.TypeList,
				Description: "GPU cluster servers list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Description: "GPU Server ID",
							Computed:    true,
						},
						"flavor_id": {
							Type:        schema.TypeString,
							Description: "Flavor id",
							Computed:    true,
						},
						"flavor": {
							Type:        schema.TypeMap,
							Description: "GPU flavor map",
							Computed:    true,
						},
						"metadata": {
							Type:        schema.TypeMap,
							Description: "GPU server metadata",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"instance_name": {
							Type:        schema.TypeString,
							Description: "GPU server name",
							Computed:    true,
						},
						"instance_description": {
							Type:        schema.TypeString,
							Description: "GPU server description",
							Computed:    true,
						},
						"instance_created": {
							Type:        schema.TypeString,
							Description: "Datetime when GPU server was created",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "GPU server status",
							Computed:    true,
						},
						"vm_state": {
							Type:        schema.TypeString,
							Description: "GPU server state (active)",
							Computed:    true,
						},
						"task_state": {
							Type:        schema.TypeString,
							Description: "Task state",
							Computed:    true,
						},
						"volumes": {
							Type:        schema.TypeSet,
							Description: "List of volumes (for virtual GPU clusters)",
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
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "List of server addresses",
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
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IP address of the interface",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type of IP address",
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

func resourceAIClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start AI cluster creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, AIClusterGPUPoint, versionPointV1)
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

	if instancesCount, ok := d.GetOk("instances_count"); ok {
		createOpts.InstancesCount = instancesCount.(int)
	}

	if userData, ok := d.GetOk("user_data"); ok {
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

	ifs := d.Get("interface").(*schema.Set).List()
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
		// on create task, the cluster_id matches the task_id
		clusterID := taskInfo.ID
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
		} else {
			for _, ipAssignment := range instanceIface.IPAssignments {
				allDetachMap[ipAssignment.SubnetID] = ai.DetachInterfaceOpts{
					PortID:    instanceIface.PortID,
					IpAddress: ipAssignment.IPAddress.String(),
				}
			}
			// check if we already found the needed (same subnetID) detach options
			if detachOpts, found := allDetachMap[detachIface.SubnetID]; found {
				return &detachOpts, nil
			}
		}

		// if we haven't found it yet, check subports
		for _, subport := range instanceIface.SubPorts {
			if detachIface.Type == types.ExternalInterfaceType && subport.NetworkDetails.External {
				return &ai.DetachInterfaceOpts{
					PortID:    subport.PortID,
					IpAddress: subport.IPAssignments[0].IPAddress.String(),
				}, nil
			} else {
				for _, ipAssignment := range subport.IPAssignments {
					allDetachMap[ipAssignment.SubnetID] = ai.DetachInterfaceOpts{
						PortID:    subport.PortID,
						IpAddress: ipAssignment.IPAddress.String(),
					}
				}
			}
		}
	}
	if detachOpts, found := allDetachMap[detachIface.SubnetID]; found {
		return &detachOpts, nil
	} else {
		return nil, fmt.Errorf("couldn't find detach options for interface: %+v", detachIface)
	}
}

var IsResize = false

func resourceAIClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start AI cluster updating")
	config := m.(*Config)
	provider := config.Provider
	clientV1, err := CreateClient(provider, d, AIClusterPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	clientV1GPU, err := CreateClient(provider, d, AIClusterGPUPoint, versionPointV1)
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
	if d.HasChanges("instances_count") {
		IsResize = true
		_, newSGs := d.GetChange("security_group")
		securityGroupList := newSGs.(*schema.Set).List()
		securityGroupIDs := make([]gcorecloud.ItemID, len(securityGroupList))
		for sgIndex, sgID := range securityGroupList {
			securityGroupIDs[sgIndex] = gcorecloud.ItemID{ID: sgID.(map[string]interface{})["id"].(string)}
		}

		instancesCount, ok := d.GetOk("instances_count")
		if !ok || instancesCount.(int) == 0 {
			// if the number of instances has been specified before, then it cannot be removed or set to 0
			return diag.FromErr(errors.New("cannot resize cluster to 0 instances"))
		}

		resizeOpts := ai.ResizeGPUAIClusterOpts{
			InstancesCount: instancesCount.(int),
		}
		log.Printf("[DEBUG] AI cluster resize options: %+v", resizeOpts)
		results, err := ai.Resize(clientV1GPU, clusterID, resizeOpts).Extract()
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
			return diag.FromErr(errors.New("baremetal servers don't support external volumes"))
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
		newAttachInterfaces := map2AttachInterfaceOpts(newIfaces.(*schema.Set).List())
		if len(newAttachInterfaces) == 0 {
			return diag.FromErr(errors.New("no interfaces is configured, at least one is required"))
		}
		oldAttachInterfaces := map2AttachInterfaceOpts(oldIfaces.(*schema.Set).List())
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

	// if only the image_id has changed, then we need to rebuild the cluster
	if d.HasChange("image_id") && !d.HasChangesExcept("image_id") {
		_, newImageID := d.GetChange("image_id")
		cluster, err := ai.Get(clientV2, clusterID).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		var nodesIDs []string
		for _, instance := range cluster.PoplarServer {
			nodesIDs = append(nodesIDs, instance.ID)
		}
		rebuildOpts := ai.RebuildGPUAIClusterOpts{ImageID: newImageID.(string), Nodes: nodesIDs}

		log.Printf("[DEBUG] GPU cluster rebuild options: %+v", rebuildOpts)
		results, err := ai.RebuildGPUAICluster(clientV1GPU, clusterID, rebuildOpts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		taskClient, err := CreateClient(provider, d, TaskPoint, versionPointV1)
		if err != nil {
			return diag.FromErr(err)
		}
		taskID := results.Tasks[0]
		log.Printf("[DEBUG] GPU cluster task_id (%s)", taskID)

		_, err = tasks.WaitTaskAndReturnResult(taskClient, taskID, true, AIClusterCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
			return nil, nil
		})
		if err != nil {
			return diag.FromErr(err)
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
