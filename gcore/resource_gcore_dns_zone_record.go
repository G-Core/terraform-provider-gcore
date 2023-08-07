package gcore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	dnssdk "github.com/G-Core/gcore-dns-sdk-go"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DNSZoneRecordResource = "gcore_dns_zone_record"

	DNSZoneRecordSchemaZone   = "zone"
	DNSZoneRecordSchemaDomain = "domain"
	DNSZoneRecordSchemaType   = "type"
	DNSZoneRecordSchemaTTL    = "ttl"
	DNSZoneRecordSchemaFilter = "filter"

	DNSZoneRecordSchemaFilterLimit  = "limit"
	DNSZoneRecordSchemaFilterType   = "type"
	DNSZoneRecordSchemaFilterStrict = "strict"

	DNSZoneRecordSchemaResourceRecord = "resource_record"
	DNSZoneRecordSchemaContent        = "content"
	DNSZoneRecordSchemaEnabled        = "enabled"

	// DNSZoneRecordSchemaMeta when adding this list, add also on dnsZoneRecordSchemaMetaList below
	DNSZoneRecordSchemaMeta           = "meta"
	DNSZoneRecordSchemaMetaAsn        = "asn"
	DNSZoneRecordSchemaMetaIP         = "ip"
	DNSZoneRecordSchemaMetaCountries  = "countries"
	DNSZoneRecordSchemaMetaContinents = "continents"
	DNSZoneRecordSchemaMetaLatLong    = "latlong"
	DNSZoneRecordSchemaMetaNotes      = "notes"
	DNSZoneRecordSchemaMetaDefault    = "default"
	DNSZoneRecordSchemaMetaFailover   = "failover"
	DNSZoneRecordSchemaMetaWeight     = "weight"
	DNSZoneRecordSchemaMetaBackup     = "backup"
	DNSZoneRecordSchemaMetaFallback   = "fallback"

	// DNSZoneRRSetSchemaMeta failover meta is inside rrset, not inside resource record
	DNSZoneRRSetSchemaMeta = "meta"

	DNSZoneRRSetSchemaMetaGeodnsLink = "geodns_link"

	// DNSZoneRRSetSchemaMetaFailover - 10. `meta` (map)
	DNSZoneRRSetSchemaMetaFailover               = "failover" // backward compatibility
	DNSZoneRRSetSchemaMetaHealthchecks           = "healthchecks"
	DNSZoneRRSetSchemaMetaFailoverFrequency      = "frequency"
	DNSZoneRRSetSchemaMetaFailoverHost           = "host"
	DNSZoneRRSetSchemaMetaFailoverHTTPStatusCode = "http_status_code"
	DNSZoneRRSetSchemaMetaFailoverCommand        = "command"
	DNSZoneRRSetSchemaMetaFailoverMethod         = "method"
	DNSZoneRRSetSchemaMetaFailoverPort           = "port"
	DNSZoneRRSetSchemaMetaFailoverProtocol       = "protocol"
	DNSZoneRRSetSchemaMetaFailoverRegexp         = "regexp"
	DNSZoneRRSetSchemaMetaFailoverTimeout        = "timeout"
	DNSZoneRRSetSchemaMetaFailoverTLS            = "tls"
	DNSZoneRRSetSchemaMetaFailoverURL            = "url"
)

var dnsZoneRecordSchemaMetaList = []string{
	DNSZoneRecordSchemaMetaAsn,
	DNSZoneRecordSchemaMetaIP,
	DNSZoneRecordSchemaMetaCountries,
	DNSZoneRecordSchemaMetaContinents,
	DNSZoneRecordSchemaMetaLatLong,
	DNSZoneRecordSchemaMetaNotes,
	DNSZoneRecordSchemaMetaDefault,
	DNSZoneRecordSchemaMetaFailover,
	DNSZoneRecordSchemaMetaWeight,
	DNSZoneRecordSchemaMetaBackup,
	DNSZoneRecordSchemaMetaFallback,
}

