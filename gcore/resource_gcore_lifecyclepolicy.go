package gcore

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/G-Core/gcorelabscloud-go/gcore/lifecyclepolicy/v1/lifecyclepolicy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	lifecyclePolicyPoint = "lifecycle_policy"
	// Maybe move to utils and use for other resources
	nameRegexString = `^[a-zA-Z0-9][a-zA-Z 0-9._\-]{1,61}[a-zA-Z0-9._]$`
)

var (
	// Maybe move to utils and use for other resources
	nameRegex = regexp.MustCompile(nameRegexString)
)

func resourceLifecyclePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLifecyclePolicyCreate,
		ReadContext:   resourceLifecyclePolicyRead,
		UpdateContext: resourceLifecyclePolicyUpdate,
		DeleteContext: resourceLifecyclePolicyDelete,
		Description:   "Represent lifecycle policy. Use to periodically take snapshots",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, lcpID, err := ImportStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(lcpID)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				ExactlyOneOf:     []string{"project_id", "project_name"},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"region_id": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				ExactlyOneOf:     []string{"region_id", "region_name"},
				DiffSuppressFunc: suppressDiffRegionID,
			},
			"project_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"project_id", "project_name"},
			},
			"region_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"region_id", "region_name"},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(nameRegex, ""),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      lifecyclepolicy.PolicyStatusActive.String(),
				ValidateFunc: validation.StringInSlice(lifecyclepolicy.PolicyStatus("").StringList(), false),
			},
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      lifecyclepolicy.PolicyActionVolumeSnapshot.String(),
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(lifecyclepolicy.PolicyAction("").StringList(), false),
			},
			"volume": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of managed volumes",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"schedule": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_quantity": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 10000),
							Description:  "Maximum number of stored resources",
						},
						"interval": {
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Description: "Use for taking actions with equal time intervals between them. Exactly one of interval and cron blocks should be provided",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"weeks": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     0,
										Description: intervalScheduleParamDescription("week"),
									},
									"days": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     0,
										Description: intervalScheduleParamDescription("day"),
									},
									"hours": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     0,
										Description: intervalScheduleParamDescription("hour"),
									},
									"minutes": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     0,
										Description: intervalScheduleParamDescription("minute"),
									},
								},
							},
							Optional: true,
						},
						"cron": {
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Description: "Use for taking actions at specified moments of time. Exactly one of interval and cron blocks should be provided",
							Elem: &schema.Resource{ // TODO: validate?
								Schema: map[string]*schema.Schema{
									"timezone": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "UTC",
									},
									"month": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "*",
										Description: cronScheduleParamDescription(1, 12),
									},
									"week": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "*",
										Description: cronScheduleParamDescription(1, 53),
									},
									"day": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "*",
										Description: cronScheduleParamDescription(1, 31),
									},
									"day_of_week": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "*",
										Description: cronScheduleParamDescription(0, 6),
									},
									"hour": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "*",
										Description: cronScheduleParamDescription(0, 23),
									},
									"minute": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "0",
										Description: cronScheduleParamDescription(0, 59),
									},
								},
							},
							Optional: true,
						},
						"resource_name_template": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "reserve snap of the volume {volume_id}",
							Description: "Used to name snapshots. {volume_id} is substituted with volume.id on creation",
						},
						"retention_time": {
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Description: "If it is set, new resource will be deleted after time",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"weeks": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     0,
										Description: retentionTimerParamDescription("week"),
									},
									"days": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     0,
										Description: retentionTimerParamDescription("day"),
									},
									"hours": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     0,
										Description: retentionTimerParamDescription("hour"),
									},
									"minutes": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     0,
										Description: retentionTimerParamDescription("minute"),
									},
								},
							},
							Optional: true,
						},
						"id": {
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
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceLifecyclePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := CreateClient(m.(*Config).Provider, d, lifecyclePolicyPoint, versionPointV1)
	if err != nil {
		return diag.Errorf("Error creating client: %s", err)
	}

	log.Printf("[DEBUG] Start of LifecyclePolicy creating")
	opts, err := buildLifecyclePolicyCreateOpts(d)
	if err != nil {
		return diag.FromErr(err)
	}
	policy, err := lifecyclepolicy.Create(client, *opts).Extract()
	if err != nil {
		return diag.Errorf("Error creating lifecycle policy: %s", err)
	}
	d.SetId(strconv.Itoa(policy.ID))
	log.Printf("[DEBUG] Finish of LifecyclePolicy %s creating", d.Id())
	return resourceLifecyclePolicyRead(ctx, d, m)
}

func resourceLifecyclePolicyRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := CreateClient(m.(*Config).Provider, d, lifecyclePolicyPoint, versionPointV1)
	if err != nil {
		return diag.Errorf("Error creating client: %s", err)
	}
	id := d.Id()
	integerId, err := strconv.Atoi(id)
	if err != nil {
		return diag.Errorf("Error converting lifecycle policy ID to integer: %s", err)
	}

	log.Printf("[DEBUG] Start of LifecyclePolicy %s reading", id)
	policy, err := lifecyclepolicy.Get(client, integerId, lifecyclepolicy.GetOpts{NeedVolumes: true}).Extract()
	if err != nil {
		return diag.Errorf("Error getting lifecycle policy: %s", err)
	}

	_ = d.Set("name", policy.Name)
	_ = d.Set("status", policy.Status)
	_ = d.Set("action", policy.Action)
	_ = d.Set("user_id", policy.UserID)
	if err = d.Set("volume", flattenVolumes(policy.Volumes)); err != nil {
		return diag.Errorf("error setting lifecycle policy volumes: %s", err)
	}
	if err = d.Set("schedule", flattenSchedules(policy.Schedules)); err != nil {
		return diag.Errorf("error setting lifecycle policy schedules: %s", err)
	}

	log.Printf("[DEBUG] Finish of LifecyclePolicy %s reading", id)
	return nil
}

func resourceLifecyclePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := CreateClient(m.(*Config).Provider, d, lifecyclePolicyPoint, versionPointV1)
	if err != nil {
		return diag.Errorf("Error creating client: %s", err)
	}
	id := d.Id()
	integerId, err := strconv.Atoi(id)
	if err != nil {
		return diag.Errorf("Error converting lifecycle policy ID to integer: %s", err)
	}

	log.Printf("[DEBUG] Start of LifecyclePolicy updating")
	_, err = lifecyclepolicy.Update(client, integerId, buildLifecyclePolicyUpdateOpts(d)).Extract()
	if err != nil {
		return diag.Errorf("Error updating lifecycle policy: %s", err)
	}

	if d.HasChange("volume") {
		oldVolumes, newVolumes := d.GetChange("volume")
		toRemove, toAdd := volumeSymmetricDifference(oldVolumes.(*schema.Set), newVolumes.(*schema.Set))
		_, err = lifecyclepolicy.RemoveVolumes(client, integerId, lifecyclepolicy.RemoveVolumesOpts{VolumeIds: toRemove}).Extract()
		if err != nil {
			return diag.Errorf("Error removing volumes from lifecycle policy: %s", err)
		}
		_, err = lifecyclepolicy.AddVolumes(client, integerId, lifecyclepolicy.AddVolumesOpts{VolumeIds: toAdd}).Extract()
		if err != nil {
			return diag.Errorf("Error adding volumes to lifecycle policy: %s", err)
		}
	}
	log.Printf("[DEBUG] Finish of LifecyclePolicy %v updating", integerId)
	return resourceLifecyclePolicyRead(ctx, d, m)
}

func resourceLifecyclePolicyDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := CreateClient(m.(*Config).Provider, d, lifecyclePolicyPoint, versionPointV1)
	if err != nil {
		return diag.Errorf("Error creating client: %s", err)
	}
	id := d.Id()
	integerId, err := strconv.Atoi(id)
	if err != nil {
		return diag.Errorf("Error converting lifecycle policy ID to integer: %s", err)
	}

	log.Printf("[DEBUG] Start of LifecyclePolicy %s deleting", id)
	err = lifecyclepolicy.Delete(client, integerId)
	if err != nil {
		return diag.Errorf("Error deleting lifecycle policy: %s", err)
	}
	d.SetId("")
	log.Printf("[DEBUG] Finish of LifecyclePolicy %s deleting", id)
	return nil
}

