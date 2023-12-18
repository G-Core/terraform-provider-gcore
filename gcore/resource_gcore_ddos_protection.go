package gcore

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/ddos/v1/ddos"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ddosProfileCreatingTimeout int = 1200
const ddosProfileDeletingTimeout int = 1200
const ddosProfileUpdatingTimeout int = 1200
const ddosProfilePoint = "ddos/profiles"

func resourceDDoSProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDDoSProtectionCreate,
		ReadContext:   resourceDDoSProtectionRead,
		UpdateContext: resourceDDoSProtectionUpdate,
		DeleteContext: resourceDDoSProtectionDelete,
		Description:   "Represents DDoS protection profile",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, profileID, err := ImportStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(profileID)

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
			"ip_address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IP address",
			},
			"profile_template": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Profile template ID",
			},
			"site": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Activate profile",
				Optional:    true,
				Default:     true,
			},
			"bgp": {
				Type:        schema.TypeBool,
				Description: "Activate BGP protocol",
				Optional:    true,
				Default:     true,
			},
			"fields": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"field_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"base_field": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Basic type value. Only one of 'value' or 'field_value' must be specified.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Field description",
						},
						"default": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"required": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"field_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Complex value. Only one of 'value' or 'field_value' must be specified.",
						},
						"validation_schema": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Json schema to validate field_values",
						},
					},
				},
			},
			"bm_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocols": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of protocols",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocols": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
					},
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

func resourceDDoSProtectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start DDoS protection profile creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, ddosProfilePoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := ddos.CreateProfileOpts{}

	createOpts.IPAddress = d.Get("ip_address").(string)
	createOpts.ProfileTemplate = d.Get("profile_template").(int)
	createOpts.BaremetalInstanceID = d.Get("bm_instance_id").(string)
	fields := d.Get("fields")
	if len(fields.([]interface{})) > 0 {
		createOpts.Fields, err = extractProfileFieldsMap(fields.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		createOpts.Fields = make([]ddos.ProfileField, 0)
	}

	results, err := ddos.CreateProfile(client, createOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	profileID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, ddosProfileCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		Profile, err := ddos.ExtractProfileIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrive DDoS protection profile ID from task info %w", err)
		}

		return Profile, nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	var options ddos.ActivateProfileOpts
	options.Active = d.Get("active").(bool)
	options.BGP = d.Get("bgp").(bool)
	if !options.Active || !options.BGP {
		id, _ := strconv.Atoi(profileID.(string))
		if results, err = ddos.ActivateProfile(client, id, options).Extract(); err != nil {
			return diag.FromErr(err)
		}

		taskID = results.Tasks[0]
		if err = tasks.WaitForStatus(client, string(taskID),
			tasks.TaskStateFinished, ddosProfileCreatingTimeout, true); err != nil {
			log.Printf("[DEBUG] failed to activate/disactivate profile: %v", err)
		}
	}

	if err != nil {
		id, _ := strconv.Atoi(profileID.(string))
		results, err = ddos.DeleteProfile(client, id).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		taskID = results.Tasks[0]
		if err = tasks.WaitForStatus(client, string(taskID),
			tasks.TaskStateFinished, ddosProfileCreatingTimeout, true); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[DEBUG] DDoS protection profile id (%s)", profileID)
	d.SetId(profileID.(string))
	resourceDDoSProtectionRead(ctx, d, m)

	log.Println("[DEBUG] Finish DDoS protection profile creating")
	return diags
}

func resourceDDoSProtectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start DDoS protection profile reading")
	log.Printf("[DEBUG] Start DDoS protection profile reading %s", d.State())
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	profileID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] DDoS profile id = %d", profileID)

	client, err := CreateClient(provider, d, ddosProfilePoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	profiles, err := ddos.ListAllProfiles(client)
	if err != nil {
		return diag.FromErr(err)
	}

	var profile ddos.Profile
	var found bool
	for _, f := range profiles {
		if f.ID == profileID {
			profile = f
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("DDoS protection profile id = %d not found", profileID)
	}

	d.Set("ip_address", profile.IPAddress)
	d.Set("profile_template", profile.ProfileTemplate)
	d.Set("bgp", profile.Options.BGP)
	d.Set("active", profile.Options.Active)
	d.Set("site", profile.Site)
	fields := make([]map[string]interface{}, len(profile.Fields))
	for i, f := range profile.Fields {
		field := map[string]interface{}{
			"id":                f.ID,
			"name":              f.Name,
			"field_type":        f.FieldType,
			"base_field":        f.BaseField,
			"value":             f.Value,
			"description":       f.Description,
			"default":           f.Default,
			"required":          f.Required,
			"field_value":       string(f.FieldValue),
			"validation_schema": string(f.ValidationSchema),
		}

		fields[i] = field
	}
	if err := d.Set("fields", fields); err != nil {
		return diag.FromErr(err)
	}

	protocols := make([]map[string]interface{}, len(profile.Protocols))
	for i, p := range profile.Protocols {
		protocol := map[string]interface{}{
			"port":      p.Port,
			"protocols": p.Protocols,
		}

		protocols[i] = protocol
	}
	if err := d.Set("protocols", protocols); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish DDoS protection profile reading")
	return diags
}

func resourceDDoSProtectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start DDoS protection profile updating")
	profileID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] DDoS protection profile id = %d", profileID)
	config := m.(*Config)
	provider := config.Provider
	client, err := CreateClient(provider, d, ddosProfilePoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	updateOpts := ddos.UpdateProfileOpts{}

	var profileChanged bool

	updateOpts.ProfileTemplate = d.Get("profile_template").(int)
	if d.HasChange("profile_template") {
		profileChanged = true
	}

	updateOpts.BaremetalInstanceID = d.Get("bm_instance_id").(string)
	if d.HasChange("bm_instance_id") {
		profileChanged = true
	}

	updateOpts.IPAddress = d.Get("ip_address").(string)
	if d.HasChange("ip_address") {
		profileChanged = true
	}

	fs := d.Get("fields")
	updateOpts.Fields = make([]ddos.ProfileField, 0)
	if len(fs.([]interface{})) > 0 {
		fields, err := extractProfileFieldsMap(fs.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateOpts.Fields = fields
	}
	if d.HasChange("fields") {
		profileChanged = true
	}

	var (
		activateOpts   ddos.ActivateProfileOpts
		optionsChanged bool
	)

	if d.HasChange("bgp") || d.HasChange("active") {
		optionsChanged = true
		activateOpts.BGP = d.Get("bgp").(bool)
		activateOpts.Active = d.Get("active").(bool)
	}

	if profileChanged {
		log.Println("[DEBUG] Profile changed: updating profile")
		results, err := ddos.UpdateProfile(client, profileID, updateOpts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		taskID := results.Tasks[0]
		log.Printf("[DEBUG] Task id (%s)", taskID)
		if err := tasks.WaitForStatus(client, string(taskID),
			tasks.TaskStateFinished, ddosProfileUpdatingTimeout, true); err != nil {
			return diag.FromErr(err)
		}
	}

	if optionsChanged {
		log.Println("[DEBUG] Profile options changed: updating profile options")
		results, err := ddos.ActivateProfile(client, profileID, activateOpts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		taskID := results.Tasks[0]
		log.Printf("[DEBUG] Task id (%s)", taskID)
		if err := tasks.WaitForStatus(client, string(taskID),
			tasks.TaskStateFinished, ddosProfileUpdatingTimeout, true); err != nil {
			return diag.FromErr(err)
		}
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))
	log.Println("[DEBUG] Finish DDoS protection profile updating")
	return resourceDDoSProtectionRead(ctx, d, m)
}

func resourceDDoSProtectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start DDoS protection profile deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	profileID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] DDoS profile id = %d", profileID)

	client, err := CreateClient(provider, d, ddosProfilePoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	results, err := ddos.DeleteProfile(client, profileID).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	err = tasks.WaitForStatus(client, string(taskID), tasks.TaskStateFinished, ddosProfileDeletingTimeout, true)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Println("[DEBUG] Finish DDoS protection profile deleting")
	return diags
}
