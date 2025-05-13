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
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the domain.",
			},
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
			"policies": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: generatePolicySchema(),
				},
			},
		},
	}
}

func resourceWaapDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).WaapClient
	domainName := d.Get("name").(string)

	resp, err := client.GetDomainsV1DomainsGetWithResponse(
		context.Background(),
		nil,
		func(ctx context.Context, req *http.Request) error {
			req.URL.RawQuery = "name=" + domainName
			return nil
		},
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing domains: %v", err))
	}

	if resp.JSON200 != nil {
		if resp.JSON200.Count != 0 {
			id := fmt.Sprintf("%v", resp.JSON200.Results[0].Id)
			d.SetId(id)
		} else {
			return diag.FromErr(fmt.Errorf("domain with name '%s' not found", domainName))
		}
	}

	return resourceWaapDomainUpdate(ctx, d, m)
}

func resourceWaapDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Config).WaapClient

	domainID, _ := strconv.Atoi(d.Get("id").(string))

	// Get domain details
	resp, err := client.GetDomainV1DomainsDomainIdGetWithResponse(
		context.Background(),
		domainID,
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting domain details: %v", err))
	}

	if resp.JSON200 != nil {
		_ = d.Set("status", string(resp.JSON200.Status))

		// Get domain settings
		settingsResp, err := client.GetDomainSettingsV1DomainsDomainIdSettingsGetWithResponse(
			context.Background(),
			domainID,
		)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting domain settings: %v", err))
		}

		if settingsResp.JSON200 != nil {
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
				_ = d.Set("settings", []interface{}{settings})
			}
		}
	}

	// Get domain policies
	policiesResp, err := client.GetRuleSetListV1DomainsDomainIdRuleSetsGetWithResponse(
		context.Background(),
		domainID,
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting domain policies: %v", err))
	}

	if policiesResp.JSON200 != nil {
		policies := make(map[string]interface{})
		for resourceSlug := range policiesMap {
			policies[resourceSlug] = []interface{}{}
		}

		for _, policy := range *policiesResp.JSON200 {
			resourceSlug := strings.ReplaceAll(*policy.ResourceSlug, "-", "_")
			ruleValues := make(map[string]interface{})

			for _, rule := range *policy.Rules {
				for ruleName, ruleID := range policiesMap[resourceSlug] {
					if ruleID == rule.Id {
						ruleValues[ruleName] = strconv.FormatBool(rule.Mode)
						break
					}
				}
			}
			if len(ruleValues) > 0 {
				policies[resourceSlug] = []interface{}{ruleValues}
			}
		}

		_ = d.Set("policies", []interface{}{policies})
	}

	return diags
}

func resourceWaapDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).WaapClient
	domainID, _ := strconv.Atoi(d.Get("id").(string))

	// Get domain details
	resp, err := client.GetDomainV1DomainsDomainIdGetWithResponse(
		context.Background(),
		domainID,
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting domain details: %v", err))
	}

	var domainStatus string
	if resp.JSON200 != nil {
		domainStatus = string(resp.JSON200.Status)

		if newStatus, ok := d.GetOk("status"); ok && newStatus != domainStatus {
			domainStatusUpdate := waap.DomainUpdateStatus(newStatus.(string))
			updateRequest := waap.UpdateDomainV1DomainsDomainIdPatchJSONRequestBody{
				Status: &domainStatusUpdate,
			}

			// Update domain status
			updateResp, err := client.UpdateDomainV1DomainsDomainIdPatchWithResponse(
				context.Background(),
				domainID,
				updateRequest,
			)

			if err != nil || updateResp.StatusCode() != 204 {
				return diag.FromErr(fmt.Errorf("failed to update domain. Status code: %d with error: %v", updateResp.StatusCode(), err))
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

						domainSettingsUpdate.Ddos = &waap.UpdateDomainDdosSettings{
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

							domainSettingsUpdate.Api = &waap.AppModelsDomainSettingsUpdateApiUrls{
								ApiUrls: &urls,
							}
						}
					}

					// Update domain settings
					updateSettingsResp, err := client.UpdateDomainSettingsV1DomainsDomainIdSettingsPatchWithResponse(
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

	// Process policies
	definedPolicies := make(map[string]map[string]string)
	if v, ok := d.GetOk("policies"); ok {
		rawList := v.([]interface{})
		if len(rawList) > 0 && rawList[0] != nil {
			rawMap := rawList[0].(map[string]interface{})
			for resourceSlug, rulesMap := range rawMap {
				definedPolicies[resourceSlug] = make(map[string]string)
				rulesList := rulesMap.([]interface{})
				if len(rulesList) > 0 && rulesList[0] != nil {
					rules := rulesList[0].(map[string]interface{})
					for ruleName, ruleValue := range rules {
						definedPolicies[resourceSlug][ruleName] = ruleValue.(string)
					}
				}
			}
		}
	}

	// Get domain policies
	policiesResp, err := client.GetRuleSetListV1DomainsDomainIdRuleSetsGetWithResponse(
		context.Background(),
		domainID,
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting domain policies: %v", err))
	}

	if policiesResp.JSON200 != nil {
		policies := make(map[string]interface{})
		for resourceSlug := range policiesMap {
			policies[resourceSlug] = []interface{}{}
		}

		for _, policy := range *policiesResp.JSON200 {
			resourceSlug := strings.ReplaceAll(*policy.ResourceSlug, "-", "_")
			ruleValues := make(map[string]interface{})

			for _, rule := range *policy.Rules {
				for ruleName, ruleID := range policiesMap[resourceSlug] {
					if ruleID == rule.Id {
						ruleValues[ruleName] = strconv.FormatBool(rule.Mode)
						break
					}
				}
			}
			if len(ruleValues) > 0 {
				policies[resourceSlug] = []interface{}{ruleValues}
			}
		}

		// Compare existing policies with defined policies
		for resourceSlug, rules := range definedPolicies {
			existingRulesList, ok := policies[resourceSlug].([]interface{})
			if !ok || len(existingRulesList) == 0 {
				continue
			}

			existingRules := existingRulesList[0].(map[string]interface{})

			for ruleName, ruleValue := range rules {
				existingRuleValue, exists := existingRules[ruleName]
				if !exists {
					continue
				}

				if existingRuleValue != ruleValue {
					ruleID, ok := policiesMap[resourceSlug][ruleName].(string)
					if !ok {
						continue
					}

					updateResp, err := client.ToggleDomainPolicyV1DomainsDomainIdPoliciesPolicyIdTogglePatchWithResponse(
						context.Background(),
						domainID,
						ruleID,
					)

					if err != nil || updateResp.StatusCode() != 200 {
						return diag.FromErr(fmt.Errorf("failed to update policy rule %s. Status code: %d with error: %v", ruleName, updateResp.StatusCode(), err))
					}
				}
			}
		}
	}

	return resourceWaapDomainRead(ctx, d, m)
}

func resourceWaapDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func generatePolicySchema() map[string]*schema.Schema {
	policySchema := make(map[string]*schema.Schema)

	for resourceSlug, rules := range policiesMap {
		ruleSchema := make(map[string]*schema.Schema)

		for ruleName := range rules {
			ruleSchema[ruleName] = &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: fmt.Sprintf("Enable %s rule", ruleName),
			}
		}

		policySchema[resourceSlug] = &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: ruleSchema,
			},
		}
	}

	return policySchema
}

func replaceHyphenWithUnderscore(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}

