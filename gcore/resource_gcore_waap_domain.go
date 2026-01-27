package gcore

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	waap "github.com/G-Core/gcore-waap-sdk-go"
)

func resourceWaapDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWaapDomainCreate,
		ReadContext:   resourceWaapDomainRead,
		UpdateContext: resourceWaapDomainUpdate,
		DeleteContext: resourceWaapDomainDelete,
		Description:   "Represent WAAP domain",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the domain.",
				ForceNew:    true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"active",
					"monitor",
				}, false),
				Description: "Status of the domain. It must be one of these values {active, monitor}.",
			},
			"settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ddos": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"global_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Global threshold for DDoS protection.",
									},
									"burst_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Burst threshold for DDoS protection.",
									},
								},
							},
						},
						"api": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_urls": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "List of API URL patterns.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"is_api": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
										Description: "Indicates if the domain is an API domain. " +
											"All requests to an API domain are treated as API requests. " +
											"If this is set to true then the api_urls field is ignored.",
									},
								},
							},
						},
					},
				},
			},
			"api_discovery_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description_file_location": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The URL of the API description file. This will be periodically scanned if `description_file_scan_enabled` is enabled. Supported formats are YAML and JSON, and it must adhere to OpenAPI versions 2, 3, or 3.1.",
						},
						"description_file_scan_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates if periodic scan of the description file is enabled.",
						},
						"description_file_scan_interval_hours": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The interval in hours for scanning the description file.",
						},
						"traffic_scan_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates if traffic scan is enabled.",
						},
						"traffic_scan_interval_hours": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The interval in hours for scanning the traffic.",
						},
					},
				},
			},
		},
	}
}

func resourceWaapDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).WaapClient
	domainName := d.Get("name").(string)
	params := waap.GetDomainsV1DomainsGetParams{
		Name: &domainName,
	}

	resp, err := client.GetDomainsV1DomainsGetWithResponse(ctx, &params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing domains: %v", err))
	}

	if resp.StatusCode() != http.StatusOK {
		return diag.Errorf("Failed to read Domains. Status code: %d with error: %s", resp.StatusCode(), resp.Body)
	}

	domain := findDomainByName(*resp.JSON200, domainName)

	if domain == nil {
		return diag.Errorf("Domain with name '%s' not found.", domainName)
	}

	status := string(domain.Status)

	// Compare domain status and update if needed
	if newStatus, ok := d.GetOk("status"); ok && newStatus != domain.Status {
		domainStatusUpdate := waap.DomainUpdateStatus(newStatus.(string))
		updateReq := waap.UpdateDomain{
			Status: &domainStatusUpdate,
		}
		updateResp, err := client.UpdateDomainV1DomainsDomainIdPatchWithResponse(ctx, domain.Id, updateReq)

		if err != nil {
			return diag.Errorf("Failed to update Domain status: %s", err)
		}

		if updateResp.StatusCode() != http.StatusNoContent {
			return diag.Errorf("Failed to update Domain status. Status code: %d with error: %s", updateResp.StatusCode(), updateResp.Body)
		}

		status = newStatus.(string)
	}

	// Update domain settings
	if settings, ok := d.GetOk("settings"); ok {
		updateSettingsResp, err := updateDomainSettings(ctx, client, settings, domain.Id)

		if err != nil {
			return diag.Errorf("Failed to update Domain settings: %s", err)
		}

		if updateSettingsResp.StatusCode() != http.StatusNoContent {
			return diag.Errorf("Failed to update Domain settings. Status code: %d with error: %s", updateSettingsResp.StatusCode(), updateSettingsResp.Body)
		}
	}

	// Update API Discovery settings
	if apiDiscoverySettings, ok := d.GetOk("api_discovery_settings"); ok {
		// Check if API Discovery feature is available
		hasFeature, diags := hasApiDiscoveryFeature(ctx, client)
		if diags.HasError() {
			return diags
		}

		if !hasFeature {
			return diag.Errorf("API Discovery feature is not available for your account. Please contact support or upgrade your plan to use api_discovery_settings.")
		}

		updateApiDiscoverySettingsResp, err := updateApiDiscoverySettings(ctx, client, apiDiscoverySettings, domain.Id)

		if err != nil {
			return diag.Errorf("Failed to update API Discovery settings: %s", err)
		}

		if updateApiDiscoverySettingsResp.StatusCode() != http.StatusOK {
			return diag.Errorf("Failed to update API Discovery settings. Status code: %d with error: %s", updateApiDiscoverySettingsResp.StatusCode(), updateApiDiscoverySettingsResp.Body)
		}
	}

	// Update state
	d.SetId(fmt.Sprintf("%d", domain.Id))
	d.Set("status", status)

	return resourceWaapDomainRead(ctx, d, m)
}

func resourceWaapDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).WaapClient

	domainID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Get domain details
	resp, err := client.GetDomainV1DomainsDomainIdGetWithResponse(ctx, domainID)
	if err != nil {
		return diag.Errorf("Failed to read Domain details: %s", err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		d.SetId("") // Resource not found, remove from state
		return diag.Diagnostics{
			{Severity: diag.Warning, Summary: fmt.Sprintf("Domain (%s) was not found, removed from TF state", d.Id())},
		}
	}

	if resp.StatusCode() != http.StatusOK {
		return diag.Errorf("Failed to read Domain details. Status code: %d with error: %s", resp.StatusCode(), resp.Body)
	}

	d.Set("status", string(resp.JSON200.Status))

	// Get domain settings
	settingsResp, err := client.GetDomainSettingsV1DomainsDomainIdSettingsGetWithResponse(ctx, domainID)
	if err != nil {
		return diag.Errorf("Failed to read Domain settings: %s", err)
	}

	if settingsResp.StatusCode() != http.StatusOK {
		return diag.Errorf("Failed to read Domain settings. Status code: %d with error: %s", settingsResp.StatusCode(), settingsResp.Body)
	}

	settings := make(map[string]interface{})
	ddosSettings := make(map[string]interface{})

	if settingsResp.JSON200.Ddos.GlobalThreshold != nil {
		ddosSettings["global_threshold"] = *settingsResp.JSON200.Ddos.GlobalThreshold
	}

	if settingsResp.JSON200.Ddos.BurstThreshold != nil {
		ddosSettings["burst_threshold"] = *settingsResp.JSON200.Ddos.BurstThreshold
	}

	if len(ddosSettings) > 0 {
		settings["ddos"] = []interface{}{ddosSettings}
	}

	if settingsResp.JSON200.Api.ApiUrls != nil {
		apiSettings := make(map[string]interface{})
		apiSettings["api_urls"] = *settingsResp.JSON200.Api.ApiUrls
		apiSettings["is_api"] = *settingsResp.JSON200.Api.IsApi
		settings["api"] = []interface{}{apiSettings}
	}

	if len(settings) > 0 {
		d.Set("settings", []interface{}{settings})
	}

	// Check if API Discovery feature is available before attempting to fetch settings
	hasApiDiscovery, diags := hasApiDiscoveryFeature(ctx, client)
	if diags.HasError() {
		return diags
	}

	if hasApiDiscovery {
		// Get API Discovery settings
		apiDiscoverySettingsResp, err := client.GetApiDiscoverySettingsV1DomainsDomainIdApiDiscoverySettingsGetWithResponse(ctx, domainID)
		if err != nil {
			return diag.Errorf("Failed to read API Discovery settings: %s", err)
		}

		if apiDiscoverySettingsResp.StatusCode() != http.StatusOK {
			return diag.Errorf("Failed to read API Discovery settings. Status code: %d with error: %s", apiDiscoverySettingsResp.StatusCode(), apiDiscoverySettingsResp.Body)
		}

		apiDiscoverySettings := make(map[string]interface{})

		if apiDiscoverySettingsResp.JSON200.DescriptionFileLocation != nil {
			apiDiscoverySettings["description_file_location"] = *apiDiscoverySettingsResp.JSON200.DescriptionFileLocation
		}

		if apiDiscoverySettingsResp.JSON200.DescriptionFileScanEnabled != nil {
			apiDiscoverySettings["description_file_scan_enabled"] = *apiDiscoverySettingsResp.JSON200.DescriptionFileScanEnabled
		}

		if apiDiscoverySettingsResp.JSON200.DescriptionFileScanIntervalHours != nil {
			apiDiscoverySettings["description_file_scan_interval_hours"] = *apiDiscoverySettingsResp.JSON200.DescriptionFileScanIntervalHours
		}

		if apiDiscoverySettingsResp.JSON200.TrafficScanEnabled != nil {
			apiDiscoverySettings["traffic_scan_enabled"] = *apiDiscoverySettingsResp.JSON200.TrafficScanEnabled
		}

		if apiDiscoverySettingsResp.JSON200.TrafficScanIntervalHours != nil {
			apiDiscoverySettings["traffic_scan_interval_hours"] = *apiDiscoverySettingsResp.JSON200.TrafficScanIntervalHours
		}

		if len(apiDiscoverySettings) > 0 {
			d.Set("api_discovery_settings", []interface{}{apiDiscoverySettings})
		}
	} else {
		// If API Discovery feature is not available, clear the field
		d.Set("api_discovery_settings", nil)
	}

	return nil
}

func resourceWaapDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).WaapClient
	domainID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Update domain status
	if d.HasChange("status") {
		if status, ok := d.GetOk("status"); ok {
			domainStatusUpdate := waap.DomainUpdateStatus(status.(string))
			updateStatusReq := waap.UpdateDomain{
				Status: &domainStatusUpdate,
			}
			updateResp, err := client.UpdateDomainV1DomainsDomainIdPatchWithResponse(ctx, domainID, updateStatusReq)

			if err != nil {
				return diag.Errorf("Failed to update Domain status: %s", err)
			}

			if updateResp.StatusCode() != http.StatusNoContent {
				return diag.Errorf("Failed to update Domain status. Status code: %d with error: %s", updateResp.StatusCode(), updateResp.Body)
			}
		}
	}

	// Update domain settings
	if d.HasChange("settings") {
		if settings, ok := d.GetOk("settings"); ok {
			updateSettingsResp, err := updateDomainSettings(ctx, client, settings, domainID)

			if err != nil {
				return diag.Errorf("Failed to update Domain settings: %s", err)
			}

			if updateSettingsResp.StatusCode() != http.StatusNoContent {
				return diag.Errorf("Failed to update Domain settings. Status code: %d with error: %s", updateSettingsResp.StatusCode(), updateSettingsResp.Body)
			}
		}
	}

	// Update API Discovery settings
	if d.HasChange("api_discovery_settings") {
		if apiDiscoverySettings, ok := d.GetOk("api_discovery_settings"); ok {
			// Check if API Discovery feature is available
			hasFeature, diags := hasApiDiscoveryFeature(ctx, client)
			if diags.HasError() {
				return diags
			}

			if !hasFeature {
				return diag.Errorf("API Discovery feature is not available for your account. Please contact support or upgrade your plan to use api_discovery_settings.")
			}

			updateApiDiscoverySettingsResp, err := updateApiDiscoverySettings(ctx, client, apiDiscoverySettings, domainID)

			if err != nil {
				return diag.Errorf("Failed to update API Discovery settings: %s", err)
			}

			if updateApiDiscoverySettingsResp.StatusCode() != http.StatusOK {
				return diag.Errorf("Failed to update API Discovery settings. Status code: %d with error: %s", updateApiDiscoverySettingsResp.StatusCode(), updateApiDiscoverySettingsResp.Body)
			}
		}
	}

	return resourceWaapDomainRead(ctx, d, m)
}

func resourceWaapDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func findDomainByName(response waap.PaginatedResponseSummaryDomainResponse, name string) *waap.SummaryDomainResponse {
	for _, domain := range response.Results {
		if strings.EqualFold(domain.Name, name) {
			return &domain
		}
	}
	return nil
}

func hasApiDiscoveryFeature(ctx context.Context, client *waap.ClientWithResponses) (bool, diag.Diagnostics) {
	// Check if the client has the API Discovery feature. This needs to be checked during runtime
	// because the feature availability depends on the account plan.
	clientResp, err := client.GetClientInfoV1ClientsMeGetWithResponse(ctx)
	if err != nil {
		return false, diag.Errorf("Failed to check API Discovery feature availability: %s", err)
	}

	if clientResp.StatusCode() != http.StatusOK {
		return false, diag.Errorf("Failed to check API Discovery feature availability. Status code: %d with error: %s", clientResp.StatusCode(), clientResp.Body)
	}

	return slices.Contains(clientResp.JSON200.Features, "api-discovery"), nil
}

