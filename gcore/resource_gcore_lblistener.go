package gcore

import (
	"context"
	"fmt"
	"log"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/listeners"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	LBListenersPoint        = "lblisteners"
	LBListenerCreateTimeout = 2400
)

func resourceLbListener() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLBListenerCreate,
		ReadContext:   resourceLBListenerRead,
		UpdateContext: resourceLBListenerUpdate,
		DeleteContext: resourceLBListenerDelete,
		Description:   "Represent load balancer listener. Can not be created without load balancer. A listener is a process that checks for connection requests, using the protocol and port that you configure",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
			"loadbalancer_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Available values is 'HTTP', 'HTTPS', 'TCP', 'UDP', 'TERMINATED_HTTPS'",
				ValidateDiagFunc: func(val interface{}, key cty.Path) diag.Diagnostics {
					v := val.(string)
					switch types.ProtocolType(v) {
					case types.ProtocolTypeHTTP, types.ProtocolTypeHTTPS, types.ProtocolTypeTCP, types.ProtocolTypeUDP, types.ProtocolTypeTerminatedHTTPS:
						return diag.Diagnostics{}
					}
					return diag.Errorf("wrong protocol %s, available values is 'HTTP', 'HTTPS', 'TCP', 'UDP', 'TERMINATED_HTTPS", v)
				},
			},
			"protocol_port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"insert_x_forwarded": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Insert *-forwarded headers",
				ForceNew:    true,
			},
			"pool_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"operating_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioning_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"sni_secret_id": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
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
	sniSecretIDRaw := d.Get("sni_secret_id").([]interface{})
	if len(sniSecretIDRaw) != 0 {
		sniSecretID := make([]string, len(sniSecretIDRaw))
		for i, s := range sniSecretIDRaw {
			sniSecretID[i] = s.(string)
		}
		opts.SNISecretID = sniSecretID
	}

	results, err := listeners.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	listenerID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, LBListenerCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
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

	fields := []string{"project_id", "region_id", "loadbalancer_id", "insert_x_forwarded"}
	revertState(d, &fields)

	log.Println("[DEBUG] Finish LBListener reading")
	return diags
}

func resourceLBListenerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LBListener updating")
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, LBListenersPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	var changed bool
	opts := listeners.UpdateOpts{}
	if d.HasChange("name") {
		opts.Name = d.Get("name").(string)
		changed = true
	}

	if d.HasChange("secret_id") {
		opts.SecretID = d.Get("secret_id").(string)
		changed = true
	}

	if d.HasChange("sni_secret_id") {
		sniSecretIDRaw := d.Get("sni_secret_id").([]interface{})
		sniSecretID := make([]string, len(sniSecretIDRaw))
		for i, s := range sniSecretIDRaw {
			sniSecretID[i] = s.(string)
		}
		opts.SNISecretID = sniSecretID
		changed = true
	}

	if changed {
		_, err = listeners.Update(client, d.Id(), opts).Extract()
		if err != nil {
			return diag.FromErr(err)
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
	results, err := listeners.Delete(client, id).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, LBListenerCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
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
