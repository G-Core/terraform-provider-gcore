package gcore

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/G-Core/gcorelabscdn-go/clients"
)

func resourceCDNClientConfig() *schema.Resource {
	return &schema.Resource{
		Description: "Manage your CDN client (only utilization_level can be updated).",

		ReadContext:   resourceClientRead,
		UpdateContext: resourceClientUpdate,
		CreateContext: resourceClientCreate,
		DeleteContext: resourceClientDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account ID.",
			},
			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain zone to which a CNAME record of your CDN resources should be pointed.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date of the first synchronization with the Platform (ISO 8601/RFC 3339 format, UTC.)",
			},
			"updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date of the last update of information about CDN service (ISO 8601/RFC 3339 format, UTC.)",
			},
			"utilization_level": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "CDN traffic usage limit in gigabytes.",
			},
			"use_balancer": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Defines whether custom balancing is used for content delivery.",
			},
			"auto_suspend_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Defines whether resources will be deactivated automatically by inactivity.",
			},
			"cdn_resources_rules_max_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Limit on the number of rules for each CDN resource.",
			},
			"service": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the CDN service status.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Defines whether the CDN service is activated.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CDN service status.",
						},
						"updated": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date of the last CDN service status update (ISO 8601/RFC 3339 format, UTC).",
						},
					},
				},
			},
		},
	}
}

func resourceClientCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceClientUpdate(ctx, d, m)
}

func resourceClientRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.CDNClient

	result, err := client.ClientsMe().Get(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", result.ID))
	d.Set("cname", result.Cname)
	d.Set("created", result.Created.Format(time.RFC3339))
	d.Set("updated", result.Updated.Format(time.RFC3339))
	d.Set("utilization_level", result.UtilizationLevel)
	d.Set("use_balancer", result.UseBalancer)
	d.Set("auto_suspend_enabled", result.AutoSuspendEnabled)
	d.Set("cdn_resources_rules_max_count", result.CDNResourcesRulesMaxCount)

	d.Set("service", []interface{}{
		map[string]interface{}{
			"enabled": result.Service.Enabled,
			"status":  result.Service.Status,
			"updated": result.Service.Updated.Format(time.RFC3339),
		},
	})

	return nil
}

func resourceClientUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.CDNClient

	req := &clients.ClientsMeUpdateRequest{
		UtilizationLevel: d.Get("utilization_level").(int),
	}

	_, err := client.ClientsMe().Update(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceClientRead(ctx, d, m)
}

func resourceClientDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
