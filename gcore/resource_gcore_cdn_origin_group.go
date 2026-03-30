package gcore

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/G-Core/gcorelabscdn-go/origingroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCDNOriginGroup() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the origin group",
			},
			"use_next": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "This options have two possible values: true — The option is active. In case the origin responds with 4XX or 5XX codes, use the next origin from the list. false — The option is disabled.",
			},
			"proxy_next_upstream": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Available values: error, timeout, invalid_header, http_403, http_404, http_429, http_500, http_502, http_503, http_504.",
			},
			"origin": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Contains information about origins in the group. Each origin can be a host origin or an S3 origin. Host origins require `source`, S3 origins require `origin_type = \"s3\"` and a `config` block.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP address or Domain name of your origin and the port if custom. Required for host origins.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "The setting allows to enable or disable an Origin source in the Origins group. Default value is true.",
						},
						"backup": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Defines whether the origin is a backup, meaning that it will not be used until one of active origins become unavailable. Default value is false.",
						},
						"origin_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "host",
							ValidateFunc: validation.StringInSlice([]string{"host", "s3"}, false),
							Description:  "Type of the origin: 'host' (default) or 's3'.",
						},
						"host_header_override": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Per-origin host header override.",
						},
						"config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "S3 configuration for the origin. Required when origin_type is 's3'.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"s3_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"other", "amazon"}, false),
										Description:  "Type of S3 storage: 'amazon' or 'other'.",
									},
									"s3_bucket_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "S3 bucket name.",
									},
									"s3_region": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "S3 region. Required when s3_type is 'amazon'.",
									},
									"s3_storage_hostname": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "S3 storage hostname. Required when s3_type is 'other'.",
									},
									"s3_access_key_id": {
										Type:        schema.TypeString,
										Required:    true,
										Sensitive:   true,
										Description: "S3 access key ID.",
									},
									"s3_secret_access_key": {
										Type:        schema.TypeString,
										Required:    true,
										Sensitive:   true,
										Description: "S3 secret access key.",
									},
									"s3_auth_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "awsSignatureV4",
										ValidateFunc: validation.StringInSlice([]string{"awsSignatureV4"}, false),
										Description:  "S3 authentication type. Default: 'awsSignatureV4'.",
									},
								},
							},
						},
					},
				},
			},
			"auth": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Deprecated:  "Use `origin` blocks with `origin_type = \"s3\"` instead. The `auth` block will be removed in a future version.",
				Description: "Deprecated: Authentication configuration for S3 storage. Use `origin` blocks with `origin_type = \"s3\"` instead.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"s3_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"other", "amazon"}, false),
							Description:  "Type of the S3 storage, accepted values: 'other' or 'amazon'",
						},
						"s3_storage_hostname": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Hostname of the S3 storage, required if s3_type is 'other'",
						},
						"s3_access_key_id": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Access key ID for the S3 storage",
						},
						"s3_secret_access_key": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Secret access key for the S3 storage",
						},
						"s3_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Region of the S3 storage, required if s3_type is 'amazon'",
						},
						"s3_bucket_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Bucket name of the S3 storage",
						},
					},
				},
			},
		},
		CreateContext: resourceCDNOriginGroupCreate,
		ReadContext:   resourceCDNOriginGroupRead,
		UpdateContext: resourceCDNOriginGroupUpdate,
		DeleteContext: resourceCDNOriginGroupDelete,
		Description:   "Represent origin group",
		CustomizeDiff: validateCDNOriginGroupConfig,
	}
}