var policiesMap = map[string]map[string]interface{}{
	"protocol_validation": {
		"invalid_user_agent":                "S76363426",
		"unknown_user_agent":                "S76363427",
		"service_protocol_validation":       "S76363428",
		"prevent_malformed_request_methods": "S76363429",
	},
	"core_waf_owasp_top_threats": {
		"sql_injection":                          "S76363430",
		"xss":                                    "S76363431",
		"shellshock_exploit":                     "S76363432",
		"rfi":                                    "S76363433",
		"apache_struts_exploit":                  "S76363434",
		"lfi":                                    "S76363435",
		"common_web_application_vulnerabilities": "S76363436",
		"web_shell_execution_prevention":         "S76363437",
		"protocol_attack":                        "S76363438",
		"csrf":                                   "S76363439",
		"open_redirect":                          "S76363440",
		"shell_injection":                        "S76363441",
		"code_injection":                         "S76363442",
		"sensitive_data_exposure":                "S76363443",
		"xxe":                                    "S76363444",
		"personally_identifiable_information":    "S76363445",
		"server_side_template_injection":         "S76363446",
	},
	"ip_reputation": {
		"traffic_via_tor_network":            "S76363447",
		"traffic_via_proxy_networks":         "S76363448",
		"traffic_from_hosting_services":      "S76363449",
		"traffic_via_vpns":                   "S76363450",
		"bot_traffic":                        "S76363451",
		"traffic_from_suspicious_nat_ranges": "S76363452",
		"external_reputation_block_list":     "S76363453",
		"traffic_via_cdns":                   "S76363454",
	},
	"behavioral_waf": {
		"anti_spam":                                  "S76363455",
		"probing_and_forced_browsing":                "S76363456",
		"obfuscated_attacks_and_zero_day_mitigation": "S76363457",
		"repeated_violations":                        "S76363458",
		"brute_force_protection":                     "S76363459",
	},
	"anti_automation_bot_protection": {
		"traffic_anomaly":   "S76363460",
		"automated_clients": "S76363461",
		"headless_browsers": "S76363462",
		"anti_scraping":     "S76363463",
	},
	"cms_protection": {
		"wordpress_waf_ruleset":              "S76363464",
		"logged_in_wordpress_admins":         "S76363465",
		"logged_in_modx_admins":              "S76363466",
		"logged_in_drupal_admins":            "S76363467",
		"logged_in_joomla_admins":            "S76363468",
		"logged_in_allowlist_magento_admins": "S76363469",
		"logged_in_umbraco_admins":           "S76363470",
		"logged_in_pimcore_admins":           "S76363471",
	},
	"common_automated_services": {
		"microsoft_msn_bot":                     "S76363472",
		"microsoft_bing_bot":                    "S76363473",
		"facebook_external_hit_bot":             "S76363474",
		"twitter_bot":                           "S76363475",
		"yahoo_inktomi_slurp_bot":               "S76363476",
		"yahoo_slurp_bot":                       "S76363477",
		"yandex_bot":                            "S76363478",
		"baidu_spider_bot":                      "S76363479",
		"baidu_spider_japan_bot":                "S76363480",
		"naver_yeti_bot":                        "S76363481",
		"seznam_bot":                            "S76363482",
		"blekko_scoutjet_bot":                   "S76363483",
		"ask_jeeves_bot":                        "S76363484",
		"linkedin_bot":                          "S76363485",
		"alexa_ia_archiver":                     "S76363486",
		"vkontakte_external_hit_bot":            "S76363487",
		"soso_spider_bot":                       "S76363488",
		"yodao_bot":                             "S76363489",
		"sogou_bot":                             "S76363490",
		"jikespider_bot":                        "S76363491",
		"yahoo_seeker_bot":                      "S76363492",
		"pingdom":                               "S76363493",
		"sitelock_spider":                       "S76363494",
		"new_relic_bot":                         "S76363495",
		"applebot":                              "S76363496",
		"gomez":                                 "S76363497",
		"chrome_compression_proxy":              "S76363498",
		"kakao_useragent":                       "S76363499",
		"yahoo_link_preview":                    "S76363500",
		"daumoa_bot":                            "S76363501",
		"yahoo_japan_bot":                       "S76363502",
		"goo_japan_bot":                         "S76363503",
		"jword_japan_bot":                       "S76363504",
		"line_japan_bot":                        "S76363505",
		"mobage_japan_bot":                      "S76363506",
		"mixi_japan_bot":                        "S76363507",
		"gree_japan_bot":                        "S76363508",
		"zoho_bot":                              "S76363509",
		"ocn_japan_bot":                         "S76363510",
		"livedoor_japan_bot":                    "S76363511",
		"microsoft_skype_bot":                   "S76363512",
		"paypal_ipn":                            "S76363513",
		"hipay":                                 "S76363514",
		"statuscake_bot":                        "S76363515",
		"thefind_crawler":                       "S76363516",
		"cybersource":                           "S76363517",
		"ias_crawler":                           "S76363518",
		"yisouspider":                           "S76363519",
		"dotmic_dotbot":                         "S76363520",
		"coccocbot":                             "S76363521",
		"microsoft_bing_preview_bot":            "S76363522",
		"qwantify_bot":                          "S76363523",
		"slack_bot":                             "S76363524",
		"uptime_robot":                          "S76363525",
		"panopta_bot":                           "S76363526",
		"server_density_service_monitoring_bot": "S76363527",
		"sagepay":                               "S76363528",
		"zum_bot":                               "S76363529",
		"ahrefs_bot":                            "S76363530",
		"requests_from_origins_ip":              "S76363531",
		"sucuri_uptime_monitor_bot":             "S76363532",
		"semrush_bot":                           "S76363533",
		"mail_ru_bot":                           "S76363534",
		"telegram_bot":                          "S76363535",
		"internet_archive_bot":                  "S76363536",
		"stripe":                                "S76363537",
		"pinterest_bot":                         "S76363538",
		"jetpack_bot":                           "S76363539",
		"alerta_bot":                            "S76363540",
		"hyperspin_bot":                         "S76363541",
		"bitbucket_webhook":                     "S76363542",
		"managewp":                              "S76363543",
		"zendesk_bot":                           "S76363544",
		"amazon_route53_health_check_service":   "S76363545",
		"lets_encrypt":                          "S76363546",
		"hetrix_tools":                          "S76363547",
		"alexa_technologies":                    "S76363548",
		"addsearch_bot":                         "S76363549",
		"site24x7_bot":                          "S76363550",
		"wordfence_central":                     "S76363551",
		"xml_sitemaps":                          "S76363552",
		"applenewsbot":                          "S76363553",
		"roger_bot":                             "S76363554",
		"duckduckgo_bot":                        "S76363555",
		"cookiebot":                             "S76363556",
		"detectify_scanner":                     "S76363557",
		"digicert_dcv_bot":                      "S76363558",
		"workato":                               "S76363559",
		"ghostinspector":                        "S76363560",
		"freshping_monitoring":                  "S76363561",
		"binarycanary":                          "S76363562",
		"adestra_bot":                           "S76363563",
		"acquia_uptime":                         "S76363564",
		"spring_bot":                            "S76363565",
		"parse_ly_scraper":                      "S76363566",
		"landau_media_spider":                   "S76363567",
		"geckoboard":                            "S76363568",
		"audisto_bot":                           "S76363569",
		"feedpress":                             "S76363570",
		"feeder_co":                             "S76363571",
		"shareaholic_bot":                       "S76363572",
		"adjust_servers":                        "S76363573",
		"kyoto_tohoku_crawler":                  "S76363574",
		"spatineo":                              "S76363575",
		"w3c":                                   "S76363576",
		"stackify":                              "S76363577",
		"sectigo_bot":                           "S76363578",
		"testomato_bot":                         "S76363579",
		"siteimprove_bot":                       "S76363580",
		"petal_bot":                             "S76363581",
		"google_cloud_monitoring":               "S76363582",
		"smart_plugin_manager_bot":              "S76363583",
		"outbrain_bot":                          "S76363584",
		"comscore_crawler":                      "S76363585",
		"google_bot":                            "S76363586",
		"google_services":                       "S76363587",
		"google_crawler":                        "S76363588",
		"google_user_triggered_fetchers":        "S76363589",
		"apple_private_relay":                   "S76363590",
		"grafana_services":                      "S76363591",
		"test_domain_static_rule":               "S76363592",
	},
	"advanced_api_protection": {
		"auth_token_protection":      "S76363593",
		"sensitive_data_exposure":    "S76363594",
		"invalid_api_traffic":        "S76363595",
		"api_level_authorization":    "S76363596",
		"non_baselined_api_requests": "S76363597",
	},
}
