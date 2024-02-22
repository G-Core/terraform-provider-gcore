package gcore

import (
	"context"
	"fmt"
	"log"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata/v1/metadata"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const networkDeleting int = 1200
const networkCreatingTimeout int = 1200
const networksPoint = "networks"
const sharedNetworksPoint = "availablenetworks"

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Description:   "Represent network. A network is a software-defined network in a cloud computing infrastructure",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, NetworkID, err := ImportStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(NetworkID)

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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"mtu": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "'vlan' or 'vxlan' network type is allowed. Default value is 'vxlan'",
			},
			"create_router": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Create external router to the network, default true",
			},
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata_map": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
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
		},
	}
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Network creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, networksPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := networks.CreateOpts{
		Name:         d.Get("name").(string),
		Type:         d.Get("type").(string),
		CreateRouter: d.Get("create_router").(bool),
	}

	if metadataRaw, ok := d.GetOk("metadata_map"); ok {
		meta, err := utils.MapInterfaceToMapString(metadataRaw)
		if err != nil {
			return diag.FromErr(err)
		}

		createOpts.Metadata = meta
	}

	log.Printf("Create network ops: %+v", createOpts)
	results, err := networks.Create(client, createOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	networkID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, networkCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		NetworkID, err := networks.ExtractNetworkIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve Network ID from task info: %w", err)
		}
		return NetworkID, nil
	},
	)
	log.Printf("[DEBUG] Network id (%s)", networkID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(networkID.(string))
	resourceNetworkRead(ctx, d, m)

	log.Printf("[DEBUG] Finish Network creating (%s)", networkID)
	return diags
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start network reading")
	log.Printf("[DEBUG] Start network reading%s", d.State())
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	networkID := d.Id()
	log.Printf("[DEBUG] Network id = %s", networkID)

	client, err := CreateClient(provider, d, networksPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := networks.Get(client, networkID).Extract()
	if err != nil {
		return diag.Errorf("cannot get network with ID: %s. Error: %s", networkID, err)
	}

	d.Set("name", network.Name)
	d.Set("mtu", network.MTU)
	d.Set("type", network.Type)
	d.Set("region_id", network.RegionID)
	d.Set("project_id", network.ProjectID)

	metadataMap := make(map[string]string)
	metadataReadOnly := make([]map[string]interface{}, 0, len(network.Metadata))

	if len(network.Metadata) > 0 {
		for _, metadataItem := range network.Metadata {
			if !metadataItem.ReadOnly {
				metadataMap[metadataItem.Key] = metadataItem.Value
			}
			metadataReadOnly = append(metadataReadOnly, map[string]interface{}{
				"key":       metadataItem.Key,
				"value":     metadataItem.Value,
				"read_only": metadataItem.ReadOnly,
			})
		}
	}

	if err := d.Set("metadata_map", metadataMap); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata_read_only", metadataReadOnly); err != nil {
		return diag.FromErr(err)
	}

	fields := []string{"create_router"}
	revertState(d, &fields)

	log.Println("[DEBUG] Finish network reading")
	return diags
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start network updating")
	networkID := d.Id()
	log.Printf("[DEBUG] Volume id = %s", networkID)
	config := m.(*Config)
	provider := config.Provider
	client, err := CreateClient(provider, d, networksPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newName := d.Get("name").(string)
		_, err := networks.Update(client, networkID, networks.UpdateOpts{Name: newName}).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("metadata_map") {
		_, nmd := d.GetChange("metadata_map")

		meta, err := utils.MapInterfaceToMapString(nmd.(map[string]interface{}))
		if err != nil {
			return diag.Errorf("cannot get metadata. Error: %s", err)
		}

		err = metadata.MetadataReplace(client, networkID, meta).Err
		if err != nil {
			return diag.Errorf("cannot update metadata. Error: %s", err)
		}
	}
	d.Set("last_updated", time.Now().Format(time.RFC850))
	log.Println("[DEBUG] Finish network updating")
	return resourceNetworkRead(ctx, d, m)
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start network deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	networkID := d.Id()
	log.Printf("[DEBUG] Network id = %s", networkID)

	client, err := CreateClient(provider, d, networksPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	results, err := networks.Delete(client, networkID).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, networkDeleting, func(task tasks.TaskID) (interface{}, error) {
		_, err := networks.Get(client, networkID).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete network with ID: %s", networkID)
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
	log.Printf("[DEBUG] Finish of network deleting")
	return diags
}
