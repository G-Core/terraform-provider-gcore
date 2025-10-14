package gcore

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/G-Core/gcorelabscloud-go/client/utils"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/file_share/v1/file_shares"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const fileSharePoint = "file_shares"
const fileShareCreatingTimeout = 1200
const fileShareDeletingTimeout = 1200

func resourceFileShare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFileShareCreate,
		ReadContext:   resourceFileShareRead,
		UpdateContext: resourceFileShareUpdate,
		DeleteContext: resourceFileShareDelete,
		Description:   "Represents a file share (NFS) in Gcore Cloud.",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, fileShareID, err := ImportStringParser(d.Id())
				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(fileShareID)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				Description:      "Project ID, only one of project_id or project_name should be set",
				ExactlyOneOf:     []string{"project_id", "project_name"},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"region_id": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				Description:      "Region ID, only one of region_id or region_name should be set",
				ExactlyOneOf:     []string{"region_id", "region_name"},
				DiffSuppressFunc: suppressDiffRegionID,
			},
			"project_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "Project name, only one of project_id or project_name should be set",
				ExactlyOneOf: []string{"project_id", "project_name"},
			},
			"region_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "Region name, only one of region_id or region_name should be set",
				ExactlyOneOf: []string{"region_id", "region_name"},
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the file share. It must be unique within the project and region.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The protocol used by the file share. Currently, only 'NFS' is supported.`,
				ValidateFunc: func(v interface{}, k string) ([]string, []error) {
					if v.(string) != "NFS" {
						return nil, []error{fmt.Errorf("only 'NFS' protocol is supported")}
					}
					return nil, nil
				},
			},
			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The size of the file share in GB. It must be a positive integer.`,
			},
			"type_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The type of the file share. Must be one of 'standard' or 'vast'.`,
				ValidateFunc: func(v interface{}, k string) ([]string, []error) {
					if v.(string) != "standard" && v.(string) != "vast" {
						return nil, []error{fmt.Errorf("type_name must be 'standard' or 'vast'")}
					}
					return nil, nil
				},
			},
			"network": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				ForceNew:    true,
				Description: "Network configuration for the file share. It must include a network ID and optionally a subnet ID. (Only required for type_name: 'standard')",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The ID of the network to which the file share will be connected. This is required for 'standard'.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "The ID of the subnet within the network. This is optional and can be used to specify a particular subnet for the file share.",
						},
					},
				},
			},
			"access": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The IP address of the file share.`,
						},
						"access_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The access mode of the file share (ro/rw).`,
							ValidateFunc: func(v interface{}, k string) ([]string, []error) {
								if v.(string) != "ro" && v.(string) != "rw" {
									return nil, []error{fmt.Errorf("access_mode must be 'ro' or 'rw'")}
								}
								return nil, nil
							}},
					},
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Tags to associate with the file share. Tags are key-value pairs.",
			},
			"connection_point": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The connection point of the file share.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the file share.`,
			},
			"created_at": {
				Type: schema.TypeString, Computed: true,
				Description: `The creation time of the file share in ISO 8601 format.`,
			},
			"creator_task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the task that created the file share.`,
			},
			"network_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the network associated with the file share.`,
			},
			"share_network_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the share network associated with the file share. This is only applicable for 'standard'.`,
			},
			"subnet_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the subnet associated with the file share`,
			},
			"share_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Share settings for the file share.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_characters": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Allowed characters in file names.",
							ValidateFunc: func(v interface{}, k string) ([]string, []error) {
								val := v.(string)
								allowedCharacters := file_shares.FileShareAllowedCharacters("").StringList()
								for _, v := range allowedCharacters {
									if v == val {
										return nil, nil
									}
								}
								return nil, []error{fmt.Errorf("%s must be one of: %s", val,
									strings.Join(allowedCharacters, ", "))}
							},
						},
						"path_length": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Affects the maximum limit of file path component name length.",
							ValidateFunc: func(v interface{}, k string) ([]string, []error) {
								val := v.(string)
								pathLengths := file_shares.FileSharePathLength("").StringList()
								for _, v := range pathLengths {
									if v == val {
										return nil, nil
									}
								}
								return nil, []error{fmt.Errorf("%s must be one of: %s", val,
									strings.Join(pathLengths, ", "))}
							},
						},
						"root_squash": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates if root squash is enabled.",
						},
					},
				},
			},
		},
	}
}

func resourceFileShareCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start file share creating")
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, fileSharePoint, "v1")
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts, err := expandFileShareCreateOpts(d)
	if err != nil {
		return diag.FromErr(err)
	}
	result := file_shares.Create(client, createOpts)
	if result.Err != nil {
		return diag.FromErr(result.Err)
	}
	taskResults, err := result.Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	taskID := taskResults.Tasks[0]

	fileShareID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, fileShareCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		return file_shares.ExtractFileShareIDFromTask(taskInfo)
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fileShareID.(string))
	return resourceFileShareRead(ctx, d, m)
}

func resourceFileShareRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start file share reading")
	config := m.(*Config)
	provider := config.Provider
	fileShareID := d.Id()
	client, err := CreateClient(provider, d, fileSharePoint, "v1")
	if err != nil {
		return diag.FromErr(err)
	}
	result := file_shares.Get(client, fileShareID)
	fileShare, err := result.Extract()
	if err != nil {
		var errDefault404 gcorecloud.ErrDefault404
		if errors.As(err, &errDefault404) {
			// removing from state because it doesn't exist anymore
			d.SetId("")
			return nil
		}
		return diag.Errorf("cannot get file share with ID: %s. Error: %s", fileShareID, err)
	}
	d.Set("name", fileShare.Name)
	d.Set("protocol", fileShare.Protocol)
	d.Set("size", fileShare.Size)
	d.Set("status", string(fileShare.Status))
	d.Set("connection_point", fileShare.ConnectionPoint)
	d.Set("created_at", fileShare.CreatedAt.String())
	d.Set("type_name", fileShare.TypeName)
	if fileShare.TypeName == "standard" {
		networkMap := map[string]interface{}{}
		if fileShare.NetworkID != "" {
			networkMap["network_id"] = fileShare.NetworkID
		}
		if fileShare.SubnetID != "" {
			networkMap["subnet_id"] = fileShare.SubnetID
		}
		if len(networkMap) > 0 {
			d.Set("network", []interface{}{networkMap})
		}
	}
	d.Set("network_name", fileShare.NetworkName)
	d.Set("subnet_name", fileShare.SubnetName)
	d.Set("creator_task_id", fileShare.CreatorTaskID)
	if fileShare.ShareNetworkName != nil {
		d.Set("share_network_name", *fileShare.ShareNetworkName)
	}
	if fileShare.Tags != nil {
		tags := make(map[string]string)
		for _, tag := range fileShare.Tags {
			if !tag.ReadOnly {
				tags[tag.Key] = tag.Value
			}
		}
		d.Set("tags", tags)
	}
	shareSettingsMap := map[string]interface{}{}
	if fileShare.ShareSettings.AllowedCharacters != nil && fileShare.ShareSettings.AllowedCharacters.String() != "" {
		shareSettingsMap["allowed_characters"] = fileShare.ShareSettings.AllowedCharacters.String()
	}
	if fileShare.ShareSettings.PathLength != nil && fileShare.ShareSettings.PathLength.String() != "" {
		shareSettingsMap["path_length"] = fileShare.ShareSettings.PathLength.String()
	}
	if fileShare.ShareSettings.RootSquash != nil {
		shareSettingsMap["root_squash"] = *fileShare.ShareSettings.RootSquash
	}
	if len(shareSettingsMap) > 0 {
		d.Set("share_settings", []interface{}{shareSettingsMap})
	}
	return nil
}

func resourceFileShareUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start file share updating")
	config := m.(*Config)
	provider := config.Provider
	fileShareID := d.Id()

	clientV1, err := CreateClient(provider, d, fileSharePoint, "v1")
	if err != nil {
		return diag.FromErr(err)
	}
	clientV3, err := CreateClient(provider, d, fileSharePoint, "v3")
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("size") {
		newSize := d.Get("size").(int)
		if newSize > 0 {
			extendOpts := file_shares.ExtendOpts{Size: newSize}
			result := file_shares.Extend(clientV1, fileShareID, extendOpts)
			if result.Err != nil {
				return diag.FromErr(result.Err)
			}

			//TODO: replace code below with waitForTaskResult() when PR #273 is merged
			taskResults, err := result.Extract()
			if err != nil {
				return diag.FromErr(err)
			}
			if len(taskResults.Tasks) == 0 {
				return diag.FromErr(errors.New("no task IDs returned"))
			}
			taskID := taskResults.Tasks[0]

			if err := tasks.WaitForFinishedTask(clientV1, taskID, fileShareCreatingTimeout); err != nil {
				return diag.FromErr(fmt.Errorf("error while waiting for task %s to finish: %w", taskID, err))
			}

		}
	}

	if d.HasChange("name") || d.HasChange("tags") || d.HasChange("share_settings") {
		// Tags needs to be initialized to avoid sending null to the API and removing all tags.
		updateOpts := file_shares.UpdateWithTagsOpts{Tags: make(map[string]*string)}
		newName := d.Get("name").(string)
		if d.HasChange("name") && newName != "" {
			updateOpts.Name = newName
		}

		if d.HasChange("tags") {
			log.Println("[DEBUG] Updating tags")
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
			updateOpts.Tags = newTagsMap
		}

		if d.HasChange("share_settings") {
			shareSettingsList := d.Get("share_settings").([]interface{})
			if len(shareSettingsList) > 0 {
				// schema restricts share_settings to a single block
				ssMap := shareSettingsList[0].(map[string]interface{})
				var shareSettingsOpts file_shares.ShareSettingsOpts
				if v, ok := ssMap["allowed_characters"]; ok && v != nil && v.(string) != "" {
					shareSettingsOpts.AllowedCharacters = file_shares.FileShareAllowedCharacters(v.(string))
				}
				if v, ok := ssMap["path_length"]; ok && v != nil && v.(string) != "" {
					shareSettingsOpts.PathLength = file_shares.FileSharePathLength(v.(string))
				}
				if v, ok := ssMap["root_squash"]; ok && v != nil {
					rootSquash := v.(bool)
					shareSettingsOpts.RootSquash = &rootSquash
				}
				updateOpts.ShareSettings = &shareSettingsOpts
			}
		}

		result := file_shares.UpdateWithTags(clientV3, fileShareID, updateOpts)
		if result.Err != nil {
			return diag.FromErr(result.Err)
		}
		//TODO: replace code below with waitForTaskResult() when PR #273 is merged
		taskResults, err := result.Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		if len(taskResults.Tasks) == 0 {
			return diag.FromErr(errors.New("no task IDs returned"))
		}
		taskID := taskResults.Tasks[0]

		if err := tasks.WaitForFinishedTask(clientV1, taskID, fileShareCreatingTimeout); err != nil {
			return diag.FromErr(fmt.Errorf("error while waiting for task %s to finish: %w", taskID, err))
		}
	}

	// Handle access rules update: remove all and re-create if changed
	if d.HasChange("access") {
		log.Println("[DEBUG] Updating access rules for file share")
		// List all current access rules
		pager := file_shares.ListAccessRules(clientV1, fileShareID)
		pages, err := pager.AllPages()
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to get all access rule pages for file share %s: %w", fileShareID, err))
		}
		rules, err := file_shares.ExtractAccessRule(pages)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to extract access rules for file share %s: %w", fileShareID, err))
		}
		// Delete all existing access rules
		for _, rule := range rules {
			result := file_shares.DeleteAccessRule(clientV1, fileShareID, rule.ID)
			if result.Err != nil {
				return diag.FromErr(fmt.Errorf("failed to delete access rule %s: %w", rule.ID, result.Err))
			}
		}
		// Add new access rules from config
		accessList := d.Get("access").([]interface{})
		for _, a := range accessList {
			amap := a.(map[string]interface{})
			createOpts := file_shares.CreateAccessRuleOpts{
				IPAddress:  amap["ip_address"].(string),
				AccessMode: amap["access_mode"].(string),
			}
			result := file_shares.CreateAccessRule(clientV1, fileShareID, createOpts)
			if result.Err != nil {
				return diag.FromErr(fmt.Errorf("failed to create access rule for file share %s: %w", fileShareID, result.Err))
			}
		}
	}
	return resourceFileShareRead(ctx, d, m)
}

func resourceFileShareDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start file share deleting")
	config := m.(*Config)
	provider := config.Provider
	fileShareID := d.Id()
	client, err := CreateClient(provider, d, fileSharePoint, "v1")
	if err != nil {
		return diag.FromErr(err)
	}
	result := file_shares.Delete(client, fileShareID)
	if result.Err != nil {
		return diag.FromErr(result.Err)
	}
	taskResults, err := result.Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	taskID := taskResults.Tasks[0]

	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, fileShareDeletingTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := file_shares.Get(client, fileShareID).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete file share with ID: %s", fileShareID)
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
	log.Printf("[DEBUG] Finish of file share deleting")
	return nil
}

func expandFileShareCreateOpts(d *schema.ResourceData) (*file_shares.CreateOpts, error) {
	name := d.Get("name").(string)
	protocol := d.Get("protocol").(string)
	size := d.Get("size").(int)

	tags := make(map[string]string)
	if v, ok := d.GetOk("tags"); ok {
		for k, val := range v.(map[string]interface{}) {
			tags[k] = val.(string)
		}
	}
	// determine file share type name (new API: 'standard' or 'vast')
	typeNameRaw, hasTypeName := d.GetOk("type_name")
	var typeName string
	if hasTypeName {
		typeName = typeNameRaw.(string)
		if typeName != "standard" && typeName != "vast" {
			return nil, fmt.Errorf("type_name must be 'standard' or 'vast'")
		}
	} else {
		return nil, fmt.Errorf("type_name is required")
	}

	// check that network and access are set only for 'standard'
	if typeName == "vast" {
		networkList := d.Get("network").([]interface{})
		if len(networkList) > 0 {
			return nil, fmt.Errorf("network block is not allowed for 'vast'")
		}
		accessList := d.Get("access").([]interface{})
		if len(accessList) > 0 {
			return nil, fmt.Errorf("access block is not allowed for 'vast'")
		}
	}

	// check that share_settings block is only set for 'vast'
	shareSettingsList := d.Get("share_settings").([]interface{})

	opts := file_shares.CreateOpts{
		Name:     name,
		Protocol: protocol,
		Size:     size,
		Tags:     tags,
		TypeName: typeName,
	}

	if typeName == "standard" {
		networkList := d.Get("network").([]interface{})
		var networkOpts file_shares.FileShareNetworkOpts
		if len(networkList) > 0 {
			netMap := networkList[0].(map[string]interface{})
			networkOpts.NetworkID = netMap["network_id"].(string)
			if v, ok := netMap["subnet_id"]; ok && v != nil && v.(string) != "" {
				networkOpts.SubnetID = v.(string)
			}
		} else {
			return nil, fmt.Errorf("network block is required")
		}
		opts.Network = &networkOpts

		accessList := d.Get("access").([]interface{})
		access := make([]file_shares.CreateAccessRuleOpts, 0, len(accessList))
		for _, a := range accessList {
			amap := a.(map[string]interface{})
			access = append(access, file_shares.CreateAccessRuleOpts{
				IPAddress:  amap["ip_address"].(string),
				AccessMode: amap["access_mode"].(string),
			})
		}
		opts.Access = access
	}

	if typeName == "vast" {
		var shareSettingsOpts file_shares.ShareSettingsOpts
		if len(shareSettingsList) > 0 {
			// schema restricts share_settings to a single block
			ssMap := shareSettingsList[0].(map[string]interface{})
			if v, ok := ssMap["allowed_characters"]; ok && v != nil && v.(string) != "" {
				shareSettingsOpts.AllowedCharacters = file_shares.FileShareAllowedCharacters(v.(string))
			}
			if v, ok := ssMap["path_length"]; ok && v != nil && v.(string) != "" {
				shareSettingsOpts.PathLength = file_shares.FileSharePathLength(v.(string))
			}
			if v, ok := ssMap["root_squash"]; ok && v != nil {
				rootSquash := v.(bool)
				shareSettingsOpts.RootSquash = &rootSquash
			}
			opts.ShareSettings = &shareSettingsOpts
		}
	}

	return &opts, nil
}
