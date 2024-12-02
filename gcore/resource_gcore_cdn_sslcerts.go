package gcore

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/G-Core/gcorelabscdn-go/sslcerts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCDNCert() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the SSL certificate. Must be unique.",
			},
			"cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "The public part of the SSL certificate. All chain of the SSL certificate should be added.",
			},
			"private_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "The private key of the SSL certificate.",
			},
			"has_related_resources": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "It shows if the SSL certificate is used by a CDN resource.",
			},
			"automated": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "The way SSL certificate was issued.",
			},
		},
		CreateContext: resourceCDNCertCreate,
		ReadContext:   resourceCDNCertRead,
		UpdateContext: resourceCDNCertUpdate,
		DeleteContext: resourceCDNCertDelete,
	}
}

func resourceCDNCertCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start CDN Cert creating")
	config := m.(*Config)
	client := config.CDNClient

	var req sslcerts.CreateRequest
	req.Name = d.Get("name").(string)
	req.Cert = d.Get("cert").(string)
	req.PrivateKey = d.Get("private_key").(string)

	if d.Get("automated") != nil {
		req.Automated = d.Get("automated").(bool)
	}

	result, err := client.SSLCerts().Create(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", result.ID))
	resourceCDNCertRead(ctx, d, m)

	log.Printf("[DEBUG] Finish CDN Cert creating (id=%d)\n", result.ID)
	return nil
}

func resourceCDNCertRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	certID := d.Id()
	log.Printf("[DEBUG] Start CDN Cert reading (id=%s)\n", certID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(certID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.SSLCerts().Get(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", result.Name)
	d.Set("has_related_resources", result.HasRelatedResources)
	d.Set("automated", result.Automated)

	log.Println("[DEBUG] Finish CDN Cert reading")
	return nil
}

func resourceCDNCertUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	certID := d.Id()
	log.Printf("[DEBUG] Start CDN Cert updating (id=%s)\n", certID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(certID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var req sslcerts.UpdateRequest
	req.Name = d.Get("name").(string)

	if _, err := client.SSLCerts().Update(ctx, id, &req); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish CDN Cert updating")
	return resourceCDNCertRead(ctx, d, m)
}

func resourceCDNCertDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	certID := d.Id()
	log.Printf("[DEBUG] Start CDN Cert deleting (id=%s)\n", certID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(certID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.SSLCerts().Delete(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Println("[DEBUG] Finish CDN Cert deleting")
	return nil
}
