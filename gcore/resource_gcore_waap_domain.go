package gcore

import (
	"context"
	"fmt"
	"net/http"
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
										Description: "Global threshold for DDoS protection",
									},
									"burst_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Burst threshold for DDoS protection",
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
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of API URL patterns",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
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
			return diag.Errorf("Failed to update Domain status: %w", err)
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
			return diag.Errorf("Failed to update Domain settings: %w", err)
		}

		if updateSettingsResp.StatusCode() != http.StatusNoContent {
			return diag.Errorf("Failed to update Domain settings. Status code: %d with error: %s", updateSettingsResp.StatusCode(), updateSettingsResp.Body)
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
	resp, err := client.GetDomainV1DomainsDomainIdGetWithResponse(
		context.Background(),
		domainID,
	)
	if err != nil {
		return diag.Errorf("Failed to read Domain details: %w", err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		d.SetId("") // Resource not found, remove from state
		return diag.Diagnostics{
			{Summary: fmt.Sprintf("Domain (%s) was not found, removed from TF state", d.Id())},
		}
	}

	if resp.StatusCode() != http.StatusOK {
		return diag.Errorf("Failed to read Domain details. Status code: %d with error: %s", resp.StatusCode(), resp.Body)
	}

	d.Set("status", string(resp.JSON200.Status))

		// Get domain settings
	settingsResp, err := client.GetDomainSettingsV1DomainsDomainIdSettingsGetWithResponse(ctx, domainID)
		if err != nil {
		return diag.Errorf("Failed to read Domain settings: %w", err)
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

				if len(apiSettings) > 0 {
					settings["api"] = []interface{}{apiSettings}
				}
			}

			if len(settings) > 0 {
				d.Set("settings", []interface{}{settings})
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
				return diag.Errorf("Failed to update Domain status: %w", err)
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
				return diag.Errorf("Failed to update Domain settings: %w", err)
			}

			if updateSettingsResp.StatusCode() != http.StatusNoContent {
				return diag.Errorf("Failed to update Domain settings. Status code: %d with error: %s", updateSettingsResp.StatusCode(), updateSettingsResp.Body)
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
		ddosSettings := struct {
			GlobalThreshold *int `json:"global_threshold,omitempty"`
			BurstThreshold  *int `json:"burst_threshold,omitempty"`
		}{}

		if v, ok := ddosMap["global_threshold"]; ok {
			val := v.(int)
			ddosSettings.GlobalThreshold = &val
		}

		if v, ok := ddosMap["burst_threshold"]; ok {
			val := v.(int)
			ddosSettings.BurstThreshold = &val
		}

		updateReq.Ddos = &waap.UpdateDomainDdosSettings{
			GlobalThreshold: ddosSettings.GlobalThreshold,
			BurstThreshold:  ddosSettings.BurstThreshold,
		}
	}

	// Process API settings
	if apiList, ok := settingsMap["api"].([]interface{}); ok && len(apiList) > 0 {
		apiMap := apiList[0].(map[string]interface{})

		if apiUrls, ok := apiMap["api_urls"].([]interface{}); ok {
			urls := make([]string, len(apiUrls))
			for i, url := range apiUrls {
				urls[i] = url.(string)
			}

			updateReq.Api = &waap.AppModelsDomainSettingsUpdateApiUrls{
				ApiUrls: &urls,
			}
		}
	}

	// Update domain settings
	return waapClient.UpdateDomainSettingsV1DomainsDomainIdSettingsPatchWithResponse(ctx, domainID, updateReq)
}