func expandIntervalSchedule(flat map[string]interface{}) *lifecyclepolicy.CreateIntervalScheduleOpts {
	return &lifecyclepolicy.CreateIntervalScheduleOpts{
		Weeks:   flat["weeks"].(int),
		Days:    flat["days"].(int),
		Hours:   flat["hours"].(int),
		Minutes: flat["minutes"].(int),
	}
}

func expandCronSchedule(flat map[string]interface{}) *lifecyclepolicy.CreateCronScheduleOpts {
	return &lifecyclepolicy.CreateCronScheduleOpts{
		Timezone:  flat["timezone"].(string),
		Week:      flat["week"].(string),
		DayOfWeek: flat["day_of_week"].(string),
		Month:     flat["month"].(string),
		Day:       flat["day"].(string),
		Hour:      flat["hour"].(string),
		Minute:    flat["minute"].(string),
	}
}

func expandRetentionTimer(flat []interface{}) *lifecyclepolicy.RetentionTimer {
	if len(flat) > 0 {
		rawRetention := flat[0].(map[string]interface{})
		return &lifecyclepolicy.RetentionTimer{
			Weeks:   rawRetention["weeks"].(int),
			Days:    rawRetention["days"].(int),
			Hours:   rawRetention["hours"].(int),
			Minutes: rawRetention["minutes"].(int),
		}
	}
	return nil
}

func expandSchedule(flat map[string]interface{}) (expanded lifecyclepolicy.CreateScheduleOpts, err error) {
	t := lifecyclepolicy.ScheduleType("")
	intervalSlice := flat["interval"].([]interface{})
	cronSlice := flat["cron"].([]interface{})
	if len(intervalSlice)+len(cronSlice) != 1 {
		return nil, fmt.Errorf("exactly one of interval and cron blocks should be provided")
	}
	if len(intervalSlice) > 0 {
		t = lifecyclepolicy.ScheduleTypeInterval
		expanded = expandIntervalSchedule(intervalSlice[0].(map[string]interface{}))
	} else {
		t = lifecyclepolicy.ScheduleTypeCron
		expanded = expandCronSchedule(cronSlice[0].(map[string]interface{}))
	}
	expanded.SetCommonCreateScheduleOpts(lifecyclepolicy.CommonCreateScheduleOpts{
		Type:                 t,
		ResourceNameTemplate: flat["resource_name_template"].(string),
		MaxQuantity:          flat["max_quantity"].(int),
		RetentionTime:        expandRetentionTimer(flat["retention_time"].([]interface{})),
	})
	return
}

