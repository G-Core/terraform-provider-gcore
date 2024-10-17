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
				Default:     "HTTP",
				Description: "This option defines the protocol that will be used by CDN servers to request content from an origin source. If not specified, we will use HTTP to connect to an origin server. Possible values are: HTTPS, HTTP, MATCH.",
			},
			"secondary_hostnames": {
				Type:     schema.TypeSet,
				Optional: true,
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
			"options": resourceOptionsSchema,
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
	req.SecondaryHostnames = make([]string, 0)
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

	if opt, ok := getOptByName(fields, "allowed_http_methods"); ok {
		opts.AllowedHTTPMethods = &gcdn.AllowedHTTPMethods{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.AllowedHTTPMethods.Value = append(opts.AllowedHTTPMethods.Value, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "brotli_compression"); ok {
		opts.BrotliCompression = &gcdn.BrotliCompression{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.BrotliCompression.Value = append(opts.BrotliCompression.Value, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "browser_cache_settings"); ok {
		opts.BrowserCacheSettings = &gcdn.BrowserCacheSettings{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "cors"); ok {
		opts.Cors = &gcdn.Cors{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.Cors.Value = append(opts.Cors.Value, v.(string))
		}
		if _, ok := opt["always"]; ok {
			opts.Cors.Always = opt["always"].(bool)
		}
	}
	if opt, ok := getOptByName(fields, "country_acl"); ok {
		opts.CountryACL = &gcdn.CountryACL{
			Enabled:    opt["enabled"].(bool),
			PolicyType: opt["policy_type"].(string),
		}
		for _, v := range opt["excepted_values"].(*schema.Set).List() {
			opts.CountryACL.ExceptedValues = append(opts.CountryACL.ExceptedValues, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "disable_proxy_force_ranges"); ok {
		opts.DisableProxyForceRanges = &gcdn.DisableProxyForceRanges{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
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
	if opt, ok := getOptByName(fields, "fastedge"); ok {
		opts.FastEdge = &gcdn.FastEdge{
			Enabled: opt["enabled"].(bool),
		}
		if onRequestHeaders, ok := getOptByName(opt, "on_request_headers"); ok {
			opts.FastEdge.OnRequestHeaders = &gcdn.FastEdgeAppConfig{
				Enabled:          onRequestHeaders["enabled"].(bool),
				AppID:            onRequestHeaders["app_id"].(string),
				InterruptOnError: onRequestHeaders["interrupt_on_error"].(bool),
			}
		}
	}
	if opt, ok := getOptByName(fields, "fetch_compressed"); ok {
		opts.FetchCompressed = &gcdn.FetchCompressed{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "follow_origin_redirect"); ok {
		opts.FollowOriginRedirect = &gcdn.FollowOriginRedirect{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["codes"].(*schema.Set).List() {
			opts.FollowOriginRedirect.Codes = append(opts.FollowOriginRedirect.Codes, v.(int))
		}
	}
	if opt, ok := getOptByName(fields, "force_return"); ok {
		opts.ForceReturn = &gcdn.ForceReturn{
			Enabled: opt["enabled"].(bool),
			Code:    opt["code"].(int),
			Body:    opt["body"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "forward_host_header"); ok {
		opts.ForwardHostHeader = &gcdn.ForwardHostHeader{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "gzip_on"); ok {
		opts.GzipOn = &gcdn.GzipOn{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "host_header"); ok {
		opts.HostHeader = &gcdn.HostHeader{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "http3_enabled"); ok {
		opts.HTTP3Enabled = &gcdn.HTTP3Enabled{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "ignore_cookie"); ok {
		opts.IgnoreCookie = &gcdn.IgnoreCookie{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "ignore_query_string"); ok {
		opts.IgnoreQueryString = &gcdn.IgnoreQueryString{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "image_stack"); ok {
		opts.ImageStack = &gcdn.ImageStack{
			Enabled: opt["enabled"].(bool),
			Quality: opt["quality"].(int),
		}
		if _, ok := opt["avif_enabled"]; ok {
			opts.ImageStack.AvifEnabled = opt["avif_enabled"].(bool)
		}
		if _, ok := opt["webp_enabled"]; ok {
			opts.ImageStack.WebpEnabled = opt["webp_enabled"].(bool)
		}
		if _, ok := opt["png_lossless"]; ok {
			opts.ImageStack.PngLossless = opt["png_lossless"].(bool)
		}
	}
	if opt, ok := getOptByName(fields, "ip_address_acl"); ok {
		opts.IPAddressACL = &gcdn.IPAddressACL{
			Enabled:    opt["enabled"].(bool),
			PolicyType: opt["policy_type"].(string),
		}
		for _, v := range opt["excepted_values"].(*schema.Set).List() {
			opts.IPAddressACL.ExceptedValues = append(opts.IPAddressACL.ExceptedValues, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "limit_bandwidth"); ok {
		opts.LimitBandwidth = &gcdn.LimitBandwidth{
			Enabled:   opt["enabled"].(bool),
			LimitType: opt["limit_type"].(string),
		}
		if _, ok := opt["speed"]; ok {
			opts.LimitBandwidth.Speed = opt["speed"].(int)
		}
		if _, ok := opt["buffer"]; ok {
			opts.LimitBandwidth.Buffer = opt["buffer"].(int)
		}
	}
	if opt, ok := getOptByName(fields, "proxy_cache_key"); ok {
		opts.ProxyCacheKey = &gcdn.ProxyCacheKey{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "proxy_cache_methods_set"); ok {
		opts.ProxyCacheMethodsSet = &gcdn.ProxyCacheMethodsSet{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "proxy_connect_timeout"); ok {
		opts.ProxyConnectTimeout = &gcdn.ProxyConnectTimeout{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "proxy_read_timeout"); ok {
		opts.ProxyReadTimeout = &gcdn.ProxyReadTimeout{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "query_params_blacklist"); ok {
		opts.QueryParamsBlacklist = &gcdn.QueryParamsBlacklist{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.QueryParamsBlacklist.Value = append(opts.QueryParamsBlacklist.Value, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "query_params_whitelist"); ok {
		opts.QueryParamsWhitelist = &gcdn.QueryParamsWhitelist{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.QueryParamsWhitelist.Value = append(opts.QueryParamsWhitelist.Value, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "redirect_https_to_http"); ok {
		opts.RedirectHttpsToHttp = &gcdn.RedirectHttpsToHttp{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "redirect_http_to_https"); ok {
		opts.RedirectHttpToHttps = &gcdn.RedirectHttpToHttps{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "referrer_acl"); ok {
		opts.ReferrerACL = &gcdn.ReferrerACL{
			Enabled:    opt["enabled"].(bool),
			PolicyType: opt["policy_type"].(string),
		}
		for _, v := range opt["excepted_values"].(*schema.Set).List() {
			opts.ReferrerACL.ExceptedValues = append(opts.ReferrerACL.ExceptedValues, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "request_limiter"); ok {
		opts.RequestLimiter = &gcdn.RequestLimiter{
			Enabled:  opt["enabled"].(bool),
			Rate:     opt["rate"].(int),
			Burst:    opt["burst"].(int),
			RateUnit: opt["rate_unit"].(string),
			Delay:    opt["delay"].(int),
		}
	}
	if opt, ok := getOptByName(fields, "response_headers_hiding_policy"); ok {
		opts.ResponseHeadersHidingPolicy = &gcdn.ResponseHeadersHidingPolicy{
			Enabled: opt["enabled"].(bool),
			Mode:    opt["mode"].(string),
		}
		for _, v := range opt["excepted"].(*schema.Set).List() {
			opts.ResponseHeadersHidingPolicy.Excepted = append(opts.ResponseHeadersHidingPolicy.Excepted, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "rewrite"); ok {
		opts.Rewrite = &gcdn.Rewrite{
			Enabled: opt["enabled"].(bool),
			Body:    opt["body"].(string),
			Flag:    opt["flag"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "secure_key"); ok {
		opts.SecureKey = &gcdn.SecureKey{
			Enabled: opt["enabled"].(bool),
			Key:     opt["key"].(string),
			Type:    opt["type"].(int),
		}
	}
	if opt, ok := getOptByName(fields, "slice"); ok {
		opts.Slice = &gcdn.Slice{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "sni"); ok {
		opts.SNI = &gcdn.SNIOption{
			Enabled:        opt["enabled"].(bool),
			SNIType:        opt["sni_type"].(string),
			CustomHostname: opt["custom_hostname"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "stale"); ok {
		opts.Stale = &gcdn.Stale{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.Stale.Value = append(opts.Stale.Value, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "static_headers"); ok {
		opts.StaticHeaders = &gcdn.StaticHeaders{
			Enabled: opt["enabled"].(bool),
			Value:   map[string]string{},
		}
		for k, v := range opt["value"].(map[string]interface{}) {
			opts.StaticHeaders.Value[k] = v.(string)
		}
	}
	if opt, ok := getOptByName(fields, "static_request_headers"); ok {
		opts.StaticRequestHeaders = &gcdn.StaticRequestHeaders{
			Enabled: opt["enabled"].(bool),
			Value:   map[string]string{},
		}
		for k, v := range opt["value"].(map[string]interface{}) {
			opts.StaticRequestHeaders.Value[k] = v.(string)
		}
	}
	if opt, ok := getOptByName(fields, "static_response_headers"); ok {
		opts.StaticResponseHeaders = &gcdn.StaticResponseHeaders{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["value"].([]interface{}) {
			item_data := v.(map[string]interface{})
			item := &gcdn.StaticResponseHeadersItem{
				Name: item_data["name"].(string),
			}
			for _, val := range item_data["value"].(*schema.Set).List() {
				item.Value = append(item.Value, val.(string))
			}
			if _, ok := item_data["always"]; ok {
				item.Always = item_data["always"].(bool)
			}
			opts.StaticResponseHeaders.Value = append(opts.StaticResponseHeaders.Value, *item)
		}
	}
	if opt, ok := getOptByName(fields, "tls_versions"); ok {
		opts.TLSVersions = &gcdn.TLSVersions{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.TLSVersions.Value = append(opts.TLSVersions.Value, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "use_default_le_chain"); ok {
		opts.UseDefaultLEChain = &gcdn.UseDefaultLEChain{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "user_agent_acl"); ok {
		opts.UserAgentACL = &gcdn.UserAgentACL{
			Enabled:    opt["enabled"].(bool),
			PolicyType: opt["policy_type"].(string),
		}
		for _, v := range opt["excepted_values"].(*schema.Set).List() {
			opts.UserAgentACL.ExceptedValues = append(opts.UserAgentACL.ExceptedValues, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "use_rsa_le_cert"); ok {
		opts.UseRSALECert = &gcdn.UseRSALECert{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "waf"); ok {
		opts.WAF = &gcdn.WAF{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "websockets"); ok {
		opts.WebSockets = &gcdn.WebSockets{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
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
	if options.AllowedHTTPMethods != nil {
		m := structToMap(options.AllowedHTTPMethods)
		result["allowed_http_methods"] = []interface{}{m}
	}
	if options.BrotliCompression != nil {
		m := structToMap(options.BrotliCompression)
		result["brotli_compression"] = []interface{}{m}
	}
	if options.BrowserCacheSettings != nil {
		m := structToMap(options.BrowserCacheSettings)
		result["browser_cache_settings"] = []interface{}{m}
	}
	if options.Cors != nil {
		m := structToMap(options.Cors)
		result["cors"] = []interface{}{m}
	}
	if options.CountryACL != nil {
		m := structToMap(options.CountryACL)
		result["country_acl"] = []interface{}{m}
	}
	if options.DisableProxyForceRanges != nil {
		m := structToMap(options.DisableProxyForceRanges)
		result["disable_proxy_force_ranges"] = []interface{}{m}
	}
	if options.EdgeCacheSettings != nil {
		m := structToMap(options.EdgeCacheSettings)
		result["edge_cache_settings"] = []interface{}{m}
	}
	if options.FastEdge != nil {
		m := structToMap(options.FastEdge)
		if options.FastEdge.OnRequestHeaders != nil {
			m["on_request_headers"] = []interface{}{structToMap(options.FastEdge.OnRequestHeaders)}
		}
		result["fastedge"] = []interface{}{m}
	}
	if options.FetchCompressed != nil {
		m := structToMap(options.FetchCompressed)
		result["fetch_compressed"] = []interface{}{m}
	}
	if options.FollowOriginRedirect != nil {
		m := structToMap(options.FollowOriginRedirect)
		result["follow_origin_redirect"] = []interface{}{m}
	}
	if options.ForceReturn != nil {
		m := structToMap(options.ForceReturn)
		result["force_return"] = []interface{}{m}
	}
	if options.ForwardHostHeader != nil {
		m := structToMap(options.ForwardHostHeader)
		result["forward_host_header"] = []interface{}{m}
	}
	if options.GzipOn != nil {
		m := structToMap(options.GzipOn)
		result["gzip_on"] = []interface{}{m}
	}
	if options.HostHeader != nil {
		m := structToMap(options.HostHeader)
		result["host_header"] = []interface{}{m}
	}
	if options.HTTP3Enabled != nil {
		m := structToMap(options.HTTP3Enabled)
		result["http3_enabled"] = []interface{}{m}
	}
	if options.IgnoreCookie != nil {
		m := structToMap(options.IgnoreCookie)
		result["ignore_cookie"] = []interface{}{m}
	}
	if options.IgnoreQueryString != nil {
		m := structToMap(options.IgnoreQueryString)
		result["ignore_query_string"] = []interface{}{m}
	}
	if options.ImageStack != nil {
		m := structToMap(options.ImageStack)
		result["image_stack"] = []interface{}{m}
	}
	if options.IPAddressACL != nil {
		m := structToMap(options.IPAddressACL)
		result["ip_address_acl"] = []interface{}{m}
	}
	if options.LimitBandwidth != nil {
		m := structToMap(options.LimitBandwidth)
		result["limit_bandwidth"] = []interface{}{m}
	}
	if options.ProxyCacheKey != nil {
		m := structToMap(options.ProxyCacheKey)
		result["proxy_cache_key"] = []interface{}{m}
	}
	if options.ProxyCacheMethodsSet != nil {
		m := structToMap(options.ProxyCacheMethodsSet)
		result["proxy_cache_methods_set"] = []interface{}{m}
	}
	if options.ProxyConnectTimeout != nil {
		m := structToMap(options.ProxyConnectTimeout)
		result["proxy_connect_timeout"] = []interface{}{m}
	}
	if options.ProxyReadTimeout != nil {
		m := structToMap(options.ProxyReadTimeout)
		result["proxy_read_timeout"] = []interface{}{m}
	}
	if options.QueryParamsBlacklist != nil {
		m := structToMap(options.QueryParamsBlacklist)
		result["query_params_blacklist"] = []interface{}{m}
	}
	if options.QueryParamsWhitelist != nil {
		m := structToMap(options.QueryParamsWhitelist)
		result["query_params_whitelist"] = []interface{}{m}
	}
	if options.RedirectHttpsToHttp != nil {
		m := structToMap(options.RedirectHttpsToHttp)
		result["redirect_https_to_http"] = []interface{}{m}
	}
	if options.RedirectHttpToHttps != nil {
		m := structToMap(options.RedirectHttpToHttps)
		result["redirect_http_to_https"] = []interface{}{m}
	}
	if options.ReferrerACL != nil {
		m := structToMap(options.ReferrerACL)
		result["referrer_acl"] = []interface{}{m}
	}
	if options.RequestLimiter != nil {
		m := structToMap(options.RequestLimiter)
		result["request_limiter"] = []interface{}{m}
	}
	if options.ResponseHeadersHidingPolicy != nil {
		m := structToMap(options.ResponseHeadersHidingPolicy)
		result["response_headers_hiding_policy"] = []interface{}{m}
	}
	if options.Rewrite != nil {
		m := structToMap(options.Rewrite)
		result["rewrite"] = []interface{}{m}
	}
	if options.SecureKey != nil {
		m := structToMap(options.SecureKey)
		result["secure_key"] = []interface{}{m}
	}
	if options.Slice != nil {
		m := structToMap(options.Slice)
		result["slice"] = []interface{}{m}
	}
	if options.SNI != nil {
		m := structToMap(options.SNI)
		result["sni"] = []interface{}{m}
	}
	if options.Stale != nil {
		m := structToMap(options.Stale)
		result["stale"] = []interface{}{m}
	}
	if options.StaticHeaders != nil {
		m := structToMap(options.StaticHeaders)
		result["static_headers"] = []interface{}{m}
	}
	if options.StaticRequestHeaders != nil {
		m := structToMap(options.StaticRequestHeaders)
		result["static_request_headers"] = []interface{}{m}
	}
	if options.StaticResponseHeaders != nil {
		m := structToMap(options.StaticResponseHeaders)
		items := []interface{}{}
		for _, v := range m["value"].([]gcdn.StaticResponseHeadersItem) {
			items = append(items, structToMap(v))
		}
		m["value"] = items
		result["static_response_headers"] = []interface{}{m}
	}
	if options.TLSVersions != nil {
		m := structToMap(options.TLSVersions)
		result["tls_versions"] = []interface{}{m}
	}
	if options.UseDefaultLEChain != nil {
		m := structToMap(options.UseDefaultLEChain)
		result["use_default_le_chain"] = []interface{}{m}
	}
	if options.UserAgentACL != nil {
		m := structToMap(options.UserAgentACL)
		result["user_agent_acl"] = []interface{}{m}
	}
	if options.UseRSALECert != nil {
		m := structToMap(options.UseRSALECert)
		result["use_rsa_le_cert"] = []interface{}{m}
	}
	if options.WAF != nil {
		m := structToMap(options.WAF)
		result["waf"] = []interface{}{m}
	}
	if options.WebSockets != nil {
		m := structToMap(options.WebSockets)
		result["websockets"] = []interface{}{m}
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