func resourceDNSZoneRecord() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			DNSZoneRecordSchemaZone: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					val := i.(string)
					if strings.TrimSpace(val) == "" || len(val) > 255 {
						return diag.Errorf("dns record zone can't be empty, it also should be less than 256 symbols")
					}
					return nil
				},
				Description: "A zone of DNS Zone Record resource.",
			},
			DNSZoneRecordSchemaDomain: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					val := i.(string)
					if strings.TrimSpace(val) == "" || len(val) > 255 {
						return diag.Errorf("dns record domain can't be empty, it also should be less than 256 symbols")
					}
					return nil
				},
				Description: "A domain of DNS Zone Record resource.",
			},
			DNSZoneRecordSchemaType: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					val := strings.TrimSpace(i.(string))
					types := []string{"A", "AAAA", "MX", "CNAME", "TXT", "CAA", "NS", "SRV", "HTTPS", "SVCB"}
					for _, t := range types {
						if strings.EqualFold(t, val) {
							return nil
						}
					}
					return diag.Errorf("dns record type should be one of %v", types)

				},
				Description: "A type of DNS Zone Record resource.",
			},
			DNSZoneRecordSchemaTTL: {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					val := i.(int)
					if val < 0 {
						return diag.Errorf("dns record ttl can't be less than 0")
					}
					return nil
				},
				Description: "A ttl of DNS Zone Record resource.",
			},
			DNSZoneRecordSchemaFilter: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						DNSZoneRecordSchemaFilterLimit: {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "A DNS Zone Record filter option that describe how many records will be percolated.",
						},
						DNSZoneRecordSchemaFilterStrict: {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A DNS Zone Record filter option that describe possibility to return answers if no records were percolated through filter.",
						},
						DNSZoneRecordSchemaFilterType: {
							Type:     schema.TypeString,
							Required: true,
							ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
								names := []string{"geodns", "geodistance", "default", "first_n", "is_healthy"}
								name := i.(string)
								for _, n := range names {
									if n == name {
										return nil
									}
								}
								return diag.Errorf("dns record filter type should be one of %v", names)
							},
							Description: "A DNS Zone Record filter option that describe a name of filter.",
						},
					},
				},
			},
			DNSZoneRecordSchemaResourceRecord: {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						DNSZoneRecordSchemaContent: {
							Type:        schema.TypeString,
							Required:    true,
							Description: `A content of DNS Zone Record resource. (TXT: 'anyString', MX: '50 mail.company.io.', CAA: '0 issue "company.org; account=12345"')`,
						},
						DNSZoneRecordSchemaEnabled: {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Manage of public appearing of DNS Zone Record resource.",
						},
						DNSZoneRecordSchemaMeta: {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									DNSZoneRecordSchemaMetaAsn: {
										Type: schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
											ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
												if i.(int) < 0 {
													return diag.Errorf("asn cannot be less then 0")
												}
												return nil
											},
										},
										Optional:    true,
										Description: "An asn meta (eg. 12345) of DNS Zone Record resource.",
									},
									DNSZoneRecordSchemaMetaIP: {
										Type: schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
												val := i.(string)
												_, _, err := net.ParseCIDR(val)
												if err != nil {
													if ip := net.ParseIP(val); ip == nil {
														return diag.Errorf("dns record meta ip has wrong format: %s: %v", val, err)
													}
												}
												return nil
											},
										},
										Optional:    true,
										Description: "An ip meta (eg. 127.0.0.0) of DNS Zone Record resource.",
									},
									DNSZoneRecordSchemaMetaLatLong: {
										Optional: true,
										Type:     schema.TypeList,
										MaxItems: 2,
										MinItems: 2,
										Elem: &schema.Schema{
											Type: schema.TypeFloat,
										},
										Description: "A latlong meta (eg. 27.988056, 86.925278) of DNS Zone Record resource.",
									},
									DNSZoneRecordSchemaMetaNotes: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A notes meta (eg. Miami DC) of DNS Zone Record resource.",
									},
									DNSZoneRecordSchemaMetaContinents: {
										Type: schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Continents meta (eg. Asia) of DNS Zone Record resource.",
									},
									DNSZoneRecordSchemaMetaCountries: {
										Type: schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Countries meta (eg. USA) of DNS Zone Record resource.",
									},
									DNSZoneRecordSchemaMetaDefault: {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Fallback meta equals true marks records which are used as a default answer (when nothing was selected by specified meta fields).",
									},
									DNSZoneRecordSchemaMetaFailover: {
										Type:        schema.TypeMap,
										Optional:    true,
										Computed:    true,
										Description: "Computed UUID of failover healtcheck property",
									},
									DNSZoneRecordSchemaMetaWeight: {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "A weight for this record",
									},
									DNSZoneRecordSchemaMetaFallback: {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Set as fallback record",
									},
									DNSZoneRecordSchemaMetaBackup: {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Set as backup record",
									},
								},
							},
						},
					},
				},
				Description: "An array of contents with meta of DNS Zone Record resource.",
			},

			DNSZoneRRSetSchemaMeta: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						DNSZoneRRSetSchemaMetaGeodnsLink: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Geodns link (domain, or cl-) of DNS Zone RRSet resource.",
						},
						// - 1. `meta` (map)
						DNSZoneRRSetSchemaMetaHealthchecks: {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: `Failover meta (eg. {"frequency": 60,"host": "www.gcore.com","http_status_code": null,"method": "GET","port": 80,"protocol": "HTTP","regexp": "","timeout": 10,"tls": false,"url": "/"}).`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									DNSZoneRRSetSchemaMetaFailoverFrequency: {
										Type:        schema.TypeInt,
										Description: "Frequency in seconds (10-3600).",
										Required:    true,
										ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
											v, dtv := toInt(i)
											if dtv != dtvInt {
												return diag.Errorf("dns record meta failover frequency %v must be integer, got: %v", path, i)
											}
											if v < 10 || v > 3600 {
												return diag.Errorf("dns record meta failover frequency %v must be in range 10-3600, got: %v", path, i)
											}
											return nil
										},
									},
									DNSZoneRRSetSchemaMetaFailoverHost: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Request host/virtualhost to send if protocol=HTTP, must be empty for non-HTTP",
									},
									DNSZoneRRSetSchemaMetaFailoverCommand: {
										Type:        schema.TypeString,
										Description: "Command to send if protocol=TCP/UDP, maximum length: 255.",
										Optional:    true,
										ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
											v, _ := toString(i)
											if len(v) > 255 {
												return diag.Errorf("dns record meta failover command %v must be less then 255, got: %v", path, i)
											}
											return nil
										},
									},
									DNSZoneRRSetSchemaMetaFailoverHTTPStatusCode: {
										Type:        schema.TypeInt,
										Description: "Expected status code if protocol=HTTP, must be empty for non-HTTP.",
										Optional:    true,
									},
									DNSZoneRRSetSchemaMetaFailoverMethod: {
										Type:        schema.TypeString,
										Description: "HTTP Method required if protocol=HTTP, must be empty for non-HTTP.",
										Optional:    true,
									},
									DNSZoneRRSetSchemaMetaFailoverPort: {
										Type:        schema.TypeInt,
										Description: "Port to check (1-65535).",
										Optional:    true,
									},
									DNSZoneRRSetSchemaMetaFailoverProtocol: {
										Type:        schema.TypeString,
										Description: "Protocol, possible value: HTTP, TCP, UDP, ICMP.",
										Required:    true,
										ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
											v, dtv := toString(i)
											if dtv != dtvString {
												return diag.Errorf("dns record meta failover protocol %v must be string, got: %v", path, i)
											}
											if v != "HTTP" && v != "TCP" && v != `UDP` && v != "ICMP" {
												return diag.Errorf("dns record meta failover protocol %v must be one of HTTP, TCP, UDP, ICMP, got: %v", path, i)
											}
											return nil
										},
									},
									DNSZoneRRSetSchemaMetaFailoverRegexp: {
										Type:        schema.TypeString,
										Description: "HTTP body or response payload to check if protocol<>ICMP, must be empty for ICMP.",
										Optional:    true,
									},
									DNSZoneRRSetSchemaMetaFailoverTimeout: {
										Type:        schema.TypeInt,
										Description: "Timeout in seconds (1-10).",
										Required:    true,
										ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
											v, dtv := toInt(i)
											if dtv != dtvInt {
												return diag.Errorf("dns record meta failover timeout %v must be an integer, got %v", path, i)
											}
											if v < 1 || v > 10 {
												return diag.Errorf("dns record meta failover timeout %v must be between 1 and 10, got %v", path, i)
											}
											return nil
										},
									},
									DNSZoneRRSetSchemaMetaFailoverTLS: {
										Type:        schema.TypeBool,
										Description: "TLS/HTTPS enabled if protocol=HTTP, must be empty for non-HTTP.",
										Optional:    true,
									},
									DNSZoneRRSetSchemaMetaFailoverURL: {
										Type:        schema.TypeString,
										Description: "URL path to check required if protocol=HTTP, must be empty for non-HTTP.",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		CreateContext: checkDNSDependency(resourceDNSZoneRecordCreate),
		UpdateContext: checkDNSDependency(resourceDNSZoneRecordUpdate),
		ReadContext:   checkDNSDependency(resourceDNSZoneRecordRead),
		DeleteContext: checkDNSDependency(resourceDNSZoneRecordDelete),
		Description:   "Represent DNS Zone Record resource. https://dns.gcore.com/zones",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), ":")
				if len(parts) != 3 {
					return nil, fmt.Errorf("format must be as zone:domain:type")
				}
				_ = d.Set(DNSZoneRecordSchemaZone, parts[0])
				d.SetId(parts[0])
				_ = d.Set(DNSZoneRecordSchemaDomain, parts[1])
				_ = d.Set(DNSZoneRecordSchemaType, parts[2])

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceDNSZoneRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	zone := strings.TrimSpace(d.Get(DNSZoneRecordSchemaZone).(string))
	domain := strings.TrimSpace(d.Get(DNSZoneRecordSchemaDomain).(string))
	rType := strings.TrimSpace(d.Get(DNSZoneRecordSchemaType).(string))
	log.Println("[DEBUG] Start DNS Zone Record Resource creating")
	defer log.Printf("[DEBUG] Finish DNS Zone Record Resource creating (id=%s %s %s)\n", zone, domain, rType)

	ttl := d.Get(DNSZoneRecordSchemaTTL).(int)
	rrSet := dnssdk.RRSet{TTL: ttl, Records: make([]dnssdk.ResourceRecord, 0)}
	err := fillRRSet(d, rType, &rrSet)
	if err != nil {
		return diag.FromErr(err)
	}

	config := m.(*Config)
	client := config.DNSClient

	_, err = client.Zone(ctx, zone)
	if err != nil {
		return diag.FromErr(fmt.Errorf("find zone: %w", err))
	}

	err = client.CreateRRSet(ctx, zone, domain, rType, rrSet)
	if err != nil {
		return diag.FromErr(fmt.Errorf("create zone rrset: %v", err))
	}
	d.SetId(zone)

	return resourceDNSZoneRecordRead(ctx, d, m)
}

func resourceDNSZoneRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Id() == "" {
		return diag.Errorf("empty id")
	}
	zone := strings.TrimSpace(d.Get(DNSZoneRecordSchemaZone).(string))
	domain := strings.TrimSpace(d.Get(DNSZoneRecordSchemaDomain).(string))
	rType := strings.TrimSpace(d.Get(DNSZoneRecordSchemaType).(string))
	log.Println("[DEBUG] Start DNS Zone Record Resource updating")
	defer log.Printf("[DEBUG] Finish DNS Zone Record Resource updating (id=%s %s %s)\n", zone, domain, rType)

	ttl := d.Get(DNSZoneRecordSchemaTTL).(int)
	rrSet := dnssdk.RRSet{TTL: ttl, Records: make([]dnssdk.ResourceRecord, 0)}
	err := fillRRSet(d, rType, &rrSet)
	if err != nil {
		return diag.FromErr(err)
	}

	config := m.(*Config)
	client := config.DNSClient

	err = client.UpdateRRSet(ctx, zone, domain, rType, rrSet)
	if err != nil {
		return diag.FromErr(fmt.Errorf("update zone rrset: %v", err))
	}
	d.SetId(zone)

	return resourceDNSZoneRecordRead(ctx, d, m)
}

func resourceDNSZoneRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Id() == "" {
		return diag.Errorf("empty id")
	}
	zone := strings.TrimSpace(d.Get(DNSZoneRecordSchemaZone).(string))
	domain := strings.TrimSpace(d.Get(DNSZoneRecordSchemaDomain).(string))
	rType := strings.TrimSpace(d.Get(DNSZoneRecordSchemaType).(string))
	log.Println("[DEBUG] Start DNS Zone Record Resource reading")
	defer log.Printf("[DEBUG] Finish DNS Zone Record Resource reading (id=%s %s %s)\n", zone, domain, rType)

	config := m.(*Config)
	client := config.DNSClient

	result, err := client.RRSet(ctx, zone, domain, rType)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get zone rrset: %w", err))
	}
	id := struct{ Zone, Domain, Type string }{zone, domain, rType}
	bs, _ := json.Marshal(id)
	d.SetId(string(bs))
	_ = d.Set(DNSZoneRecordSchemaZone, zone)
	_ = d.Set(DNSZoneRecordSchemaDomain, domain)
	_ = d.Set(DNSZoneRecordSchemaType, rType)
	_ = d.Set(DNSZoneRecordSchemaTTL, result.TTL)

	filters := make([]map[string]interface{}, 0)
	for _, f := range result.Filters {
		filters = append(filters, map[string]interface{}{
			DNSZoneRecordSchemaFilterLimit:  f.Limit,
			DNSZoneRecordSchemaFilterType:   f.Type,
			DNSZoneRecordSchemaFilterStrict: f.Strict,
		})
	}
	if len(filters) > 0 {
		_ = d.Set(DNSZoneRecordSchemaFilter, filters)
	}

	// rr and meta of rr
	log.Printf("result.Records: %v\n", result.Records)
	rr := make([]map[string]interface{}, 0)
	rrMetaValidKeys := map[string]bool{}
	for _, metaKey := range dnsZoneRecordSchemaMetaList {
		rrMetaValidKeys[metaKey] = true
	}
	for _, rec := range result.Records {
		r := map[string]any{}
		r[DNSZoneRecordSchemaEnabled] = rec.Enabled
		log.Printf("rr: %v\n", rec)
		log.Printf("rr.ContentToString: %v\n", rec.ContentToString())
		r[DNSZoneRecordSchemaContent] = rec.ContentToString()
		meta := map[string]interface{}{}
		for key, val := range rec.Meta {
			// skip import failover since it's not configurable
			if key == DNSZoneRecordSchemaMetaFailover {
				continue
			}
			// skip if empty
			if val == nil {
				continue
			}
			log.Printf("rr.meta[%s] = %v\n", key, val)
			if !rrMetaValidKeys[key] {
				log.Printf("WARNING: rr.meta[%s] is not list of valid keys, create PR to https://github.com/G-Core/terraform-provider-gcore\n", key)
				continue
			}

			// if array
			valArr, ok := val.([]any)
			if ok {
				if len(valArr) > 0 {
					if key == DNSZoneRecordSchemaMetaNotes {
						meta[key] = valArr[0]
						continue
					}
				}
			}
			meta[key] = val
		}
		if len(meta) > 0 {
			r[DNSZoneRecordSchemaMeta] = []map[string]interface{}{meta}
		} else {
			r[DNSZoneRecordSchemaMeta] = nil
		}
		rr = append(rr, r)
	}
	if len(rr) > 0 {
		_ = d.Set(DNSZoneRecordSchemaResourceRecord, rr)
	}

	// meta of RRSet
	rrMeta := map[string]any{}
	for key, val := range result.Meta {
		v, ok := val.(map[string]any)
		if ok {
			if key == DNSZoneRRSetSchemaMetaFailover { // copy to healthcheck
				rrMeta[DNSZoneRRSetSchemaMetaHealthchecks] = []map[string]any{v}
			} else {
				rrMeta[key] = []map[string]any{v}
			}
			continue
		}
		s, ok := val.(string)
		if ok {
			if key == DNSZoneRRSetSchemaMetaGeodnsLink {
				rrMeta[key] = s
				continue
			}
		}
	}

	rrsm := make([]map[string]interface{}, 0)
	if len(rrMeta) > 0 {
		rrsm = append(rrsm, rrMeta)
		_ = d.Set(DNSZoneRRSetSchemaMeta, rrsm)
	}

	return nil
}

func resourceDNSZoneRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Id() == "" {
		return diag.Errorf("empty id")
	}
	zone := strings.TrimSpace(d.Get(DNSZoneRecordSchemaZone).(string))
	domain := strings.TrimSpace(d.Get(DNSZoneRecordSchemaDomain).(string))
	rType := strings.TrimSpace(d.Get(DNSZoneRecordSchemaType).(string))
	log.Println("[DEBUG] Start DNS Zone Record Resource deleting")
	defer log.Printf("[DEBUG] Finish DNS Zone Record Resource deleting (id=%s %s %s)\n", zone, domain, rType)

	config := m.(*Config)
	client := config.DNSClient

	err := client.DeleteRRSet(ctx, zone, domain, rType)
	if err != nil {
		return diag.FromErr(fmt.Errorf("delete zone rrset: %w", err))
	}

	d.SetId("")

	return nil
}

func fillRRSet(d *schema.ResourceData, rType string, rrSet *dnssdk.RRSet) error {
	// set filters
	for _, resource := range d.Get(DNSZoneRecordSchemaFilter).([]any) {
		filter := dnssdk.RecordFilter{}
		filterData := resource.(map[string]interface{})
		name := filterData[DNSZoneRecordSchemaFilterType].(string)
		filter.Type = name
		limit, ok := filterData[DNSZoneRecordSchemaFilterLimit].(int)
		if ok {
			filter.Limit = uint(limit)
		}
		strict, ok := filterData[DNSZoneRecordSchemaFilterStrict].(bool)
		if ok {
			filter.Strict = strict
		}
		rrSet.AddFilter(filter)
	}
	// set meta of RRSet
	for _, kv := range d.Get(DNSZoneRRSetSchemaMeta).(*schema.Set).List() {
		if rrSet.Meta == nil {
			rrSet.Meta = map[string]any{}
		}
		metaErrs := make([]error, 0)
		for k, v := range kv.(map[string]any) {
			switch k {
			// DNSZoneRRSetSchemaMetaFailover - 10. `meta` (map)
			case DNSZoneRRSetSchemaMetaFailover, DNSZoneRRSetSchemaMetaHealthchecks:
				failoverObj, ok := v.(*schema.Set)
				if !ok {
					return fmt.Errorf("invalid type of rrset meta.healthchecks, expected map[string]any, got %T %v", v, v)
				}
				fMap := map[string]any{}
				for _, kv2 := range failoverObj.List() {
					m2, ok := kv2.(map[string]any)
					if !ok {
						return fmt.Errorf("invalid type of rrset meta.healthchecks.*, expected map[string]any, got %T %v", kv2, kv2)
					}
					for k2, v2 := range m2 {
						fMap[k2] = v2
					}
				}
				rrSet.Meta[DNSZoneRRSetSchemaMetaHealthchecks] = fMap
				rrSet.Meta[DNSZoneRRSetSchemaMetaFailover] = fMap
			case DNSZoneRRSetSchemaMetaGeodnsLink:
				var ok bool
				rrSet.Meta[k], ok = v.(string)
				if !ok {
					return fmt.Errorf("invalid type of rrset meta.geodns_link, expected string, got %T %v", v, v)
				}
			default:
				return fmt.Errorf("unsupported rrset meta key %s with value %v", k, v)
			}
			if len(metaErrs) > 0 {
				return fmt.Errorf("invalid meta for rrset %#v: %v", kv, metaErrs)
			}
		}
	}
	// set meta of ResourceRecord
	for _, resource := range d.Get(DNSZoneRecordSchemaResourceRecord).(*schema.Set).List() {
		data := resource.(map[string]interface{})
		content := data[DNSZoneRecordSchemaContent].(string)
		rr := (&dnssdk.ResourceRecord{}).SetContent(rType, content)
		enabled := data[DNSZoneRecordSchemaEnabled].(bool)
		rr.Enabled = enabled
		metaErrs := make([]error, 0)

		for _, dataMeta := range data[DNSZoneRecordSchemaMeta].(*schema.Set).List() {
			meta := dataMeta.(map[string]interface{})
			log.Println("dataMeta: ", dataMeta)
			validWrap := func(rm dnssdk.ResourceMeta) dnssdk.ResourceMeta {
				if rm.Valid() != nil {
					metaErrs = append(metaErrs, rm.Valid())
				}
				return rm
			}

			val := meta[DNSZoneRecordSchemaMetaIP].([]any)
			ips := make([]string, len(val))
			for i, v := range val {
				ips[i] = v.(string)
			}
			if len(ips) > 0 {
				rr.AddMeta(dnssdk.NewResourceMetaIP(ips...))
			}

			val = meta[DNSZoneRecordSchemaMetaCountries].([]any)
			countries := make([]string, len(val))
			for i, v := range val {
				countries[i] = v.(string)
			}
			if len(countries) > 0 {
				rr.AddMeta(dnssdk.NewResourceMetaCountries(countries...))
			}

			val = meta[DNSZoneRecordSchemaMetaContinents].([]any)
			continents := make([]string, len(val))
			for i, v := range val {
				continents[i] = v.(string)
			}
			if len(continents) > 0 {
				rr.AddMeta(dnssdk.NewResourceMetaContinents(continents...))
			}

			valStr := meta[DNSZoneRecordSchemaMetaNotes].(string)
			rr.AddMeta(dnssdk.NewResourceMetaNotes(valStr))

			latLongVal := meta[DNSZoneRecordSchemaMetaLatLong].([]any)
			if len(latLongVal) == 2 {
				rr.AddMeta(
					validWrap(
						dnssdk.NewResourceMetaLatLong(
							fmt.Sprintf("%f,%f", latLongVal[0].(float64), latLongVal[1].(float64)))))
			}

			val = meta[DNSZoneRecordSchemaMetaAsn].([]any)
			asn := make([]uint64, len(val))
			for i, v := range val {
				asn[i] = uint64(v.(int))
			}
			if len(asn) > 0 {
				rr.AddMeta(dnssdk.NewResourceMetaAsn(asn...))
			}

			valDefault := meta[DNSZoneRecordSchemaMetaDefault].(bool)
			if valDefault {
				rr.AddMeta(validWrap(dnssdk.NewResourceMetaDefault()))
			}

			valBool, ok := meta[DNSZoneRecordSchemaMetaBackup].(bool)
			if ok && valBool {
				rr.AddMeta(validWrap(dnssdk.NewResourceMetaBackup()))
			}

			valBool, ok = meta[DNSZoneRecordSchemaMetaFallback].(bool)
			if ok && valBool {
				rr.AddMeta(validWrap(dnssdk.NewResourceMetaFallback()))
			}

			valInt, ok := meta[DNSZoneRecordSchemaMetaWeight].(int)
			if ok && valInt > 0 {
				rr.AddMeta(validWrap(dnssdk.NewResourceMetaWeight(valInt)))
			}
		}

		if len(metaErrs) > 0 {
			return fmt.Errorf("invalid meta for zone rrset with content %s: %v", content, metaErrs)
		}
		rrSet.Records = append(rrSet.Records, *rr)
	}
	return nil
}
