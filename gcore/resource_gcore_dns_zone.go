package gcore

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DNSZoneResource = "gcore_dns_zone"

	DNSZoneSchemaName       = "name"
	DNSZoneSchemaDNSSECName = "dnssec"
)

func resourceDNSZone() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			DNSZoneSchemaName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					zoneName := i.(string)
					if strings.TrimSpace(zoneName) == "" || len(zoneName) > 255 {
						return diag.Errorf("dns name can't be empty, it also should be less than 256 symbols")
					}
					return nil
				},
				Description: "A name of DNS Zone resource.",
			},
			DNSZoneSchemaDNSSECName: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "Activation or deactivation of DNSSEC for the zone. By default, DNSSEC is disabled. " +
					"However, if this is set to true, DNSSEC will be enabled.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		CreateContext: checkDNSDependency(resourceDNSZoneCreate),
		ReadContext:   checkDNSDependency(resourceDNSZoneRead),
		UpdateContext: checkDNSDependency(resourceDNSZoneUpdate),
		DeleteContext: checkDNSDependency(resourceDNSZoneDelete),
		Description:   "Represent DNS zone resource. https://dns.gcore.com/zones",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func checkDNSDependency(next func(context.Context, *schema.ResourceData,
	interface{}) diag.Diagnostics) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {

	return func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
		config := i.(*Config)
		client := config.DNSClient
		if client == nil {
			return diag.Errorf("dns api client is null. make sure that you defined gcore_dns_api var in gcore provider section.")
		}
		return next(ctx, data, i)
	}
}

func resourceDNSZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := strings.TrimSpace(d.Get(DNSZoneSchemaName).(string))
	log.Println("[DEBUG] Start DNS Zone Resource creating")
	defer log.Printf("[DEBUG] Finish DNS Zone Resource creating (id=%s)\n", name)

	config := m.(*Config)
	client := config.DNSClient

	_, err := client.CreateZone(ctx, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("create zone: %v", err))
	}

	enableDnssec := d.Get(DNSZoneSchemaDNSSECName).(bool)
	if enableDnssec {
		_, err = client.ToggleDnssec(ctx, name, true)
		if err != nil {
			return diag.FromErr(fmt.Errorf("enable dnssec: %v", err))
		}
	}

	d.SetId(name)
	return resourceDNSZoneRead(ctx, d, m)
}

func resourceDNSZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := strings.TrimSpace(d.Get(DNSZoneSchemaName).(string))
	if d.Id() == "" {
		return diag.Errorf("empty id")
	}

	log.Printf("[DEBUG] Start DNS Zone Resource updating (id=%s)\n", name)
	defer log.Printf("[DEBUG] Finish DNS Zone Resource updating (id=%s)\n", name)

	config := m.(*Config)
	client := config.DNSClient

	enableDnssec := d.Get(DNSZoneSchemaDNSSECName).(bool)
	_, err := client.ToggleDnssec(ctx, name, enableDnssec)
	if err != nil {
		return diag.FromErr(fmt.Errorf("enable dnssec: %v", err))
	}

	return resourceDNSZoneRead(ctx, d, m)
}

func resourceDNSZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	zoneName := dnsZoneResourceID(d)
	log.Printf("[DEBUG] Start DNS Zone Resource reading (id=%s)\n", zoneName)
	defer log.Println("[DEBUG] Finish DNS Zone Resource reading")

	config := m.(*Config)
	client := config.DNSClient

	result, err := client.Zone(ctx, zoneName)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get zone: %w", err))
	}

	enableDnssec := d.Get(DNSZoneSchemaDNSSECName).(bool)
	if enableDnssec {
		_, errDnssecDS := client.DNSSecDS(ctx, zoneName)
		if errDnssecDS != nil {
			return diag.FromErr(fmt.Errorf("verify dnssec created: %w", errDnssecDS))
		}
	}

	d.SetId(result.Name)
	_ = d.Set(DNSZoneSchemaName, result.Name)

	return nil
}

func resourceDNSZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	zoneName := dnsZoneResourceID(d)
	log.Printf("[DEBUG] Start DNS Zone Resource deleting (id=%s)\n", zoneName)
	defer log.Println("[DEBUG] Finish DNS Zone Resource deleting")
	if zoneName == "" {
		return diag.Errorf("empty zone name")
	}

	config := m.(*Config)
	client := config.DNSClient

	err := client.DeleteZone(ctx, zoneName)
	if err != nil {
		return diag.FromErr(fmt.Errorf("delete zone: %w", err))
	}
	d.SetId("")

	return nil
}

func dnsZoneResourceID(d *schema.ResourceData) string {
	resourceID := d.Id()
	if resourceID == "" {
		resourceID = d.Get(DNSZoneSchemaName).(string)
	}
	return resourceID
}