func expandSchedules(flat []interface{}) ([]lifecyclepolicy.CreateScheduleOpts, error) {
	expanded := make([]lifecyclepolicy.CreateScheduleOpts, len(flat))
	for i, x := range flat {
		exp, err := expandSchedule(x.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		expanded[i] = exp
	}
	return expanded, nil
}

func expandVolumeIds(flat []interface{}) []string {
	expanded := make([]string, len(flat))
	for i, x := range flat {
		expanded[i] = x.(map[string]interface{})["id"].(string)
	}
	return expanded
}

func buildLifecyclePolicyCreateOpts(d *schema.ResourceData) (*lifecyclepolicy.CreateOpts, error) {
	schedules, err := expandSchedules(d.Get("schedule").([]interface{}))
	if err != nil {
		return nil, err
	}
	opts := &lifecyclepolicy.CreateOpts{
		Name:      d.Get("name").(string),
		Status:    lifecyclepolicy.PolicyStatus(d.Get("status").(string)),
		Schedules: schedules,
		VolumeIds: expandVolumeIds(d.Get("volume").(*schema.Set).List()),
	}

	// Action is required field from API point of view, but optional for us
	if action, ok := d.GetOk("action"); ok {
		opts.Action = lifecyclepolicy.PolicyAction(action.(string))
	} else {
		opts.Action = lifecyclepolicy.PolicyActionVolumeSnapshot
	}
	return opts, nil
}

func volumeSymmetricDifference(oldVolumes, newVolumes *schema.Set) ([]string, []string) {
	toRemove := make([]string, 0)
	for _, v := range oldVolumes.List() {
		if !newVolumes.Contains(v) {
			toRemove = append(toRemove, v.(map[string]interface{})["id"].(string))
		}
	}
	toAdd := make([]string, 0)
	for _, v := range newVolumes.List() {
		if !oldVolumes.Contains(v) {
			toAdd = append(toAdd, v.(map[string]interface{})["id"].(string))
		}
	}
	return toRemove, toAdd
}

func buildLifecyclePolicyUpdateOpts(d *schema.ResourceData) lifecyclepolicy.UpdateOpts {
	opts := lifecyclepolicy.UpdateOpts{
		Name:   d.Get("name").(string),
		Status: lifecyclepolicy.PolicyStatus(d.Get("status").(string)),
	}
	return opts
}

func flattenIntervalSchedule(expanded lifecyclepolicy.IntervalSchedule) interface{} {
	return []map[string]int{{
		"weeks":   expanded.Weeks,
		"days":    expanded.Days,
		"hours":   expanded.Hours,
		"minutes": expanded.Minutes,
	}}
}

func flattenCronSchedule(expanded lifecyclepolicy.CronSchedule) interface{} {
	return []map[string]string{{
		"timezone":    expanded.Timezone,
		"week":        expanded.Week,
		"day_of_week": expanded.DayOfWeek,
		"month":       expanded.Month,
		"day":         expanded.Day,
		"hour":        expanded.Hour,
		"minute":      expanded.Minute,
	}}
}

func flattenRetentionTimer(expanded *lifecyclepolicy.RetentionTimer) interface{} {
	if expanded != nil {
		return []map[string]int{{
			"weeks":   expanded.Weeks,
			"days":    expanded.Days,
			"hours":   expanded.Hours,
			"minutes": expanded.Minutes,
		}}
	}
	return []interface{}{}
}

func flattenSchedule(expanded lifecyclepolicy.Schedule) map[string]interface{} {
	common := expanded.GetCommonSchedule()
	flat := map[string]interface{}{
		"max_quantity":           common.MaxQuantity,
		"resource_name_template": common.ResourceNameTemplate,
		"retention_time":         flattenRetentionTimer(common.RetentionTime),
		"id":                     common.ID,
		"type":                   common.Type,
	}
	switch common.Type {
	case lifecyclepolicy.ScheduleTypeInterval:
		flat["interval"] = flattenIntervalSchedule(expanded.(lifecyclepolicy.IntervalSchedule))
	case lifecyclepolicy.ScheduleTypeCron:
		flat["cron"] = flattenCronSchedule(expanded.(lifecyclepolicy.CronSchedule))
	}
	return flat
}

func flattenSchedules(expanded []lifecyclepolicy.Schedule) []map[string]interface{} {
	flat := make([]map[string]interface{}, len(expanded))
	for i, x := range expanded {
		flat[i] = flattenSchedule(x)
	}
	return flat
}

func flattenVolumes(expanded []lifecyclepolicy.Volume) []map[string]string {
	flat := make([]map[string]string, len(expanded))
	for i, volume := range expanded {
		flat[i] = map[string]string{"id": volume.ID, "name": volume.Name}
	}
	return flat
}

func cronScheduleParamDescription(min, max int) string {
	return fmt.Sprintf("Either single asterisk or comma-separated list of integers (%v-%v)", min, max)
}

func intervalScheduleParamDescription(unit string) string {
	return fmt.Sprintf("Number of %ss to wait between actions", unit)
}

func retentionTimerParamDescription(unit string) string {
	return fmt.Sprintf("Number of %ss to wait before deleting snapshot", unit)
}
