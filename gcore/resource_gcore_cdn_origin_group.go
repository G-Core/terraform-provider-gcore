package gcore

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
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
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Contains information about all IP address or Domain names of your origin and the port if custom. This field is required unless `auth` is specified. `origin` and `auth` cannot both be specified simultaneously.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP address or Domain name of your origin and the port if custom",
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
					},
				},
			},
			"auth": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Authentication configuration for S3 storage. This field is required unless `origin` is specified. `auth` and `origin` cannot both be specified simultaneously.",
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
		return fmt.Errorf("One of `origin` or `auth` must be specified")
	}

	if originExists && authExists {
		return fmt.Errorf("Both `origin` and `auth` cannot be specified at the same time")
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

	return nil
}

func resourceCDNOriginGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start CDN OriginGroup creating")
	config := m.(*Config)
	client := config.CDNClient

	var req origingroups.GroupRequest
	req.Name = d.Get("name").(string)
	req.UseNext = d.Get("use_next").(bool)

	if originSet, ok := d.GetOk("origin"); ok {
		req.AuthType = "none"
		req.Sources = setToSourceRequests(originSet.(*schema.Set))
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
	if err := d.Set("origin", originsToSet(result.Sources)); err != nil {
		return diag.FromErr(err)
	}
	d.Set("proxy_next_upstream", result.ProxyNextUpstream)

	// keep s3_secret_access_key and s3_access_key_id unchanged by API response
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

	if originSet, ok := d.GetOk("origin"); ok {
		req.AuthType = "none"
		req.Sources = setToSourceRequests(originSet.(*schema.Set))
	} else {
		req.AuthType = "awsSignatureV4"
		req.Auth = listToAuthS3(d.Get("auth").([]interface{}))
		req.Sources = nil
	}

	if req.UseNext == true {
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

func setToSourceRequests(s *schema.Set) (origins []origingroups.SourceRequest) {
	for _, fields := range s.List() {
		var originReq origingroups.SourceRequest

		for key, val := range fields.(map[string]interface{}) {
			switch key {
			case "source":
				originReq.Source = val.(string)
			case "enabled":
				originReq.Enabled = val.(bool)
			case "backup":
				originReq.Backup = val.(bool)
			}
		}

		origins = append(origins, originReq)
	}

	return origins
}

func originsToSet(origins []origingroups.Source) *schema.Set {
	s := &schema.Set{F: originSetIDFunc}

	for _, origin := range origins {
		fields := make(map[string]interface{})
		fields["source"] = origin.Source
		fields["enabled"] = origin.Enabled
		fields["backup"] = origin.Backup

		s.Add(fields)
	}

	return s
}

func originSetIDFunc(i interface{}) int {
	fields := i.(map[string]interface{})
	h := md5.New()

	key := fmt.Sprintf("%d-%s-%t", fields["source"], fields["enabled"], fields["backup"])
	log.Printf("[DEBUG] Origin Set ID = %s\n", key)

	io.WriteString(h, key)

	return int(binary.BigEndian.Uint64(h.Sum(nil)))
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
