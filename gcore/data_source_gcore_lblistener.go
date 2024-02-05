package gcore

import (
	"context"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/listeners"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLBListener() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBListenerRead,
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
			},
			"loadbalancer_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protocol": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Available values is 'HTTP', 'HTTPS', 'TCP', 'UDP'",
			},
			"protocol_port": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
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
			"timeout_client_data": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"timeout_member_connect": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"timeout_member_data": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"connection_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func dataSourceLBListenerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LBListener reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, LBListenersPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	var opts listeners.ListOpts
	name := d.Get("name").(string)
	lbID := d.Get("loadbalancer_id").(string)
	if lbID != "" {
		opts.LoadBalancerID = &lbID
	}

	ls, err := listeners.ListAll(client, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	var found bool
	var lb listeners.Listener
	for _, l := range ls {
		if l.Name == name {
			lb = l
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("lb listener with name %s not found", name)
	}

	d.SetId(lb.ID)
	d.Set("name", lb.Name)
	d.Set("protocol", lb.Protocol.String())
	d.Set("protocol_port", lb.ProtocolPort)
	d.Set("pool_count", lb.PoolCount)
	d.Set("operating_status", lb.OperationStatus.String())
	d.Set("provisioning_status", lb.ProvisioningStatus.String())
	d.Set("loadbalancer_id", lbID)
	d.Set("project_id", d.Get("project_id").(int))
	d.Set("region_id", d.Get("region_id").(int))
	d.Set("timeout_client_data", lb.TimeoutClientData)
	d.Set("timeout_member_connect", lb.TimeoutMemberConnect)
	d.Set("timeout_member_data", lb.TimeoutMemberData)
	d.Set("connection_limit", lb.ConnectionLimit)

	log.Println("[DEBUG] Finish LBListener reading")
	return diags
}
