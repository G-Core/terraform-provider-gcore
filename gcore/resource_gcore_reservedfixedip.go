package gcore

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/port/v1/ports"
	"github.com/G-Core/gcorelabscloud-go/gcore/reservedfixedip/v1/reservedfixedips"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	reservedFixedIPsPoint        = "reserved_fixed_ips"
	portsPoint                   = "ports"
	ReservedFixedIPCreateTimeout = 1200
)

func resourceReservedFixedIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReservedFixedIPCreate,
		ReadContext:   resourceReservedFixedIPRead,
		UpdateContext: resourceReservedFixedIPUpdate,
		DeleteContext: resourceReservedFixedIPDelete,
		Description:   "Represent reserved ips",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, ipID, err := ImportStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(ipID)

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
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: fmt.Sprintf("Available values is '%s', '%s', '%s', '%s'", reservedfixedips.External, reservedfixedips.Subnet, reservedfixedips.AnySubnet, reservedfixedips.IPAddress),
				ValidateDiagFunc: func(val interface{}, key cty.Path) diag.Diagnostics {
					v := val.(string)
					switch reservedfixedips.ReservedFixedIPType(v) {
					case reservedfixedips.External, reservedfixedips.Subnet, reservedfixedips.AnySubnet, reservedfixedips.IPAddress:
						return diag.Diagnostics{}
					}
					return diag.Errorf("wrong type %s, available values is '%s', '%s', '%s', '%s'", v, reservedfixedips.External, reservedfixedips.Subnet, reservedfixedips.AnySubnet, reservedfixedips.IPAddress)
				},
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"fixed_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateDiagFunc: func(val interface{}, key cty.Path) diag.Diagnostics {
					v := val.(string)
					ip := net.ParseIP(v)
					if ip != nil {
						return diag.Diagnostics{}
					}

					return diag.FromErr(fmt.Errorf("%q must be a valid ip, got: %s", key, v))
				},
			},
			"subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"network_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"is_vip": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"port_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "ID of the port_id underlying the reserved fixed IP",
				Computed:    true,
			},
			"allowed_address_pairs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Group of IP addresses that share the current IP as VIP",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Optional: true,
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

func resourceReservedFixedIPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start ReservedFixedIP creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, reservedFixedIPsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := reservedfixedips.CreateOpts{
		IsVip: d.Get("is_vip").(bool),
	}

	portType := d.Get("type").(string)
	switch reservedfixedips.ReservedFixedIPType(portType) {
	case reservedfixedips.External:
	case reservedfixedips.Subnet:
		subnetID := d.Get("subnet_id").(string)
		if subnetID == "" {
			return diag.Errorf("'subnet_id' required if the type is 'subnet'")
		}

		opts.SubnetID = subnetID
	case reservedfixedips.AnySubnet:
		networkID := d.Get("network_id").(string)
		if networkID == "" {
			return diag.Errorf("'network_id' required if the type is 'any_subnet'")
		}
		opts.NetworkID = networkID
	case reservedfixedips.IPAddress:
		networkID := d.Get("network_id").(string)
		ipAddress := d.Get("fixed_ip_address").(string)
		if networkID == "" || ipAddress == "" {
			return diag.Errorf("'network_id' and 'fixed_ip_address' required if the type is 'ip_address'")
		}

		opts.NetworkID = networkID
		opts.IPAddress = net.ParseIP(ipAddress)
	default:
		return diag.Errorf("wrong type %s, available values is 'external', 'subnet', 'any_subnet', 'ip_address'", portType)
	}

	opts.Type = reservedfixedips.ReservedFixedIPType(portType)
	results, err := reservedfixedips.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	reservedFixedIPID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, ReservedFixedIPCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		reservedFixedIPID, err := reservedfixedips.ExtractReservedFixedIPIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve reservedFixedIP ID from task info: %w", err)
		}
		return reservedFixedIPID, nil
	})

	log.Printf("[DEBUG] ReservedFixedIP id (%s)", reservedFixedIPID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(reservedFixedIPID.(string))
	resourceReservedFixedIPRead(ctx, d, m)

	log.Printf("[DEBUG] Finish ReservedFixedIP creating (%s)", reservedFixedIPID)
	return diags
}

func resourceReservedFixedIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start ReservedFixedIP reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, reservedFixedIPsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	reservedFixedIP, err := reservedfixedips.Get(client, d.Id()).Extract()
	if err != nil {
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			log.Printf("[WARN] Removing reserved fixed ip %s because resource doesn't exist anymore", d.Id())
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("project_id", reservedFixedIP.ProjectID)
	d.Set("region_id", reservedFixedIP.RegionID)
	d.Set("status", reservedFixedIP.Status)
	d.Set("fixed_ip_address", reservedFixedIP.FixedIPAddress.String())
	d.Set("subnet_id", reservedFixedIP.SubnetID)
	d.Set("network_id", reservedFixedIP.NetworkID)
	d.Set("is_vip", reservedFixedIP.IsVip)
	d.Set("port_id", reservedFixedIP.PortID)

	allowedPairs := make([]map[string]interface{}, len(reservedFixedIP.AllowedAddressPairs))
	for i, p := range reservedFixedIP.AllowedAddressPairs {
		pair := make(map[string]interface{})

		pair["ip_address"] = p.IPAddress
		pair["mac_address"] = p.MacAddress

		allowedPairs[i] = pair
	}
	if err := d.Set("allowed_address_pairs", allowedPairs); err != nil {
		return diag.FromErr(err)
	}
	fields := []string{"type"}
	revertState(d, &fields)

	log.Println("[DEBUG] Finish ReservedFixedIP reading")
	return diags
}

func resourceReservedFixedIPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start ReservedFixedIP updating")
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, reservedFixedIPsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	if d.HasChange("is_vip") {
		opts := reservedfixedips.SwitchVIPOpts{IsVip: d.Get("is_vip").(bool)}
		_, err := reservedfixedips.SwitchVIP(client, id, opts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("allowed_address_pairs") {
		aap := d.Get("allowed_address_pairs").([]interface{})
		allowedAddressPairs := make([]reservedfixedips.AllowedAddressPairs, len(aap))
		for i, p := range aap {
			pair := p.(map[string]interface{})
			allowedAddressPairs[i] = reservedfixedips.AllowedAddressPairs{
				IPAddress:  pair["ip_address"].(string),
				MacAddress: pair["mac_address"].(string),
			}
		}

		clientPort, err := CreateClient(provider, d, portsPoint, versionPointV1)
		if err != nil {
			return diag.FromErr(err)
		}

		opts := ports.AllowAddressPairsOpts{AllowedAddressPairs: allowedAddressPairs}
		if _, err := ports.AllowAddressPairs(clientPort, id, opts).Extract(); err != nil {
			return diag.FromErr(err)
		}
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))
	log.Println("[DEBUG] Finish ReservedFixedIP updating")
	return resourceReservedFixedIPRead(ctx, d, m)
}

func resourceReservedFixedIPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start ReservedFixedIP deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, reservedFixedIPsPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	// only is_vip == false
	isVip := d.Get("is_vip").(bool)
	if isVip {
		return diag.Errorf("could not delete reserved fixed ip with is_vip=true")
	}

	id := d.Id()
	results, err := reservedfixedips.Delete(client, id).Extract()
	if err != nil {
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			d.SetId("")
			log.Printf("[DEBUG] Finish of ReservedFixedIP deleting")
			return diags
		default:
			return diag.FromErr(err)
		}
	}

	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, ReservedFixedIPCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := reservedfixedips.Get(client, id).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete reserved fixed ip with ID: %s", id)
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
	log.Printf("[DEBUG] Finish of ReservedFixedIP deleting")
	return diags
}
