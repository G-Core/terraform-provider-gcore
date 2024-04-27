package gcore

import (
	"context"
	"fmt"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/lbpools"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLBPool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBPoolRead,
		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "ID of the project in which load balancer pool was created.",
				Optional:    true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"region_id": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "ID of the region in which load balancer pool was created.",
				Optional:    true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				DiffSuppressFunc: suppressDiffRegionID,
			},
			"project_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the project in which load balancer pool was created.",
				Optional:    true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the region in which load balancer pool was created.",
				Optional:    true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the load balancer pool.",
				Required:    true,
			},
			"loadbalancer_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "ID of the load balancer to which pool was attached.",
				Optional:    true,
				Computed:    true,
			},
			"listener_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "ID of the load balancer listener to which pool was attached.",
				Optional:    true,
				Computed:    true,
			},
			"lb_algorithm": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: fmt.Sprintf("Available values is '%s', '%s', '%s'", types.LoadBalancerAlgorithmRoundRobin, types.LoadBalancerAlgorithmLeastConnections, types.LoadBalancerAlgorithmSourceIP),
			},
			"protocol": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: fmt.Sprintf("Available values are '%s', '%s', '%s', '%s', '%s', '%s'", types.ProtocolTypeHTTP, types.ProtocolTypeHTTPS, types.ProtocolTypeTCP, types.ProtocolTypeUDP, types.ProtocolTypePROXY, types.ProtocolTypePROXYV2),
			},
			"health_monitor": &schema.Schema{
				Type:        schema.TypeList,
				Description: "Health Monitor settings for defining health state of members inside this pool.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Description: "Health Monitor ID.",
							Computed:    true,
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Description: fmt.Sprintf("Available values is '%s', '%s', '%s', '%s', '%s', '%s", types.HealthMonitorTypeHTTP, types.HealthMonitorTypeHTTPS, types.HealthMonitorTypePING, types.HealthMonitorTypeTCP, types.HealthMonitorTypeTLSHello, types.HealthMonitorTypeUDPConnect),
							Computed:    true,
						},
						"delay": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "The time, in seconds, between sending probes to members.",
							Computed:    true,
						},
						"max_retries": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "The number of successful checks before changing the operating status of the member to ONLINE.",
							Computed:    true,
						},
						"timeout": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "The maximum time, in seconds, that a monitor waits to connect before it times out.",
							Computed:    true,
						},
						"max_retries_down": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "The number of allowed check failures before changing the operating status of the member to ERROR.",
							Computed:    true,
						},
						"http_method": &schema.Schema{
							Type:        schema.TypeString,
							Description: "The HTTP method that the health monitor uses for requests.",
							Computed:    true,
						},
						"url_path": &schema.Schema{
							Type:        schema.TypeString,
							Description: "The HTTP URL path of the request sent by the monitor to test the health of a backend member.",
							Computed:    true,
						},
						"expected_codes": &schema.Schema{
							Type:        schema.TypeString,
							Description: "The list of HTTP status codes expected in response from the member to declare it healthy.",
							Computed:    true,
						},
					},
				},
			},
			"session_persistence": &schema.Schema{
				Type:        schema.TypeList,
				Description: "Pool session persistence tells the load balancer to attempt to send future requests from a client to the same backend member as the initial request.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Description: "One of: 'APP_COOKIE' (an application supplied cookie), 'HTTP_COOKIE' (a cookie created by the load balancer), 'SOURCE_IP' (the source IP address).",
							Computed:    true,
						},
						"cookie_name": &schema.Schema{
							Type:        schema.TypeString,
							Description: "The name of the application cookie to use for session persistence.",
							Computed:    true,
						},
						"persistence_granularity": &schema.Schema{
							Type:        schema.TypeString,
							Description: "The netmask used to determine SCTP or UDP SOURCE_IP session persistence.",
							Computed:    true,
						},
						"persistence_timeout": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "The timeout, in seconds, after which a SCTP or UDP flow may be rescheduled to a different member.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceLBPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LBPool reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, LBPoolsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	var opts lbpools.ListOpts
	name := d.Get("name").(string)
	lbID := d.Get("loadbalancer_id").(string)
	if lbID != "" {
		opts.LoadBalancerID = &lbID
	}
	lID := d.Get("listener_id").(string)
	if lbID != "" {
		opts.ListenerID = &lID
	}

	pools, err := lbpools.ListAll(client, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	var found bool
	var lb lbpools.Pool
	for _, p := range pools {
		if p.Name == name {
			lb = p
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("lb listener with name %s not found", name)
	}

	d.SetId(lb.ID)
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

	d.Set("project_id", d.Get("project_id").(int))
	d.Set("region_id", d.Get("region_id").(int))

	log.Println("[DEBUG] Finish LBPool reading")
	return diags
}
