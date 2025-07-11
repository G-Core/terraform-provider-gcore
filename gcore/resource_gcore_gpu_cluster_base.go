package gcore

import (
	"context"
	"fmt"
	"log"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/gpu/v3/clusters"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const GPUClusterCreationTimeout = 1800

type GPUNodeType string

const (
	GPUNodeTypeVirtual   GPUNodeType = "virtual"
	GPUNodeTypeBaremetal GPUNodeType = "baremetal"
)

func resourceGPUCluster(gpuNodeType GPUNodeType) *schema.Resource {
	return &schema.Resource{
		CreateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
			return resourceGPUClusterCreate(ctx, d, m, gpuNodeType)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
			return resourceGPUClusterRead(ctx, d, m, gpuNodeType)
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
			return resourceGPUClusterUpdate(ctx, d, m, gpuNodeType)
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
			return resourceGPUClusterDelete(ctx, d, m, gpuNodeType)
		},
		Description: fmt.Sprintf("Manages a %s GPU cluster", gpuNodeType),
		Schema:      resourceGPUClusterSchema(),
	}
}

func getGPUServicePath(gpuNodeType GPUNodeType) string {
	return fmt.Sprintf("gpu/%s", gpuNodeType)
}

func resourceGPUClusterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			ExactlyOneOf: []string{
				"project_id",
				"project_name",
			},
			DiffSuppressFunc: suppressDiffProjectID,
		},
		"region_id": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			ExactlyOneOf: []string{
				"region_id",
				"region_name",
			},
			DiffSuppressFunc: suppressDiffRegionID,
		},
		"project_name": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ExactlyOneOf: []string{
				"project_id",
				"project_name",
			},
		},
		"region_name": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ExactlyOneOf: []string{
				"region_id",
				"region_name",
			},
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the GPU cluster",
		},
		"flavor": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Flavor name for the GPU cluster",
		},

		"tags": {
			Type:        schema.TypeMap,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Tags to associate with the GPU cluster",
		},
		"servers_count": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Number of servers in the GPU cluster",
		},
		"servers_settings": {
			Type:        schema.TypeList,
			Required:    true,
			ForceNew:    true,
			MaxItems:    1,
			Description: "Settings for the GPU cluster servers (immutable, changes force resource recreation)",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"interface": {
						Type:        schema.TypeSet,
						Set:         instanceInterfaceUniqueIDByName,
						Required:    true,
						Description: "List of interfaces to attach to the instance",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"type": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Interface type (subnet, any_subnet, external)",
								},
								"name": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Name of interface, should be unique for the instance",
								},
								"network_id": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Required if type is 'subnet' or 'any_subnet'",
								},
								"subnet_id": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Required if type is 'subnet'",
								},
								"ip_family": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "IP family for the interface (dual, ipv4, ipv6)",
								},
								"ip_address": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "IP address for the interface",
								},
								"floating_ip": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "Floating IP configuration",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"source": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "Source of the floating IP",
											},
										},
									},
								},
							},
						},
					},
					"security_groups": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "List of security group IDs to associate with the cluster",
						Elem:        &schema.Schema{Type: schema.TypeString},
					},
					"volume": {
						Type:        schema.TypeList,
						Required:    true,
						Description: "Volumes to attach to the cluster servers",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"source": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Volume source (new, image, snapshot)",
								},
								"boot_index": {
									Type:        schema.TypeInt,
									Required:    true,
									Description: "Boot order for the volume",
								},
								"delete_on_termination": {
									Type:        schema.TypeBool,
									Optional:    true,
									Default:     false,
									Description: "Whether to delete the volume when the cluster is terminated",
								},
								"name": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Name of the volume",
								},
								"size": {
									Type:        schema.TypeInt,
									Required:    true,
									Description: "Size of the volume in GB",
								},
								"image_id": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "ID of the image to use (required if source is 'image')",
									Computed:    true,
								},
								"snapshot_id": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "ID of the snapshot to use (required if source is 'snapshot')",
									Computed:    true,
								},
								"tags": {
									Type:        schema.TypeMap,
									Optional:    true,
									Elem:        &schema.Schema{Type: schema.TypeString},
									Description: "Tags to associate with the volume",
								},
								"type": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Type of volume",
								},
							},
						},
					},
					"user_data": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "User data to provide to the instance for cloud-init",
						Computed:    true,
					},
					"credentials": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Credentials for accessing the instances",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"username": {
									Type:        schema.TypeString,
									Optional:    true,
									Computed:    true,
									Description: "Username for the instance",
								},
								"password": {
									Type:        schema.TypeString,
									Optional:    true,
									Computed:    true,
									Description: "Password for the instance",
								},
								"ssh_key_name": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Name of the keypair to use for SSH access",
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceGPUClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}, gpuNodeType GPUNodeType) diag.Diagnostics {
	log.Printf("[DEBUG] Start %s GPU cluster creation", gpuNodeType)
	log.Printf("[DEBUG] resource data: %+v", d)
	config := m.(*Config)
	provider := config.Provider
	client, err := CreateClient(provider, d, getGPUServicePath(gpuNodeType), "v3")
	if err != nil {
		return diag.FromErr(err)
	}

	// parse cluster creation options from schema
	tags := make(map[string]string)
	if v, ok := d.GetOk("tags"); ok {
		for k, val := range v.(map[string]interface{}) {
			tags[k] = val.(string)
		}
	}
	serverOpts, err := extractServerSettings(d)
	if err != nil {
		return diag.FromErr(err)
	}
	opts := &clusters.CreateClusterOpts{
		Name:            d.Get("name").(string),
		Flavor:          d.Get("flavor").(string),
		ServersCount:    d.Get("servers_count").(int),
		Tags:            tags,
		ServersSettings: serverOpts,
	}

	// create a cluster and wait for task completion
	result := clusters.Create(client, opts)
	if result.Err != nil {
		return diag.FromErr(result.Err)
	}
	taskID, err := result.Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	firstTaskID := taskID.Tasks[0]

	taskClient, err := CreateClient(provider, d, "tasks", "v1")
	if err != nil {
		return diag.FromErr(err)
	}

	waitResult, err := tasks.WaitTaskAndReturnResult(taskClient, firstTaskID, true, GPUClusterCreationTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(taskClient, string(task)).Extract()
		if err != nil {
			return nil, err
		}
		return taskInfo, nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	// extract cluster ID from the created results section of task results
	resultTask := waitResult.(*tasks.Task)
	clusterID, err := clusters.ExtractGPUClusterIDFromTask(resultTask)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(clusterID)
	return resourceGPUClusterRead(ctx, d, m, gpuNodeType)
}

func resourceGPUClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}, gpuNodeType GPUNodeType) diag.Diagnostics {
	log.Printf("[DEBUG] Start %s GPU cluster update", gpuNodeType)
	config := m.(*Config)
	provider := config.Provider
	client, err := CreateClient(provider, d, getGPUServicePath(gpuNodeType), "v3")
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("servers_count") {
		log.Printf("[DEBUG] Updating servers_count for %s GPU cluster", gpuNodeType)
		serversCount := d.Get("servers_count").(int)

		taskResult, err := clusters.Resize(client, d.Id(), serversCount).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		taskID := taskResult.Tasks[0]

		taskClient, err := CreateClient(provider, d, tasksPoint, versionPointV1)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = tasks.WaitTaskAndReturnResult(taskClient, taskID, true, GPUClusterCreationTimeout, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(taskClient, string(task)).Extract()
			if err != nil {
				return nil, err
			}
			return taskInfo, nil
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("name") {
		log.Printf("[DEBUG] Updating name for %s GPU cluster", gpuNodeType)
		name := d.Get("name").(string)
		opts := clusters.RenameClusterOpts{Name: name}
		getResult := clusters.Rename(client, d.Id(), opts)
		if getResult.Err != nil {
			return diag.FromErr(getResult.Err)
		}
	}

	if d.HasChange("tags") {
		log.Printf("[DEBUG] Updating tags for %s GPU cluster", gpuNodeType)

		// Get old and new tag values
		oldTags, newTags := d.GetChange("tags")

		// Create a new map with new tags
		newTagsMap := make(map[string]*string)
		if newTags != nil {
			for k, val := range newTags.(map[string]interface{}) {
				newTagsMap[k] = utils.StringToPointer(val.(string))
			}
		}

		// Convert old tags to map[string]string for comparison
		oldTagsMap := make(map[string]string)
		if oldTags != nil {
			for k, val := range oldTags.(map[string]interface{}) {
				oldTagsMap[k] = val.(string)
			}
		}

		// Identify tags to remove (present in old but not in new)
		for oldKey := range oldTagsMap {
			if _, exists := newTagsMap[oldKey]; !exists {
				newTagsMap[oldKey] = nil // nil value indicates removal
			}
		}
		log.Printf("[DEBUG] New tags sent for update: %v", newTagsMap)

		taskResult, err := clusters.UpdateAndRemoveTags(client, d.Id(), newTagsMap).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		taskID := taskResult.Tasks[0]

		taskClient, err := CreateClient(provider, d, tasksPoint, versionPointV1)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = tasks.WaitTaskAndReturnResult(taskClient, taskID, true, GPUClusterCreationTimeout, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(taskClient, string(task)).Extract()
			if err != nil {
				return nil, err
			}
			return taskInfo, nil
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGPUClusterRead(ctx, d, m, gpuNodeType)
}

func resourceGPUClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}, gpuNodeType GPUNodeType) diag.Diagnostics {
	log.Printf("[DEBUG] Start %s GPU image reading", gpuNodeType)
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	client, err := CreateClient(provider, d, getGPUServicePath(gpuNodeType), "v3")
	if err != nil {
		return diag.FromErr(err)
	}

	var cluster *clusters.Cluster
	err = retry.RetryContext(ctx, 20*time.Second, func() *retry.RetryError {
		getResult := clusters.Get(client, d.Id())
		if getResult.Err != nil {
			if _, ok := getResult.Err.(gcorecloud.ErrDefault404); ok {
				return retry.RetryableError(getResult.Err)
			}
			return retry.NonRetryableError(getResult.Err)
		}
		var err error
		cluster, err = getResult.Extract()
		if err != nil {
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	err = setGPUClusterResourceData(d, cluster)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] Finish GPU cluster reading")
	return diags
}

func resourceGPUClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}, gpuNodeType GPUNodeType) diag.Diagnostics {
	log.Printf("[DEBUG] Start %s GPU cluster deletion", gpuNodeType)
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	client, err := CreateClient(provider, d, getGPUServicePath(gpuNodeType), "v3")
	if err != nil {
		return diag.FromErr(err)
	}

	opts := &clusters.DeleteClusterOpts{AllVolumes: true}
	results, err := clusters.Delete(client, d.Id(), opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] GPU cluster deletion Task id: (%s)", taskID)

	tasksClient, err := CreateClient(provider, d, tasksPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = tasks.WaitTaskAndReturnResult(tasksClient, taskID, true, GPUClusterCreationTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := clusters.Get(client, d.Id()).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete GPU cluster with id: %s", d.Id())
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
	log.Printf("[DEBUG] Finish %s GPU cluster deletion", gpuNodeType)
	return diags
}

func setGPUClusterResourceData(d *schema.ResourceData, cluster *clusters.Cluster) error {
	if cluster == nil {
		return fmt.Errorf("cluster is nil")
	}

	log.Printf("[DEBUG] ==== tags: %v", cluster.Tags)

	d.Set("name", cluster.Name)
	d.Set("flavor", cluster.Flavor)
	d.Set("servers_count", cluster.ServersCount)
	// Ensure tags are properly set as a map[string]interface{} for Terraform
	tagsMap := make(map[string]interface{})
	for _, v := range cluster.Tags {
		tagsMap[v.Key] = v.Value
	}
	log.Printf("[DEBUG] tagsMap: %v", tagsMap)
	if err := d.Set("tags", tagsMap); err != nil {
		return fmt.Errorf("error setting tags: %v", err)
	}

	serverSettings, err := extractServerSettingsFromCluster(cluster)
	if err != nil {
		return err
	}

	// Preserve credentials from state since API doesn't return them
	if v, ok := d.GetOk("servers_settings"); ok {
		settingsList := v.([]interface{})
		if len(settingsList) > 0 {
			oldSettings := settingsList[0].(map[string]interface{})
			log.Printf("[DEBUG] Setting from old settings: (%s)", oldSettings)
			if creds, ok := oldSettings["credentials"]; ok && creds != nil {
				credsList := creds.([]interface{})
				if len(credsList) > 0 {
					credsMap := credsList[0].(map[string]interface{})
					log.Printf("[DEBUG] Setting credentials (before): (%s)", credsMap)

					// Log each field's value and type
					for k, v := range credsMap {
						log.Printf("[DEBUG] Credential field %s: value=%v, type=%T", k, v, v)
					}

					// Create a new map with only the non-empty values
					newCredsMap := make(map[string]interface{})

					// Only add non-empty values
					for k, v := range credsMap {
						if strVal, ok := v.(string); ok && strVal != "" {
							newCredsMap[k] = strVal
						}
					}

					// Only set credentials if we have at least one non-empty value
					if len(newCredsMap) > 0 {
						serverSettings["credentials"] = []interface{}{newCredsMap}
						log.Printf("[DEBUG] Setting credentials (after): (%s)", newCredsMap)
					}
				}
			}
		}
	}

	log.Printf("[DEBUG] Setting servers_settings: (%s)", serverSettings)
	if err := d.Set("servers_settings", []interface{}{serverSettings}); err != nil {
		return fmt.Errorf("error setting servers_settings: %v", err)
	}

	d.SetId(cluster.ID)
	return nil
}

func extractServerSettingsFromCluster(cluster *clusters.Cluster) (map[string]interface{}, error) {
	if cluster == nil {
		return nil, fmt.Errorf("cluster is nil")
	}

	serverSettings := make(map[string]interface{})

	// Extract interfaces
	if len(cluster.ServersSettings.Interfaces) > 0 {
		interfaces := make([]interface{}, 0, len(cluster.ServersSettings.Interfaces))
		for _, iface := range cluster.ServersSettings.Interfaces {
			ifaceMap := make(map[string]interface{})
			ifaceType := iface.InterfaceType()
			ifaceMap["type"] = ifaceType

			switch clusters.InterfaceType(ifaceType) {
			case clusters.External:
				ifaceMap["name"] = iface.ExternalInterface.Name
				ifaceMap["ip_family"] = string(iface.ExternalInterface.IPFamily)
			case clusters.Subnet:
				ifaceMap["name"] = iface.SubnetInterface.Name
				ifaceMap["network_id"] = iface.SubnetInterface.NetworkID
				ifaceMap["subnet_id"] = iface.SubnetInterface.SubnetID
				if iface.SubnetInterface.FloatingIP != nil {
					ifaceMap["floating_ip"] = map[string]string{"source": iface.SubnetInterface.FloatingIP.Source}
				}
			case clusters.AnySubnet:
				ifaceMap["name"] = iface.AnySubnetInterface.Name
				ifaceMap["network_id"] = iface.AnySubnetInterface.NetworkID
				ifaceMap["ip_address"] = iface.AnySubnetInterface.IPAddress
				ifaceMap["ip_family"] = string(iface.AnySubnetInterface.IPFamily)
				if iface.AnySubnetInterface.FloatingIP != nil {
					ifaceMap["floating_ip"] = map[string]string{"source": iface.AnySubnetInterface.FloatingIP.Source}
				}
			default:
				return nil, fmt.Errorf("unknown interface type: %s", ifaceType)
			}
			interfaces = append(interfaces, ifaceMap)
		}
		serverSettings["interface"] = interfaces
	}

	// Extract volumes
	if len(cluster.ServersSettings.Volumes) > 0 {
		volumes := make([]map[string]interface{}, 0)
		for _, vol := range cluster.ServersSettings.Volumes {
			volMap := map[string]interface{}{
				"boot_index":            vol.BootIndex,
				"delete_on_termination": vol.DeleteOnTermination,
				"name":                  vol.Name,
				"size":                  vol.Size,
				"type":                  string(vol.Type),
			}

			if vol.ImageID != nil && *vol.ImageID != "" {
				volMap["image_id"] = *vol.ImageID
				volMap["source"] = string(clusters.Image)
			}
			if vol.SnapshotID != nil && *vol.SnapshotID != "" {
				volMap["snapshot_id"] = *vol.SnapshotID
				volMap["source"] = string(clusters.Snapshot)
			}
			if vol.ImageID == nil && vol.SnapshotID == nil {
				volMap["source"] = string(clusters.NewVolume)
			}
			if vol.Tags != nil && len(vol.Tags) > 0 {
				volMap["tags"] = vol.Tags
			}

			volumes = append(volumes, volMap)
		}
		serverSettings["volume"] = volumes
	}

	// Extract security groups if available
	if len(cluster.ServersSettings.SecurityGroups) > 0 {
		serverSettings["security_groups"] = cluster.ServersSettings.SecurityGroups
	}
	// Don't set user_data in the resource data if it's nil or empty string from the API
	// This ensures consistency between nil and empty string representation
	if cluster.ServersSettings.UserData != nil && *cluster.ServersSettings.UserData != "" {
		serverSettings["user_data"] = *cluster.ServersSettings.UserData
	}

	return serverSettings, nil
}

func extractServerSettings(d *schema.ResourceData) (clusters.ServerSettingsOpts, error) {
	var serverSettings clusters.ServerSettingsOpts
	if v, ok := d.GetOk("servers_settings"); ok {
		settingsList := v.([]interface{})
		if len(settingsList) == 0 {
			return serverSettings, fmt.Errorf("servers_settings cannot be empty")
		}
		settingsMap := settingsList[0].(map[string]interface{})
		serverSettings.Interfaces = extractInterfaces(settingsMap["interface"])
		serverSettings.Volumes = extractVolumes(settingsMap["volume"])

		if secGroups, ok := settingsMap["security_groups"]; ok {
			serverSettings.SecurityGroups = extractSecurityGroups(secGroups.([]interface{}))
		}
		if userData, ok := settingsMap["user_data"].(string); ok {
			serverSettings.UserData = utils.StringToPointer(userData)
		}
		if creds, ok := settingsMap["credentials"]; ok {
			serverSettings.Credentials = extractCredentials(creds.([]interface{}))
		}
	}
	return serverSettings, nil
}

func extractInterfaces(interfaces any) []clusters.InterfaceOpts {
	if interfaces == nil {
		return nil
	}
	interfaceSet := interfaces.(*schema.Set)
	result := make([]clusters.InterfaceOpts, 0, interfaceSet.Len())
	for _, item := range interfaceSet.List() {
		ifaceMap := item.(map[string]interface{})
		ifaceType := ifaceMap["type"].(string)
		var iface clusters.InterfaceOpts

		switch clusters.InterfaceType(ifaceType) {
		case clusters.External:
			iface = clusters.ExternalInterfaceOpts{
				Name:     utils.StringToPointer(ifaceMap["name"].(string)),
				Type:     ifaceType,
				IPFamily: clusters.IPFamilyType(ifaceMap["ip_family"].(string)),
			}
		case clusters.Subnet:
			floatingIPOpts := extractFloatingIP(ifaceMap)
			iface = clusters.SubnetInterfaceOpts{
				NetworkID:  ifaceMap["network_id"].(string),
				Name:       utils.StringToPointer(ifaceMap["name"].(string)),
				Type:       ifaceType,
				SubnetID:   ifaceMap["subnet_id"].(string),
				FloatingIP: floatingIPOpts,
			}
		case clusters.AnySubnet:
			floatingIPOpts := extractFloatingIP(ifaceMap)
			iface = clusters.AnySubnetInterfaceOpts{
				NetworkID:  ifaceMap["network_id"].(string),
				Name:       utils.StringToPointer(ifaceMap["name"].(string)),
				Type:       ifaceType,
				IPFamily:   clusters.IPFamilyType(ifaceMap["ip_family"].(string)),
				IPAddress:  utils.StringToPointer(ifaceMap["ip_address"].(string)),
				FloatingIP: floatingIPOpts,
			}
		}
		result = append(result, iface)
	}
	return result
}

func extractFloatingIP(ifaceMap map[string]interface{}) *clusters.FloatingIPOpts {
	var floatingIPOpts *clusters.FloatingIPOpts
	if floatingIP, ok := ifaceMap["floating_ip"]; ok && floatingIP != nil {
		if floatingIPMaps, ok := floatingIP.([]interface{}); ok && len(floatingIPMaps) > 0 {
			floatingIPMap := floatingIPMaps[0].(map[string]interface{})
			floatingIPOpts = &clusters.FloatingIPOpts{
				Source: floatingIPMap["source"].(string),
			}
		}
	}
	return floatingIPOpts
}

func extractVolumes(volume any) []clusters.VolumeOpts {
	if volume == nil {
		return nil
	}
	volumeList := volume.([]interface{})
	result := make([]clusters.VolumeOpts, 0, len(volumeList))
	for _, item := range volumeList {
		volMap := item.(map[string]interface{})
		source := volMap["source"].(string)

		volOpts := clusters.VolumeOpts{
			Source:              clusters.VolumeSource(source),
			BootIndex:           volMap["boot_index"].(int),
			DeleteOnTermination: volMap["delete_on_termination"].(bool),
			Name:                volMap["name"].(string),
			Size:                volMap["size"].(int),
			Type:                clusters.VolumeType(volMap["type"].(string)),
		}
		if imageID, ok := volMap["image_id"]; ok && imageID != "" {
			volOpts.ImageID = imageID.(string)
		}
		if snapshotID, ok := volMap["snapshot_id"]; ok && snapshotID != "" {
			volOpts.SnapshotID = snapshotID.(string)
		}
		// Handle optional fields
		if deleteOnTermination, ok := volMap["delete_on_termination"]; ok {
			volOpts.DeleteOnTermination = deleteOnTermination.(bool)
		}
		if tags, ok := volMap["tags"]; ok && tags != nil {
			volOpts.Tags = make(map[string]string)
			for k, val := range tags.(map[string]interface{}) {
				volOpts.Tags[k] = val.(string)
			}
		}
		result = append(result, volOpts)
	}
	return result
}

func extractSecurityGroups(securityGroups []interface{}) []string {
	if len(securityGroups) == 0 {
		return nil
	}
	result := make([]string, 0, len(securityGroups))
	for _, sg := range securityGroups {
		if sg != nil {
			result = append(result, sg.(string))
		}
	}
	return result
}

func extractCredentials(credentials []interface{}) *clusters.ServerCredentialsOpts {
	if len(credentials) == 0 {
		return nil
	}

	credsMap := credentials[0].(map[string]interface{})
	result := &clusters.ServerCredentialsOpts{}
	if username, ok := credsMap["username"].(string); ok && username != "" {
		result.Username = username
	}
	if password, ok := credsMap["password"].(string); ok && password != "" {
		result.Password = password
	}
	if sshKeyName, ok := credsMap["ssh_key_name"].(string); ok && sshKeyName != "" {
		result.SSHKeyName = sshKeyName
	}
	return result
}
