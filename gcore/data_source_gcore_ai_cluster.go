package gcore

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	ai "github.com/G-Core/gcorelabscloud-go/gcore/ai/v1/ais"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	BM_FLAVOR_PREFIX = "bm"
)

func dataSourceAICluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAIClusterRead,
		Description: "Represent GPU Cluster",
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
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "AI Cluster ID",
				Required:    true,
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Description: "AI Cluster Name",
				Computed:    true,
			},
			"cluster_status": {
				Type:        schema.TypeString,
				Description: "AI Cluster status",
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
				Computed:    true,
			},
			"image_id": {
				Type:        schema.TypeString,
				Description: "Image ID",
				Computed:    true,
			},
			"image_name": {
				Type:        schema.TypeString,
				Description: "Image name",
				Computed:    true,
			},
			"volume": {
				Type:        schema.TypeSet,
				Description: "List of volumes attached to the cluster",
				Computed:    true,
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
							Computed:    true,
							Optional:    true,
						},
						"attachments": {
							Type:        schema.TypeSet,
							Description: "Attachment list",
							Computed:    true,
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
				Computed:    true,
				Set:         aiSgHashID,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Security group ID",
							Computed:    true,
						},
					},
				},
			},
			"interface": {
				Type:        schema.TypeSet,
				Description: "Networks managed by user and associated with the cluster",
				Computed:    true,
				Set:         aiInterfaceHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Description: "Network type",
							Computed:    true,
						},
						"network_id": {
							Type:        schema.TypeString,
							Description: "Network ID",
							Computed:    true,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "Network ID the subnet belongs to. Port will be plugged in this network",
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
				Description: "Ssh keypair name",
				Computed:    true,
			},
			"user_data": {
				Type:        schema.TypeString,
				Description: "String in base64 format. Must not be passed together with 'username' or 'password'. Examples of the user_data: https://cloudinit.readthedocs.io/en/latest/topics/examples.html",
				Computed:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "A name of a new user in the Linux instance. It may be passed with a 'password' parameter",
				Computed:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "A password for baremetal instance. This parameter is used to set a password for the Admin user on a Windows instance, a default user or a new user on a Linux instance",
				Computed:    true,
			},
			"poplar_servers": {
				Type:        schema.TypeList,
				Description: "GPU cluster servers list",
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
			"cluster_metadata": {
				Type:        schema.TypeMap,
				Description: "Cluster metadata (simple key-value pairs)",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func aiSgHashID(i interface{}) int {
	var buf bytes.Buffer
	sg := i.(map[string]interface{})
	buf.WriteString(sg["id"].(string))
	return schema.HashString(buf.String())
}

func aiInterfaceHash(i interface{}) int {
	var buf bytes.Buffer
	iface := i.(map[string]interface{})
	switch {
	case iface["type"].(string) == string(types.ExternalInterfaceType):
		buf.WriteString(string(types.ExternalInterfaceType))
	case iface["type"].(string) == string(types.SubnetInterfaceType):
		buf.WriteString(string(types.SubnetInterfaceType))
		buf.WriteString(iface["network_id"].(string))
		buf.WriteString(iface["subnet_id"].(string))
	case iface["type"].(string) == string(types.AnySubnetInterfaceType):
		buf.WriteString(string(types.AnySubnetInterfaceType))
		buf.WriteString(iface["network_id"].(string))
	case iface["type"].(string) == string(types.ReservedFixedIpType):
		buf.WriteString(string(types.ReservedFixedIpType))
		buf.WriteString(iface["port_id"].(string))
	}
	return schema.HashString(buf.String())
}

func aiVolumeHash(i interface{}) int {
	var buf bytes.Buffer
	volume := i.(map[string]interface{})
	switch {
	case volume["source"].(string) == string(types.Image):
		buf.WriteString(string(types.Image))
		buf.WriteString(volume["image_id"].(string))
		buf.WriteString(volume["volume_type"].(string))
		buf.WriteString(strconv.Itoa(volume["size"].(int)))
	case volume["source"].(string) == string(types.ExistingVolume):
		buf.WriteString(string(types.ExistingVolume))
		buf.WriteString(volume["volume_id"].(string))
	}
	return schema.HashString(buf.String())
}

func flattenSecurityGroup(securityGroups []ai.PoplarInterfaceSecGrop) []interface{} {
	if len(securityGroups) == 0 {
		return nil
	}
	sgIDs := make([]interface{}, len(securityGroups[0].SecurityGroups))
	for index, sgID := range securityGroups[0].SecurityGroups {
		sgMap := make(map[string]interface{})
		sgMap["id"] = sgID
		sgIDs[index] = sgMap
	}
	return sgIDs
}

func flattenInterfaces(interfaces []ai.AIClusterInterface) []map[string]interface{} {
	clusterInterfaces := make([]map[string]interface{}, len(interfaces))
	for index, iface := range interfaces {
		interfaceMap := make(map[string]interface{})
		interfaceMap["type"] = iface.Type
		interfaceMap["port_id"] = iface.PortID
		interfaceMap["network_id"] = iface.NetworkID
		interfaceMap["subnet_id"] = iface.SubnetID
		clusterInterfaces[index] = interfaceMap

	}

	return clusterInterfaces
}

func flattenPoplarServers(servers []instances.Instance) []map[string]interface{} {
	serverList := make([]map[string]interface{}, len(servers))
	for index, instance := range servers {
		serverMap := make(map[string]interface{})
		serverMap["instance_id"] = instance.ID
		serverMap["instance_name"] = instance.Name
		serverMap["flavor_id"] = instance.Flavor.FlavorID
		serverMap["status"] = instance.Status
		serverMap["vm_state"] = instance.VMState
		serverMap["metadata"] = instance.Metadata
		serverMap["task_id"] = instance.TaskID
		serverMap["creator_task_id"] = instance.CreatorTaskID
		serverMap["task_state"] = instance.TaskState
		serverMap["instance_created"] = instance.CreatedAt.Format(gcorecloud.RFC3339NoZ)
		serverMap["instance_description"] = instance.Description

		flavor := make(map[string]interface{}, 4)
		flavor["flavor_id"] = instance.Flavor.FlavorID
		flavor["flavor_name"] = instance.Flavor.FlavorName
		flavor["ram"] = strconv.Itoa(instance.Flavor.RAM)
		flavor["vcpus"] = strconv.Itoa(instance.Flavor.VCPUS)
		serverMap["flavor"] = flavor

		volumes := make([]map[string]interface{}, len(instance.Volumes))
		for volIndex, vol := range instance.Volumes {
			volMap := make(map[string]interface{})
			volMap["id"] = vol.ID
			volMap["delete_on_termination"] = vol.DeleteOnTermination
			volumes[volIndex] = volMap
		}
		serverMap["volumes"] = volumes

		addresses := make([]map[string]interface{}, len(instance.Addresses))
		var addressIndex int
		for netName, addrs := range instance.Addresses {
			ifacesMap := make(map[string]interface{})
			addrsList := make([]map[string]interface{}, len(addrs))
			for addrIndex, addr := range addrs {
				addrMap := make(map[string]interface{})
				addrMap["addr"] = addr.Address.String()
				addrMap["type"] = addr.Type
				addrsList[addrIndex] = addrMap
			}
			ifacesMap["network_name"] = netName
			ifacesMap["address"] = addrsList
			addresses[addressIndex] = ifacesMap
			addressIndex++
		}
		serverMap["addresses"] = addresses

		secGrps := make([]map[string]interface{}, len(instance.SecurityGroups))
		for sgIndex, sg := range instance.SecurityGroups {
			sgMap := make(map[string]interface{})
			sgMap["name"] = sg.Name
			secGrps[sgIndex] = sgMap
		}
		serverMap["security_groups"] = secGrps

		serverList[index] = serverMap
	}
	return serverList
}

func flattenClusterVolumes(volumes []volumes.Volume) ([]interface{}, error) {
	volumeList := make([]interface{}, len(volumes))
	for index, volume := range volumes {
		volMap := make(map[string]interface{})
		volMap["volume_id"] = volume.ID
		volMap["name"] = volume.Name
		volMap["status"] = volume.Status
		volMap["size"] = volume.Size
		volMap["created_at"] = volume.CreatedAt.Format(gcorecloud.RFC3339NoZ)
		volMap["updated_at"] = volume.UpdatedAt.Format(gcorecloud.RFC3339NoZ)
		volMap["volume_type"] = volume.VolumeType
		volMap["creator_task_id"] = volume.CreatorTaskID
		volMap["image_id"] = volume.VolumeImageMetadata.ImageID
		if volume.VolumeImageMetadata.ImageID != "" {
			volMap["source"] = types.Image
		} else {
			volMap["source"] = types.ExistingVolume
		}

		imageMetadataMap := make(map[string]interface{})
		data, err := json.Marshal(volume.VolumeImageMetadata)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &imageMetadataMap)
		if err != nil {
			return nil, err
		}

		volMap["volume_image_metadata"] = imageMetadataMap

		attachments := make([]map[string]interface{}, len(volume.Attachments))
		for attachIndex, attach := range volume.Attachments {
			attachMap := make(map[string]interface{})
			attachMap["server_id"] = attach.ServerID
			attachMap["instance_name"] = attach.InstanceName
			attachMap["attachment_id"] = attach.AttachmentID
			attachMap["volume_id"] = attach.VolumeID
			attachMap["device"] = attach.Device
			attachMap["attached_at"] = attach.AttachedAt.Format(gcorecloud.RFC3339NoZ)
			attachments[attachIndex] = attachMap
		}
		volMap["attachments"] = attachments

		volumeList[index] = volMap
	}
	return volumeList, nil
}

func isBmFlavor(flavor string) bool {
	return strings.HasPrefix(flavor, "bm")
}

func setAIClusterResourcerData(d *schema.ResourceData, provider *gcorecloud.ProviderClient, cluster *ai.AICluster) error {
	d.Set("region_id", cluster.RegionID)
	d.Set("region_name", cluster.Region)
	d.Set("project_id", cluster.ProjectID)
	d.Set("cluster_name", cluster.ClusterName)
	d.Set("cluster_id", cluster.ClusterID)
	d.Set("flavor", cluster.Flavor)
	d.Set("cluster_status", cluster.ClusterStatus)
	d.Set("task_id", cluster.TaskID)
	d.Set("task_status", cluster.TaskStatus)
	d.Set("created_at", cluster.CreatedAt.Format(gcorecloud.RFC3339NoZ))
	d.Set("image_id", cluster.ImageID)
	d.Set("image_name", cluster.ImageName)
	d.Set("password", cluster.Password)
	d.Set("username", cluster.Username)
	d.Set("keypair_name", cluster.KeypairName)
	d.Set("user_data", cluster.UserData)
	d.Set("security_group", flattenSecurityGroup(cluster.SecurityGroups))
	client, err := CreateClient(provider, d, AIClusterPoint, versionPointV1)
	if err != nil {
		return err
	}
	clusterInterfaces, err := ai.ListInterfacesAll(client, cluster.ClusterID)
	if err != nil {
		return err
	}

	// we don't know how many interfaces we will have
	var aiClusterInterfaces []ai.AIClusterInterface
	for _, iface := range clusterInterfaces {
		// get parent interface (pub or private and private might have multiple subnets)
		// we add one interface per subnet
		var ifaceType string
		if iface.NetworkDetails.External {
			ifaceType = string(types.ExternalInterfaceType)
			// external network has multiple subnets, leaving SubnetID empty
			ifaceParent := ai.AIClusterInterface{
				Type:      ifaceType,
				PortID:    iface.PortID,
				NetworkID: iface.NetworkID,
			}
			aiClusterInterfaces = append(aiClusterInterfaces, ifaceParent)
		} else {
			ifaceType = string(types.SubnetInterfaceType)
			for _, subnet := range iface.NetworkDetails.Subnets {
				ifaceParent := ai.AIClusterInterface{
					Type:      ifaceType,
					PortID:    iface.PortID,
					NetworkID: iface.NetworkID,
					SubnetID:  subnet.ID,
				}
				aiClusterInterfaces = append(aiClusterInterfaces, ifaceParent)
			}
		}

		// check if there are more interfaces, if yes, they are inside subports
		if len(iface.SubPorts) > 0 {
			for _, subPort := range iface.SubPorts {
				if subPort.NetworkDetails.External {
					ifaceType = string(types.ExternalInterfaceType)
					// external network has multiple subnets, leaving SubnetID empty
					ifaceSubPort := ai.AIClusterInterface{
						Type:      ifaceType,
						PortID:    subPort.PortID,
						NetworkID: subPort.NetworkID,
					}
					aiClusterInterfaces = append(aiClusterInterfaces, ifaceSubPort)
				} else {
					ifaceType = string(types.SubnetInterfaceType)
					for _, subnet := range subPort.NetworkDetails.Subnets {
						ifaceSubPort := ai.AIClusterInterface{
							Type:      ifaceType,
							PortID:    subPort.PortID,
							NetworkID: subPort.NetworkID,
							SubnetID:  subnet.ID,
						}
						aiClusterInterfaces = append(aiClusterInterfaces, ifaceSubPort)
					}
				}
			}
		}
	}

	d.Set("interface", flattenInterfaces(aiClusterInterfaces))
	d.Set("cluster_metadata", cluster.Metadata)
	d.Set("poplar_servers", flattenPoplarServers(cluster.PoplarServer))

	volumes, err := flattenClusterVolumes(cluster.Volumes)
	if err != nil {
		return err
	}
	d.Set("volume", volumes)

	return nil
}

func dataSourceAIClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start AI Cluster reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	clusterID := d.Get("cluster_id").(string)
	log.Printf("[DEBUG] Getting AI cluster id = %s", clusterID)

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
	d.SetId(cluster.ClusterID)
	err = setAIClusterResourcerData(d, provider, cluster)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish AI cluster reading")
	return diags
}
