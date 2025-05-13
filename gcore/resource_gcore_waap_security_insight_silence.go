package gcore

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	waap "github.com/G-Core/gcore-waap-sdk-go"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWaapSecurityInsightSilence() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityInsightSilenceCreate,
		ReadContext:   resourceSecurityInsightSilenceRead,
		UpdateContext: resourceSecurityInsightSilenceUpdate,
		DeleteContext: resourceSecurityInsightSilenceDelete,
		Description:   "Represent Security Insight Silence for a specific WAAP domain",

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A generated unique identifier for the silence.",
			},
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The WAAP domain ID for which the insight silence is configured.",
			},
			"insight_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The slug of the insight type.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "A hash table of label names and values that apply to the insight silence.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"comment": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A comment explaining the reason for the silence.",
			},
			"author": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The author of the silence.",
			},
			"expire_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The date and time the silence expires in ISO 8601 format.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					strVal := val.(string)
					if strVal != "" {
						_, err := time.Parse(time.RFC3339, strVal)
						if err != nil {
							errs = append(errs, fmt.Errorf("%q must be a valid RFC3339 date: %s", key, err))
						}
					}
					return
				},
			},
		},
	}
}

func resourceSecurityInsightSilenceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start WAAP Security Insight Silence creating")

	client := m.(*Config).WaapClient

	req := waap.CreateInsightSilencePayload{
		InsightType: d.Get("insight_type").(string),
		Comment:     d.Get("comment").(string),
		Author:      d.Get("author").(string),
	}

	rawLabels := d.Get("labels").(map[string]interface{})
	labels := make(map[string]string, len(rawLabels))
	for key, value := range rawLabels {
		labels[key] = value.(string)
	}
	req.Labels = labels

	if expireAtStr, ok := d.GetOk("expire_at"); ok {
		expireAt, err := time.Parse(time.RFC3339, expireAtStr.(string))
		if err != nil {
			return diag.Errorf("Error parsing time: %s", err)
		}
		req.ExpireAt = &expireAt
	}

	result, err := client.CreateInsightSilenceV1DomainsDomainIdInsightSilencesPostWithResponse(ctx, d.Get("domain_id").(int), req)

	if err != nil {
		return diag.FromErr(err)
	}

	if result.StatusCode() != http.StatusOK {
		return diag.Errorf("Failed to create Insight Silence. Status code: %d with error: %s", result.StatusCode(), result.Body)
	}

	d.SetId(result.JSON200.Id.String())

	log.Printf("[DEBUG] Finish WAAP Security Insight Silence creating (id=%d)\n", result.JSON200.Id)
	return resourceSecurityInsightSilenceRead(ctx, d, m)
}

func resourceSecurityInsightSilenceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Start WAAP Security Insight Silence reading (id=%s)\n", d.Id())

	client := m.(*Config).WaapClient

	insightSilenceUUID, err := uuid.Parse(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing UUID: %v", err)
	}

	result, err := client.GetInsightSilenceV1DomainsDomainIdInsightSilencesSilenceIdGetWithResponse(ctx, d.Get("domain_id").(int), insightSilenceUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	if result.StatusCode() == http.StatusNotFound {
		d.SetId("") // Resource not found, remove from state
		return diag.Diagnostics{
			{Summary: fmt.Sprintf("Insight Silence (%s) was not found, removed from TF state", d.Id())},
		}
	}

	if result.StatusCode() != http.StatusOK {
		return diag.Errorf("Failed to read Insight Silence. Status code: %d with error: %s", result.StatusCode(), result.Body)
	}

	d.Set("insight_type", result.JSON200.InsightType)
	d.Set("labels", result.JSON200.Labels)
	d.Set("comment", result.JSON200.Comment)
	d.Set("author", result.JSON200.Author)
	d.Set("expire_at", result.JSON200.ExpireAt)

	log.Printf("[DEBUG] Finish WAAP Security Insight Silence reading (id=%s)\n", d.Id())
	return nil
}

func resourceSecurityInsightSilenceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Start WAAP Security Insight Silence updating (id=%s)\n", d.Id())

	client := m.(*Config).WaapClient

	insightSilenceUUID, err := uuid.Parse(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing UUID: %v", err)
	}

	req := waap.UpdateInsightSilencePayload{
		Comment: d.Get("comment").(string),
		Author:  d.Get("author").(string),
	}

	if expireAtStr, ok := d.GetOk("expire_at"); ok {
		expireAt, err := time.Parse(time.RFC3339, expireAtStr.(string))
		if err != nil {
			return diag.Errorf("Error parsing time: %s", err)
		}
		req.ExpireAt = &expireAt
	}

	rawLabels := d.Get("labels").(map[string]interface{})
	labels := make(map[string]string, len(rawLabels))
	for key, value := range rawLabels {
		labels[key] = value.(string)
	}
	req.Labels = &labels

	result, err := client.UpdateInsightSilenceV1DomainsDomainIdInsightSilencesSilenceIdPatchWithResponse(ctx, d.Get("domain_id").(int), insightSilenceUUID, req)
	if err != nil {
		return diag.FromErr(err)
	}

	if result.StatusCode() != http.StatusOK {
		return diag.Errorf("Failed to update Insight Silence. Status code: %d with error: %s", result.StatusCode(), result.Body)
	}

	log.Printf("[DEBUG] Finish WAAP Security Insight Silence updating (id=%s)", d.Id())
	return resourceSecurityInsightSilenceRead(ctx, d, m)
}

func resourceSecurityInsightSilenceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Start WAAP Security Insight Silence deleting (id=%s)\n", d.Id())

	client := m.(*Config).WaapClient

	insightSilenceUUID, err := uuid.Parse(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing UUID: %v", err)
	}

	result, err := client.DeleteInsightSilenceV1DomainsDomainIdInsightSilencesSilenceIdDeleteWithResponse(ctx, d.Get("domain_id").(int), insightSilenceUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	if result.StatusCode() != http.StatusNoContent {
		return diag.Errorf("Failed to delete Insight Silence. Status code: %d with error: %s", result.StatusCode(), result.Body)
	}

	log.Printf("[DEBUG] Finish WAAP Security Insight Silence deleting (id=%s)\n", d.Id())
	d.SetId("")

	return nil
}
