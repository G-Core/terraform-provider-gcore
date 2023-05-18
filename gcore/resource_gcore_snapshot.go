package gcore

import (
	"context"
	"fmt"
	"log"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/snapshot/v1/snapshots"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const snapshotDeleting int = 1200
const snapshotCreatingTimeout int = 1200
const snapshotsPoint = "snapshots"

func resourceSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSnapshotCreate,
		ReadContext:   resourceSnapshotRead,
		UpdateContext: resourceSnapshotUpdate,
		DeleteContext: resourceSnapshotDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, snapshotID, err := ImportStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(snapshotID)

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
				ForceNew: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSnapshotCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start snapshot creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, snapshotsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := getSnapshotData(d)
	results, err := snapshots.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	SnapshotID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, snapshotCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		snapshotID, err := snapshots.ExtractSnapshotIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve snapshot ID from task info: %w", err)
		}
		return snapshotID, nil
	},
	)
	log.Printf("[DEBUG] Snapshot id (%s)", SnapshotID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(SnapshotID.(string))
	resourceSnapshotRead(ctx, d, m)

	log.Printf("[DEBUG] Finish snapshot creating (%s)", SnapshotID)
	return diags
}

func resourceSnapshotRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start snapshot reading")
	log.Printf("[DEBUG] Start snapshot reading %s", d.State())
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	snapshotID := d.Id()
	log.Printf("[DEBUG] Snapshot id = %s", snapshotID)

	client, err := CreateClient(provider, d, snapshotsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	snapshot, err := snapshots.Get(client, snapshotID).Extract()
	if err != nil {
		return diag.Errorf("cannot get snapshot with ID: %s. Error: %s", snapshotID, err)
	}

	d.Set("name", snapshot.Name)
	d.Set("description", snapshot.Description)
	d.Set("status", snapshot.Status)
	d.Set("size", snapshot.Size)
	d.Set("volume_id", snapshot.VolumeID)
	d.Set("region_id", snapshot.RegionID)
	d.Set("project_id", snapshot.ProjectID)
	if err := d.Set("metadata", snapshot.Metadata); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish snapshot reading")
	return diags
}

func resourceSnapshotUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start snapshot updating")
	snapshotID := d.Id()
	if d.HasChange("metadata") {
		config := m.(*Config)
		provider := config.Provider
		client, err := CreateClient(provider, d, snapshotsPoint, versionPointV1)
		if err != nil {
			return diag.FromErr(err)
		}

		newMeta := prepareRawMetadata(d.Get("metadata").(map[string]interface{}))
		metadata := make([]snapshots.MetadataOpts, 0, len(newMeta))
		for k, v := range newMeta {
			metadata = append(metadata, snapshots.MetadataOpts{Key: k, Value: v})
		}
		opts := snapshots.MetadataSetOpts{Metadata: metadata}
		if _, err := snapshots.MetadataReplace(client, snapshotID, opts).Extract(); err != nil {
			return diag.FromErr(err)
		}
	}
	d.Set("last_updated", time.Now().Format(time.RFC850))
	log.Println("[DEBUG] Finish snapshot updating")
	return resourceSnapshotRead(ctx, d, m)
}

func resourceSnapshotDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start snapshot deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	snapshotID := d.Id()
	log.Printf("[DEBUG] Snapshot id = %s", snapshotID)

	client, err := CreateClient(provider, d, snapshotsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	results, err := snapshots.Delete(client, snapshotID).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, snapshotDeleting, func(task tasks.TaskID) (interface{}, error) {
		_, err := snapshots.Get(client, snapshotID).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete snapshot with ID: %s", snapshotID)
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
	log.Printf("[DEBUG] Finish of snapshot deleting")
	return diags
}

func getSnapshotData(d *schema.ResourceData) *snapshots.CreateOpts {
	snapshotData := snapshots.CreateOpts{}
	snapshotData.Name = d.Get("name").(string)
	snapshotData.VolumeID = d.Get("volume_id").(string)
	snapshotData.Description = d.Get("description").(string)
	metadataRaw := d.Get("metadata").(map[string]interface{})
	if len(metadataRaw) > 0 {
		snapshotData.Metadata = prepareRawMetadata(metadataRaw)
	}

	return &snapshotData
}

func prepareRawMetadata(raw map[string]interface{}) map[string]string {
	meta := make(map[string]string, len(raw))
	for k, v := range raw {
		meta[k] = v.(string)
	}
	return meta
}
