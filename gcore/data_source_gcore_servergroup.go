package gcore

import (
	"context"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/servergroup/v1/servergroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceServerGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServerGroupRead,
		Description: "Represent server group data",
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
			"name": {
				Type:        schema.TypeString,
				Description: "Displayed server group name",
				Required:    true,
			},
			"policy": {
				Type:        schema.TypeString,
				Description: "Server group policy. Available value is 'affinity', 'anti-affinity'",
				Computed:    true,
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

func dataSourceServerGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start ServerGroup reading")
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, serverGroupsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	var serverGroup servergroups.ServerGroup
	serverGroups, err := servergroups.ListAll(client)
	if err != nil {
		return diag.FromErr(err)
	}

	var found bool
	name := d.Get("name").(string)
	for _, sg := range serverGroups {
		if sg.Name == name {
			serverGroup = sg
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("server group with name %s not found", name)
	}

	d.SetId(serverGroup.ServerGroupID)
	d.Set("name", name)
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
	return nil
}
