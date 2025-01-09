package gcore

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/G-Core/gcorelabscdn-go/cacerts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCDNCACert() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the CA certificate. Must be unique.",
			},
			"cert": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "The public part of the CA certificate. It must be in the PEM format.",
			},
			"has_related_resources": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "It shows if the CA certificate is used by a CDN resource.",
			},
		},
		CreateContext: resourceCDNCACertCreate,
		ReadContext:   resourceCDNCACertRead,
		UpdateContext: resourceCDNCACertUpdate,
		DeleteContext: resourceCDNCACertDelete,
		Description:   "Represent CDN CA Certificate",
	}
}

func resourceCDNCACertCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start CDN CA Cert creating")
	config := m.(*Config)
	client := config.CDNClient

	var req cacerts.CreateRequest
	req.Name = d.Get("name").(string)
	req.Cert = d.Get("cert").(string)

	result, err := client.CACerts().Create(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", result.ID))
	resourceCDNCACertRead(ctx, d, m)

	log.Printf("[DEBUG] Finish CDN CA Cert creating (id=%d)\n", result.ID)
	return nil
}

func resourceCDNCACertRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	certID := d.Id()
	log.Printf("[DEBUG] Start CDN CA Cert reading (id=%s)\n", certID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(certID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.CACerts().Get(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", result.Name)
	d.Set("has_related_resources", result.HasRelatedResources)

	log.Println("[DEBUG] Finish CDN CA Cert reading")
	return nil
}

func resourceCDNCACertUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	certID := d.Id()
	log.Printf("[DEBUG] Start CDN CA Cert updating (id=%s)\n", certID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(certID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var req cacerts.UpdateRequest
	req.Name = d.Get("name").(string)

	if _, err := client.CACerts().Update(ctx, id, &req); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish CDN CA Cert updating")
	return resourceCDNCACertRead(ctx, d, m)
}

func resourceCDNCACertDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	certID := d.Id()
	log.Printf("[DEBUG] Start CDN CA Cert deleting (id=%s)\n", certID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(certID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.CACerts().Delete(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Println("[DEBUG] Finish CDN CA Cert deleting")
	return nil
}