func updateDomainSettings(ctx context.Context, waapClient *waap.ClientWithResponses, settings any, domainID int) (*waap.UpdateDomainSettingsV1DomainsDomainIdSettingsPatchResponse, error) {
	settingsList := settings.([]interface{})

	if len(settingsList) <= 0 || settingsList[0] == nil {
		return nil, nil
	}

	settingsMap := settingsList[0].(map[string]interface{})
	var updateReq waap.UpdateDomainSettings

	// Process DDOS settings
	if ddosList, ok := settingsMap["ddos"].([]interface{}); ok && len(ddosList) > 0 {
		ddosMap := ddosList[0].(map[string]interface{})
		ddosSettings := waap.UpdateDomainDdosSettings{}

		if v, ok := ddosMap["global_threshold"]; ok {
			val := v.(int)
			ddosSettings.GlobalThreshold = &val
		}

		if v, ok := ddosMap["burst_threshold"]; ok {
			val := v.(int)
			ddosSettings.BurstThreshold = &val
		}

		updateReq.Ddos = &ddosSettings
	}

	// Process API settings
	if apiList, ok := settingsMap["api"].([]interface{}); ok && len(apiList) > 0 {
		apiMap := apiList[0].(map[string]interface{})
		apiSettings := waap.UpdateApiSettings{}

		if v, ok := apiMap["api_urls"].(*schema.Set); ok {
			urls := convertSchemaSetToStringList(v)
			apiSettings.ApiUrls = &urls
		}

		if v, ok := apiMap["is_api"]; ok {
			val := v.(bool)
			apiSettings.IsApi = &val
		}

		updateReq.Api = &apiSettings
	}

	// Update domain settings
	return waapClient.UpdateDomainSettingsV1DomainsDomainIdSettingsPatchWithResponse(ctx, domainID, updateReq)
}

func updateApiDiscoverySettings(ctx context.Context, waapClient *waap.ClientWithResponses, apiDiscoverySettings any, domainID int) (*waap.UpdateApiDiscoverySettingsV1DomainsDomainIdApiDiscoverySettingsPatchResponse, error) {
	settingsList := apiDiscoverySettings.([]interface{})

	if len(settingsList) <= 0 || settingsList[0] == nil {
		return nil, nil
	}

	settingsMap := settingsList[0].(map[string]interface{})

	var descriptionFileLocation *string
	if v, exists := settingsMap["description_file_location"]; exists {
		value := v.(string)
		descriptionFileLocation = &value
	}

	var descriptionFileScanEnabled *bool
	if v, exists := settingsMap["description_file_scan_enabled"]; exists {
		value := v.(bool)
		descriptionFileScanEnabled = &value
	}

	var descriptionFileScanIntervalHours *int
	if v, exists := settingsMap["description_file_scan_interval_hours"]; exists {
		value := v.(int)
		descriptionFileScanIntervalHours = &value
	}

	var trafficScanEnabled *bool
	if v, exists := settingsMap["traffic_scan_enabled"]; exists {
		value := v.(bool)
		trafficScanEnabled = &value
	}

	var trafficScanIntervalHours *int
	if v, exists := settingsMap["traffic_scan_interval_hours"]; exists {
		value := v.(int)
		trafficScanIntervalHours = &value
	}

	updateReq := waap.UpdateApiDiscoverySettings{
		DescriptionFileLocation:          descriptionFileLocation,
		DescriptionFileScanEnabled:       descriptionFileScanEnabled,
		DescriptionFileScanIntervalHours: descriptionFileScanIntervalHours,
		TrafficScanEnabled:               trafficScanEnabled,
		TrafficScanIntervalHours:         trafficScanIntervalHours,
	}

	return waapClient.UpdateApiDiscoverySettingsV1DomainsDomainIdApiDiscoverySettingsPatchWithResponse(ctx, domainID, updateReq)
}
