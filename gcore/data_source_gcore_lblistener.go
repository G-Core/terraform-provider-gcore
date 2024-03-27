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
				Type:        schema.TypeInt,
				Description: "ID of the project in which load balancer listener was created.",
				Optional:    true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"region_id": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "ID of the region in which load balancer listener was created.",
				Optional:    true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				DiffSuppressFunc: suppressDiffRegionID,
			},
			"project_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the project in which load balancer listener was created.",
				Optional:    true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the region in which load balancer listener was created.",
				Optional:    true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the load balancer listener.",
				Required:    true,
			},
			"loadbalancer_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "ID of the load balancer to which listener was attached.",
				Optional:    true,
				Computed:    true,
			},
			"protocol": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Available values are 'HTTP', 'HTTPS', 'TCP', 'UDP', 'TERMINATED_HTTPS', 'PROMETHEUS'",
			},
			"protocol_port": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Port number to listen, between 1 and 65535.",
				Computed:    true,
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
							Type:     schema.TypeString,
							Required: true,
						},
						"encrypted_password": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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

	userList := make([]map[string]string, len(lb.UserList))
	for i, userData := range lb.UserList {
		u := map[string]string{"username": userData.Username, "encrypted_password": userData.EncryptedPassword}
		userList[i] = u
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
	d.Set("user_list", userList)

	log.Println("[DEBUG] Finish LBListener reading")
	return diags
}
