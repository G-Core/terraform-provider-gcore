package gcore

import (
	"context"
	"fmt"
	"log"
	"strings"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	portsv1 "github.com/G-Core/gcorelabscloud-go/gcore/port/v1/ports"
	portsv2 "github.com/G-Core/gcorelabscloud-go/gcore/port/v2/ports"
	reservedfixedips "github.com/G-Core/gcorelabscloud-go/gcore/reservedfixedip/v1/reservedfixedips"
	tasksv1 "github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGcorePortAllowedAddressPairs() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePortAAPCreate,
		ReadContext:   resourcePortAAPRead,
		UpdateContext: resourcePortAAPUpdate,
		DeleteContext: resourcePortAAPDelete,
		Description:   "Manages allowed address pairs of any port resource.",
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "Region ID where the port is located.",
				ConflictsWith: []string{"region_name"},
			},
			"region_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Region name where the port is located.",
				ConflictsWith: []string{"region_id"},
			},
			"project_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "Project ID in which the port exists.",
				ConflictsWith: []string{"project_name"},
			},
			"project_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Project name in which the port exists.",
				ConflictsWith: []string{"project_id"},
			},
			"port_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the port where allowed address pairs will be managed.",
			},
			"allowed_address_pair": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of allowed address pairs to associate with the port.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP address or subnet of the allowed address pair (e.g., 192.0.2.10/32).",
						},
						"mac_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "MAC address of the allowed address pair. If omitted, port's MAC address will be used.",
						},
					},
				},
			},
		},
	}
}

func resourcePortAAPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Port AllowedAddressPairs creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	clientPort, err := CreateClient(provider, d, "ports", "v2")
	if err != nil {
		return diag.FromErr(err)
	}
	taskClient, err := CreateClient(provider, d, "tasks", "v1")
	if err != nil {
		return diag.FromErr(err)
	}

	portID := d.Get("port_id").(string)
	pairs := expandAAP(d.Get("allowed_address_pair").([]interface{}))

	if err := putAllowedAddressPairs(ctx, clientPort, taskClient, portID, pairs); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(portID)
	return diags
}

func resourcePortAAPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourcePortAAPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Port AllowedAddressPairs updating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	clientPort, err := CreateClient(provider, d, "ports", "v2")
	if err != nil {
		return diag.FromErr(err)
	}
	taskClient, err := CreateClient(provider, d, "tasks", "v1")
	if err != nil {
		return diag.FromErr(err)
	}

	portID := d.Get("port_id").(string)
	pairs := expandAAP(d.Get("allowed_address_pair").([]interface{}))

	if err := putAllowedAddressPairs(ctx, clientPort, taskClient, portID, pairs); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourcePortAAPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Port AllowedAddressPairs deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	clientPort, err := CreateClient(provider, d, "ports", "v2")
	if err != nil {
		return diag.FromErr(err)
	}
	taskClient, err := CreateClient(provider, d, "tasks", "v1")
	if err != nil {
		return diag.FromErr(err)
	}

	portID := d.Id()
	empty := make([]reservedfixedips.AllowedAddressPairs, 0)
	if err := putAllowedAddressPairs(ctx, clientPort, taskClient, portID, empty); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func putAllowedAddressPairs(ctx context.Context, clientPort *gcorecloud.ServiceClient, taskClient *gcorecloud.ServiceClient, portID string, pairs []reservedfixedips.AllowedAddressPairs) error {
	opts := portsv1.AllowAddressPairsOpts{AllowedAddressPairs: pairs}
	res := portsv2.AllowAddressPairs(clientPort, portID, opts)
	tr, err := res.Extract()
	if err != nil {
		if strings.Contains(err.Error(), "already has allowed address pairs set as requested") {
			return nil
		}
		return fmt.Errorf("allow_address_pairs PUT failed: %w", err)
	}
	if tr != nil && len(tr.Tasks) > 0 {
		for _, tid := range tr.Tasks {
			if _, err := tasksv1.WaitTaskAndReturnResult(taskClient, tid, true, 600, func(t tasksv1.TaskID) (interface{}, error) { return nil, nil }); err != nil {
				return fmt.Errorf("task %s did not finish successfully: %w", string(tid), err)
			}
		}
	}
	return nil
}

func expandAAP(list []interface{}) []reservedfixedips.AllowedAddressPairs {
	if len(list) == 0 {
		return nil
	}
	out := make([]reservedfixedips.AllowedAddressPairs, 0, len(list))
	for _, raw := range list {
		m := raw.(map[string]interface{})
		aap := reservedfixedips.AllowedAddressPairs{IPAddress: m["ip_address"].(string)}
		if mac, ok := m["mac_address"].(string); ok && mac != "" {
			aap.MacAddress = mac
		}
		out = append(out, aap)
	}
	return out
}
