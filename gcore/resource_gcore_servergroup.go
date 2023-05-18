package gcore

import (
	"context"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/servergroup/v1/servergroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	serverGroupsPoint = "servergroups"
)

func resourceServerGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerGroupCreate,
		ReadContext:   resourceServerGroupRead,
		DeleteContext: resourceServerGroupDelete,
		Description:   "Represent server group resource",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, sgID, err := ImportStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(sgID)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"region_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				DiffSuppressFunc: suppressDiffRegionID,
			},
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": &schema.Schema{
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
				Description: "Displayed server group name",
				Required:    true,
				ForceNew:    true,
			},
			"policy": {
				Type:        schema.TypeString,
				Description: "Server group policy. Available value is 'affinity', 'anti-affinity'",
				Required:    true,
				ForceNew:    true,
			},
			"instances": {
				Type:        schema.TypeList,
				Description: "Instances in this server group",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceServerGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start ServerGroup creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, serverGroupsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := servergroups.CreateOpts{
		Name:   d.Get("name").(string),
		Policy: servergroups.ServerGroupPolicy(d.Get("policy").(string)),
	}

	serverGroup, err := servergroups.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serverGroup.ServerGroupID)
	resourceServerGroupRead(ctx, d, m)
	log.Println("[DEBUG] Finish ServerGroup creating")
	return diags
}

func resourceServerGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start ServerGroup reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, serverGroupsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	serverGroup, err := servergroups.Get(client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", serverGroup.Name)
	d.Set("project_id", serverGroup.ProjectID)
	d.Set("region_id", serverGroup.RegionID)
	d.Set("policy", serverGroup.Policy.String())

	instances := make([]map[string]string, len(serverGroup.Instances))
	for i, instance := range serverGroup.Instances {
		rawInstance := make(map[string]string)
		rawInstance["instance_id"] = instance.InstanceID
		rawInstance["instance_name"] = instance.InstanceName
		instances[i] = rawInstance
	}
	if err := d.Set("instances", instances); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish ServerGroup reading")
	return diags
}

func resourceServerGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start ServerGroup deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, serverGroupsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	err = servergroups.Delete(client, d.Id()).Err
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Println("[DEBUG] Finish ServerGroup deleting")
	return diags
}
