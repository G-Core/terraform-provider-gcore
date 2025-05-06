package gcore

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	waap "github.com/G-Core/gcore-waap-sdk-go"
)

func resourceWaapDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGcoreDomainCreate,
		ReadContext:   resourceGcoreDomainRead,
		UpdateContext: resourceGcoreDomainUpdate,
		DeleteContext: resourceGcoreDomainDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the domain.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the domain.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Status of the domain. If not provided, the current status will be read from the API.",
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
									"sub_second_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Sub-second threshold for DDoS protection",
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

func resourceGcoreDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceGcoreDomainUpdate(ctx, d, m)
}

func resourceGcoreDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var domainID int
	domainName := d.Get("name").(string)

	authToken := m.(*Config).Provider.APIToken
	apiURL := m.(*Config).Provider.IdentityEndpoint

	// Create a new client
	client, err := waap.NewClient(
		apiURL,
		waap.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Add("Authorization", authToken)
			return nil
		}),
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating client: %v", err))
	}

	clientWithResponses := &waap.ClientWithResponses{ClientInterface: client}

	// List domains
	resp, err := clientWithResponses.GetDomainsV1DomainsGetWithResponse(context.Background(), &waap.GetDomainsV1DomainsGetParams{})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing domains: %v", err))
	}

	if resp.JSON200 != nil {
		found := false
		for _, domain := range resp.JSON200.Results {
			id := fmt.Sprintf("%v", domain.Id)
			name := domain.Name

			if name == domainName {
				d.SetId(id)
				_ = d.Set("status", domain.Status)
				domainID = domain.Id
				found = true
				break
			}
		}

		if !found {
			d.SetId("")
			return diag.FromErr(fmt.Errorf("domain with name '%s' not found", domainName))
		}

		// Get domain settings
		settingsResp, err := clientWithResponses.GetDomainSettingsV1DomainsDomainIdSettingsGetWithResponse(
			context.Background(),
			domainID,
		)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting domain settings: %v", err))
		}

		if settingsResp.JSON200 != nil {
			settings := make(map[string]interface{})

			if settingsResp.JSON200 != nil {
				ddosSettings := make(map[string]interface{})

				if settingsResp.JSON200.Ddos.GlobalThreshold != nil {
					ddosSettings["global_threshold"] = *settingsResp.JSON200.Ddos.GlobalThreshold
				}

				if settingsResp.JSON200.Ddos.BurstThreshold != nil {
					ddosSettings["burst_threshold"] = *settingsResp.JSON200.Ddos.BurstThreshold
				}

				if settingsResp.JSON200.Ddos.SubSecondThreshold != nil {
					ddosSettings["sub_second_threshold"] = *settingsResp.JSON200.Ddos.SubSecondThreshold
				}

				if len(ddosSettings) > 0 {
					settings["ddos"] = []interface{}{ddosSettings}
				}
			}

			if settingsResp.JSON200 != nil && settingsResp.JSON200.Api.ApiUrls != nil {
				apiSettings := make(map[string]interface{})
				apiSettings["api_urls"] = *settingsResp.JSON200.Api.ApiUrls

				if len(apiSettings) > 0 {
					settings["api"] = []interface{}{apiSettings}
				}
			}

			if len(settings) > 0 {
				_ = d.Set("settings", []interface{}{settings})
			}
		}
	}
	return diags
}

func resourceGcoreDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	domainName := d.Get("name").(string)
	authToken := m.(*Config).Provider.APIToken
	apiURL := m.(*Config).Provider.IdentityEndpoint

	client, err := waap.NewClient(
		apiURL,
		waap.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Add("Authorization", authToken)
			return nil
		}),
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating client: %v", err))
	}

	clientWithResponses := &waap.ClientWithResponses{ClientInterface: client}

	// List domains
	resp, err := clientWithResponses.GetDomainsV1DomainsGetWithResponse(context.Background(), &waap.GetDomainsV1DomainsGetParams{})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing domains: %v", err))
	}

	var domainID int
	var domainStatus string
	if resp.JSON200 != nil {

		for _, domain := range resp.JSON200.Results {
			if domain.Name == domainName {
				id := fmt.Sprintf("%v", domain.Id)
				domainID = domain.Id
				domainStatus = string(domain.Status)
				d.SetId(id)
				break
			}
		}

		if domainID == 0 {
			return diag.FromErr(fmt.Errorf("domain with name '%s' not found", domainName))
		}

		newStatus := d.Get("status")
		if newStatus == "protect" {
			newStatus = "active"
		}

		if newStatus == "active" || newStatus == "monitor" {
			if newStatus != domainStatus {
				newStatusStr, _ := newStatus.(waap.DomainUpdateStatus)
				updateRequest := waap.UpdateDomainV1DomainsDomainIdPatchJSONRequestBody{
					Status: &struct {
						waap.DomainUpdateStatus `yaml:",inline"`
					}{
						DomainUpdateStatus: newStatusStr,
					},
				}

				// Update domain status
				updateResp, err := clientWithResponses.UpdateDomainV1DomainsDomainIdPatchWithResponse(
					context.Background(),
					domainID,
					updateRequest,
				)

				if err != nil || updateResp.StatusCode() != 204 {
					return diag.FromErr(fmt.Errorf("failed to update domain. Status code: %d with error: %v", updateResp.StatusCode(), err))
				}
			}
		}

		if d.HasChange("settings") {
			var domainSettingsUpdate waap.UpdateDomainSettingsV1DomainsDomainIdSettingsPatchJSONRequestBody

			if v, ok := d.GetOk("settings"); ok {
				settingsList := v.([]interface{})
				if len(settingsList) > 0 && settingsList[0] != nil {
					settingsMap := settingsList[0].(map[string]interface{})

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

						domainSettingsUpdate.Ddos = &struct {
							waap.UpdateDomainDdosSettings `yaml:",inline"`
						}{
							UpdateDomainDdosSettings: waap.UpdateDomainDdosSettings{
								GlobalThreshold: ddosSettings.GlobalThreshold,
								BurstThreshold:  ddosSettings.BurstThreshold,
							},
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

							domainSettingsUpdate.Api = &struct {
								waap.AppModelsDomainSettingsUpdateApiUrls `yaml:",inline"`
							}{
								AppModelsDomainSettingsUpdateApiUrls: waap.AppModelsDomainSettingsUpdateApiUrls{
									ApiUrls: &urls,
								},
							}
						}
					}

					// Update domain settings
					updateSettingsResp, err := clientWithResponses.UpdateDomainSettingsV1DomainsDomainIdSettingsPatchWithResponse(
						context.Background(),
						domainID,
						domainSettingsUpdate,
					)

					if err != nil || updateSettingsResp.StatusCode() != 204 {
						return diag.FromErr(fmt.Errorf("failed to update domain settings. Status code: %d with error: %v", updateSettingsResp.StatusCode(), err))
					}
				}
			}
		}
	}
	return resourceGcoreDomainRead(ctx, d, m)
}

func resourceGcoreDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
