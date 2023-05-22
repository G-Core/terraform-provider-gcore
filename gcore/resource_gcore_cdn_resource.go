package gcore

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"

	gcdn "github.com/G-Core/gcorelabscdn-go/gcore"
	"github.com/G-Core/gcorelabscdn-go/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	optionsSchema = &schema.Schema{
		Type:        schema.TypeList,
		MaxItems:    1,
		Optional:    true,
		Computed:    true,
		Description: "Each option in CDN resource settings. Each option added to CDN resource settings should have the following mandatory request fields: enabled, value.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"edge_cache_settings": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Computed:    true,
					Description: "The cache expiration time for CDN servers.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Caching time for a response with codes 200, 206, 301, 302. Responses with codes 4xx, 5xx will not be cached. Use '0s' disable to caching. Use custom_values field to specify a custom caching time for a response with specific codes.",
							},
							"default": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Content will be cached according to origin cache settings. The value applies for a response with codes 200, 201, 204, 206, 301, 302, 303, 304, 307, 308 if an origin server does not have caching HTTP headers. Responses with other codes will not be cached.",
							},
							"custom_values": {
								Type:     schema.TypeMap,
								Optional: true,
								Computed: true,
								DefaultFunc: func() (interface{}, error) {
									return map[string]interface{}{}, nil
								},
								Elem:        schema.TypeString,
								Description: "Caching time for a response with specific codes. These settings have a higher priority than the value field. Response code ('304', '404' for example). Use 'any' to specify caching time for all response codes. Caching time in seconds ('0s', '600s' for example). Use '0s' to disable caching for a specific response code.",
							},
						},
					},
				},
				"browser_cache_settings": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Computed:    true,
					Description: "",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "",
							},
						},
					},
				},
				"host_header": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Specify the Host header that CDN servers use when request content from an origin server. Your server must be able to process requests with the chosen header. If the option is in NULL state Host Header value is taken from the CNAME field.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},
				"redirect_http_to_https": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Sets redirect from HTTP protocol to HTTPS for all resource requests.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"request_limiter": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "It allows to limit the amount of HTTP requests",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"rate": {
								Type:     schema.TypeInt,
								Required: true,
							},
							"burst": {
								Type:     schema.TypeInt,
								Required: true,
							},
							"rate_unit": {
								Type:     schema.TypeInt,
								Optional: true,
								Default:  "r/s",
							},
							"delay": {
								Type:     schema.TypeInt,
								Optional: true,
								Default:  0,
							},
						},
					},
				},
				"gzip_on": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"cors": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeSet,
								Elem:     &schema.Schema{Type: schema.TypeString},
								Required: true,
							},
						},
					},
				},
				"rewrite": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"body": {
								Type:     schema.TypeString,
								Required: true,
							},
							"flag": {
								Type:     schema.TypeString,
								Optional: true,
								Default:  "break",
							},
						},
					},
				},
				"webp": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"jpg_quality": {
								Type:     schema.TypeInt,
								Required: true,
							},
							"png_quality": {
								Type:     schema.TypeInt,
								Required: true,
							},
							"png_lossless": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  false,
							},
						},
					},
				},
				"sni": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Computed:    true,
					Description: "",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"sni_type": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Available values 'dynamic' or 'custom'",
							},
							"custom_hostname": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Required to set custom hostname in case sni-type='custom'",
							},
						},
					},
				},
				"ignore_query_string": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"query_params_whitelist": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeSet,
								Elem:     &schema.Schema{Type: schema.TypeString},
								Required: true,
							},
						},
					},
				},
				"query_params_blacklist": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeSet,
								Elem:     &schema.Schema{Type: schema.TypeString},
								Required: true,
							},
						},
					},
				},
				"static_request_headers": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeMap,
								Elem:     &schema.Schema{Type: schema.TypeString},
								Required: true,
							},
						},
					},
				},
				"static_headers": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeMap,
								Elem:     &schema.Schema{Type: schema.TypeString},
								Required: true,
							},
						},
					},
				},
				"websockets": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"tls_versions": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeSet,
								Elem:     &schema.Schema{Type: schema.TypeString},
								Required: true,
							},
						},
					},
				},
				"use_rsa_le_cert": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"force_return": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Allows to apply custom HTTP code to the CDN content.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"code": {
								Type:        schema.TypeInt,
								Required:    true,
								Description: "HTTP response status code. Available codes: 100 <= value <= 599. Reserved codes: 408, 444, 477, 494, 495, 496, 497, 499",
							},
							"body": {
								Type:        schema.TypeString,
								Optional:    true,
								Default:     "",
								Description: "Response text or URL if you're going to set up redirection. Max length = 100.",
							},
						},
					},
				},
			},
		},
	}
)

