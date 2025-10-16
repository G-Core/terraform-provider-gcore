package gcore

import (
	"context"
	"errors"
	"log"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/file_share/v1/file_shares"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFileShare() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFileShareRead,
		Description: "Get information about a file share (NFS) in Gcore Cloud.",
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:             schema.TypeInt,
				Optional:         true,
				ExactlyOneOf:     []string{"project_id", "project_name"},
				DiffSuppressFunc: suppressDiffProjectID,
				Description:      "Project ID, only one of project_id or project_name should be set",
			},
			"region_id": {
				Type:             schema.TypeInt,
				Optional:         true,
				ExactlyOneOf:     []string{"region_id", "region_name"},
				DiffSuppressFunc: suppressDiffRegionID,
				Description:      "Region ID, only one of region_id or region_name should be set",
			},
			"project_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"project_id", "project_name"},
				Description:  "Project name, only one of project_id or project_name should be set",
			},
			"region_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"region_id", "region_name"},
				Description:  "Region name, only one of region_id or region_name should be set",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the file share. It must be unique within the project and region.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protocol used by the file share. Currently, only 'NFS' is supported.",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the file share in GB. It must be a positive integer.",
			},
			// Deprecated and removed: volume_type. Use type_name instead.
			"connection_point": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The connection point of the file share.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the file share.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the file share in ISO 8601 format.",
			},
			"type_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the file share (standard or vast).",
			},
			"network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the network to which the file share will be connected.",
			},
			"network_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the network associated with the file share.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the subnet within the network. This is optional and can be used to specify a particular subnet for the file share.",
			},
			"subnet_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the subnet associated with the file share",
			},
			"share_network_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the share network associated with the file share. This is only applicable for 'standard'.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Tags associated with the file share. Tags are key-value pairs.",
			},
			"share_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Share settings for the file share.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The type of the file share (standard or vast).",
						},
						"allowed_characters": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Allowed characters in file names.",
						},
						"path_length": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Affects the maximum limit of file path component name length.",
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

func dataSourceFileShareRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FileShare data source reading")
	config := m.(*Config)
	provider := config.Provider

	projectID, regionID, err := getProjectAndRegionID(provider, d)
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("name").(string)

	client, err := CreateClient(provider, d, fileSharePoint, "v1")
	if err != nil {
		return diag.FromErr(err)
	}

	fileShares, err := file_shares.ListAll(client)
	if err != nil {
		return diag.FromErr(err)
	}

	var found bool
	var fs file_shares.FileShare
	for _, f := range fileShares {
		if f.Name == name && f.ProjectID == projectID && f.RegionID == regionID {
			fs = f
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("file share with name '%s' not found in project %d and region %d", name, projectID, regionID)
	}

	d.SetId(fs.ID)
	d.Set("name", fs.Name)
	d.Set("protocol", fs.Protocol)
	d.Set("size", fs.Size)
	// volume_type removed in API; expose type_name only
	d.Set("connection_point", fs.ConnectionPoint)
	d.Set("status", string(fs.Status))
	d.Set("created_at", fs.CreatedAt.String())
	d.Set("type_name", fs.TypeName)
	d.Set("network_id", fs.NetworkID)
	d.Set("network_name", fs.NetworkName)
	d.Set("subnet_id", fs.SubnetID)
	d.Set("subnet_name", fs.SubnetName)
	d.Set("region_name", fs.Region)
	if fs.ShareNetworkName != nil {
		d.Set("share_network_name", *fs.ShareNetworkName)
	}
	if fs.Tags != nil {
		tags := make(map[string]string)
		for _, tag := range fs.Tags {
			if !tag.ReadOnly {
				tags[tag.Key] = tag.Value
			}
		}
		d.Set("tags", tags)
	}
	shareSettingsMap := map[string]interface{}{
		"type_name": fs.ShareSettings.TypeName,
	}
	if fs.ShareSettings.AllowedCharacters != nil && fs.ShareSettings.AllowedCharacters.String() != "" {
		shareSettingsMap["allowed_characters"] = fs.ShareSettings.AllowedCharacters.String()
	}
	if fs.ShareSettings.PathLength != nil && fs.ShareSettings.PathLength.String() != "" {
		shareSettingsMap["path_length"] = fs.ShareSettings.PathLength.String()
	}
	if fs.ShareSettings.RootSquash != nil {
		shareSettingsMap["root_squash"] = *fs.ShareSettings.RootSquash
	}
	d.Set("share_settings", []interface{}{shareSettingsMap})
	return nil
}

// getProjectAndRegionID is a helper to resolve project/region from id or name
func getProjectAndRegionID(provider *gcorecloud.ProviderClient, d *schema.ResourceData) (int, int, error) {
	var projectID, regionID int
	var err error
	if v, ok := d.GetOk("project_id"); ok {
		projectID = v.(int)
	} else if v, ok := d.GetOk("project_name"); ok {
		projectID, err = GetProject(provider, 0, v.(string))
		if err != nil {
			return 0, 0, err
		}
	} else {
		return 0, 0, errors.New("project_id or project_name must be set")
	}
	if v, ok := d.GetOk("region_id"); ok {
		regionID = v.(int)
	} else if v, ok := d.GetOk("region_name"); ok {
		regionID, err = GetRegion(provider, 0, v.(string))
		if err != nil {
			return 0, 0, err
		}
	} else {
		return 0, 0, errors.New("region_id or region_name must be set")
	}
	return projectID, regionID, nil
}
