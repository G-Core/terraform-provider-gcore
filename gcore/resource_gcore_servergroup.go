package gcore

import (
	"context"
	"errors"
	"fmt"
	"log"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	gcloudInstance "github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/servergroup/v1/servergroups"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	serverGroupsPoint            = "servergroups"
	putToServerGroupTimeout      = 1200
	removeFromServerGroupTimeout = 1200
)

func resourceServerGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerGroupCreate,
		ReadContext:   resourceServerGroupRead,
		UpdateContext: resourceServerGroupUpdate,
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
				Type: schema.TypeString,
				Description: `
Server group policy. Available value is 'affinity' 'anti-affinity' 'soft-anti-affinity'. 

**Affinity Rule**:

Strictly isolates instances across different physical servers, minimizing simultaneous failures for
critical applications. This rule is recommended for enhanced fault tolerance.

**Soft Anti-affinity Rule**:

Attempts to place instances belonging to the same server group on a single host. If this is not
possible, instances will be scheduled on the minimum number of hosts required.

**Anti-affinity Rule**:

Improves performance by grouping related instances on the same physical server to enhance communication,
resource sharing, and faster interaction.`,
				Required: true,
				ForceNew: true,
			},
			"instance": {
				Type:        schema.TypeSet,
				Description: "Instances in this server group",
				Optional:    true,
				Set:         serverGroupInstanceUniqueID,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Optional: true,
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

	instanceClient, err := CreateClient(provider, d, InstancePoint, versionPointV1)
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

	instancesSet := d.Get("instance").(*schema.Set)
	instances := instancesSet.List()
	if len(instances) > 0 {
		for _, instance := range instances {
			instanceMap := instance.(map[string]interface{})
			instanceID := instanceMap["instance_id"].(string)

			err := putInstanceToServerGroup(instanceClient, instanceID, serverGroup.ServerGroupID)
			if err != nil {
				// we try to delete server group if we can't add instance to it on server group creation
				errDel := servergroups.Delete(client, serverGroup.ServerGroupID).Err
				if errDel != nil {
					return diag.FromErr(errDel)
				}
				return diag.FromErr(err)
			}

		}
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

	instances := make([]interface{}, len(serverGroup.Instances))
	for i, instance := range serverGroup.Instances {
		rawInstance := make(map[string]interface{})
		rawInstance["instance_id"] = instance.InstanceID
		rawInstance["instance_name"] = instance.InstanceName
		instances[i] = rawInstance
	}
	if err := d.Set("instance", schema.NewSet(serverGroupInstanceUniqueID, instances)); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish ServerGroup reading")
	return diags
}

func resourceServerGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start ServerGroup updating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, InstancePoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("instance") {
		oldInstances, newInstances := d.GetChange("instance")
		oldInstancesSet := oldInstances.(*schema.Set)
		newInstancesSet := newInstances.(*schema.Set)

		for _, i := range oldInstancesSet.Difference(newInstancesSet).List() {
			instance := i.(map[string]interface{})
			instanceID := instance["instance_id"].(string)

			log.Printf("[DEBUG] try to delete instance %s from server group", instanceID)
			if err := removeInstanceFromServerGroup(client, instanceID); err != nil {
				log.Printf("[DEBUG] error while removing instance from server group: %s", err)
				var errDefault404 gcorecloud.ErrDefault404
				if errors.As(err, &errDefault404) {
					continue
				}

				return diag.FromErr(err)
			}
		}

		for _, i := range newInstancesSet.Difference(oldInstancesSet).List() {
			instance := i.(map[string]interface{})
			instanceID := instance["instance_id"].(string)

			log.Printf("[DEBUG] try to add instance %s into server group", instanceID)
			if err := putInstanceToServerGroup(client, instanceID, d.Id()); err != nil {
				return diag.FromErr(err)
			}
		}
	}

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

func removeInstanceFromServerGroup(client *gcorecloud.ServiceClient, instanceID string) error {
	results, err := gcloudInstance.RemoveFromServerGroup(client, instanceID).Extract()
	if err != nil {
		return err
	}

	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, removeFromServerGroupTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}

		return nil, nil
	},
	)
	return err
}

func putInstanceToServerGroup(client *gcorecloud.ServiceClient, instanceID, serverGroupID string) error {
	opts := gcloudInstance.PutServerGroupOpts{
		ServerGroupID: serverGroupID,
	}
	results, err := gcloudInstance.PutToServerGroup(client, instanceID, opts).Extract()
	if err != nil {
		return err
	}

	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, putToServerGroupTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}

		return nil, nil
	},
	)
	return err
}