func resourceCDNResource() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"cname": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A CNAME that will be used to deliver content though a CDN. If you update this field new resource will be created.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom client description of the resource.",
			},
			"origin_group": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ExactlyOneOf: []string{
					"origin_group",
					"origin",
				},
				Description: "ID of the Origins Group. Use one of your Origins Group or create a new one. You can use either 'origin' parameter or 'originGroup' in the resource definition.",
			},
			"origin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ExactlyOneOf: []string{
					"origin_group",
					"origin",
				},
				Description: "A domain name or IP of your origin source. Specify a port if custom. You can use either 'origin' parameter or 'originGroup' in the resource definition.",
			},
			"origin_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "This option defines the protocol that will be used by CDN servers to request content from an origin source. If not specified, we will use HTTP to connect to an origin server. Possible values are: HTTPS, HTTP, MATCH.",
			},
			"secondary_hostnames": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				DefaultFunc: func() (interface{}, error) {
					return []string{}, nil
				},
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of additional CNAMEs.",
			},
			"ssl_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Use HTTPS protocol for content delivery.",
			},
			"ssl_data": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"ssl_enabled"},
				Description:  "Specify the SSL Certificate ID which should be used for the CDN Resource.",
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "The setting allows to enable or disable a CDN Resource",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of a CDN resource content availability. Possible values are: Active, Suspended, Processed.",
			},
			"options": optionsSchema,
		},
		CreateContext: resourceCDNResourceCreate,
		ReadContext:   resourceCDNResourceRead,
		UpdateContext: resourceCDNResourceUpdate,
		DeleteContext: resourceCDNResourceDelete,
		Description:   "Represent CDN resource",
	}
}

func resourceCDNResourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start CDN Resource creating")
	config := m.(*Config)
	client := config.CDNClient

	var req resources.CreateRequest
	req.Cname = d.Get("cname").(string)
	req.Description = d.Get("description").(string)
	req.Origin = d.Get("origin").(string)
	req.OriginGroup = d.Get("origin_group").(int)
	req.OriginProtocol = resources.Protocol(d.Get("origin_protocol").(string))
	req.SSlEnabled = d.Get("ssl_enabled").(bool)
	req.SSLData = d.Get("ssl_data").(int)

	req.Options = listToOptions(d.Get("options").([]interface{}))

	for _, hostname := range d.Get("secondary_hostnames").(*schema.Set).List() {
		req.SecondaryHostnames = append(req.SecondaryHostnames, hostname.(string))
	}

	result, err := client.Resources().Create(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", result.ID))
	resourceCDNResourceRead(ctx, d, m)

	log.Printf("[DEBUG] Finish CDN Resource creating (id=%d)\n", result.ID)
	return nil
}

func resourceCDNResourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceID := d.Id()
	log.Printf("[DEBUG] Start CDN Resource reading (id=%s)\n", resourceID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.Resources().Get(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("cname", result.Cname)
	d.Set("description", result.Description)
	d.Set("origin_group", result.OriginGroup)
	d.Set("origin_protocol", result.OriginProtocol)
	d.Set("secondary_hostnames", result.SecondaryHostnames)
	d.Set("ssl_enabled", result.SSlEnabled)
	d.Set("ssl_data", result.SSLData)
	d.Set("status", result.Status)
	d.Set("active", result.Active)
	if err := d.Set("options", optionsToList(result.Options)); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish CDN Resource reading")
	return nil
}

func resourceCDNResourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceID := d.Id()
	log.Printf("[DEBUG] Start CDN Resource updating (id=%s)\n", resourceID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var req resources.UpdateRequest
	req.Active = d.Get("active").(bool)
	req.Description = d.Get("description").(string)
	req.OriginGroup = d.Get("origin_group").(int)
	req.SSlEnabled = d.Get("ssl_enabled").(bool)
	req.SSLData = d.Get("ssl_data").(int)
	req.OriginProtocol = resources.Protocol(d.Get("origin_protocol").(string))
	req.Options = listToOptions(d.Get("options").([]interface{}))
	for _, hostname := range d.Get("secondary_hostnames").(*schema.Set).List() {
		req.SecondaryHostnames = append(req.SecondaryHostnames, hostname.(string))
	}

	if _, err := client.Resources().Update(ctx, id, &req); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish CDN Resource updating")
	return resourceCDNResourceRead(ctx, d, m)
}

func resourceCDNResourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceID := d.Id()
	log.Printf("[DEBUG] Start CDN Resource deleting (id=%s)\n", resourceID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.Resources().Delete(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Println("[DEBUG] Finish CDN Resource deleting")
	return nil
}

func listToOptions(l []interface{}) *gcdn.Options {
	if len(l) == 0 {
		return nil
	}

	var opts gcdn.Options
	fields := l[0].(map[string]interface{})
	if opt, ok := getOptByName(fields, "edge_cache_settings"); ok {
		rawCustomVals := opt["custom_values"].(map[string]interface{})
		customVals := make(map[string]string, len(rawCustomVals))
		for key, value := range rawCustomVals {
			customVals[key] = value.(string)
		}

		opts.EdgeCacheSettings = &gcdn.EdgeCacheSettings{
			Enabled:      opt["enabled"].(bool),
			Value:        opt["value"].(string),
			CustomValues: customVals,
			Default:      opt["default"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "browser_cache_settings"); ok {
		enabled := true
		if _, ok := opt["enabled"]; ok {
			enabled = opt["enabled"].(bool)
		}
		opts.BrowserCacheSettings = &gcdn.BrowserCacheSettings{
			Enabled: enabled,
			Value:   opt["value"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "host_header"); ok {
		opts.HostHeader = &gcdn.HostHeader{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "redirect_http_to_https"); ok {
		enabled := true
		if _, ok := opt["enabled"]; ok {
			enabled = opt["enabled"].(bool)
		}
		opts.RedirectHttpToHttps = &gcdn.RedirectHttpToHttps{
			Enabled: enabled,
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "request_limiter"); ok {
		enabled := true
		if _, ok := opt["enabled"]; ok {
			enabled = opt["enabled"].(bool)
		}
		opts.RequestLimiter = &gcdn.RequestLimiter{
			Enabled:  enabled,
			Rate:     opt["rate"].(int),
			Burst:    opt["burst"].(int),
			RateUnit: opt["rate_unit"].(string),
			Delay:    opt["delay"].(int),
		}
	}
	if opt, ok := getOptByName(fields, "gzip_on"); ok {
		enabled := true
		if _, ok := opt["enabled"]; ok {
			enabled = opt["enabled"].(bool)
		}
		opts.GzipOn = &gcdn.GzipOn{
			Enabled: enabled,
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "cors"); ok {
		enabled := true
		if _, ok := opt["enabled"]; ok {
			enabled = opt["enabled"].(bool)
		}
		opts.Cors = &gcdn.Cors{
			Enabled: enabled,
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.Cors.Value = append(opts.Cors.Value, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "rewrite"); ok {
		enabled := true
		if _, ok := opt["enabled"]; ok {
			enabled = opt["enabled"].(bool)
		}
		opts.Rewrite = &gcdn.Rewrite{
			Enabled: enabled,
			Body:    opt["body"].(string),
			Flag:    opt["flag"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "webp"); ok {
		enabled := true
		if _, ok := opt["enabled"]; ok {
			enabled = opt["enabled"].(bool)
		}
		opts.Webp = &gcdn.Webp{
			Enabled:     enabled,
			JPGQuality:  opt["jpg_quality"].(int),
			PNGQuality:  opt["png_quality"].(int),
			PNGLossless: opt["png_lossless"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "sni"); ok {
		enabled := true
		if _, ok := opt["enabled"]; ok {
			enabled = opt["enabled"].(bool)
		}
		opts.SNI = &gcdn.SNIOption{
			Enabled:        enabled,
			SNIType:        opt["sni_type"].(string),
			CustomHostname: opt["custom_hostname"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "ignore_query_string"); ok {
		enabled := true
		if _, ok := opt["enabled"]; ok {
			enabled = opt["enabled"].(bool)
		}
		opts.IgnoreQueryString = &gcdn.IgnoreQueryString{
			Enabled: enabled,
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "query_params_whitelist"); ok {
		enabled := true
		if _, ok := opt["enabled"]; ok {
			enabled = opt["enabled"].(bool)
		}
		opts.QueryParamsWhitelist = &gcdn.QueryParamsWhitelist{
			Enabled: enabled,
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.QueryParamsWhitelist.Value = append(opts.QueryParamsWhitelist.Value, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "query_params_blacklist"); ok {
		enabled := true
		if _, ok := opt["enabled"]; ok {
			enabled = opt["enabled"].(bool)
		}
		opts.QueryParamsBlacklist = &gcdn.QueryParamsBlacklist{
			Enabled: enabled,
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.QueryParamsBlacklist.Value = append(opts.QueryParamsBlacklist.Value, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "static_request_headers"); ok {
		enabled := true
		if _, ok := opt["enabled"]; ok {
			enabled = opt["enabled"].(bool)
		}
		opts.StaticRequestHeaders = &gcdn.StaticRequestHeaders{
			Enabled: enabled,
			Value:   map[string]string{},
		}
		for k, v := range opt["value"].(map[string]interface{}) {
			opts.StaticRequestHeaders.Value[k] = v.(string)
		}
	}
	if opt, ok := getOptByName(fields, "static_headers"); ok {
		enabled := true
		if _, ok := opt["enabled"]; ok {
			enabled = opt["enabled"].(bool)
		}
		opts.StaticHeaders = &gcdn.StaticHeaders{
			Enabled: enabled,
			Value:   map[string]string{},
		}
		for k, v := range opt["value"].(map[string]interface{}) {
			opts.StaticHeaders.Value[k] = v.(string)
		}
	}
	if opt, ok := getOptByName(fields, "websockets"); ok {
		enabled := true
		if _, ok := opt["enabled"]; ok {
			enabled = opt["enabled"].(bool)
		}
		opts.WebSockets = &gcdn.WebSockets{
			Enabled: enabled,
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "tls_versions"); ok {
		enabled := true
		if _, ok := opt["enabled"]; ok {
			enabled = opt["enabled"].(bool)
		}
		opts.TLSVersions = &gcdn.TLSVersions{
			Enabled: enabled,
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.TLSVersions.Value = append(opts.TLSVersions.Value, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "use_rsa_le_cert"); ok {
		enabled := true
		if _, ok := opt["enabled"]; ok {
			enabled = opt["enabled"].(bool)
		}
		opts.UseRSALECert = &gcdn.UseRSALECert{
			Enabled: enabled,
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "force_return"); ok {
		enabled := true
		if _, ok := opt["enabled"]; ok {
			enabled = opt["enabled"].(bool)
		}
		opts.ForceReturn = &gcdn.ForceReturn{
			Enabled: enabled,
			Code:    opt["code"].(int),
			Body:    opt["body"].(string),
		}
	}
	return &opts
}

func getOptByName(fields map[string]interface{}, name string) (map[string]interface{}, bool) {
	if _, ok := fields[name]; !ok {
		return nil, false
	}

	container, ok := fields[name].([]interface{})
	if !ok {
		return nil, false
	}

	if len(container) == 0 {
		return nil, false
	}

	opt, ok := container[0].(map[string]interface{})
	if !ok {
		return nil, false
	}

	return opt, true
}

func optionsToList(options *gcdn.Options) []interface{} {
	result := make(map[string][]interface{})
	if options.EdgeCacheSettings != nil {
		m := structToMap(options.EdgeCacheSettings)
		result["edge_cache_settings"] = []interface{}{m}
	}
	if options.BrowserCacheSettings != nil {
		m := structToMap(options.BrowserCacheSettings)
		result["browser_cache_settings"] = []interface{}{m}
	}
	if options.HostHeader != nil {
		m := structToMap(options.HostHeader)
		result["host_header"] = []interface{}{m}
	}
	if options.RedirectHttpToHttps != nil {
		m := structToMap(options.RedirectHttpToHttps)
		result["redirect_http_to_https"] = []interface{}{m}
	}
	if options.RequestLimiter != nil {
		m := structToMap(options.RequestLimiter)
		result["request_limiter"] = []interface{}{m}
	}
	if options.GzipOn != nil {
		m := structToMap(options.GzipOn)
		result["gzip_on"] = []interface{}{m}
	}
	if options.Cors != nil {
		m := structToMap(options.Cors)
		result["cors"] = []interface{}{m}
	}
	if options.Rewrite != nil {
		m := structToMap(options.Rewrite)
		result["rewrite"] = []interface{}{m}
	}
	if options.Webp != nil {
		m := structToMap(options.Webp)
		result["webp"] = []interface{}{m}
	}
	if options.SNI != nil {
		m := structToMap(options.SNI)
		result["sni"] = []interface{}{m}
	}
	if options.IgnoreQueryString != nil {
		m := structToMap(options.IgnoreQueryString)
		result["ignore_query_string"] = []interface{}{m}
	}
	if options.QueryParamsWhitelist != nil {
		m := structToMap(options.QueryParamsWhitelist)
		result["query_params_whitelist"] = []interface{}{m}
	}
	if options.QueryParamsBlacklist != nil {
		m := structToMap(options.QueryParamsBlacklist)
		result["query_params_blacklist"] = []interface{}{m}
	}
	if options.StaticRequestHeaders != nil {
		m := structToMap(options.StaticRequestHeaders)
		result["static_request_headers"] = []interface{}{m}
	}
	if options.StaticHeaders != nil {
		m := structToMap(options.StaticHeaders)
		result["static_headers"] = []interface{}{m}
	}
	if options.WebSockets != nil {
		m := structToMap(options.WebSockets)
		result["websockets"] = []interface{}{m}
	}
	if options.TLSVersions != nil {
		m := structToMap(options.TLSVersions)
		result["tls_versions"] = []interface{}{m}
	}
	if options.UseRSALECert != nil {
		m := structToMap(options.UseRSALECert)
		result["use_rsa_le_cert"] = []interface{}{m}
	}
	if options.ForceReturn != nil {
		m := structToMap(options.ForceReturn)
		result["force_return"] = []interface{}{m}
	}
	return []interface{}{result}
}

func structToMap(item interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = structToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}
