package gcore

import (
	"context"
	"fmt"
	"log"
	"net/http"

	waap "github.com/G-Core/gcore-waap-sdk-go"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceWaapApiPath() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWaapApiPathCreate,
		ReadContext:   resourceWaapApiPathRead,
		UpdateContext: resourceWaapApiPathUpdate,
		DeleteContext: resourceWaapApiPathDelete,
		Description:   "Represent API Paths for a specific WAAP domain",

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The WAAP domain ID for which the API Path is configured.",
			},
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The API path, locations that are saved for resource IDs will be put in curly brackets.",
			},
			"method": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The different methods an API path can have. It must be one of these values: GET, POST, PUT, PATCH, DELETE, TRACE, HEAD, OPTIONS.",
				ValidateFunc: validation.StringInSlice([]string{
					"GET",
					"POST",
					"PUT",
					"PATCH",
					"DELETE",
					"TRACE",
					"HEAD",
					"OPTIONS",
				}, false),
			},
			"http_scheme": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The different HTTP schemes an API path can have. It must be one of these values: HTTP, HTTPS.",
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP",
					"HTTPS",
				}, false),
			},
			"api_version": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The API version.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "An array of tags associated with the API path.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"api_groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "An array of API groups associated with the API path.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				Description:  "The status of the discovered API path. It must be one of these values: CONFIRMED_API, POTENTIAL_API, NOT_API, DELISTED_API",
				ValidateFunc: validation.StringInSlice([]string{"CONFIRMED_API", "POTENTIAL_API", "NOT_API", "DELISTED_API"}, false),
			},
		},
	}
}

func resourceWaapApiPathCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start WAAP API Path creating")

	client := m.(*Config).WaapClient

	req := waap.CreateApiPath{
		Path:       d.Get("path").(string),
		Method:     waap.ApiPathMethod(d.Get("method").(string)),
		HttpScheme: waap.ApiPathHttpScheme(d.Get("http_scheme").(string)),
	}

	if v, ok := d.GetOk("api_version"); ok {
		version := v.(string)
		req.ApiVersion = &version
	}

	if tags := convertStringList(d.Get("tags").([]interface{})); tags != nil && len(tags) > 0 {
		req.Tags = &tags
	}

	if apiGroups := convertStringList(d.Get("api_groups").([]interface{})); apiGroups != nil && len(apiGroups) > 0 {
		req.ApiGroups = &apiGroups
	}

	result, err := client.CreateApiPathV1DomainsDomainIdApiPathsPostWithResponse(ctx, d.Get("domain_id").(int), req)

	if err != nil {
		return diag.Errorf("Failed to create API Path: %w", err)
	}

	if result.StatusCode() != http.StatusCreated {
		return diag.Errorf("Failed to create API Path. Status code: %d with error: %s", result.StatusCode(), result.Body)
	}

	d.SetId(result.JSON201.Id.String())

	// Check if the status should be updated
	if status, ok := d.GetOk("status"); ok && status != "CONFIRMED_API" {
		pathId, _ := uuid.Parse(d.Id())
		pathStatus := waap.ApiPathStatus(status.(string))
		updateReq := waap.UpdateApiPath{
			Status: &pathStatus,
		}

		result, err := client.UpdateApiPathV1DomainsDomainIdApiPathsPathIdPatchWithResponse(ctx, d.Get("domain_id").(int), pathId, updateReq)

		if err != nil {
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Failed to update API Path status: %w", err),
				},
			}
		}

		if result.StatusCode() != http.StatusNoContent {
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Failed to update API Pathstatus. Status code: %d with error: %s", result.StatusCode(), result.Body),
				},
			}
		}
	}

	log.Printf("[DEBUG] Finish WAAP API Path creating (id=%s)\n", d.Id())

	return resourceWaapApiPathRead(ctx, d, m)
}

func resourceWaapApiPathRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start WAAP API Path reading")

	client := m.(*Config).WaapClient

	apiPathUUID, err := uuid.Parse(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing UUID: %v", err)
	}

	result, err := client.GetApiPathV1DomainsDomainIdApiPathsPathIdGetWithResponse(ctx, d.Get("domain_id").(int), apiPathUUID)
	if err != nil {
		return diag.Errorf("Failed to read API Path: %w", apiPathUUID, err)
	}

	if result.StatusCode() == http.StatusNotFound {
		d.SetId("") // Resource not found, remove from state
		return diag.Diagnostics{
			{Summary: fmt.Sprintf("API Path (%s) was not found, removed from TF state", apiPathUUID)},
		}
	}

	if result.StatusCode() != http.StatusOK {
		return diag.Errorf("Failed to read API Path. Status code: %d with error: %s", result.StatusCode(), result.Body)
	}

	d.Set("path", result.JSON200.Path)
	d.Set("method", result.JSON200.Method)
	d.Set("http_scheme", result.JSON200.HttpScheme)
	d.Set("api_version", result.JSON200.ApiVersion)
	d.Set("tags", result.JSON200.Tags)
	d.Set("api_groups", result.JSON200.ApiGroups)
	d.Set("status", result.JSON200.Status)

	log.Printf("[DEBUG] Finish WAAP API Path reading (id=%s)\n", apiPathUUID)
	return nil
}

func resourceWaapApiPathUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Start WAAP API Path updating (id=%s)\n", d.Id())

	client := m.(*Config).WaapClient

	pathUUID, err := uuid.Parse(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing UUID: %v", err)
	}

	path := d.Get("path").(string)
	emptyList := []string{}
	req := waap.UpdateApiPath{
		Path:      &path,
		ApiGroups: &emptyList,
		Tags:      &emptyList,
	}

	if v, ok := d.GetOk("status"); ok {
		status := waap.ApiPathStatus(v.(string))
		req.Status = &status
	}

	if apiGroups := convertStringList(d.Get("api_groups").([]interface{})); apiGroups != nil && len(apiGroups) > 0 {
		req.ApiGroups = &apiGroups
	}

	if tags := convertStringList(d.Get("tags").([]interface{})); tags != nil && len(tags) > 0 {
		req.Tags = &tags
	}

	result, err := client.UpdateApiPathV1DomainsDomainIdApiPathsPathIdPatchWithResponse(ctx, d.Get("domain_id").(int), pathUUID, req)

	if err != nil {
		return diag.Errorf("Failed to update API Path: %w", err)
	}

	if result.StatusCode() != http.StatusNoContent {
		return diag.Errorf("Failed to update API Path. Status code: %d with error: %s", result.StatusCode(), result.Body)
	}

	log.Printf("[DEBUG] Finish WAAP API Path updating (id=%s)", d.Id())
	return resourceWaapApiPathRead(ctx, d, m)
}

func resourceWaapApiPathDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Start WAAP API Path deleting (id=%s)\n", d.Id())

	client := m.(*Config).WaapClient

	apiPathUUID, err := uuid.Parse(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing UUID: %v", err)
	}

	result, err := client.DeleteApiPathV1DomainsDomainIdApiPathsPathIdDeleteWithResponse(ctx, d.Get("domain_id").(int), apiPathUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	if result.StatusCode() != http.StatusNoContent {
		return diag.Errorf("Failed to delete API PAth. Status code: %d with error: %s", result.StatusCode(), result.Body)
	}

	log.Printf("[DEBUG] Finish WAAP API Path deleting (id=%s)\n", d.Id())
	d.SetId("")

	return nil
}

func convertStringList(v []interface{}) []string {
	var result []string
	for _, item := range v {
		result = append(result, item.(string))
	}
	return result
}
