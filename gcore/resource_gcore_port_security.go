package gcore

import (
	"context"
	"log"

	ports "github.com/G-Core/gcorelabscloud-go/gcore/port/v1/ports"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePortSecurity() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePortSecurityUpdate,
		ReadContext:   resourcePortSecurityRead,
		UpdateContext: resourcePortSecurityUpdate,
		DeleteContext: resourcePortSecurityDelete,
		Description:   `Manages security for any private port attached to a Virtual Instance.`,
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
			"enable_port_security": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable port security for the port.",
				Default:     true,
			},
		},
	}
}

func resourcePortSecurityUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Port Security updating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	clientPort, err := CreateClient(provider, d, "ports", "v1")
	if err != nil {
		return diag.FromErr(err)
	}

	portID := d.Get("port_id").(string)
	enablePortSecurity := d.Get("enable_port_security").(bool)
	log.Println("[DEBUG] Enable port security: ", enablePortSecurity)
	if enablePortSecurity {
		_, err := ports.EnablePortSecurity(clientPort, portID).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		_, err := ports.DisablePortSecurity(clientPort, portID).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(portID)

	return diags
}

func resourcePortSecurityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourcePortSecurityDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}
