package gcore

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	clientUtils "github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata/v1/metadata"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/floatingip/v1/floatingips"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	floatingIPsPoint        = "floatingips"
	FloatingIPCreateTimeout = 1200
)

func resourceFloatingIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFloatingIPCreate,
		ReadContext:   resourceFloatingIPRead,
		UpdateContext: resourceFloatingIPUpdate,
		DeleteContext: resourceFloatingIPDelete,
		Description:   "A floating IP is a static IP address that points to one of your Instances. It allows you to redirect network traffic to any of your Instances in the same datacenter.",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, fipID, err := ImportStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(fipID)

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
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"fixed_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateDiagFunc: func(val interface{}, key cty.Path) diag.Diagnostics {
					v := val.(string)
					ip := net.ParseIP(v)
					if ip != nil {
						return diag.Diagnostics{}
					}

					return diag.FromErr(fmt.Errorf("%q must be a valid ip, got: %s", key, v))
				},
			},
			"floating_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"router_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"updated_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Tags to associate with the floating IP. Tags are key-value pairs.",
			},
		},
	}
}

func resourceFloatingIPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FloatingIP creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, floatingIPsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := floatingips.CreateOpts{
		PortID:         d.Get("port_id").(string),
		FixedIPAddress: net.ParseIP(d.Get("fixed_ip_address").(string)),
	}

	metadataRaw, metadataOk := d.GetOk("metadata_map")
	tagsRaw, tagsOk := d.GetOk("tags")

	if metadataOk && tagsOk {
		return diag.FromErr(fmt.Errorf("only one of metadata_map or tags can be set"))
	}

	if metadataOk {
		meta, err := utils.MapInterfaceToMapString(metadataRaw)
		if err != nil {
			return diag.FromErr(err)
		}
		opts.Metadata = meta
	}

	if tagsOk {
		tags, err := utils.MapInterfaceToMapString(tagsRaw)
		if err != nil {
			return diag.FromErr(err)
		}
		opts.Metadata = tags
	}

	results, err := floatingips.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	floatingIPID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, FloatingIPCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		floatingIPID, err := floatingips.ExtractFloatingIPIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve FloatingIP ID from task info: %w", err)
		}
		return floatingIPID, nil
	})

	log.Printf("[DEBUG] FloatingIP id (%s)", floatingIPID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(floatingIPID.(string))
	resourceFloatingIPRead(ctx, d, m)

	log.Printf("[DEBUG] Finish FloatingIP creating (%s)", floatingIPID)
	return diags
}

func resourceFloatingIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FloatingIP reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, floatingIPsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	floatingIP, err := floatingips.Get(client, d.Id()).Extract()
	if err != nil {
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			log.Printf("[WARN] Removing floating ip %s because resource doesn't exist anymore", d.Id())
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	if floatingIP.FixedIPAddress != nil {
		d.Set("fixed_ip_address", floatingIP.FixedIPAddress.String())
	} else {
		d.Set("fixed_ip_address", "")
	}

	d.Set("project_id", floatingIP.ProjectID)
	d.Set("region_id", floatingIP.RegionID)
	d.Set("status", floatingIP.Status)
	d.Set("port_id", floatingIP.PortID)
	d.Set("router_id", floatingIP.RouterID)
	d.Set("floating_ip_address", floatingIP.FloatingIPAddress.String())

	metadataMap := make(map[string]string)
	metadataReadOnly := make([]map[string]interface{}, 0, len(floatingIP.Metadata))

	if len(floatingIP.Metadata) > 0 {
		for _, metadataItem := range floatingIP.Metadata {
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

	log.Println("[DEBUG] Finish FloatingIP reading")
	return diags
}

func resourceFloatingIPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FloatingIP updating")
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, floatingIPsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("fixed_ip_address", "port_id") {
		oldFixedIP, newFixedIP := d.GetChange("fixed_ip_address")
		oldPortID, newPortID := d.GetChange("port_id")
		if oldPortID.(string) != "" || oldFixedIP.(string) != "" {
			_, err := floatingips.UnAssign(client, d.Id()).Extract()
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if newPortID.(string) != "" || newFixedIP.(string) != "" {
			opts := floatingips.CreateOpts{
				PortID:         d.Get("port_id").(string),
				FixedIPAddress: net.ParseIP(d.Get("fixed_ip_address").(string)),
			}

			_, err = floatingips.Assign(client, d.Id(), opts).Extract()
			if err != nil {
				return diag.FromErr(err)
			}
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	if d.HasChange("metadata_map") {
		_, nmd := d.GetChange("metadata_map")

		meta, err := utils.MapInterfaceToMapString(nmd.(map[string]interface{}))
		if err != nil {
			return diag.Errorf("cannot get metadata. Error: %s", err)
		}

		err = metadata.MetadataReplace(client, d.Id(), meta).Err
		if err != nil {
			return diag.Errorf("cannot update metadata. Error: %s", err)
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		newTagsMap := make(map[string]*string)
		if newTags != nil {
			for k, val := range newTags.(map[string]interface{}) {
				newTagsMap[k] = clientUtils.StringToPointer(val.(string))
			}
		}

		oldTagsMap := make(map[string]string)
		if oldTags != nil {
			for k, val := range oldTags.(map[string]interface{}) {
				oldTagsMap[k] = val.(string)
			}
		}

		for oldKey := range oldTagsMap {
			if _, exists := newTagsMap[oldKey]; !exists {
				newTagsMap[oldKey] = nil
			}
		}

		_, err := floatingips.Update(client, d.Id(), floatingips.UpdateOpts{
			Tags: newTagsMap,
		}).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceFloatingIPRead(ctx, d, m)
}

func resourceFloatingIPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FloatingIP deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, floatingIPsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	results, err := floatingips.Delete(client, id).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, FloatingIPCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := floatingips.Get(client, id).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete floating ip with ID: %s", id)
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
	log.Printf("[DEBUG] Finish of FloatingIP deleting")
	return diags
}
