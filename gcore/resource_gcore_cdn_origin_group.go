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
				Required:    true,
				Description: "This options have two possible values: true — The option is active. In case the origin responds with 4XX or 5XX codes, use the next origin from the list. false — The option is disabled.",
			},
			"origin": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Contains information about all IP address or Domain names of your origin and the port if custom",
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
							Computed:    true,
							Description: "The setting allows to enable or disable an Origin source in the Origins group",
						},
						"backup": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "true — The option is active. The origin will not be used until one of active origins become unavailable. false — The option is disabled.",
						},
					},
				},
			},
			"proxy_next_upstream": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Available values: error, timeout, invalid_header, http_403, http_404, http_429, http_500, http_502, http_503, http_504.",
			},
		},
		CreateContext: resourceCDNOriginGroupCreate,
		ReadContext:   resourceCDNOriginGroupRead,
		UpdateContext: resourceCDNOriginGroupUpdate,
		DeleteContext: resourceCDNOriginGroupDelete,
		Description:   "Represent origin group",
	}
}

func resourceCDNOriginGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start CDN OriginGroup creating")
	config := m.(*Config)
	client := config.CDNClient

	var req origingroups.GroupRequest
	req.Name = d.Get("name").(string)
	req.UseNext = d.Get("use_next").(bool)
	req.Sources = setToSourceRequests(d.Get("origin").(*schema.Set))

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
	resourceCDNOriginGroupRead(ctx, d, m)

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
		return diag.FromErr(err)
	}

	d.Set("name", result.Name)
	d.Set("use_next", result.UseNext)
	if err := d.Set("origin", originsToSet(result.Sources)); err != nil {
		return diag.FromErr(err)
	}
	d.Set("proxy_next_upstream", result.ProxyNextUpstream)

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
	req.Sources = setToSourceRequests(d.Get("origin").(*schema.Set))

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

	log.Println("[DEBUG] Finish CDN OriginGroup updating")

	return resourceCDNOriginGroupRead(ctx, d, m)
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
