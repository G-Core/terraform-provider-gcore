package gcore

import (
	"context"
	"fmt"
	"log"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/lbpools"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/hashicorp/go-cty/cty"

	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	LBPoolsPoint         = "lbpools"
	LBPoolsCreateTimeout = 2400
)

func resourceLBPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLBPoolCreate,
		ReadContext:   resourceLBPoolRead,
		UpdateContext: resourceLBPoolUpdate,
		DeleteContext: resourceLBPoolDelete,
		Description:   "Represent load balancer listener pool. A pool is a list of virtual machines to which the listener will redirect incoming traffic",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, lbPoolID, err := ImportStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(lbPoolID)

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "ID of the desired project to create load balancer pool in. Alternative for `project_name`. One of them should be specified.",
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
				Description: "ID of the desired region to create load balancer pool in. Alternative for `region_name`. One of them should be specified.",
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
				Description: "Name of the desired project to create load balancer pool in. Alternative for `project_id`. One of them should be specified.",
				Optional:    true,
				ForceNew:    true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the desired region to create load balancer pool in. Alternative for `region_id`. One of them should be specified.",
				Optional:    true,
				ForceNew:    true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Pool name.",
				Required:    true,
			},
			"lb_algorithm": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: fmt.Sprintf("Available values is '%s', '%s', '%s', '%s'", types.LoadBalancerAlgorithmRoundRobin, types.LoadBalancerAlgorithmLeastConnections, types.LoadBalancerAlgorithmSourceIP, types.LoadBalancerAlgorithmSourceIPPort),
				ValidateDiagFunc: func(val interface{}, key cty.Path) diag.Diagnostics {
					v := val.(string)
					switch types.LoadBalancerAlgorithm(v) {
					case types.LoadBalancerAlgorithmRoundRobin, types.LoadBalancerAlgorithmLeastConnections, types.LoadBalancerAlgorithmSourceIP, types.LoadBalancerAlgorithmSourceIPPort:
						return diag.Diagnostics{}
					}
					return diag.Errorf("wrong type %s, available values is '%s', '%s', '%s', '%s'", v, types.LoadBalancerAlgorithmRoundRobin, types.LoadBalancerAlgorithmLeastConnections, types.LoadBalancerAlgorithmSourceIP, types.LoadBalancerAlgorithmSourceIPPort)
				},
			},
			"protocol": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: fmt.Sprintf("Available values are '%s', '%s', '%s', '%s', '%s', '%s'", types.ProtocolTypeHTTP, types.ProtocolTypeHTTPS, types.ProtocolTypeTCP, types.ProtocolTypeUDP, types.ProtocolTypePROXY, types.ProtocolTypePROXYV2),
				ValidateDiagFunc: func(val interface{}, key cty.Path) diag.Diagnostics {
					v := val.(string)
					switch types.ProtocolType(v) {
					case types.ProtocolTypeHTTP, types.ProtocolTypeHTTPS, types.ProtocolTypeTCP, types.ProtocolTypeUDP, types.ProtocolTypePROXY, types.ProtocolTypePROXYV2:
						return diag.Diagnostics{}
					}
					return diag.Errorf("wrong type %s, available values are '%s', '%s', '%s', '%s', '%s', '%s'", v, types.ProtocolTypeHTTP, types.ProtocolTypeHTTPS, types.ProtocolTypeTCP, types.ProtocolTypeUDP, types.ProtocolTypePROXY, types.ProtocolTypePROXYV2)
				},
			},
			"loadbalancer_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "ID of the target load balancer to attach newly created pool.",
				Optional:    true,
				ForceNew:    true,
			},
			"listener_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "ID of the target listener associated with load balancer to attach newly created pool.",
				Optional:    true,
				ForceNew:    true,
			},
			"health_monitor": &schema.Schema{
				Type:        schema.TypeList,
				Description: "Health Monitor settings for defining health state of members inside this pool.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Description: "Health Monitor ID.",
							Optional:    true,
							Computed:    true,
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: fmt.Sprintf("Available values is '%s', '%s', '%s', '%s', '%s', '%s", types.HealthMonitorTypeHTTP, types.HealthMonitorTypeHTTPS, types.HealthMonitorTypePING, types.HealthMonitorTypeTCP, types.HealthMonitorTypeTLSHello, types.HealthMonitorTypeUDPConnect),
							ValidateDiagFunc: func(val interface{}, key cty.Path) diag.Diagnostics {
								v := val.(string)
								switch types.HealthMonitorType(v) {
								case types.HealthMonitorTypeHTTP, types.HealthMonitorTypeHTTPS, types.HealthMonitorTypePING, types.HealthMonitorTypeTCP, types.HealthMonitorTypeTLSHello, types.HealthMonitorTypeUDPConnect:
									return diag.Diagnostics{}
								}
								return diag.Errorf("wrong type %s, available values is '%s', '%s', '%s', '%s', '%s', '%s", v, types.HealthMonitorTypeHTTP, types.HealthMonitorTypeHTTPS, types.HealthMonitorTypePING, types.HealthMonitorTypeTCP, types.HealthMonitorTypeTLSHello, types.HealthMonitorTypeUDPConnect)
							},
						},
						"delay": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "The time, in seconds, between sending probes to members.",
							Required:    true,
						},
						"max_retries": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "The number of successful checks before changing the operating status of the member to ONLINE.",
							Required:    true,
						},
						"timeout": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "The maximum time, in seconds, that a monitor waits to connect before it times out.",
							Required:    true,
						},
						"max_retries_down": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "The number of allowed check failures before changing the operating status of the member to ERROR.",
							Optional:    true,
							Computed:    true,
						},
						"http_method": &schema.Schema{
							Type:        schema.TypeString,
							Description: "The HTTP method that the health monitor uses for requests.",
							Optional:    true,
							Computed:    true,
						},
						"url_path": &schema.Schema{
							Type:        schema.TypeString,
							Description: "The HTTP URL path of the request sent by the monitor to test the health of a backend member.",
							Optional:    true,
							Computed:    true,
						},
						"expected_codes": &schema.Schema{
							Type:        schema.TypeString,
							Description: "The list of HTTP status codes expected in response from the member to declare it healthy.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"session_persistence": &schema.Schema{
				Type:        schema.TypeList,
				Description: "Pool session persistence tells the load balancer to attempt to send future requests from a client to the same backend member as the initial request.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Description: "One of: 'APP_COOKIE' (an application supplied cookie), 'HTTP_COOKIE' (a cookie created by the load balancer), 'SOURCE_IP' (the source IP address).",
							Required:    true,
						},
						"cookie_name": &schema.Schema{
							Type:        schema.TypeString,
							Description: "The name of the application cookie to use for session persistence.",
							Optional:    true,
							Computed:    true,
						},
						"persistence_granularity": &schema.Schema{
							Type:        schema.TypeString,
							Description: "The netmask used to determine SCTP or UDP SOURCE_IP session persistence.",
							Optional:    true,
							Computed:    true,
						},
						"persistence_timeout": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "The timeout, in seconds, after which a SCTP or UDP flow may be rescheduled to a different member.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"last_updated": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Datetime when load balancer pool was updated at the last time.",
				Computed:    true,
			},
		},
	}
}

func resourceLBPoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LBPool creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, LBPoolsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	healthOpts := extractHealthMonitorMap(d)
	sessionOpts := extractSessionPersistenceMap(d)
	opts := lbpools.CreateOpts{
		Name:               d.Get("name").(string),
		Protocol:           types.ProtocolType(d.Get("protocol").(string)),
		LBPoolAlgorithm:    types.LoadBalancerAlgorithm(d.Get("lb_algorithm").(string)),
		LoadBalancerID:     d.Get("loadbalancer_id").(string),
		ListenerID:         d.Get("listener_id").(string),
		HealthMonitor:      healthOpts,
		SessionPersistence: sessionOpts,
	}

	results, err := lbpools.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	lbPoolID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, LBPoolsCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		lbPoolID, err := lbpools.ExtractPoolIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve LBPool ID from task info: %w", err)
		}
		return lbPoolID, nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(lbPoolID.(string))
	resourceLBPoolRead(ctx, d, m)

	log.Printf("[DEBUG] Finish LBPool creating (%s)", lbPoolID)
	return diags
}

func resourceLBPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LBPool reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, LBPoolsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	lb, err := lbpools.Get(client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", lb.Name)
	d.Set("lb_algorithm", lb.LoadBalancerAlgorithm.String())
	d.Set("protocol", lb.Protocol.String())

	if len(lb.LoadBalancers) > 0 {
		d.Set("loadbalancer_id", lb.LoadBalancers[0].ID)
	}

	if len(lb.Listeners) > 0 {
		d.Set("listener_id", lb.Listeners[0].ID)
	}

	if lb.HealthMonitor != nil {
		healthMonitor := map[string]interface{}{
			"id":               lb.HealthMonitor.ID,
			"type":             lb.HealthMonitor.Type.String(),
			"delay":            lb.HealthMonitor.Delay,
			"timeout":          lb.HealthMonitor.Timeout,
			"max_retries":      lb.HealthMonitor.MaxRetries,
			"max_retries_down": lb.HealthMonitor.MaxRetriesDown,
			"url_path":         lb.HealthMonitor.URLPath,
			"expected_codes":   lb.HealthMonitor.ExpectedCodes,
		}
		if lb.HealthMonitor.HTTPMethod != nil {
			healthMonitor["http_method"] = lb.HealthMonitor.HTTPMethod.String()
		}

		if err := d.Set("health_monitor", []interface{}{healthMonitor}); err != nil {
			return diag.FromErr(err)
		}
	}

	if lb.SessionPersistence != nil {
		sessionPersistence := map[string]interface{}{
			"type":                    lb.SessionPersistence.Type.String(),
			"cookie_name":             lb.SessionPersistence.CookieName,
			"persistence_granularity": lb.SessionPersistence.PersistenceGranularity,
			"persistence_timeout":     lb.SessionPersistence.PersistenceTimeout,
		}

		if err := d.Set("session_persistence", []interface{}{sessionPersistence}); err != nil {
			return diag.FromErr(err)
		}
	}

	fields := []string{"project_id", "region_id"}
	revertState(d, &fields)

	log.Println("[DEBUG] Finish LBPool reading")
	return diags
}

func resourceLBPoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LBPool updating")
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, LBPoolsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	var change bool
	opts := lbpools.UpdateOpts{Name: d.Get("name").(string)}

	if d.HasChange("lb_algorithm") {
		opts.LBPoolAlgorithm = types.LoadBalancerAlgorithm(d.Get("lb_algorithm").(string))
		change = true
	}

	if d.HasChange("health_monitor") {
		opts.HealthMonitor = extractHealthMonitorMap(d)
		if opts.HealthMonitor == nil {
			lbpools.DeleteHealthMonitor(client, d.Id())
		} else {
			change = true
		}
	}

	if d.HasChange("session_persistence") {
		opts.SessionPersistence = extractSessionPersistenceMap(d)
		if opts.SessionPersistence == nil {
			results, err := lbpools.Unset(client, d.Id(), lbpools.UnsetOpts{SessionPersistence: true}).Extract()
			if err != nil {
				return diag.FromErr(err)
			}
			taskID := results.Tasks[0]
			_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, LBPoolsCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
				_, err := tasks.Get(client, string(task)).Extract()
				if err != nil {
					return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
				}
				return nil, nil
			})
		} else {
			change = true
		}
	}

	if !change {
		log.Println("[DEBUG] Finish LBPool updating")
		return resourceLBPoolRead(ctx, d, m)
	}

	results, err := lbpools.Update(client, d.Id(), opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, LBPoolsCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		return nil, nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))
	log.Println("[DEBUG] Finish LBPool updating")
	return resourceLBPoolRead(ctx, d, m)
}

func resourceLBPoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LBPool deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, LBPoolsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	results, err := lbpools.Delete(client, id).Extract()
	if err != nil {
		switch err.(type) {
		case gcorecloud.ErrDefault404:
		default:
			return diag.FromErr(err)
		}
	}

	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, LBPoolsCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := lbpools.Get(client, id).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete LBPool with ID: %s", id)
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
	log.Printf("[DEBUG] Finish of LBPool deleting")
	return diags
}