func validateCDNOriginGroupConfig(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	_, originExists := diff.GetOk("origin")
	authRaw, authExists := diff.GetOk("auth")

	if !originExists && !authExists {
		return fmt.Errorf("one of `origin` or `auth` must be specified")
	}

	if originExists && authExists {
		return fmt.Errorf("`origin` and `auth` cannot both be specified at the same time")
	}

	if authExists {
		authList := authRaw.([]interface{})

		if len(authList) > 0 {
			auth := authList[0].(map[string]interface{})
			s3Type := auth["s3_type"].(string)

			if s3Type == "other" {
				if storageHostname, ok := auth["s3_storage_hostname"].(string); !ok || storageHostname == "" {
					return fmt.Errorf("`s3_storage_hostname` is required when `s3_type` is 'other'")
				}
			}

			if s3Type == "amazon" {
				if s3Region, ok := auth["s3_region"].(string); !ok || s3Region == "" {
					return fmt.Errorf("`s3_region` is required when `s3_type` is 'amazon'")
				}
			}
		}
	}

	if originExists {
		originList := diff.Get("origin").([]interface{})
		for i, originRaw := range originList {
			if originRaw == nil {
				continue
			}
			origin := originRaw.(map[string]interface{})
			originType, _ := origin["origin_type"].(string)
			if originType == "" {
				originType = "host"
			}

			if originType == "host" {
				source, _ := origin["source"].(string)
				if source == "" {
					return fmt.Errorf("origin.%d: `source` is required for host origins", i)
				}
			}

			if originType == "s3" {
				configList, _ := origin["config"].([]interface{})
				if len(configList) > 0 && configList[0] != nil {
					if cfg, ok := configList[0].(map[string]interface{}); ok {
						if err := validateS3ConfigFields(i, cfg); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

func validateS3ConfigFields(index int, cfg map[string]interface{}) error {
	s3Type, _ := cfg["s3_type"].(string)

	if s3Type == "other" {
		if val, ok := cfg["s3_storage_hostname"].(string); !ok || val == "" {
			return fmt.Errorf("origin.%d.config: `s3_storage_hostname` is required when `s3_type` is 'other'", index)
		}
	}

	if s3Type == "amazon" {
		if val, ok := cfg["s3_region"].(string); !ok || val == "" {
			return fmt.Errorf("origin.%d.config: `s3_region` is required when `s3_type` is 'amazon'", index)
		}
	}

	return nil
}

func resourceCDNOriginGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start CDN OriginGroup creating")
	config := m.(*Config)
	client := config.CDNClient

	var req origingroups.GroupRequest
	req.Name = d.Get("name").(string)
	req.UseNext = d.Get("use_next").(bool)

	if originList, ok := d.GetOk("origin"); ok {
		req.AuthType = "none"
		req.Sources = listToSourceRequests(originList.([]interface{}))
	} else {
		req.AuthType = "awsSignatureV4"
		req.Auth = listToAuthS3(d.Get("auth").([]interface{}))
		req.Sources = nil
	}

	proxyNextUpstream, ok := d.Get("proxy_next_upstream").(*schema.Set)
	if ok && proxyNextUpstream.Len() > 0 {
		req.ProxyNextUpstream = make([]string, 0)
		for _, upstreamError := range d.Get("proxy_next_upstream").(*schema.Set).List() {
			req.ProxyNextUpstream = append(req.ProxyNextUpstream, upstreamError.(string))
		}
	}

	result, err := client.OriginGroups().Create(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", result.ID))
	diags := resourceCDNOriginGroupRead(ctx, d, m)

	if diags.HasError() {
		return diags
	}

	if _, ok := d.GetOk("auth"); ok {
		d.Set("auth.0.s3_secret_access_key", req.Auth.S3SecretAccessKey)
		d.Set("auth.0.s3_access_key_id", req.Auth.S3AccessKeyID)
	}

	log.Printf("[DEBUG] Finish CDN OriginGroup creating (id=%d)\n", result.ID)
	return nil
}

func resourceCDNOriginGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Id()
	log.Printf("[DEBUG] Start CDN OriginGroup reading (id=%s)\n", groupID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(groupID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.OriginGroups().Get(ctx, id)
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("") // Resource not found, remove from state
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Origin Group \"%s\" not found, removing from state", d.Get("name")),
				},
			}
		}

		return diag.FromErr(err)
	}

	d.Set("name", result.Name)
	d.Set("use_next", result.UseNext)
	d.Set("proxy_next_upstream", result.ProxyNextUpstream)

	// Preserve S3 credentials from current state before setting origins from API.
	// The API masks sensitive values, so we keep the user-provided ones.
	s3Credentials := preserveS3OriginCredentials(d)

	originsList := sourcesToList(result.Sources)
	restoreS3OriginCredentials(originsList, s3Credentials)

	if err := d.Set("origin", originsList); err != nil {
		return diag.FromErr(err)
	}

	// Legacy auth block: keep s3_secret_access_key and s3_access_key_id unchanged by API response
	currentSecretAccessKey, keyExists := d.GetOk("auth.0.s3_secret_access_key")
	currentAccessKeyID, keyIDExists := d.GetOk("auth.0.s3_access_key_id")
	if err := d.Set("auth", authToList(result.Auth)); err != nil {
		return diag.FromErr(err)
	}
	if keyExists && keyIDExists {
		authList := d.Get("auth").([]interface{})
		if len(authList) > 0 {
			authMap := authList[0].(map[string]interface{})
			authMap["s3_secret_access_key"] = currentSecretAccessKey
			authMap["s3_access_key_id"] = currentAccessKeyID
			if err := d.Set("auth", []interface{}{authMap}); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	log.Println("[DEBUG] Finish CDN OriginGroup reading")
	return nil
}

func resourceCDNOriginGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Id()
	log.Printf("[DEBUG] Start CDN OriginGroup updating (id=%s)\n", groupID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(groupID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var req origingroups.GroupRequest
	req.Name = d.Get("name").(string)
	req.UseNext = d.Get("use_next").(bool)

	if originList, ok := d.GetOk("origin"); ok {
		req.AuthType = "none"
		req.Sources = listToSourceRequests(originList.([]interface{}))
	} else {
		req.AuthType = "awsSignatureV4"
		req.Auth = listToAuthS3(d.Get("auth").([]interface{}))
		req.Sources = nil
	}

	if req.UseNext {
		proxyNextUpstream, ok := d.Get("proxy_next_upstream").(*schema.Set)
		if ok && proxyNextUpstream.Len() > 0 {
			req.ProxyNextUpstream = make([]string, 0)
			for _, upstreamError := range d.Get("proxy_next_upstream").(*schema.Set).List() {
				req.ProxyNextUpstream = append(req.ProxyNextUpstream, upstreamError.(string))
			}
		}
	}

	if _, err := client.OriginGroups().Update(ctx, id, &req); err != nil {
		return diag.FromErr(err)
	}

	diags := resourceCDNOriginGroupRead(ctx, d, m)

	if diags.HasError() {
		return diags
	}

	if authList, ok := d.GetOk("auth"); ok {
		if len(authList.([]interface{})) > 0 {
			authConfig := authList.([]interface{})[0].(map[string]interface{})
			if secretAccessKey, ok := authConfig["s3_secret_access_key"].(string); ok {
				d.Set("s3_secret_access_key", secretAccessKey)
			}
			if accessKeyID, ok := authConfig["s3_access_key_id"].(string); ok {
				d.Set("s3_access_key_id", accessKeyID)
			}
		}
	}

	log.Println("[DEBUG] Finish CDN OriginGroup updating")
	return nil
}

func resourceCDNOriginGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceID := d.Id()
	log.Printf("[DEBUG] Start CDN OriginGroup deleting (id=%s)\n", resourceID)

	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.OriginGroups().Delete(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Println("[DEBUG] Finish CDN Resource deleting")
	return nil
}

func listToSourceRequests(origins []interface{}) []origingroups.SourceRequest {
	var result []origingroups.SourceRequest

	for _, raw := range origins {
		if raw == nil {
			continue
		}
		fields := raw.(map[string]interface{})
		originType, _ := fields["origin_type"].(string)
		if originType == "" {
			originType = "host"
		}

		originReq := origingroups.SourceRequest{
			Enabled: fields["enabled"].(bool),
			Backup:  fields["backup"].(bool),
		}

		if hostHeader, ok := fields["host_header_override"].(string); ok && hostHeader != "" {
			originReq.HostHeaderOverride = &hostHeader
		}

		if originType == "s3" {
			originReq.OriginType = "s3"
			configList, _ := fields["config"].([]interface{})
			if len(configList) > 0 && configList[0] != nil {
				cfg := configList[0].(map[string]interface{})
				s3AuthType, _ := cfg["s3_auth_type"].(string)
				if s3AuthType == "" {
					s3AuthType = "awsSignatureV4"
				}
				originReq.Config = &origingroups.S3Config{
					S3Type:            cfg["s3_type"].(string),
					S3BucketName:      cfg["s3_bucket_name"].(string),
					S3AccessKeyID:     cfg["s3_access_key_id"].(string),
					S3SecretAccessKey: cfg["s3_secret_access_key"].(string),
					S3AuthType:        s3AuthType,
				}
				if region, ok := cfg["s3_region"].(string); ok && region != "" {
					originReq.Config.S3Region = region
				}
				if hostname, ok := cfg["s3_storage_hostname"].(string); ok && hostname != "" {
					originReq.Config.S3StorageHostname = hostname
				}
			}
		} else {
			originReq.Source, _ = fields["source"].(string)
		}

		result = append(result, originReq)
	}

	return result
}

func sourcesToList(sources []origingroups.Source) []interface{} {
	var result []interface{}

	for _, origin := range sources {
		fields := map[string]interface{}{
			"enabled":              origin.Enabled,
			"backup":               origin.Backup,
			"source":               origin.Source,
			"origin_type":          "host",
			"host_header_override": "",
			"config":               []interface{}{},
		}

		if origin.HostHeaderOverride != nil {
			fields["host_header_override"] = *origin.HostHeaderOverride
		}

		if origin.OriginType == "s3" && origin.Config != nil {
			fields["origin_type"] = "s3"
			fields["source"] = ""
			cfgMap := map[string]interface{}{
				"s3_type":              origin.Config.S3Type,
				"s3_bucket_name":       origin.Config.S3BucketName,
				"s3_region":            origin.Config.S3Region,
				"s3_storage_hostname":  origin.Config.S3StorageHostname,
				"s3_access_key_id":     origin.Config.S3AccessKeyID,
				"s3_secret_access_key": origin.Config.S3SecretAccessKey,
				"s3_auth_type":         "awsSignatureV4",
			}
			if origin.Config.S3AuthType != "" {
				cfgMap["s3_auth_type"] = origin.Config.S3AuthType
			}
			fields["config"] = []interface{}{cfgMap}
		}

		result = append(result, fields)
	}

	return result
}

// s3OriginCredentials stores S3 credentials keyed by bucket name for state preservation.
type s3OriginCredentials struct {
	accessKeyID     string
	secretAccessKey string
}

// preserveS3OriginCredentials extracts S3 credentials from current state before API read overwrites them.
func preserveS3OriginCredentials(d *schema.ResourceData) map[string]s3OriginCredentials {
	creds := make(map[string]s3OriginCredentials)
	originList, ok := d.GetOk("origin")
	if !ok {
		return creds
	}
	for _, raw := range originList.([]interface{}) {
		if raw == nil {
			continue
		}
		origin := raw.(map[string]interface{})
		originType, _ := origin["origin_type"].(string)
		if originType != "s3" {
			continue
		}
		configList, ok := origin["config"].([]interface{})
		if !ok || len(configList) == 0 || configList[0] == nil {
			continue
		}
		cfg := configList[0].(map[string]interface{})
		bucketName, _ := cfg["s3_bucket_name"].(string)
		if bucketName == "" {
			continue
		}
		accessKeyID, _ := cfg["s3_access_key_id"].(string)
		secretAccessKey, _ := cfg["s3_secret_access_key"].(string)
		if accessKeyID != "" || secretAccessKey != "" {
			creds[bucketName] = s3OriginCredentials{
				accessKeyID:     accessKeyID,
				secretAccessKey: secretAccessKey,
			}
		}
	}
	return creds
}

// restoreS3OriginCredentials restores S3 credentials in the origins list after API read.
func restoreS3OriginCredentials(origins []interface{}, creds map[string]s3OriginCredentials) {
	if len(creds) == 0 {
		return
	}
	for _, raw := range origins {
		if raw == nil {
			continue
		}
		origin := raw.(map[string]interface{})
		originType, _ := origin["origin_type"].(string)
		if originType != "s3" {
			continue
		}
		configList, ok := origin["config"].([]interface{})
		if !ok || len(configList) == 0 || configList[0] == nil {
			continue
		}
		cfg := configList[0].(map[string]interface{})
		bucketName, _ := cfg["s3_bucket_name"].(string)
		if saved, ok := creds[bucketName]; ok {
			cfg["s3_access_key_id"] = saved.accessKeyID
			cfg["s3_secret_access_key"] = saved.secretAccessKey
		}
	}
}

func listToAuthS3(authList []interface{}) *origingroups.AuthS3 {
	if len(authList) == 0 {
		return nil
	}

	authConfig := authList[0].(map[string]interface{})

	auth := &origingroups.AuthS3{
		S3Type:            authConfig["s3_type"].(string),
		S3AccessKeyID:     authConfig["s3_access_key_id"].(string),
		S3SecretAccessKey: authConfig["s3_secret_access_key"].(string),
		S3BucketName:      authConfig["s3_bucket_name"].(string),
	}

	if s3StorageHostname, ok := authConfig["s3_storage_hostname"]; ok {
		auth.S3StorageHostname = s3StorageHostname.(string)
	}

	if s3Region, ok := authConfig["s3_region"]; ok {
		auth.S3Region = s3Region.(string)
	}

	return auth
}

func authToList(auth *origingroups.AuthS3) []interface{} {
	if auth == nil {
		return nil
	}

	authMap := map[string]interface{}{
		"s3_type":        auth.S3Type,
		"s3_bucket_name": auth.S3BucketName,
	}

	if auth.S3StorageHostname != "" {
		authMap["s3_storage_hostname"] = auth.S3StorageHostname
	}

	if auth.S3Region != "" {
		authMap["s3_region"] = auth.S3Region
	}

	return []interface{}{authMap}
}
