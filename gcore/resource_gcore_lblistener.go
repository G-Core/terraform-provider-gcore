package gcore

import (
	"context"
	"fmt"
	"log"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/listeners"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/loadbalancers"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	LBListenersPoint                 = "lblisteners"
	LBListenerResourceTimeoutMinutes = 30

	LoadbalancerProvisioningStatusActive = "ACTIVE"
)

func resourceLbListener() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLBListenerCreate,
		ReadContext:   resourceLBListenerRead,
		UpdateContext: resourceLBListenerUpdate,
		DeleteContext: resourceLBListenerDelete,
		Description:   "Represent load balancer listener. Can not be created without load balancer. A listener is a process that checks for connection requests, using the protocol and port that you configure",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(LBListenerResourceTimeoutMinutes * time.Minute),
			Delete: schema.DefaultTimeout(LBListenerResourceTimeoutMinutes * time.Minute),
			Update: schema.DefaultTimeout(LBListenerResourceTimeoutMinutes * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, listenerID, lbID, err := ImportStringParserExtended(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.Set("loadbalancer_id", lbID)
				d.SetId(listenerID)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "ID of the desired project to create load balancer listener in. Alternative for `project_name`. One of them should be specified.",
				Optional:    true,
				ForceNew:    true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"region_id": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "ID of the desired region to create load balancer listener in. Alternative for `region_name`. One of them should be specified.",
				Optional:    true,
				ForceNew:    true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				DiffSuppressFunc: suppressDiffRegionID,
			},
			"project_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the desired project to create load balancer listener in. Alternative for `project_id`. One of them should be specified.",
				Optional:    true,
				ForceNew:    true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the desired region to create load balancer listener in. Alternative for `region_id`. One of them should be specified.",
				Optional:    true,
				ForceNew:    true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
			},
			"loadbalancer_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "ID of the target load balancer to attach newly created listener.",
				Required:    true,
				ForceNew:    true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Listener name.",
				Required:    true,
			},
			"protocol": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Available values are 'HTTP', 'HTTPS', 'TCP', 'UDP', 'TERMINATED_HTTPS', 'PROMETHEUS'",
				ValidateDiagFunc: func(val interface{}, key cty.Path) diag.Diagnostics {
					v := val.(string)
					switch types.ProtocolType(v) {
					case types.ProtocolTypeHTTP, types.ProtocolTypeHTTPS, types.ProtocolTypeTCP, types.ProtocolTypeUDP, types.ProtocolTypeTerminatedHTTPS, types.ProtocolTypePrometheus:
						return diag.Diagnostics{}
					}
					return diag.Errorf("wrong protocol %s, available values are 'HTTP', 'HTTPS', 'TCP', 'UDP', 'TERMINATED_HTTPS', 'PROMETHEUS'", v)
				},
			},
			"protocol_port": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Port number to listen, between 1 and 65535.",
				Required:    true,
				ForceNew:    true,
			},
			"insert_x_forwarded": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Insert X-Forwarded headers for 'HTTP', 'HTTPS', 'TERMINATED_HTTPS' protocols.",
				ForceNew:    true,
			},
			"pool_count": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Number of pools in this listener.",
				Computed:    true,
			},
			"operating_status": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Operating status of this listener.",
				Computed:    true,
			},
			"provisioning_status": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Provisioning status of this listener.",
				Computed:    true,
			},
			"secret_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Secret ID to use with 'TERMINATED_HTTPS' protocol.",
				Optional:    true,
			},
			"sni_secret_id": &schema.Schema{
				Type:        schema.TypeList,
				Description: "List of additional Secret IDs to use with 'TERMINATED_HTTPS' protocol.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"allowed_cidrs": &schema.Schema{
				Type:        schema.TypeList,
				Description: "List of networks from which listener is accessible",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"last_updated": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Datetime when load balancer was updated at the last time.",
				Computed:    true,
			},
			"timeout_client_data": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Frontend client inactivity timeout in milliseconds.",
				Optional:    true,
			},
			"timeout_member_connect": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Backend member connection timeout in milliseconds.",
				Optional:    true,
			},
			"timeout_member_data": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Backend member inactivity timeout in milliseconds.",
				Optional:    true,
			},
			"connection_limit": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Number of simultaneous connections for this listener, between 1 and 1,000,000.",
				Optional:    true,
			},
			"user_list": &schema.Schema{
				Type:        schema.TypeList,
				Description: "Load balancer listener list of username and encrypted password items.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": &schema.Schema{
							Type:        schema.TypeString,
							Description: "Username to auth via Basic Authentication",
							Required:    true,
						},
						"encrypted_password": &schema.Schema{
							Type:        schema.TypeString,
							Description: "Encrypted password (hash) to auth via Basic Authentication",
							Sensitive:   true,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func LoadbalancerProvisioningStatusRefreshedFunc(client *gcorecloud.ServiceClient, loadbalancerID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		loadbalancer, err := loadbalancers.Get(client, loadbalancerID, nil).Extract()
		if err != nil {
			return loadbalancer, "", err
		}
		return loadbalancer, loadbalancer.ProvisioningStatus.String(), nil
	}
}

func resourceLBListenerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LBListener creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, LBListenersPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := listeners.CreateOpts{
		Name:             d.Get("name").(string),
		Protocol:         types.ProtocolType(d.Get("protocol").(string)),
		ProtocolPort:     d.Get("protocol_port").(int),
		LoadBalancerID:   d.Get("loadbalancer_id").(string),
		InsertXForwarded: d.Get("insert_x_forwarded").(bool),
		SecretID:         d.Get("secret_id").(string),
	}
	if tcd, ok := d.GetOkExists("timeout_client_data"); ok {
		timeoutClientData := tcd.(int)
		opts.TimeoutClientData = &timeoutClientData
	}
	if tmc, ok := d.GetOkExists("timeout_member_connect"); ok {
		timeoutMemberConnect := tmc.(int)
		opts.TimeoutMemberConnect = &timeoutMemberConnect
	}
	if tmd, ok := d.GetOkExists("timeout_member_data"); ok {
		timeoutMemberConnect := tmd.(int)
		opts.TimeoutMemberData = &timeoutMemberConnect
	}
	if cl, ok := d.GetOkExists("connection_limit"); ok {
		connectionLimit := cl.(int)
		opts.ConnectionLimit = &connectionLimit
	}

	sniSecretIDRaw := d.Get("sni_secret_id").([]interface{})
	if len(sniSecretIDRaw) != 0 {
		sniSecretID := make([]string, len(sniSecretIDRaw))
		for i, s := range sniSecretIDRaw {
			sniSecretID[i] = s.(string)
		}
		opts.SNISecretID = sniSecretID
	}

	allowedCIDRSRaw := d.Get("allowed_cidrs").([]interface{})
	if len(allowedCIDRSRaw) != 0 {
		allowedCIDRS := make([]string, len(allowedCIDRSRaw))
		for i, a := range allowedCIDRSRaw {
			allowedCIDRS[i] = a.(string)
		}
		opts.AllowedCIDRS = allowedCIDRS
	}

	u := d.Get("user_list")
	if len(u.([]interface{})) > 0 {
		userList, err := extractUserList(u.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		opts.UserList = userList
	}
	timeout := int(d.Timeout(schema.TimeoutCreate).Seconds())
	rc := GetConflictRetryConfig(timeout)
	results, err := listeners.Create(client, opts, &gcorecloud.RequestOpts{
		ConflictRetryAmount:   rc.Amount,
		ConflictRetryInterval: rc.Interval,
	}).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	listenerID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, timeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		listenerID, err := listeners.ExtractListenerIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve LBListener ID from task info: %w", err)
		}
		return listenerID, nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(listenerID.(string))
	resourceLBListenerRead(ctx, d, m)

	log.Printf("[DEBUG] Finish LBListener creating (%s)", listenerID)
	return diags
}

func resourceLBListenerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LBListener reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, LBListenersPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	lb, err := listeners.Get(client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", lb.Name)
	d.Set("protocol", lb.Protocol.String())
	d.Set("protocol_port", lb.ProtocolPort)
	d.Set("pool_count", lb.PoolCount)
	d.Set("operating_status", lb.OperationStatus.String())
	d.Set("provisioning_status", lb.ProvisioningStatus.String())
	d.Set("secret_id", lb.SecretID)
	d.Set("sni_secret_id", lb.SNISecretID)
	d.Set("allowed_cidrs", lb.AllowedCIDRS)

	fields := []string{"project_id", "region_id", "loadbalancer_id", "insert_x_forwarded"}
	revertState(d, &fields)

	log.Println("[DEBUG] Finish LBListener reading")
	return diags
}

func resourceLBListenerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LBListener updating")
	config := m.(*Config)
	provider := config.Provider

	clientV1, err := CreateClient(provider, d, LoadBalancersPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	clientV2, err := CreateClient(provider, d, LBListenersPoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	var changed bool
	var toUnset bool
	updateOpts := listeners.UpdateOpts{}
	unsetOpts := listeners.UnsetOpts{}

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
		changed = true
	}

	if d.HasChange("secret_id") {
		updateOpts.SecretID = d.Get("secret_id").(string)
		changed = true
	}

	if d.HasChange("sni_secret_id") {
		sniSecretIDRaw := d.Get("sni_secret_id").([]interface{})
		sniSecretID := make([]string, len(sniSecretIDRaw))
		for i, s := range sniSecretIDRaw {
			sniSecretID[i] = s.(string)
		}
		updateOpts.SNISecretID = sniSecretID
		changed = true
	}

	if d.HasChange("timeout_client_data") {
		timeoutClientData := d.Get("timeout_client_data").(int)
		updateOpts.TimeoutClientData = &timeoutClientData
		changed = true
	}

	if d.HasChange("timeout_member_connect") {
		timeoutMemberConnect := d.Get("timeout_member_connect").(int)
		updateOpts.TimeoutMemberConnect = &timeoutMemberConnect
		changed = true
	}

	if d.HasChange("timeout_member_data") {
		timeoutMemberData := d.Get("timeout_member_data").(int)
		updateOpts.TimeoutMemberData = &timeoutMemberData
		changed = true
	}

	if d.HasChange("connection_limit") {
		connectionLimit := d.Get("connection_limit").(int)
		updateOpts.ConnectionLimit = &connectionLimit
		changed = true
	}

	if d.HasChange("allowed_cidrs") {
		allowedCIDRSRaw := d.Get("allowed_cidrs").([]interface{})
		if len(allowedCIDRSRaw) > 0 {
			allowedCIDRS := make([]string, len(allowedCIDRSRaw))
			for i, a := range allowedCIDRSRaw {
				allowedCIDRS[i] = a.(string)
			}
			updateOpts.AllowedCIDRS = allowedCIDRS
		} else {
			unsetOpts.AllowedCIDRS = true
			toUnset = true
		}
		changed = true
	}

	if d.HasChange("user_list") {
		u := d.Get("user_list")
		updateOpts.UserList = make([]listeners.CreateUserListOpts, 0)
		if len(u.([]interface{})) > 0 {
			userList, err := extractUserList(u.([]interface{}))
			if err != nil {
				return diag.FromErr(err)
			}
			updateOpts.UserList = userList
		} else {
			unsetOpts.UserList = true
			toUnset = true
		}
		changed = true
	}

	if changed {
		rc := GetConflictRetryConfig(int(d.Timeout(schema.TimeoutUpdate).Seconds()))
		_, err = listeners.Update(clientV2, d.Id(), updateOpts, &gcorecloud.RequestOpts{
			ConflictRetryAmount:   rc.Amount,
			ConflictRetryInterval: rc.Interval,
		}).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		if toUnset {
			stopWaitConf := retry.StateChangeConf{
				Target:     []string{types.ProvisioningStatusActive.String()},
				Refresh:    LoadbalancerProvisioningStatusRefreshedFunc(clientV1, d.Get("loadbalancer_id").(string)),
				Timeout:    3 * time.Minute,
				Delay:      10 * time.Second,
				MinTimeout: 5 * time.Second,
			}
			_, err = stopWaitConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for loadbalancer (%s): %s", d.Get("loadbalancer_id").(string), err)
			}
			_, err := listeners.Unset(clientV2, d.Id(), unsetOpts, &gcorecloud.RequestOpts{
				ConflictRetryAmount:   rc.Amount,
				ConflictRetryInterval: rc.Interval,
			}).Extract()
			if err != nil {
				return diag.FromErr(err)
			}
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	log.Println("[DEBUG] Finish LBListener updating")
	return resourceLBListenerRead(ctx, d, m)
}

func resourceLBListenerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LBListener deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, LBListenersPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	timeout := int(d.Timeout(schema.TimeoutDelete).Seconds())
	rc := GetConflictRetryConfig(timeout)
	results, err := listeners.Delete(client, id, &gcorecloud.RequestOpts{
		ConflictRetryAmount:   rc.Amount,
		ConflictRetryInterval: rc.Interval,
	}).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, timeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := listeners.Get(client, id).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete LBListener with ID: %s", id)
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
	log.Printf("[DEBUG] Finish of LBListener deleting")
	return diags
}
