package gcore

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	dnssdk "github.com/G-Core/gcore-dns-sdk-go"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DNSZoneResource = "gcore_dns_zone"

	DNSZoneSchemaName           = "name"
	DNSZoneSchemaDNSSEC         = "dnssec"
	DNSZoneSchemaEnabled        = "enabled"
	DNSZoneSchemaClientID       = "client_id"
	DNSZoneSchemaContact        = "contact"
	DNSZoneSchemaExpiry         = "expiry"
	DNSZoneSchemaMeta           = "meta"
	DNSZoneSchemaNX_TTL         = "nx_ttl"
	DNSZoneSchemaPrimary_server = "primary_server"
	DNSZoneSchemaRecords        = "records"
	DNSZoneSchemaRefresh        = "refresh"
	DNSZoneSchemaRetry          = "retry"
	DNSZoneSchemaRRSetsAmount   = "rrsets_amount"
	DNSZoneSchemaSerial         = "serial"
	DNSZoneSchemaStatus         = "status"
	// DNSZoneSchemaImportFileContent = "import_file_content"
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
			// DNSZoneSchemaImportFileContent: {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "Content of the BIND file to import zone from. Zone will be imported on creation.",
			// },
			DNSZoneSchemaDNSSEC: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "Activation or deactivation of DNSSEC for the zone." +
					"Set it to true to enable DNSSEC for the zone or false to disable it." +
					"By default, DNSSEC is set to false wich means it is disabled.",
			},
			DNSZoneSchemaEnabled: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Default: true. If a zone is disabled, then its records will not be resolved on dns servers",
			},
			DNSZoneSchemaContact: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email address of the administrator responsible for this zone",
			},
			DNSZoneSchemaExpiry: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "number of seconds after which secondary name servers should stop answering request for this zone",
			},
			DNSZoneSchemaMeta: {
				Type:     schema.TypeMap,
				Optional: true,
				Description: "Arbitrary data of zone in JSON format. " +
					"You can specify webhook URL and webhook_method here. " +
					"Webhook will receive a map with three arrays: for created, updated, and deleted rrsets. " +
					"webhook_method can be omitted; POST will be used by default.",
			},
			DNSZoneSchemaNX_TTL: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Time To Live of cache",
			},
			DNSZoneSchemaPrimary_server: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Primary master name server for zone",
			},
			DNSZoneSchemaRefresh: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "number of seconds after which secondary name servers should refresh the zone",
			},
			DNSZoneSchemaRetry: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "number of seconds after which secondary name servers should retry to request the serial number",
			},
			DNSZoneSchemaSerial: {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Serial number for this zone or Timestamp of zone modification moment. " +
					"If a secondary name server slaved to this one observes an increase in this number, " +
					"the slave will assume that the zone has been updated and initiate a zone transfer.",
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
	zoneName := dnsZoneResourceID(d)
	log.Println("[DEBUG] Start DNS Zone Resource creating")
	defer log.Printf("[DEBUG] Finish DNS Zone Resource creating (id=%s)\n", zoneName)

	config := m.(*Config)
	client := config.DNSClient

	addZone := dnssdk.AddZone{
		Name:          zoneName,
		Contact:       d.Get(DNSZoneSchemaContact).(string),
		Expiry:        uint64(d.Get(DNSZoneSchemaExpiry).(int)),
		Enabled:       d.Get(DNSZoneSchemaEnabled).(bool),
		Meta:          d.Get(DNSZoneSchemaMeta).(map[string]interface{}),
		NxTTL:         uint64(d.Get(DNSZoneSchemaNX_TTL).(int)),
		PrimaryServer: d.Get(DNSZoneSchemaPrimary_server).(string),
		Refresh:       uint64(d.Get(DNSZoneSchemaRefresh).(int)),
		Retry:         uint64(d.Get(DNSZoneSchemaRetry).(int)),
		Serial:        uint64(d.Get(DNSZoneSchemaSerial).(int)),
	}

	_, err := client.CreateZone(ctx, addZone)
	if err != nil {
		return diag.FromErr(fmt.Errorf("create zone: %v", err))
	}

	// if importFileContent, ok := d.GetOk(DNSZoneSchemaImportFileContent); ok {
	// 	_, err = client.ImportZone(ctx, zoneName, importFileContent.(string))
	// 	if err != nil {
	// 		return diag.FromErr(fmt.Errorf("import zone from file content: %w", err))
	// 	}
	// }

	enableDnssec := d.Get(DNSZoneSchemaDNSSEC).(bool)
	if enableDnssec {
		_, err = client.ToggleDnssec(ctx, zoneName, true)
		if err != nil {
			return diag.FromErr(fmt.Errorf("enable dnssec: %v", err))
		}
	}

	d.SetId(zoneName)

	return resourceDNSZoneRead(ctx, d, m)
}

func resourceDNSZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	zoneName := dnsZoneResourceID(d)
	if d.Id() == "" {
		return diag.Errorf("empty id")
	}

	log.Printf("[DEBUG] Start DNS Zone Resource updating (id=%s)\n", zoneName)
	defer log.Printf("[DEBUG] Finish DNS Zone Resource updating (id=%s)\n", zoneName)

	config := m.(*Config)
	client := config.DNSClient

	if d.HasChange(DNSZoneSchemaDNSSEC) {
		enableDnssec := d.Get(DNSZoneSchemaDNSSEC).(bool)
		_, err := client.ToggleDnssec(ctx, zoneName, enableDnssec)
		if err != nil {
			return diag.FromErr(fmt.Errorf("enable dnssec: %v", err))
		}
	}

	if d.HasChange(DNSZoneSchemaEnabled) {
		enabled := d.Get(DNSZoneSchemaEnabled).(bool)
		if enabled {
			err := client.EnableZone(ctx, zoneName)
			if err != nil {
				return diag.FromErr(fmt.Errorf("enabling zone: %w", err))
			}
		} else {
			err := client.DisableZone(ctx, zoneName)
			if err != nil {
				return diag.FromErr(fmt.Errorf("disabling zone: %w", err))
			}
		}
	}
	hasChangesForUpdateZone := d.HasChange(DNSZoneSchemaContact) ||
		d.HasChange(DNSZoneSchemaExpiry) ||
		d.HasChange(DNSZoneSchemaMeta) ||
		d.HasChange(DNSZoneSchemaNX_TTL) ||
		d.HasChange(DNSZoneSchemaPrimary_server) ||
		d.HasChange(DNSZoneSchemaRefresh) ||
		d.HasChange(DNSZoneSchemaRetry) ||
		d.HasChange(DNSZoneSchemaSerial)

	if hasChangesForUpdateZone {
		updateZone := dnssdk.AddZone{
			Name:          zoneName,
			Contact:       d.Get(DNSZoneSchemaContact).(string),
			Expiry:        uint64(d.Get(DNSZoneSchemaExpiry).(int)),
			Enabled:       d.Get(DNSZoneSchemaEnabled).(bool),
			Meta:          d.Get(DNSZoneSchemaMeta).(map[string]interface{}),
			NxTTL:         uint64(d.Get(DNSZoneSchemaNX_TTL).(int)),
			PrimaryServer: d.Get(DNSZoneSchemaPrimary_server).(string),
			Refresh:       uint64(d.Get(DNSZoneSchemaRefresh).(int)),
			Retry:         uint64(d.Get(DNSZoneSchemaRetry).(int)),
			Serial:        uint64(d.Get(DNSZoneSchemaSerial).(int)),
		}

		_, err := client.UpdateZone(ctx, zoneName, updateZone)
		if err != nil {
			log.Printf("[DEBUG] Error updating zone: %v", err)
			return diag.FromErr(fmt.Errorf("update zone: %v", err))
		}
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
		return diag.FromErr(err)
	}

	enableDnssec := result.DNSSECEnabled
	if enableDnssec {
		_, errDnssecDS := client.DNSSecDS(ctx, zoneName)
		if errDnssecDS != nil {
			return diag.FromErr(fmt.Errorf("verify dnssec created: %w", errDnssecDS))
		}
	}

	d.SetId(result.Name)
	d.Set(DNSZoneSchemaName, result.Name)
	d.Set(DNSZoneSchemaDNSSEC, result.DNSSECEnabled)
	d.Set(DNSZoneSchemaClientID, result.ClientID)
	d.Set(DNSZoneSchemaContact, result.Contact)
	d.Set(DNSZoneSchemaExpiry, result.Expiry)
	d.Set(DNSZoneSchemaMeta, result.Meta)
	d.Set(DNSZoneSchemaNX_TTL, result.NxTTL)
	d.Set(DNSZoneSchemaPrimary_server, result.PrimaryServer)
	d.Set(DNSZoneSchemaRecords, result.Records)
	d.Set(DNSZoneSchemaRefresh, result.Refresh)
	d.Set(DNSZoneSchemaRetry, result.Retry)
	d.Set(DNSZoneSchemaRRSetsAmount, result.RRSetsAmount)
	d.Set(DNSZoneSchemaSerial, result.Serial)
	d.Set(DNSZoneSchemaStatus, result.Status)

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
